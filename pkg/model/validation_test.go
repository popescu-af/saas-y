package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/popescu-af/saas-y/pkg/model"
)

// ServiceCommon
// ExternalService

// API
// Service

// Spec

func TestMethodTypeValid(t *testing.T) {
	tests := []struct {
		method *model.Method
		valid  bool
	}{
		{&model.Method{Type: "dummy"}, false},
		{&model.Method{Type: "get"}, true},
		{&model.Method{Type: "post"}, true},
		{&model.Method{Type: "patch"}, true},
		{&model.Method{Type: "delete"}, true},
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

func TestMethodBodyAllowed(t *testing.T) {
	tests := []struct {
		method *model.Method
		valid  bool
	}{
		{&model.Method{Type: "get"}, true},
		{&model.Method{Type: "post"}, true},
		{&model.Method{Type: "patch"}, true},
		{&model.Method{Type: "delete"}, true},
		{&model.Method{Type: "get", InputType: "whatever"}, false},
		{&model.Method{Type: "post", InputType: "whatever"}, true},
		{&model.Method{Type: "patch", InputType: "whatever"}, true},
		{&model.Method{Type: "delete", InputType: "whatever"}, false},
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
		{&model.Method{Type: "post", InputType: "something_unknown"}, false},
		{&model.Method{Type: "post", InputType: "something_known"}, true},
		{&model.Method{Type: "post", ReturnType: "something_unknown"}, false},
		{&model.Method{Type: "post", ReturnType: "something_known"}, true},
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
		{&model.Variable{Name: "good_name_42", Type: "int", Value: "1000000000"}, true},
		{&model.Variable{Name: "good_name_42", Type: "int", Value: "37"}, true},
		{&model.Variable{Name: "good_name_42", Type: "int", Value: ""}, true},
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

func TestStructValid(t *testing.T) {
	goodVariables := []model.Variable{
		{Name: "good_name_42", Type: "string", Value: "dummy_value"},
		{Name: "good_name_42", Type: "string", Value: ""},
	}
	badVariables := []model.Variable{
		{Name: "good_name_42", Type: "float", Value: "1000000000"},
		{Name: "good_name_42", Type: "bad_type", Value: "dummy_value"},
	}

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
