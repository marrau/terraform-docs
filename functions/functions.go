package functions

import "html/template"

// Apply adds additional functions to a functionmap for sprig-templates
func Apply(funcMap template.FuncMap) template.FuncMap {
	funcMap["gitUrl"] = gitURL
	return funcMap
}
