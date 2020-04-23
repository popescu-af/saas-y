package templates

// HTTPWrapper is the template for HTTP boilerplate in go code.
const HTTPWrapper = `package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"{{.Name}}/pkg/structs"
)

// HTTPWrapper decorates the APIs with from/to HTTP code.
type HTTPWrapper struct {
	api API
}

// NewHTTPWrapper creates an HTTP wrapper for the service API.
func NewHTTPWrapper(api API) *HTTPWrapper {
	return &HTTPWrapper{api: api}
}

func encodeJSONResponse(i interface{}, status *int, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if status != nil {
		w.WriteHeader(*status)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	return json.NewEncoder(w).Encode(i)
}

func parseIntParameter(param string) (int64, error) {
	return strconv.ParseInt(param, 10, 64)
}

// Paths lists the paths that the API serves.
func (h *HTTPWrapper) Paths() Paths {
	return Paths{
		{{range $a := .API}}{{range $mname, $method := $a.Methods}}{
			strings.ToUpper("{{$method.Type}}"),
			"{{$a.Path}}",
			h.{{$mname | capitalize | symbolize}},
		},
		{{end}}{{end}}
	}
}

{{range $a := .API}}{{range $mname, $method := $a.Methods}}
func (h *HTTPWrapper) {{$mname | capitalize | symbolize}}(w http.ResponseWriter, r *http.Request) {
	{{if $method.InputType}}// Body
	body := structs.{{$method.InputType | capitalize | symbolize}}{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(500)
		return
	}

	{{end}}{{if $a.Path | pathHasParameters}}// Path params
	pathParams := mux.Vars(r)
	{{with $params := $a.Path | pathParameters}}
	{{range $pnameidx := $params | indicesParameters}}
	{{with $ptypeidx := inc $pnameidx}}
	{{index $params $pnameidx}}, err := parse{{index $params $ptypeidx | capitalize}}Parameter(params["{{index $params $pnameidx}}"])
	if err != nil {
		w.WriteHeader(500)
		return
	}

	{{end}}{{end}}{{end}}{{end}}{{if $method.HeaderParams}}// Header params
	{{range $method.HeaderParams}}{{if eq $method.Type "options"}}
	{{/* TODO */}}
	{{else}}{{.Name}}, err := parse{{.Type | capitalize}}Parameter(r.Header.Get("{{.Name}}"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	{{end}}{{end}}{{end}}{{if $method.QueryParams}}// Query params
	{{range $method.QueryParams}}{{.Name}}, err := parse{{.Type | capitalize}}Parameter(r.URL.Query.Get("{{.Name}}"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	{{end}}{{end}}
	result, err := s.api.{{$mname | capitalize | symbolize}}({{if $method.InputType}}body{{end}})
	if err != nil {
		w.WriteHeader(500)
		return
	}

	encodeJSONResponse(result, nil, w)
}
{{end}}{{end}}`