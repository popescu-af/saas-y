package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/popescu-af/saas-y/internal/model"
)

func TestMethodTypeValid(t *testing.T) {
	tests := []struct {
		method *model.Method
		valid  bool
	}{
		{&model.Method{Type: "dummy"}, false},
		{&model.Method{Type: "GET"}, true},
		{&model.Method{Type: "POST"}, true},
		{&model.Method{Type: "PATCH"}, true},
		{&model.Method{Type: "DELETE"}, true},
	}

	for _, tt := range tests {
		err := tt.method.Validate([]string{})
		if tt.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func TestMethodBodyOrReturnTypeAllowed(t *testing.T) {
	tests := []struct {
		method *model.Method
		valid  bool
	}{
		{&model.Method{Type: "GET"}, true},
		{&model.Method{Type: "POST"}, true},
		{&model.Method{Type: "PATCH"}, true},
		{&model.Method{Type: "DELETE"}, true},
		{&model.Method{Type: "WS"}, true},
		{&model.Method{Type: "GET", InputType: "whatever"}, false},
		{&model.Method{Type: "POST", InputType: "whatever"}, true},
		{&model.Method{Type: "PATCH", InputType: "whatever"}, true},
		{&model.Method{Type: "DELETE", InputType: "whatever"}, false},
		{&model.Method{Type: "WS", InputType: "whatever"}, false},
		{&model.Method{Type: "GET", ReturnType: "whatever"}, true},
		{&model.Method{Type: "POST", ReturnType: "whatever"}, true},
		{&model.Method{Type: "PATCH", ReturnType: "whatever"}, true},
		{&model.Method{Type: "DELETE", ReturnType: "whatever"}, true},
		{&model.Method{Type: "WS", ReturnType: "whatever"}, false},
	}

	for _, tt := range tests {
		err := tt.method.Validate([]string{"whatever"})
		if tt.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func TestMethodTypes(t *testing.T) {
	tests := []struct {
		method *model.Method
		valid  bool
	}{
		{&model.Method{Type: "POST", InputType: "something_unknown"}, false},
		{&model.Method{Type: "POST", InputType: "something_known"}, true},
		{&model.Method{Type: "POST", ReturnType: "something_unknown"}, false},
		{&model.Method{Type: "POST", ReturnType: "something_known"}, true},
	}

	for _, tt := range tests {
		err := tt.method.Validate([]string{"something_known"})
		if tt.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func TestPathRegex(t *testing.T) {
	tests := []struct {
		path   string
		parsed int
		valid  bool
	}{
		{"", 0, false},
		{"some/invalid/path", 0, false},
		{"/some/{invalid:bad_type}/path", 6, false},
		{"/some/invalid/path:string}", 18, false},
		{"some/invalid/path/", 0, false},
		{"/some/{invalid:bad_type}/path/", 6, false},
		{"/some/invalid/path:string}/", 18, false},
		{"//", 1, false},
		{"/some/valid/path", 16, true},
		{"/some/{valid:int}/path", 22, true},
		{"/some/valid/{path:string}", 25, true},
		{"/some/valid/path/", 17, true},
		{"/some/{valid:int}/path/", 23, true},
		{"/some/valid/{path:string}/", 26, true},
		{"/", 1, true},
	}

	for _, tt := range tests {
		parsed, err := model.ValidatePathValue(tt.path)
		require.Equal(t, tt.parsed, parsed, "wrong number of parsed characters: %d, expected: %d", parsed, tt.parsed)
		if tt.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func TestPathValid(t *testing.T) {
	parentSubdomain := &model.Subdomain{
		Name: "dummy-subdomain",
	}

	tests := []struct {
		path  *model.Path
		valid bool
	}{
		{&model.Path{Value: "some/invalid/path", Endpoint: "some_unknown_endpoint"}, false},
		{&model.Path{Value: "/some/valid/path", Endpoint: "some_unknown_endpoint"}, false},
		{&model.Path{Value: "some/invalid/path", Endpoint: "some_known_endpoint"}, false},
		{&model.Path{Value: "/some/valid/path", Endpoint: "some_known_endpoint"}, true},
	}

	for _, tt := range tests {
		err := tt.path.Validate(parentSubdomain, []string{"some_known_endpoint"})
		if tt.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func TestSubdomainValid(t *testing.T) {
	goodPaths := []model.Path{
		{Value: "/some/valid/path", Endpoint: "some_known_endpoint"},
		{Value: "/", Endpoint: "some_known_endpoint"},
	}
	badPaths := []model.Path{
		{Value: "/some/valid/path", Endpoint: "some_known_endpoint"},
		{Value: "some/invalid/path", Endpoint: "some_known_endpoint"},
	}

	tests := []struct {
		subdomain *model.Subdomain
		valid     bool
	}{
		{&model.Subdomain{Name: "bad.char.subdomain", Paths: badPaths}, false},
		{&model.Subdomain{Name: "-starts-with-hyphen-subdomain", Paths: badPaths}, false},
		{&model.Subdomain{Name: "ends-with-hyphen-subdomain-", Paths: badPaths}, false},
		{&model.Subdomain{Name: "-starts-ends-with-hyphen-subdomain-", Paths: badPaths}, false},
		{&model.Subdomain{Name: "valid-subdomain-37", Paths: badPaths}, false},
		{&model.Subdomain{Name: "bad.char.subdomain", Paths: goodPaths}, false},
		{&model.Subdomain{Name: "-starts-with-hyphen-subdomain", Paths: goodPaths}, false},
		{&model.Subdomain{Name: "ends-with-hyphen-subdomain-", Paths: goodPaths}, false},
		{&model.Subdomain{Name: "-starts-ends-with-hyphen-subdomain-", Paths: goodPaths}, false},
		{&model.Subdomain{Name: "valid-subdomain-37", Paths: goodPaths}, true},
	}

	for _, tt := range tests {
		err := tt.subdomain.Validate([]string{"some_known_endpoint"})
		if tt.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func TestNameValid(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"badName69", false},
		{"bad-name-69", false},
		{"bad+name&&", false},
		{"_bad_name_66", false},
		{"bad__name_66", false},
		{"1_bad_name", false},
		{"bad_name_", false},
		{"good_name_42", true},
	}

	for _, tt := range tests {
		err := model.ValidateName(tt.name, "dummy type")
		if tt.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func TestVariableValid(t *testing.T) {
	tests := []struct {
		variable *model.Variable
		valid    bool
	}{
		{&model.Variable{Name: "good_name_42", Type: "bad_type", Value: "dummy_value"}, false},
		{&model.Variable{Name: "good_name_42", Type: "string", Value: "dummy_value"}, true},
		{&model.Variable{Name: "good_name_42", Type: "string", Value: ""}, true},
		{&model.Variable{Name: "good_name_42", Type: "float", Value: "dummy_value"}, false},
		{&model.Variable{Name: "good_name_42", Type: "float", Value: "3.14f"}, false},
		{&model.Variable{Name: "good_name_42", Type: "float", Value: "a3.14"}, false},
		{&model.Variable{Name: "good_name_42", Type: "float", Value: "3e14"}, true},
		{&model.Variable{Name: "good_name_42", Type: "float", Value: "3.14"}, true},
		{&model.Variable{Name: "good_name_42", Type: "float", Value: "1000000000"}, true},
		{&model.Variable{Name: "good_name_42", Type: "float", Value: ""}, true},
		{&model.Variable{Name: "good_name_42", Type: "int", Value: "dummy_value"}, false},
		{&model.Variable{Name: "good_name_42", Type: "int", Value: "3.14f"}, false},
		{&model.Variable{Name: "good_name_42", Type: "int", Value: "a3.14"}, false},
		{&model.Variable{Name: "good_name_42", Type: "int", Value: "3e14"}, false},
		{&model.Variable{Name: "good_name_42", Type: "int", Value: "3.14"}, false},
		{&model.Variable{Name: "good_name_42", Type: "int", Value: "-1000000000"}, true},
		{&model.Variable{Name: "good_name_42", Type: "int", Value: "1000000000"}, true},
		{&model.Variable{Name: "good_name_42", Type: "int", Value: "37"}, true},
		{&model.Variable{Name: "good_name_42", Type: "int", Value: ""}, true},
		{&model.Variable{Name: "good_name_42", Type: "int", Value: "dummy_value"}, false},
		{&model.Variable{Name: "good_name_42", Type: "int", Value: "3.14f"}, false},
		{&model.Variable{Name: "good_name_42", Type: "int", Value: "a3.14"}, false},
		{&model.Variable{Name: "good_name_42", Type: "int", Value: "3e14"}, false},
		{&model.Variable{Name: "good_name_42", Type: "int", Value: "3.14"}, false},
		{&model.Variable{Name: "good_name_42", Type: "uint", Value: "-1000000000"}, false},
		{&model.Variable{Name: "good_name_42", Type: "uint", Value: "1000000000"}, true},
		{&model.Variable{Name: "good_name_42", Type: "uint", Value: "37"}, true},
		{&model.Variable{Name: "good_name_42", Type: "uint", Value: ""}, true},
	}

	for _, tt := range tests {
		err := tt.variable.Validate()
		if tt.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

var goodVariables = []model.Variable{
	{Name: "good_name_42", Type: "string", Value: "dummy_value"},
	{Name: "good_name_42", Type: "string", Value: ""},
}

var badVariables = []model.Variable{
	{Name: "good_name_42", Type: "float", Value: "1000000000"},
	{Name: "good_name_42", Type: "bad_type", Value: "dummy_value"},
}

func TestStructValid(t *testing.T) {
	tests := []struct {
		s     *model.Struct
		valid bool
	}{
		{&model.Struct{Name: "_bad_name", Fields: badVariables}, false},
		{&model.Struct{Name: "good_name", Fields: badVariables}, false},
		{&model.Struct{Name: "_bad_name", Fields: goodVariables}, false},
		{&model.Struct{Name: "good_name", Fields: goodVariables}, true},
	}

	for _, tt := range tests {
		err := tt.s.Validate()
		if tt.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func TestServiceCommonValid(t *testing.T) {
	knownDependencies := []string{"dep_1", "dep_2", "dep_3"}
	goodDependencies := []string{"dep_1", "dep_2"}
	badDependencies := []string{"dep_3", "dep_4"}

	tests := []struct {
		svcCommon *model.ServiceCommon
		valid     bool
	}{
		{&model.ServiceCommon{Name: "bad_service_name_", Port: "80000", Environment: badVariables, Dependencies: badDependencies}, false},
		{&model.ServiceCommon{Name: "bad_service_name_", Port: "80000", Environment: badVariables, Dependencies: goodDependencies}, false},
		{&model.ServiceCommon{Name: "bad_service_name_", Port: "80000", Environment: goodVariables, Dependencies: badDependencies}, false},
		{&model.ServiceCommon{Name: "bad_service_name_", Port: "80000", Environment: goodVariables, Dependencies: goodDependencies}, false},
		{&model.ServiceCommon{Name: "bad_service_name_", Port: "30000", Environment: badVariables, Dependencies: badDependencies}, false},
		{&model.ServiceCommon{Name: "bad_service_name_", Port: "30000", Environment: badVariables, Dependencies: goodDependencies}, false},
		{&model.ServiceCommon{Name: "bad_service_name_", Port: "30000", Environment: goodVariables, Dependencies: badDependencies}, false},
		{&model.ServiceCommon{Name: "bad_service_name_", Port: "30000", Environment: goodVariables, Dependencies: goodDependencies}, false},
		{&model.ServiceCommon{Name: "good_service_name", Port: "80000", Environment: badVariables, Dependencies: badDependencies}, false},
		{&model.ServiceCommon{Name: "good_service_name", Port: "80000", Environment: badVariables, Dependencies: goodDependencies}, false},
		{&model.ServiceCommon{Name: "good_service_name", Port: "80000", Environment: goodVariables, Dependencies: badDependencies}, false},
		{&model.ServiceCommon{Name: "good_service_name", Port: "80000", Environment: goodVariables, Dependencies: goodDependencies}, false},
		{&model.ServiceCommon{Name: "good_service_name", Port: "30000", Environment: badVariables, Dependencies: badDependencies}, false},
		{&model.ServiceCommon{Name: "good_service_name", Port: "30000", Environment: badVariables, Dependencies: goodDependencies}, false},
		{&model.ServiceCommon{Name: "good_service_name", Port: "30000", Environment: goodVariables, Dependencies: badDependencies}, false},
		{&model.ServiceCommon{Name: "good_service_name", Port: "30000", Environment: goodVariables, Dependencies: goodDependencies}, true},
	}

	for _, tt := range tests {
		err := tt.svcCommon.Validate(knownDependencies)
		if tt.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func TestExternalServiceValid(t *testing.T) {
	validSvcCommon := model.ServiceCommon{Name: "good_name", RepositoryURL: "good-name-repo", Port: "80"}

	tests := []struct {
		extSvc *model.ExternalService
		valid  bool
	}{
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "bad_url"}, false},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "/a/b/c"}, false},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: ":/a/b/c"}, false},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "http//a/b/c"}, false},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "http:/a/b/c"}, false},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker/whalesay"}, true},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker/whalesay:"}, false},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker/whalesay:lates"}, false},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker/whalesay:latest"}, true},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker/whalesay:masta"}, false},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker/whalesay:master"}, true},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker/whalesay:v2"}, true},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker/whalesay:v2."}, false},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker/whalesay:v2.3"}, true},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker/whalesay:v2.3."}, false},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker/whalesay:v2.3.255"}, true},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker/whalesay:2"}, true},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker/whalesay:2."}, false},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker/whalesay:2.3"}, true},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker/whalesay:2.3."}, false},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker/whalesay:2.3.255"}, true},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "docker.com/whalesay:v2.3"}, true},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "hub.docker.com/whalesay:v2.3"}, true},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "http://hub.docker.com:8080/whalesay:v2.3"}, true},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "https://hub.docker.com/whalesay:v2.3"}, true},
		{&model.ExternalService{ServiceCommon: validSvcCommon, ImageURL: "localhost:32000/whalesay:v2.3"}, true},
	}

	for _, tt := range tests {
		err := tt.extSvc.Validate([]string{})
		if tt.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}
