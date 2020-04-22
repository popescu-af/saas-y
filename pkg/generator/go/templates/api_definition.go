package templates

// APIDefinition is a template for an example of go implementation of the API.
const APIDefinition = `package service

import "{{.Name}}/pkg/structs"

// API defines the operations supported by the {{.Name}} service.
type API interface {
	{{range .API}}// {{.Path}}
	{{range $mname, $method := .Methods}}{{$mname | capitalize | symbolize}}({{if $method.InputType}}body structs.{{$method.InputType | capitalize | symbolize}}{{end}}) (interface{}, error)
	{{end}}
{{end}}}`
