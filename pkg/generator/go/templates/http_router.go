package templates

// HTTPRouter is a template for the HTTP router go file.
const HTTPRouter = `package service

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// A PathDefinition groups an HTTP method on a path with its handler function.
type PathDefinition struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

// Paths represents a collection of path definitions.
type Paths []PathDefinition

// NewRouter creates a new router for the given paths.
func NewRouter(paths Paths, logger *zap.Logger) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, p := range paths {
		router.
			Methods(p.Method).
			Path(p.Path).
			Handler(apiLogger(p.Handler, logger))
	}
	return router
}

func apiLogger(handler http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info(
			"serving",
			zap.String("method", r.Method),
			zap.String("path", r.RequestURI),
		)

		start := time.Now()
		handler.ServeHTTP(w, r)

		logger.Info(
			"served",
			zap.String("method", r.Method),
			zap.String("path", r.RequestURI),
			zap.String("duration", time.Since(start).String()),
		)
	})
}`
