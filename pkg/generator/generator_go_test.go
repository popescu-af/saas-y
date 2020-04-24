package generator

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	saasy_testing "github.com/popescu-af/saas-y/pkg/testing"
)

// env generation
// body(input) generation
// header params generation
// query params generation
// path params generation
// combination of all types of params & body generation + parameter passing
// path generation
// method type validation
// correct params for each method type validation
// options method correctly generated
// structs generation

func TestGeneratedEnv(t *testing.T) {
	var spec = `
    {
        "services": [
            {
                "name": "foo-service",
                "port": "80",
                "env": [
                    {
                        "name": "ENV_VAR_NAME",
                        "type": "int64",
                        "value": "42"
                    }
                ]
            }
        ]
    }
    `

	pSpec, pOutdir, err := saasy_testing.CreateJSONSpecFileAndOutdir(spec, ".", "spec.json")
	require.NoError(t, err)

	defer os.Remove(pSpec)
	defer os.RemoveAll(pOutdir)

	// TODO: test
}
