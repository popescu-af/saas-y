package golang

import (
	"github.com/popescu-af/saas-y/pkg/generator"
	gengo "github.com/popescu-af/saas-y/pkg/generator/go"
	"github.com/popescu-af/saas-y/pkg/parser"
)

// GenerateSourcesFromJSONSpec generates go code for the services from
// a JSON specification, saving it under the specified path.
func GenerateSourcesFromJSONSpec(jsonSpecFilePath, outdir string) (err error) {
	p := &parser.JSON{}
	spec, err := p.Parse(jsonSpecFilePath)
	if err != nil {
		return
	}

	g := &gengo.Generator{}
	err = generator.Do(g, spec, outdir)
	return
}
