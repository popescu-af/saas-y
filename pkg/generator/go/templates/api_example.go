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

{{range .API}}// {{.Path}}
{{range $mname, $method := .Methods}}
// {{$mname | capitalize}} example.
func (a *API) {{$mname | capitalize}}({{if $method.InputType}}body structs.{{$method.InputType | capitalize | symbolize}}{{end}}) (interface{}, error) {
	{{if $method.InputType}}a.logger.Info("{{$mname}}", zap.Any("body", body))
	{{end}}return nil, errors.New("method '{{$mname}}' not implemented")
}{{end}}
{{end}}`
