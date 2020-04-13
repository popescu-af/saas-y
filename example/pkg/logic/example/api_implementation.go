package example

import (
	"errors"

	"github.com/popescu-af/saas-y/template/pkg/service"
	"github.com/popescu-af/saas-y/template/pkg/service/structs"

	"go.uber.org/zap"
)

// API TODO.
type API struct {
	logger *zap.Logger
}

// NewAPI creates a dummy example API implementation.
func NewAPI(logger *zap.Logger) service.API {
	return &API{logger: logger}
}

// AddSomething example.
func (s *API) AddSomething(body structs.Something) (interface{}, error) {
	s.logger.Info("AddSomething", zap.Any("something", body))
	return nil, errors.New("method 'AddSomething' not implemented")
}
