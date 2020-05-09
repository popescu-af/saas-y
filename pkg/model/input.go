package model

import (
	"fmt"
)

// Spec is the saas-y specification.
type Spec struct {
	Domain           string            `json:"domain"`
	Subdomains       []Subdomain       `json:"subdomains"`
	Services         []Service         `json:"services"`
	ExternalServices []ExternalService `json:"external_services"`
}

// Validate checks a specification.
func (s *Spec) Validate() (err error) {
	// TODO: Domain validation

	// TODO: list of knonw services
	var knownServices []string

	for _, svc := range s.Services {
		if err = svc.Validate(); err != nil {
			return
		}
	}
	for _, esvc := range s.ExternalServices {
		if err = esvc.Validate(); err != nil {
			return
		}
	}
	for _, subd := range s.Subdomains {
		if err = subd.Validate(knownServices); err != nil {
			return
		}
	}
	return
}

// Subdomain is a subdomain entry in the specification.
type Subdomain struct {
	Name  string `json:"name"`
	Paths []Path `json:"paths"`
}

// Validate checks a subdomain.
func (s *Subdomain) Validate(knownServices []string) (err error) {
	// TODO: Name validation

	for _, p := range s.Paths {
		if err = p.Validate(s, knownServices); err != nil {
			return
		}
	}
	return
}

// Path represents a URL path.
type Path struct {
	Value    string `json:"value"`
	Endpoint string `json:"endpoint"`
}

// Validate checks if the path has a proper value and
// ends up at a known service, which is defined in the spec
// as a new service or as an external service.
func (p *Path) Validate(parent *Subdomain, knownServices []string) (err error) {
	// TODO: Value validation

	for _, s := range knownServices {
		if p.Endpoint == s {
			return
		}
	}

	err = fmt.Errorf("cannot validate path %s, unknown endpoint %s", p.Value, p.Endpoint)
	return
}

// Service represents a saas-y defined service.
type Service struct {
	ServiceCommon
	API     []API    `json:"api"`
	Structs []Struct `json:"structs"`
}

// Validate checks if the service is well defined.
func (s *Service) Validate() (err error) {
	// TODO: validation

	return
}

// API represents a saas-y defined API.
type API struct {
	Path    string            `json:"path"`
	Methods map[string]Method `json:"methods"`
}

// Validate checks if the API is well defined.
func (a *API) Validate() (err error) {
	// TODO: validation

	return
}

// APIMethodType is the type for saas-y API methods.
type APIMethodType string

const (
	// GET is the HTTP GET type
	GET APIMethodType = "get"
	// POST is the HTTP POST type
	POST APIMethodType = "post"
	// PATCH is the HTTP PATCH type
	PATCH APIMethodType = "patch"
	// DELETE is the HTTP DELETE type
	DELETE APIMethodType = "delete"
)

// Method represents a saas-y API method.
type Method struct {
	Type         APIMethodType `json:"type"`
	HeaderParams []Variable    `json:"header_params"`
	QueryParams  []Variable    `json:"query_params"`
	InputType    string        `json:"input_type"`
	ReturnType   string        `json:"return_type"`
}

// Validate checks if the method is well defined.
func (m *Method) Validate(knownTypes []string) (err error) {
	typeOK := false
	for _, t := range []APIMethodType{GET, POST, PATCH, DELETE} {
		if m.Type == t {
			typeOK = true
			break
		}
	}

	if !typeOK {
		err = fmt.Errorf("invalid method type %s", m.Type)
		return
	}

	for _, p := range m.HeaderParams {
		if err = p.Validate(); err != nil {
			return
		}
	}
	for _, p := range m.QueryParams {
		if err = p.Validate(); err != nil {
			return
		}
	}

	for _, t := range []APIMethodType{GET, DELETE} {
		if m.Type == t && m.InputType != "" {
			err = fmt.Errorf("body is not allowed for method type %s", m.Type)
			return
		}
	}

	foundInputType := false
	if m.InputType == "" {
		foundInputType = true
	}

	foundReturnType := false
	if m.ReturnType == "" {
		foundReturnType = true
	}

	for _, t := range knownTypes {
		if m.InputType == t {
			foundInputType = true
		}
		if m.ReturnType == t {
			foundReturnType = true
		}
	}

	if !foundInputType {
		err = fmt.Errorf("unknown type %s", m.InputType)
		return
	}

	if !foundReturnType {
		err = fmt.Errorf("unknown type %s", m.ReturnType)
		return
	}

	return
}

// Variable represents an environment / struct variable
// or a header / query param.
type Variable struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Validate checks if the variable is well defined.
func (v *Variable) Validate() (err error) {
	// TODO: validation

	return
}

// Struct represents an API struct.
type Struct struct {
	Name   string     `json:"name"`
	Fields []Variable `json:"fields"`
}

// Validate checks if the struct is well defined.
func (s *Struct) Validate() (err error) {
	// TODO: validation

	return
}

// ExternalService defines a service that is defined outside of the spec.
type ExternalService struct {
	ServiceCommon
	ImageURL string `json:"image_url"`
}

// Validate checks if the external service is well defined.
func (s *ExternalService) Validate() (err error) {
	// TODO: validation

	return
}

// ServiceCommon contains the core attributes of both saas-y and external services.
type ServiceCommon struct {
	Name         string     `json:"name"`
	Port         string     `json:"port"`
	Environment  []Variable `json:"env"`
	Dependencies []string   `json:"dependencies"`
}

// Validate checks if the service core attributes are well defined.
func (s *ServiceCommon) Validate() (err error) {
	// TODO: validation

	return
}
