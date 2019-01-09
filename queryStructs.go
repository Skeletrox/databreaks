package databreaks

/*FieldPair is a struct that holds a Function(Column) Pair*/
type FieldPair struct {
	Function string `json:"function"`
	Column   string `json:"column"`
}

/*Query is a struct that holds a query
  The length of operators is always len(Fields) - 1
*/
type Query struct {
	Measurement string      `json:"measurement"`
	Fields      []FieldPair `json:"fields"`
	Operators   []string    `json:"operators"`
}

/*A FieldComposite is a collection of FieldBranches
 */
type FieldComposite struct {
	FieldBranches []FieldBranch `json:"field_branches"`
}

/*A FieldBranch is a single column in the output that might be a single column from the DB or a composite of multiple, with an optional alias
 */
type FieldBranch struct {
	FieldUnits []FieldCompositeUnit `json:"field_units"`
	Operators  []string             `json:"operators"`
	Alias      string               `json:"alias"`
}

/*A FieldCompositeUnit is a single unit that contains a function and a column name
 */
type FieldCompositeUnit struct {
	Function string `json:"function"`
	Column   string `json:"column"`
}

/*A SQLQuery is a machine that decomposes a SQL query into language-independent constituents
 */
type SQLQuery struct {
	Fields      FieldComposite `json:"fields"`
	Measurement string         `json:"measurement"`
}

/*A SQLParser parses the Query from/to the Lang*/
type SQLParser struct {
	Query SQLQuery
	Lang  string
}
