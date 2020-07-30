package templates

// HTTPErrorHandler is the template for the go definition
// of the HTTP error handler code.
const HTTPErrorHandler = `package service

import (
	"net/http"

	"{{.RepositoryURL}}/internal/logic"
)

func writeErrorToHTTPResponse(err error, w http.ResponseWriter) {
	if err == nil {
		return
	}

	switch err.(type) {
	case *logic.NotFoundError:
		w.WriteHeader(http.StatusNotFound)
	}

	w.Write([]byte(err.Error()))
}`
