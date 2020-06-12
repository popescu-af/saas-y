package generator_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/popescu-af/saas-y/internal/generator"
	"github.com/popescu-af/saas-y/internal/generator/common/templates/k8s"
	gengo "github.com/popescu-af/saas-y/internal/generator/go"
	"github.com/popescu-af/saas-y/internal/model"
	saasytesting "github.com/popescu-af/saas-y/internal/testing"
)

func generateServiceFiles(svc model.Service) (pOutdir string, err error) {
	pOutdir, err = saasytesting.CreateOutdir()
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
		API:           []model.API{},
		Structs:       []model.Struct{},
		RepositoryURL: "foo-service",
	}

	generator.Init()

	pOutdir, err := generateServiceFiles(svc)
	require.NoError(t, err)
	defer os.RemoveAll(pOutdir)

	pOutdir = path.Join(pOutdir, "services", svc.Name)
	referenceDir := path.Join(saasytesting.GetTestingCommonDirectory(), "..", "generator", "testdata", "generated_env")
	saasytesting.CheckFilesInDirsEqual(t, path.Join(pOutdir, "internal", "config"), referenceDir, []string{"env.go"})
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

	pOutdir, err := saasytesting.CreateOutdir()
	if err != nil {
		return
	}
	require.NoError(t, err)
	defer os.RemoveAll(pOutdir)

	err = generator.CommonEntity(spec, k8s.Ingress, path.Join(pOutdir, "ingress.yaml"))
	require.NoError(t, err)

	referenceDir := path.Join(saasytesting.GetTestingCommonDirectory(), "..", "generator", "testdata")
	saasytesting.CheckFilesInDirsEqual(t, pOutdir, referenceDir, []string{"ingress.yaml"})
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
		RepositoryURL: "foo-service",
	}

	generator.Init()

	pOutdir, err := generateServiceFiles(svc)
	require.NoError(t, err)
	defer os.RemoveAll(pOutdir)

	pOutdir = path.Join(pOutdir, "services", svc.Name)
	referenceDir := path.Join(saasytesting.GetTestingCommonDirectory(), "..", "generator", "testdata", "generated_methods")
	saasytesting.CheckFilesInDirsEqual(t, path.Join(pOutdir, "internal", "logic"), referenceDir, []string{"api_example.go"})
	saasytesting.CheckFilesInDirsEqual(t, path.Join(pOutdir, "internal", "service"), referenceDir, []string{"http_wrapper.go"})
	saasytesting.CheckFilesInDirsEqual(t, path.Join(pOutdir, "pkg", "exports"), referenceDir, []string{"api_definition.go"})
}

func TestGeneratedStruct(t *testing.T) {
	svc := model.Service{
		ServiceCommon: model.ServiceCommon{
			Name: "foo-service",
			Port: "80",
		},
		API: []model.API{},
		Structs: []model.Struct{
			{
				Name: "hakuna_matata",
				Fields: []model.Variable{
					{
						Name: "whatever_int",
						Type: "int",
					},
					{
						Name: "whatever_float",
						Type: "float",
					},
					{
						Name: "whatever_string",
						Type: "string",
					},
				},
			},
		},
	}

	generator.Init()

	pOutdir, err := generateServiceFiles(svc)
	require.NoError(t, err)
	defer os.RemoveAll(pOutdir)

	pOutdir = path.Join(pOutdir, "services", svc.Name)
	referenceDir := path.Join(saasytesting.GetTestingCommonDirectory(), "..", "generator", "testdata", "generated_structs")
	saasytesting.CheckFilesInDirsEqual(t, path.Join(pOutdir, "pkg", "exports"), referenceDir, []string{"hakuna_matata.go"})
}
