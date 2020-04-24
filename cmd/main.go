package main

import (
	"log"

	"github.com/popescu-af/saas-y/pkg/engine/golang"
)

func main() {
	err := golang.GenerateSourcesFromJSONSpec(
		"/Users/alexandru/go/src/github.com/popescu-af/saas-y/example/spec.json",
		"/Users/alexandru/go/src/github.com/popescu-af/saas-y/example/_gen_test",
	)
	if err != nil {
		log.Fatalf("saas-y error: %v", err)
	}
}
