package templates

// HTTPWrapper is the template for HTTP boilerplate in go code.
const HTTPWrapper = `package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/popescu-af/saas-y/template/pkg/service/structs"
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
func (s *HTTPWrapper) Paths() Paths {
	return Paths{
		{
			strings.ToUpper("POST"),
			"/something",
			s.addSomething,
		},
	}
}

func (s *HTTPWrapper) addSomething(w http.ResponseWriter, r *http.Request) {
	body := &structs.Something{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(500)
		return
	}

	result, err := s.api.AddSomething(*body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	encodeJSONResponse(result, nil, w)
}`
