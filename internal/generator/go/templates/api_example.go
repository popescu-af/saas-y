package templates

// APIExample is the template for an example of go implementation of the API.
const APIExample = `package logic

{{range $a := .API}}{{range $mname, $method := $a.Methods}}
	{{$method.Type | print | checkIfWebSocket}}
{{end}}{{end}}

import (
	"errors"
	{{- if eq foundWebSocket "yes"}}
	"time"
	{{- end}}

	"github.com/popescu-af/saas-y/pkg/log"
	{{- if eq foundWebSocket "yes"}}
	"github.com/popescu-af/saas-y/pkg/connection"
	{{- end}}

	"{{.RepositoryURL}}/pkg/exports"

	{{range $i, $d := .DependencyInfos -}}
	exports{{$i}} "{{$d.RepositoryURL}}/pkg/exports"
	{{end}}
)

{{resetFoundWebSocket}}

// ExampleAPI is an example, trivial implementation of the API interface.
// It simply logs the request name.
type ExampleAPI struct {
	{{range $i, $d := .Dependencies -}}
	{{$d | replaceHyphens | toLower}} exports{{$i}}.API
	{{- end}}
}

// NewAPI creates an instance of the example API implementation.
func NewAPI(
	{{- range $i, $d := .Dependencies -}}
	{{$d | replaceHyphens | toLower}} exports{{$i}}.API,
	{{- end -}}
) exports.API {
	return &ExampleAPI{
		{{range $i, $d := .Dependencies -}}
		{{$d | replaceHyphens | toLower}}: {{$d | replaceHyphens | toLower}},
		{{- end}}
	}
}

{{range $a := .API}}
	// {{$a.Path}}
	{{range $mname, $method := $a.Methods}}
	{{if eq $method.Type "WS"}}
		// New{{$mname | capitalize}}Handler example.
		func (a *ExampleAPI) New{{$mname | capitalize}}Handler (connection.FullDuplexEndpoint, error) {
			log.Info("called {{$mname}}")
			return nil, errors.New("method '{{$mname}}' not implemented")
		}

		type {{$mname}}Handler struct {
		}

		// ProcessMessage implements a method of the connection.FullDuplexEndpoint interface.
		func (s *{{$mname}}Handler) ProcessMessage(m *connection.Message, write connection.WriteFn) error {
			log.Info("ProcessMessage not implemented")
			return nil
		}

		// Poll implements a method of the connection.FullDuplexEndpoint interface.
		func (s *{{$mname}}Handler) Poll(t time.Time, write connection.WriteFn) error {
			log.Info("Poll not implemented")
			return nil
		}
	{{- else -}}
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
			log.Info("called {{$mname}}")
			return nil, errors.New("method '{{$mname}}' not implemented")
		}
	{{- end}}
	{{- end -}}
{{- end -}}`
