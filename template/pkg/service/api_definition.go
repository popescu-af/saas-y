package service

import "github.com/popescu-af/saas-y/template/pkg/service/structs"

// API defines the operations supported by the HTTP service.
type API interface {
	AddSomething(structs.Something) (interface{}, error)
}
