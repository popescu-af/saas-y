package templates

// APIExample is a template for an example of go implementation of the API.
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

{{range .API}}// {{.Path}}
{{range $methodName, $method := .Methods}}
{{with $mname := $methodName | symbolName}}
// {{$mname}} example.
func (s *API) {{$mname}}({{if $method.InputType}}body structs.{{$method.InputType | symbolName}}{{end}}) (interface{}, error) {
	{{if $method.InputType}}s.logger.Info("{{$mname}}", zap.Any("body", body))
	{{end}}return nil, errors.New("method '{{$mname}}' not implemented")
}{{end}}
{{end}}
{{end}}`
