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

// 1. definition, example, wrapper tests
// (from same svc specs, three content equality expectations, one for each file)
// - body(input) generation
// - header params generation
// - query params generation
// - path params generation
// - combination of all types of params & body generation + parameter passing
// - path generation
// 2. structs generation
// 3. method type validation
// - correct params for each method type validation

func generateServiceFiles(svc model.Service, components []string) (pOutdir string, err error) {
	pOutdir, err = saasy_testing.CreateOutdir()
	if err != nil {
		return
	}

	g := &gengo.Generator{}
	for _, c := range components {
		if err = generator.ServiceComponent(g, svc, c, pOutdir); err != nil {
			os.RemoveAll(pOutdir)
			return
		}
	}
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

	pOutdir, err := generateServiceFiles(svc, []string{"env"})
	require.NoError(t, err)
	defer os.RemoveAll(pOutdir)

	referenceDir := path.Join(saasy_testing.GetTestingCommonDirectory(), "..", "generator", "go", "testing", "expected")
	saasy_testing.CheckFilesInDirsEqual(t, pOutdir, referenceDir, []string{"env.go"})
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
