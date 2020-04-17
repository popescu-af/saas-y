package main

import (
	"log"

	"github.com/popescu-af/saas-y/pkg/generator"
	"github.com/popescu-af/saas-y/pkg/parser"
)

func main() {
	outdir := "/Users/alexandru/go/src/github.com/popescu-af/saas-y/gen_test"

	p := &parser.JSON{}

	spec, err := p.Parse("/Users/alexandru/go/src/github.com/popescu-af/saas-y/example/spec.json")
	if err != nil {
		log.Fatalf("saas-y error: %v", err)
	}

	g := &generator.Go{}
	generator.Do(g, spec, outdir)
}
