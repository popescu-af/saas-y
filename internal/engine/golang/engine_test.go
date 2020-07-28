package golang

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	saasy_testing "github.com/popescu-af/saas-y/internal/testing"
)

var fullSpec = `
{
    "repository_url": "example.com/example",
    "services": [
        {
            "name": "foo-service",
            "port": "80",
            "api": [
                {
                    "path": "/foo",
                    "methods": {
                        "method_name_0": {
                            "type": "GET",
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
                            "type": "POST",
                            "input_type": "input_struct_name",
                            "return_type": "return_struct_name"
                        },
                        "cool_websocket": {
                            "type": "WS"
                        }
                    }
                },
                {
                    "path": "/foo/{rank:uint}/{price:float}",
                    "methods": {
                        "method_name_3": {
                            "type": "GET",
                            "return_type": "return_struct_name"
                        },
                        "method_name_4": {
                            "type": "DELETE",
                            "return_type": "return_struct_name"
                        },
                        "method_name_5": {
                            "type": "PATCH",
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
                },
                {
                    "name": "return_struct_name",
                    "fields": [
                        {
                            "name": "status",
                            "type": "int"
                        }
                    ]
                }
            ],
            "env": [
                {
                    "name": "env_var_name",
                    "type": "int",
                    "value": "42"
                }
            ]
        },
        {
            "name": "bar-service",
            "port": "80",
            "api": [
                {
                    "path": "/bar",
                    "methods": {
                        "method_name_0": {
                            "type": "GET",
                            "return_type": "return_struct_name"
                        },
                        "method_name_1": {
                            "type": "POST",
                            "return_type": "return_struct_name"
                        }
                    }
                }
            ],
            "structs": [
                {
                    "name": "return_struct_name",
                    "fields": [
                        {
                            "name": "status",
                            "type": "int"
                        }
                    ]
                }
            ],
            "dependencies": [ "foo-service" ]
        }
    ]
}
`

func TestGeneratedServiceCompiles(t *testing.T) {
	pSpec, err := saasy_testing.CreateJSONSpecFile(fullSpec, ".", "spec.json")
	require.NoError(t, err)

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
