package parser

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/popescu-af/saas-y/internal/model"
)

// A JSON parses the saas-y spec from a JSON file.
type JSON struct {
}

// Parse does the actual parsing, creating the Spec struct from a JSON file.
func (j *JSON) Parse(filename string) (spec *model.Spec, err error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	spec = &model.Spec{}
	err = json.NewDecoder(bytes.NewBuffer(b)).Decode(spec)
	if err != nil {
		spec = nil
	}

	spec.GenerateAdditionalInformation()
	return
}
