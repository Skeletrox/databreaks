package databreaks

import (
	"errors"
	"log"
	"regexp"
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

		var fUnits []FieldCompositeUnit
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
			function := fCMatches[0][1]
			column := fCMatches[0][2]
			fUnit := FieldCompositeUnit{Function: function, Column: column}
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
