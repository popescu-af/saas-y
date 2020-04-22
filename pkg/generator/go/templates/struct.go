package templates

// Struct is a template for API structs.
const Struct = `package structs

// {{.Name | capitalize}} - generated API structure
type {{.Name | capitalize}} struct {
	{{range .Fields}}{{.Name | capitalize}} {{.Type}} ` + "`" + `json:"{{.Name}}"` + "`" + `
	{{end}}
}`
