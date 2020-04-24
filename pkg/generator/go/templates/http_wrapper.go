package templates

// HTTPWrapper is the template for HTTP boilerplate in go code.
const HTTPWrapper = `package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"{{.Name}}/pkg/logic"
	"{{.Name}}/pkg/structs"
)

// HTTPWrapper decorates the APIs with from/to HTTP code.
type HTTPWrapper struct {
	api logic.API
}

// NewHTTPWrapper creates an HTTP wrapper for the service API.
func NewHTTPWrapper(api logic.API) *HTTPWrapper {
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

func parseUintParameter(param string) (uint64, error) {
	return strconv.ParseUint(param, 10, 64)
}

func parseFloatParameter(param string) (float64, error) {
	return strconv.ParseFloat(param, 64)
}

// Paths lists the paths that the API serves.
func (h *HTTPWrapper) Paths() Paths {
	return Paths{
		{{range $a := .API}}{{range $mname, $method := $a.Methods}}{
			strings.ToUpper("{{$method.Type}}"),
			"{{$a.Path | cleanPath}}",
			h.{{$mname | capitalize | symbolize}},
		},
		{{end}}{{end}}
	}
}

{{range $a := .API}}{{range $mname, $method := $a.Methods}}// {{$mname | capitalize | symbolize}} HTTP wrapper.
func (h *HTTPWrapper) {{$mname | capitalize | symbolize}}(w http.ResponseWriter, r *http.Request) {
	{{if $method.InputType}}// Body
	{{"body" | pushParam}} := &structs.{{$method.InputType | capitalize | symbolize}}{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(500)
		return
	}

	{{end}}{{if $a.Path | pathHasParameters}}// Path params
		pathParams := mux.Vars(r)
		{{with $params := $a.Path | pathParameters}}
			{{range $pnameidx := $params | indicesParameters}}
				{{with $ptypeidx := inc $pnameidx}}
					{{with $ptype := index $params $ptypeidx}}
						{{if eq $ptype "string"}}
							{{index $params $pnameidx | decapitalize | pushParam}} := pathParams["{{index $params $pnameidx}}"]
						{{else}}
							{{index $params $pnameidx | decapitalize | pushParam}}, err := parse{{index $params $ptypeidx | capitalize}}Parameter(pathParams["{{index $params $pnameidx}}"])
							if err != nil {
								w.WriteHeader(500)
								return
							}

						{{end}}
					{{end}}
				{{end}}
			{{end}}
		{{end}}
	{{end}}{{if $method.HeaderParams}}{{if eq $method.Type "options"}}// Response headers
	{{else}}// Header params
	{{end}}{{range $method.HeaderParams}}{{if eq $method.Type "options"}}w.Header().Set("{{.Name}}", "{{.Value}}"){{else}}{{if eq .Type "string"}}{{.Name | decapitalize | pushParam}} := r.Header.Get("{{.Name}}"){{else}}{{.Name | decapitalize | pushParam}}, err := parse{{.Type | capitalize}}Parameter(r.Header.Get("{{.Name}}"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	{{end}}{{end}}{{end}}{{end}}{{if $method.QueryParams}}// Query params
	query := r.URL.Query()
	{{range $method.QueryParams}}{{if eq .Type "string"}}{{.Name | decapitalize | pushParam}} := query.Get("{{.Name}}"){{else}}{{.Name | decapitalize | pushParam}}, err := parse{{.Type | capitalize}}Parameter(query.Get("{{.Name}}"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	{{end}}{{end}}{{end}}{{if eq $method.Type "options"}}
	{{/* NADA for "options" method */}}{{else}}// Call implementation
	result, err := h.api.{{$mname | capitalize | symbolize}}({{printParamStack}})
	if err != nil {
		w.WriteHeader(500)
		return
	}

	encodeJSONResponse(result, nil, w){{end}}
}
{{end}}{{end}}`
