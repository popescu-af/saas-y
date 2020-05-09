package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/popescu-af/saas-y/pkg/model"
)

// Path
// Subdomain
// Variable
// Struct
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
