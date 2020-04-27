package golang

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	saasy_testing "github.com/popescu-af/saas-y/pkg/testing"
)

var fullSpec = `
{
    "services": [
        {
            "name": "foo-service",
            "port": "80",
            "api": [
                {
                    "path": "/foo",
                    "methods": {
                        "method_name_0": {
                            "type": "get",
                            "header_params": [
                                {
                                    "name": "header_param_name",
                                    "type": "int"
                                }
                            ],
                            "query_params": [
                                {
                                    "name": "query_param_name",
                                    "type": "string"
                                }
                            ],
                            "return_type": "return_struct_name"
                        },
                        "method_name_1": {
                            "type": "post",
                            "input_type": "input_struct_name",
                            "return_type": "return_struct_name"
                        },
                        "method_name_2": {
                            "type": "options",
                            "header_params": [
                                {
                                    "name": "header_param_name",
                                    "value": "value"
                                }
                            ]
                        }
                    }
                },
                {
                    "path": "/bar/{rank:uint}/{price:float}",
                    "methods": {
                        "method_name_3": {
                            "type": "get",
                            "return_type": "return_struct_name"
                        },
                        "method_name_4": {
                            "type": "delete",
                            "return_type": "return_struct_name"
                        },
                        "method_name_5": {
                            "type": "patch",
                            "input_type": "input_struct_name",
                            "return_type": "return_struct_name"
                        }
                    }
                }
            ],
            "structs": [
                {
                    "name": "input_struct_name",
                    "fields": [
                        {
                            "name": "a_field_name",
                            "type": "int"
                        },
                        {
                            "name": "another_field_name",
                            "type": "string"
                        }
                    ]
                }
            ],
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

func TestGeneratedServiceCompiles(t *testing.T) {
	pSpec, err := saasy_testing.CreateJSONSpecFile(fullSpec, ".", "spec.json")
	pOutdir, err := saasy_testing.CreateOutdir()
	require.NoError(t, err)

	defer os.Remove(pSpec)
	defer os.RemoveAll(pOutdir)

	err = GenerateSourcesFromJSONSpec(pSpec, pOutdir)
	require.NoError(t, err)

	// compile
	var errout bytes.Buffer
	cmd := exec.Command("go", "build", "./cmd/main.go")
	cmd.Stderr = &errout
	cmd.Dir = path.Join(pOutdir, "services", "foo-service")

	err = cmd.Run()
	require.NoError(t, err, errout.String())
}
