package templates

// HTTPErrorHandlerExample is the template for the go definition
// of the HTTP error handler example code.
const HTTPErrorHandlerExample = `package service

import (
	"net/http"

	"{{.Name}}/pkg/logic"
)

// NotFoundError should be returned when a specific resource
// requested by the client of the logic does not exist.
type NotFoundError struct {
	message string
}

func (e *NotFoundError) Error() string {
	return e.message
}

// NewNotFoundError creates a custom NotFoundError instance.
func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{message: message}
}

func writeErrorToHTTPResponse(err error, w http.ResponseWriter) {
	if err == nil {
		return
	}

	switch err.(type) {
	case *logic.NotFoundError:
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
	}
}`
