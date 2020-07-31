package templates

// Impl is the template for the main go implementation of the API.
const Impl = `package logic

import (
	"errors"
	"time"

	"github.com/popescu-af/saas-y/pkg/log"
	"github.com/popescu-af/saas-y/pkg/connection"

	"{{.RepositoryURL}}/internal/config"
	"{{.RepositoryURL}}/pkg/exports"

	{{range $d := .DependencyInfos -}}
	{{$d.Name | cleanName | toLower}} "{{$d.RepositoryURL}}/pkg/exports"
	{{end}}
)

// Implementation is the main implementation of the API interface.
type Implementation struct {
	env config.Env
	{{range $d := .DependencyInfos -}}
	{{$d.Name | cleanName}} {{$d.Name | cleanName | toLower}}.APIClient
	{{end}}
}

// NewImpl creates an instance of the main implementation.
func NewImpl(env config.Env,
	{{- range $d := .DependencyInfos -}}
	{{$d.Name | cleanName}} {{$d.Name | cleanName | toLower}}.APIClient,
	{{- end -}}
) exports.API {
	return &Implementation{
		env: env,
		{{range $d := .DependencyInfos -}}
		{{$d.Name | cleanName}}: {{$d.Name | cleanName}},
		{{end}}
	}
}

{{range $a := .API}}
	// {{$a.Path}}
	{{range $mname, $method := $a.Methods}}
	{{if eq $method.Type "WS"}}
		// New{{$mname | capitalize}}ChannelListener implementation.
		func (i *Implementation) New{{$mname | capitalize}}ChannelListener() (connection.ChannelListener, error) {
			log.Info("called {{$mname}}")
			return nil, errors.New("method '{{$mname}}' not implemented")
		}

		type {{$mname}}ChannelListener struct {
		}

		// ProcessMessage implements a method of the connection.ChannelListener interface.
		func (s *{{$mname}}ChannelListener) ProcessMessage(m *connection.Message, write connection.WriteOnChannelFunc) error {
			log.Info("ProcessMessage not implemented")
			return nil
		}
	{{- else -}}
		// {{$mname | capitalize}} implementation.
		func (i *Implementation) {{$mname | capitalize}}(
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
