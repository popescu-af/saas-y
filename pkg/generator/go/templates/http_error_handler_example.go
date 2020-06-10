package templates

// HTTPErrorHandlerExample is the template for the go definition
// of the HTTP error handler example code.
const HTTPErrorHandlerExample = `package service

import (
	"net/http"

	"{{.Name}}/internal/logic"
)

func writeErrorToHTTPResponse(err error, w http.ResponseWriter) {
	if err == nil {
		return
	}

	switch err.(type) {
	case *logic.NotFoundError:
		w.WriteHeader(404)
	}

	w.Write([]byte(err.Error()))
}`
