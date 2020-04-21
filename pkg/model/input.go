package model

type Spec struct {
	Domain           string            `json:"domain"`
	Subdomains       []Subdomain       `json:"subdomains"`
	Services         []Service         `json:"services"`
	ExternalServices []ExternalService `json:"external_services"`
}

type Subdomain struct {
	Name  string `json:"name"`
	Paths []Path `json:"paths"`
}

type Path struct {
	Value    string `json:"value"`
	Endpoint string `json:"endpoint"`
}

type Service struct {
	serviceCommon
	API     []API    `json:"api"`
	Structs []Struct `json:"structs"`
}

type API struct {
	Path    string            `json:"path"`
	Methods map[string]Method `json:"methods"`
}

type APIMethodType string

const (
	GET     APIMethodType = "get"
	POST    APIMethodType = "post"
	PATCH   APIMethodType = "patch"
	DELETE  APIMethodType = "delete"
	OPTIONS APIMethodType = "options"
)

type Method struct {
	Type         APIMethodType `json:"type"`
	HeaderParams []Variable    `json:"header_params"`
	QueryParams  []Variable    `json:"query_params"`
	InputType    string        `json:"input_type"`
	ReturnType   string        `json:"return_type"`
}

type Variable struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Struct struct {
	Name   string     `json:"name"`
	Fields []Variable `json:"fields"`
}

type ExternalService struct {
	serviceCommon
	ImageURL string `json:"image_url"`
}

type serviceCommon struct {
	Name         string     `json:"name"`
	Port         string     `json:"port"`
	Environment  []Variable `json:"env"`
	Dependencies []string   `json:"dependencies"`
}
