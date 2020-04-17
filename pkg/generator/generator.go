package generator

import (
	"github.com/popescu-af/saas-y/pkg/model"
)

// Abstract is the interface for saas-y code & infrastructure generators.
type Abstract interface {
	structs(structs []model.Struct, outdir string) error
}

// Do generates code and infrastructure declaration for the given Spec
// and dumps it in the specified output directory.
func Do(generator Abstract, spec *model.Spec, outdir string) (err error) {
	err = generator.structs(spec.Structs, outdir)
	return
}
