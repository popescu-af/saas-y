package templates

// Struct is a template for API structs.
const Struct = `package structs

// {{.Name}} - generated API structure
type {{.Name}} struct {
{{range .Fields}}
	{{.Name}} {{.Type}} ` + "`" + `json:"{{.Name}}"` + "`" +
	`{{end}}
}`
