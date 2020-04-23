package templates

// APIExample is the template for an example of go implementation of the API.
const APIExample = `package example

import (
	"errors"

	"go.uber.org/zap"

	"{{.Name}}/pkg/service"
	"{{.Name}}/pkg/structs"
)

// API is an example implementation of the API interface.
// It simply logs the request body.
type API struct {
	logger *zap.Logger
}

// NewAPI creates an instance of the example API implementation.
func NewAPI(logger *zap.Logger) service.API {
	return &API{logger: logger}
}

{{range $a := .API}}// {{$a.Path}}
{{range $mname, $method := $a.Methods}}{{if eq $method.Type "options"}}
{{/* NADA for options method */}}
{{else}}
// {{$mname | capitalize}} example.
func (a *API) {{$mname | capitalize}}(
{{if $method.InputType}}*structs.{{$method.InputType | capitalize | symbolize}},
{{end}}{{if $a.Path | pathHasParameters}}{{with $params := $a.Path | pathParameters}}{{range $pnameidx := $params | indicesParameters}}{{with $ptypeidx := inc $pnameidx}}{{index $params $ptypeidx | typeName}},
{{end}}{{end}}{{end}}{{end}}{{if $method.HeaderParams}}{{range $method.HeaderParams}}{{.Type | typeName}},
{{end}}{{end}}{{if $method.QueryParams}}{{range $method.QueryParams}}{{.Type | typeName}},
{{end}}{{end}}) (interface{}, error) {
	a.logger.Info("called {{$mname}}")
	return nil, errors.New("method '{{$mname}}' not implemented")
}

{{end}}{{end}}{{end}}`
