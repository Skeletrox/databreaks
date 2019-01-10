package databreaks

//RunTests is the official testing function for the package, to ensure things run smoothly.
func RunTests() map[string]interface{} {
	returnable := make(map[string]interface{})
	query := "SELECT MEAN(ensemble), MAX(entropy) + MIN(enchiladas) AS extra_Es FROM E_LETTERS"
	parser := SQLParser{}
	err := parser.ParseQueryFromLang(query)
	returnable["influx_parse"] = parser.Query
	if err != nil {
		returnable["influx_parse"] = err
	}
	parsed := parser.ParseQueryToLang()
	returnable["influx_compose"] = parsed
	return returnable
}
