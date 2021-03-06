package templates

// HTTPWrapper is the template for HTTP boilerplate in go code.
const HTTPWrapper = `package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/popescu-af/saas-y/pkg/connection"
	"github.com/popescu-af/saas-y/pkg/log"

	"{{.RepositoryURL}}/pkg/exports"
)

// HTTPWrapper decorates the APIs with from/to HTTP code.
type HTTPWrapper struct {
	api exports.API
}

// NewHTTPWrapper creates an HTTP wrapper for the service API.
func NewHTTPWrapper(api exports.API) *HTTPWrapper {
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

{{range $a := .API}}{{range $mname, $method := $a.Methods}}
{{if eq $method.Type "WS"}}
// {{$mname | capitalize | symbolize}} WebSocket wrapper.
func (h *HTTPWrapper) {{$mname | capitalize | symbolize}}(w http.ResponseWriter, r *http.Request) {
	listener, err := h.api.New{{$mname | capitalize | symbolize}}ChannelListener()
	if err != nil {
		writeErrorToHTTPResponse(err, w)
		log.ErrorCtx("creating instance of {{$mname | capitalize | symbolize}}ChannelListener failed", log.Context{"error": err})
		return
	}

	conn, err := connection.NewWebSocketServer(w, r, listener)
	if err != nil {
		writeErrorToHTTPResponse(err, w)
		log.ErrorCtx("creating websocket connection failed", log.Context{"error": err})
		return
	}

	conn.Run()
}
{{else}}
// {{$mname | capitalize | symbolize}} HTTP wrapper.
func (h *HTTPWrapper) {{$mname | capitalize | symbolize}}(w http.ResponseWriter, r *http.Request) {
	{{if $method.InputType}}// Body
	{{"body" | pushParam}} := &exports.{{$method.InputType | capitalize | symbolize}}{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.ErrorCtx("decoding input failed", log.Context{"error": err})
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
								w.WriteHeader(http.StatusBadRequest)
								return
							}

						{{end}}
					{{end}}
				{{end}}
			{{end}}
		{{end}}
	{{end}}{{if $method.QueryParams}}// Query params
	query := r.URL.Query()

	{{range $method.QueryParams}}{{if eq .Type "string"}}{{.Name | decapitalize | pushParam}} := query.Get("{{.Name}}")

	{{else}}{{.Name | decapitalize | pushParam}}, err := parse{{.Type | capitalize}}Parameter(query.Get("{{.Name}}"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	{{end}}{{end}}{{end}}{{if $method.HeaderParams}}// Header params
	{{range $method.HeaderParams}}{{if eq .Type "string"}}{{.Name | decapitalize | pushParam}} := r.Header.Get("{{.Name}}")

	{{else}}{{.Name | decapitalize | pushParam}}, err := parse{{.Type | capitalize}}Parameter(r.Header.Get("{{.Name}}"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	{{end}}{{end}}{{end}}
	// Call implementation
	{{if ne $method.ReturnType ""}}result,{{end}} err := h.api.{{$mname | capitalize | symbolize}}({{printParamStack}})
	if err != nil {
		writeErrorToHTTPResponse(err, w)
		log.ErrorCtx("call to implementation failed", log.Context{"error": err})
		return
	}
	{{- if ne $method.ReturnType ""}}

	encodeJSONResponse(result, nil, w){{end}}
}
{{end}}
{{end}}{{end}}`
