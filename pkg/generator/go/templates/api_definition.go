package templates

// APIDefinition is the template for the go definition of the API.
const APIDefinition = `package logic

import "{{.Name}}/pkg/structs"

// API defines the operations supported by the {{.Name}} service.
type API interface {
	{{range $a := .API}}// {{$a.Path}}
	{{range $mname, $method := $a.Methods}}{{$mname | capitalize | symbolize}}(
	{{if $method.InputType}}*structs.{{$method.InputType | capitalize | symbolize}},
	{{end}}{{if $a.Path | pathHasParameters}}{{with $params := $a.Path | pathParameters}}{{range $pnameidx := $params | indicesParameters}}{{with $ptypeidx := inc $pnameidx}}{{index $params $ptypeidx | typeName}},
	{{end}}{{end}}{{end}}{{end}}{{if $method.HeaderParams}}{{range $method.HeaderParams}}{{.Type | typeName}},
	{{end}}{{end}}{{if $method.QueryParams}}{{range $method.QueryParams}}{{.Type | typeName}},
	{{end}}{{end}}) (*structs.{{$method.ReturnType | capitalize | symbolize}}, error)

	{{end}}{{end}}
}`
