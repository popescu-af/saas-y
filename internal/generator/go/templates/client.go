package templates

// Client is the template for the client of the service.
const Client = `package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"{{.RepositoryURL}}/pkg/exports"
)

{{with $cleanName := .Name | cleanName | capitalize}}
// {{$cleanName}}Client is the structure that encompasses a {{$.Name}} client.
type {{$cleanName}}Client struct {
	RemoteAddress string
}

// New{{$cleanName}}Client creates a new instance of {{$.Name}} client.
func New{{$cleanName}}Client(remoteAddress string) *{{$cleanName}}Client {
	return &{{$cleanName}}Client{
		RemoteAddress: remoteAddress,
	}
}

{{range $a := $.API}}
{{range $mname, $method := $a.Methods}}
// {{$mname | capitalize}} is the client function for {{$method.Type}} '{{$a.Path}}'.
func (a *{{$cleanName}}Client) {{$mname | capitalize}}(
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
	var body io.Reader

	{{if $method.InputType -}}
		b, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}

		body = bytes.NewBuffer(b)
	{{end}}

	{{with $fmtAndArgs := $a.Path | createPathWithParameterValues -}}
		url := a.RemoteAddress + fmt.Sprintf("{{index $fmtAndArgs 0}}"{{index $fmtAndArgs 1}})
	{{- end}}
	{{if $method.QueryParams -}}
		{{- range $i, $p := $method.QueryParams}}
			{{if eq $i 0 -}}
				url += fmt.Sprintf("?{{$p.Name}}={{$p.Type | typePlaceholder}}", {{$p.Name}})
			{{- else -}}
				url += fmt.Sprintf("&{{$p.Name}}={{$p.Type | typePlaceholder}}", {{$p.Name}})
			{{- end}}
		{{- end}}
	{{- end}}

	request, err := http.NewRequest("{{$method.Type}}", url, body)
	{{- if $method.HeaderParams -}}
		{{range $method.HeaderParams}}
			request.Header.Set("{{.Name}}", fmt.Sprintf("{{.Type | typePlaceholder}}", {{.Name}}))
		{{- end}}
	{{- end}}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("{{$method.Type}} %s failed with status code %d", url, response.StatusCode)
	}

	result := new(exports.{{$method.ReturnType | capitalize | symbolize}})
	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

{{end}}
{{end}}
{{end}}
`
