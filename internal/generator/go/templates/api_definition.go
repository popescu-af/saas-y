package templates

// APIDefinition is the template for the go definition of the API.
const APIDefinition = `package exports

{{range $a := .API}}{{range $mname, $method := $a.Methods}}
	{{$method.Type | print | checkIfWebSocket}}
{{end}}{{end}}

import (
	"time"

	"github.com/popescu-af/saas-y/pkg/connection"
)

// API defines the operations supported by the {{.Name}} service.
type API interface {
	{{- range $a := .API}}
		// {{$a.Path}}
		{{range $mname, $method := $a.Methods}}
		{{- if eq $method.Type "WS" -}}
			{{with $fname := $mname | capitalize -}}
			{{printf "%s%s%s" "New" $fname "ChannelListener" | symbolize}}() (connection.ChannelListener, error)
			{{- end}}
		{{else -}}
			{{- $mname | capitalize | symbolize}}(
				{{- if $method.InputType -}}
					*{{- $method.InputType | capitalize | symbolize}},
				{{- end -}}
				{{- if $a.Path | pathHasParameters -}}
					{{- with $params := $a.Path | pathParameters -}}
						{{- range $pnameidx := $params | indicesParameters -}}
							{{- with $ptypeidx := inc $pnameidx -}}
								{{- index $params $ptypeidx | typeName -}},
						{{- end -}}
					{{- end -}}
				{{- end -}}
			{{- end -}}
			{{- if $method.QueryParams -}}
				{{- range $method.QueryParams -}}
					{{- .Type | typeName -}},
				{{- end -}}
			{{- end -}}
			{{- if $method.HeaderParams -}}
				{{- range $method.HeaderParams -}}
					{{- .Type | typeName -}},
				{{- end -}}
			{{- end -}}
			)
			{{- if eq $method.ReturnType "" -}}
				error
			{{- else -}}
				(*{{$method.ReturnType | capitalize | symbolize}}, error)
			{{- end}}
		{{end -}}
		{{end -}}
	{{- end -}}
}

// APIClient defines the operations supported by the {{.Name}} service client.
type APIClient interface {
	{{- range $a := .API}}
		// {{$a.Path}}
		{{range $mname, $method := $a.Methods}}
		{{- if eq $method.Type "WS" -}}
			{{with $fname := $mname | capitalize -}}
			{{printf "%s%s%s" "New" $fname "Client" | symbolize}}(connection.ChannelListener) (*connection.FullDuplex, error)
			{{- end}}
		{{else -}}
			{{- $mname | capitalize | symbolize}}(
				{{- if $method.InputType -}}
					*{{- $method.InputType | capitalize | symbolize}},
				{{- end -}}
				{{- if $a.Path | pathHasParameters -}}
					{{- with $params := $a.Path | pathParameters -}}
						{{- range $pnameidx := $params | indicesParameters -}}
							{{- with $ptypeidx := inc $pnameidx -}}
								{{- index $params $ptypeidx | typeName -}},
						{{- end -}}
					{{- end -}}
				{{- end -}}
			{{- end -}}
			{{- if $method.QueryParams -}}
				{{- range $method.QueryParams -}}
					{{- .Type | typeName -}},
				{{- end -}}
			{{- end -}}
			{{- if $method.HeaderParams -}}
				{{- range $method.HeaderParams -}}
					{{- .Type | typeName -}},
				{{- end -}}
			{{- end -}}
			)
			{{- if eq $method.ReturnType "" -}}
				error
			{{- else -}}
				(*{{$method.ReturnType | capitalize | symbolize}}, error)
			{{- end}}
		{{end -}}
		{{end -}}
	{{- end}}
	{{- if eq foundWebSocket "yes"}}
			CloseConnections()
	{{end -}}
}
`
