package generator_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/popescu-af/saas-y/pkg/generator"
	"github.com/popescu-af/saas-y/pkg/generator/common/templates/k8s"
	gengo "github.com/popescu-af/saas-y/pkg/generator/go"
	"github.com/popescu-af/saas-y/pkg/model"
	saasy_testing "github.com/popescu-af/saas-y/pkg/testing"
)

// 1. structs generation
// 2. method type validation
// - correct params for each method type validation

func generateServiceFiles(svc model.Service) (pOutdir string, err error) {
	pOutdir, err = saasy_testing.CreateOutdir()
	if err != nil {
		return
	}

	err = generator.Service(&gengo.Generator{}, svc, pOutdir)
	return
}

func TestGeneratedEnv(t *testing.T) {
	svc := model.Service{
		ServiceCommon: model.ServiceCommon{
			Name: "foo-service",
			Port: "80",
			Environment: []model.Variable{
				{
					Name:  "ENV_VAR_NAME",
					Type:  "int64",
					Value: "42",
				},
			},
		},
		API:     []model.API{},
		Structs: []model.Struct{},
	}

	generator.Init()

	pOutdir, err := generateServiceFiles(svc)
	require.NoError(t, err)
	defer os.RemoveAll(pOutdir)

	pOutdir = path.Join(pOutdir, "services", svc.Name)
	referenceDir := path.Join(saasy_testing.GetTestingCommonDirectory(), "..", "generator", "testdata", "generated_env")
	saasy_testing.CheckFilesInDirsEqual(t, path.Join(pOutdir, "pkg", "config"), referenceDir, []string{"env.go"})
}

func TestGeneratedIngress(t *testing.T) {
	spec := model.Spec{
		Domain: "foo.bar",
		Subdomains: []model.Subdomain{
			{
				Name: "api",
				Paths: []model.Path{
					{
						Value:    "/",
						Endpoint: "api-service",
					},
				},
			},
			{
				Name: "baz",
				Paths: []model.Path{
					{
						Value:    "/hakuna",
						Endpoint: "baz0-service",
					},
					{
						Value:    "/matata",
						Endpoint: "baz1-service",
					},
				},
			},
		},
	}

	pOutdir, err := saasy_testing.CreateOutdir()
	if err != nil {
		return
	}
	require.NoError(t, err)
	defer os.RemoveAll(pOutdir)

	err = generator.CommonEntity(spec, k8s.Ingress, path.Join(pOutdir, "ingress.yaml"))
	require.NoError(t, err)

	referenceDir := path.Join(saasy_testing.GetTestingCommonDirectory(), "..", "generator", "testdata")
	saasy_testing.CheckFilesInDirsEqual(t, pOutdir, referenceDir, []string{"ingress.yaml"})
}

func TestGeneratedMethods(t *testing.T) {
	qParams := []model.Variable{
		{Name: "query_param_0", Type: "int"},
		{Name: "query_param_1", Type: "float"},
		{Name: "query_param_2", Type: "string"},
	}

	hParams := []model.Variable{
		{Name: "header_param_0", Type: "string"},
		{Name: "header_param_1", Type: "float"},
		{Name: "header_param_2", Type: "int"},
	}

	bType := "body_type"

	svc := model.Service{
		ServiceCommon: model.ServiceCommon{
			Name: "foo-service",
			Port: "80",
		},
		API: []model.API{
			{
				Path: "/method_no_path_params",
				Methods: map[string]model.Method{
					"method_no_path_params_0": {Type: model.POST, QueryParams: nil, HeaderParams: nil, InputType: "", ReturnType: "return_type"},
					"method_no_path_params_1": {Type: model.POST, QueryParams: nil, HeaderParams: nil, InputType: bType, ReturnType: "return_type"},
					"method_no_path_params_2": {Type: model.POST, QueryParams: nil, HeaderParams: hParams, InputType: "", ReturnType: "return_type"},
					"method_no_path_params_3": {Type: model.POST, QueryParams: nil, HeaderParams: hParams, InputType: bType, ReturnType: "return_type"},
					"method_no_path_params_4": {Type: model.POST, QueryParams: qParams, HeaderParams: nil, InputType: "", ReturnType: "return_type"},
					"method_no_path_params_5": {Type: model.POST, QueryParams: qParams, HeaderParams: nil, InputType: bType, ReturnType: "return_type"},
					"method_no_path_params_6": {Type: model.POST, QueryParams: qParams, HeaderParams: hParams, InputType: "", ReturnType: "return_type"},
					"method_no_path_params_7": {Type: model.POST, QueryParams: qParams, HeaderParams: hParams, InputType: bType, ReturnType: "return_type"},
				},
			},
			{
				Path: "/method/{path_param_0:int}/{path_param_1:string}",
				Methods: map[string]model.Method{
					"method_0": {Type: model.POST, QueryParams: nil, HeaderParams: nil, InputType: "", ReturnType: "return_type"},
					"method_1": {Type: model.POST, QueryParams: nil, HeaderParams: nil, InputType: bType, ReturnType: "return_type"},
					"method_2": {Type: model.POST, QueryParams: nil, HeaderParams: hParams, InputType: "", ReturnType: "return_type"},
					"method_3": {Type: model.POST, QueryParams: nil, HeaderParams: hParams, InputType: bType, ReturnType: "return_type"},
					"method_4": {Type: model.POST, QueryParams: qParams, HeaderParams: nil, InputType: "", ReturnType: "return_type"},
					"method_5": {Type: model.POST, QueryParams: qParams, HeaderParams: nil, InputType: bType, ReturnType: "return_type"},
					"method_6": {Type: model.POST, QueryParams: qParams, HeaderParams: hParams, InputType: "", ReturnType: "return_type"},
					"method_7": {Type: model.POST, QueryParams: qParams, HeaderParams: hParams, InputType: bType, ReturnType: "return_type"},
				},
			},
		},
		Structs: []model.Struct{
			{
				Name: "body_type",
				Fields: []model.Variable{
					{Name: "variable_0", Type: "int"},
					{Name: "variable_1", Type: "string"},
					{Name: "variable_2", Type: "float"},
				},
			},
			{
				Name: "return_type",
				Fields: []model.Variable{
					{Name: "return_variable_0", Type: "string"},
					{Name: "return_variable_1", Type: "string"},
					{Name: "return_variable_2", Type: "float"},
				},
			},
		},
	}

	generator.Init()

	pOutdir, err := generateServiceFiles(svc)
	require.NoError(t, err)
	defer os.RemoveAll(pOutdir)

	pOutdir = path.Join(pOutdir, "services", svc.Name)
	referenceDir := path.Join(saasy_testing.GetTestingCommonDirectory(), "..", "generator", "testdata", "generated_methods")
	saasy_testing.CheckFilesInDirsEqual(t, path.Join(pOutdir, "pkg", "logic"), referenceDir, []string{"api_definition.go", "api_example.go"})
	saasy_testing.CheckFilesInDirsEqual(t, path.Join(pOutdir, "pkg", "service"), referenceDir, []string{"http_wrapper.go"})
}
