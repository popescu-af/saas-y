package parser

import (
	"github.com/popescu-af/saas-y/pkg/model"
)

// Abstract is the interface for saas-y parsers of different input formats.
type Abstract interface {
	Parse(filename string) (*model.Spec, error)
}
