package parser

import "io"

type AbstractParser interface {
	Parse(r io.Reader) *Spec
}
