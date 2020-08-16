package templates

// Client is the template for the client of the service.
const Client = `package client

{{range $a := .API}}{{range $mname, $method := $a.Methods}}
	{{$method.Type | print | checkIfWebSocket}}
{{end}}{{end}}

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/popescu-af/saas-y/pkg/connection"

	"{{.RepositoryURL}}/pkg/exports"
)

{{with $cleanName := .Name | cleanName | capitalize}}
// {{$cleanName}}Client is the structure that encompasses a {{$.Name}} client.
type {{$cleanName}}Client struct {
	connectionManager *connection.FullDuplexManager
	remoteAddress string
}

// New{{$cleanName}}Client creates a new instance of {{$.Name}} client.
func New{{$cleanName}}Client(remoteAddress string) *{{$cleanName}}Client {
	return &{{$cleanName}}Client{
		connectionManager: connection.NewFullDuplexManager(),
		remoteAddress: remoteAddress,
	}
}

{{range $a := $.API}}
{{range $mname, $method := $a.Methods}}
{{if eq $method.Type "WS"}}
// New{{$mname | capitalize}}Client creates a client for websocket at the path '{{$a.Path}}'.
// The caller is responsible to close the returned websocket channel when done.
func (c *{{$cleanName}}Client) New{{$mname | capitalize}}Client(listener connection.ChannelListener) error {
	u := url.URL{Scheme: "ws", Host: c.remoteAddress, Path: "{{$a.Path}}"}
	conn, err := connection.NewWebSocketClient(u, listener)
	if err != nil {
		return err
	}
	c.connectionManager.AddConnection(conn)
	return nil
}
{{- else -}}
// {{$mname | capitalize}} is the client function for {{$method.Type}} '{{$a.Path}}'.
func (c *{{$cleanName}}Client) {{$mname | capitalize}}(
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
)
{{- if eq $method.ReturnType "" -}}
error {
{{- else -}}
(*exports.{{$method.ReturnType | capitalize | symbolize}}, error) {
{{- end}}
	var body io.Reader

	{{if $method.InputType -}}
		b, err := json.Marshal(input)
		if err != nil {
			return {{if ne $method.ReturnType ""}}nil,{{end}} err
		}

		body = bytes.NewBuffer(b)
	{{end}}

	{{with $fmtAndArgs := $a.Path | createPathWithParameterValues -}}
		url := "http://" + c.remoteAddress + fmt.Sprintf("{{index $fmtAndArgs 0}}"{{index $fmtAndArgs 1}})
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
		return {{if ne $method.ReturnType ""}}nil,{{end}} err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return {{if ne $method.ReturnType ""}}nil,{{end}} fmt.Errorf("{{$method.Type}} %s failed with status code %d", url, response.StatusCode)
	}

	{{if eq $method.ReturnType "" -}}
	return nil
	{{- else -}}
	result := new(exports.{{$method.ReturnType | capitalize | symbolize}})
	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
	{{- end}}
}
{{end}}

{{end}}
{{end}}
{{if eq foundWebSocket "yes"}}
	// CloseConnections closes all connections made by this client.
	func (c *{{$cleanName}}Client) CloseConnections() {
		c.connectionManager.CloseConnections()
	}
{{end -}}
{{end}}
`
