package templates

// Struct is the template for API structures in go code.
const Struct = `package structs

// {{.Name | capitalize}} - generated API structure
type {{.Name | capitalize}} struct {
	{{range .Fields}}{{.Name | capitalize}} {{.Type | typeName}} ` + "`" + `json:"{{.Name}}"` + "`" + `
	{{end}}
}`
