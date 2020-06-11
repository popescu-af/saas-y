package parser

import (
	"github.com/popescu-af/saas-y/internal/model"
)

// Abstract is the interface for saas-y parsers of different input formats.
type Abstract interface {
	Parse(filename string) (*model.Spec, error)
}
