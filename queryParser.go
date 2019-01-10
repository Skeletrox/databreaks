package databreaks

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

/*DecomposeQuery decomposes a query*/
func DecomposeQuery(query string) (SQLQuery, error) {
	/*
		Assume a query can be decomposed as SELECT << field composite >> FROM << measurement >>
	*/
	var err error
	fieldCompositeMatchRegex := "^(?i)SELECT (.+) (?i)FROM"
	measurementMatchRegex := "FROM (.+)$"
	fieldCompositeMatcher := regexp.MustCompile(fieldCompositeMatchRegex)
	measurementMatcher := regexp.MustCompile(measurementMatchRegex)
	matches := fieldCompositeMatcher.FindStringSubmatch(query)

	if len(matches) == 0 {
		err = errors.New("Cannot parse query: Improper syntax")
		log.Println(err.Error())
		return SQLQuery{}, err
	}

	measurementMatches := measurementMatcher.FindStringSubmatch(query)
	if len(measurementMatches) == 0 {
		err = errors.New("Cannot parse query: No measurement")
		log.Println(err.Error())
		return SQLQuery{}, err
	}

	fieldComposite := matches[1]
	log.Println("Obtained Field Composite:", fieldComposite)

	measurement := measurementMatches[1]
	log.Println("Obtained Measurement:", measurement)

	// A field composite is separated by commas to get field branches
	// A single branch composite has no comma, therefore there is ALWAYS >= 1 match

	fieldCompositeSplitRegex := "[^,]+"
	fieldBranchMatcher := regexp.MustCompile(fieldCompositeSplitRegex)
	branchMatches := fieldBranchMatcher.FindAllStringSubmatch(fieldComposite, -1)

	if len(matches[0]) == 0 {
		err = errors.New("Cannot parse query: Nothing to match")
		log.Println(err.Error())
		return SQLQuery{}, err
	}

	log.Println("Obtained branches:", branchMatches)

	/*
		Each branch can be denoted as a collection of functions and columns, linked to each other over functions and coming together as aliases
		We shall assume that each branch is splittable into multiple fieldUnits, identifiable by an alias
	*/
	fieldUnitMatchRegex := "([^\\+\\-\\*\\/]+\\([^\\+\\-\\*\\/\\s]+\\))"
	operatorMatchRegex := "[\\+\\-\\*\\/]"
	asMatchRegex := "(?i)AS (.+)\\b"
	functionColumnRegex := "(.+)\\((.+)\\)"

	fieldUnitMatcher := regexp.MustCompile(fieldUnitMatchRegex)
	operatorMatcher := regexp.MustCompile(operatorMatchRegex)
	asMatcher := regexp.MustCompile(asMatchRegex)
	funcColMatcher := regexp.MustCompile(functionColumnRegex)

	// Go through each branch
	var fBranches []FieldBranch
	for _, bMatch := range branchMatches {

		operatorMatches := operatorMatcher.FindAllStringSubmatch(bMatch[0], -1)
		fieldUnitMatches := fieldUnitMatcher.FindAllStringSubmatch(bMatch[0], -1)
		asMatches := asMatcher.FindAllStringSubmatch(bMatch[0], -1)

		var fUnits []FieldPair
		/*
			Break each field unit into the function and column
		*/
		for _, fU := range fieldUnitMatches {
			fCMatches := funcColMatcher.FindAllStringSubmatch(fU[0], -1)
			log.Println("For", fU[0], "found:", fCMatches)
			if len(fCMatches[0]) < 3 {
				err = errors.New("Cannot parse query: Improper field syntax")
				log.Println(err.Error())
				return SQLQuery{}, err
			}
			function := strings.Trim(fCMatches[0][1], " \t")
			column := strings.Trim(fCMatches[0][2], " \t")
			fUnit := FieldPair{Function: function, Column: column}
			fUnits = append(fUnits, fUnit)
		}
		var alias string
		var operators []string
		if len(operatorMatches) > 0 {
			operators = operatorMatches[0]
		}
		if len(asMatches) != 0 && len(asMatches[0]) == 2 {
			alias = asMatches[0][1]
		}
		fBranch := FieldBranch{FieldUnits: fUnits, Operators: operators, Alias: alias}
		fBranches = append(fBranches, fBranch)
	}
	fieldCompositeStruct := FieldComposite{FieldBranches: fBranches}
	queryParsed := SQLQuery{Measurement: measurement, Fields: fieldCompositeStruct}
	return queryParsed, nil
}

/*ParseQueryFromLang is a method that allows a SQLParser object to parse a query string and convert it into the appropriate structure based on the language
As of now it ignores the language and just converts from InfluxQL

*/
func (s *SQLParser) ParseQueryFromLang(query string) error {
	var err error
	s.Query, err = DecomposeQuery(query)
	return err
}

/*ParseQueryToLang for a SQL parser is a that composes a query from the struct based on the language.
As if now it ignores the language and just converts to InfluxQL

*/
func (s *SQLParser) ParseQueryToLang() string {
	returnable := s.Query.Stringify()
	return returnable
}

/*Stringify for a FieldPair uses the function and column to create a function(column) pair

 */
func (f *FieldPair) Stringify() string {
	returnable := fmt.Sprintf("%s(%s)", f.Function, f.Column)
	return returnable
}

/*Stringify for a FieldBranch recursively stringifies the fieldPairs, separating them by their operators
 */
func (f *FieldBranch) Stringify() string {
	var returnableBuff bytes.Buffer
	numOps := len(f.Operators)
	for i, fu := range f.FieldUnits {
		returnableBuff.WriteString(fu.Stringify())
		if i < numOps {
			returnableBuff.WriteString(fmt.Sprintf(" %s ", f.Operators[i]))
		}
	}
	if len(f.Alias) > 0 {
		returnableBuff.WriteString(fmt.Sprintf(" AS %s", f.Alias))
	}
	return returnableBuff.String()
}

/*Stringify or a FieldComposite stringifies its FieldBranches and places commas between them*/
func (f *FieldComposite) Stringify() string {
	var returnableBuff bytes.Buffer
	numBranches := len(f.FieldBranches)
	for i, fB := range f.FieldBranches {
		returnableBuff.WriteString(fB.Stringify())
		if i < numBranches-1 {
			returnableBuff.WriteString(", ")
		}
	}
	return returnableBuff.String()
}

/*Stringify for a SQLQuery builds the entire string by manufacturing a legitimate SQL query
 */
func (s *SQLQuery) Stringify() string {
	var returnableBuff bytes.Buffer
	returnableBuff.WriteString("SELECT ")
	returnableBuff.WriteString(s.Fields.Stringify())
	returnableBuff.WriteString(fmt.Sprintf(" FROM %s", s.Measurement))
	return returnableBuff.String()
}
