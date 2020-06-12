package templates

// APIExample is the template for an example of go implementation of the API.
const APIExample = `package logic

import (
	"errors"

	"github.com/popescu-af/saas-y/pkg/logutil"

	"{{.RepositoryURL}}/pkg/exports"
)

// ExampleAPI is an example, trivial implementation of the API interface.
// It simply logs the request name.
type ExampleAPI struct {
}

// NewAPI creates an instance of the example API implementation.
func NewAPI() exports.API {
	return &ExampleAPI{}
}

{{range $a := .API}}
	// {{$a.Path}}
	{{range $mname, $method := $a.Methods}}
		// {{$mname | capitalize}} example.
		func (a *ExampleAPI) {{$mname | capitalize}}(
			{{- if $method.InputType -}}
				input *exports.{{$method.InputType | capitalize | symbolize}},
			{{- end -}}
			{{- if $a.Path | pathHasParameters -}}
				{{- with $params := $a.Path | pathParameters -}}
					{{- range $pnameidx := $params | indicesParameters -}}
						{{- index $params $pnameidx}} {{with $ptypeidx := inc $pnameidx}}{{index $params $ptypeidx | typeName}},{{end}}
					{{- end -}}
				{{- end -}}
			{{- end -}}
			{{- if $method.QueryParams -}}
				{{- range $method.QueryParams -}}
					{{- .Name}} {{.Type | typeName}},
				{{- end -}}
			{{- end -}}
			{{- if $method.HeaderParams -}}
				{{- range $method.HeaderParams -}}
					{{- .Name}} {{.Type | typeName}},
				{{- end -}}
			{{- end -}}
		) (*exports.{{$method.ReturnType | capitalize | symbolize}}, error) {
			logutil.Info("called {{$mname}}")
			return nil, errors.New("method '{{$mname}}' not implemented")
		}
	{{- end -}}
{{- end -}}`
