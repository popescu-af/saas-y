package generator

import (
	"os"
	"path"

	"github.com/popescu-af/saas-y/pkg/model"
)

// Go generates go code + infrastructure specification from a saas-y Spec.
type Go struct {
}

// type Variable struct {
// 	Name  string `json:"name"`
// 	Type  string `json:"type"`
// 	Value string `json:"value"`
// }

// type Struct struct {
// 	Name   string     `json:"name"`
// 	Fields []Variable `json:"fields"`
// }

// Generate does the actual code + infrastructure specification generation.
func (g *Go) structs(structs []model.Struct, outdir string) (err error) {
	dir := path.Join(outdir, "pkg", "structs")
	if err = os.MkdirAll(dir, 0770); err != nil {
		return
	}

	// for _, s := range structs {

	// }

	return
}
