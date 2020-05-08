package logic

import (
	"errors"

	"go.uber.org/zap"

	"foo-service/pkg/structs"
)

// ExampleAPI is an example, trivial implementation of the API interface.
// It simply logs the request name.
type ExampleAPI struct {
	logger *zap.Logger
}

// NewAPI creates an instance of the example API implementation.
func NewAPI(logger *zap.Logger) API {
	return &ExampleAPI{logger: logger}
}

// /method_no_path_params

// MethodNoPathParams0 example.
func (a *ExampleAPI) MethodNoPathParams0() (*structs.ReturnType, error) {
	a.logger.Info("called method_no_path_params_0")
	return nil, errors.New("method 'method_no_path_params_0' not implemented")
}

// MethodNoPathParams1 example.
func (a *ExampleAPI) MethodNoPathParams1(
	*structs.BodyType,
) (*structs.ReturnType, error) {
	a.logger.Info("called method_no_path_params_1")
	return nil, errors.New("method 'method_no_path_params_1' not implemented")
}

// MethodNoPathParams2 example.
func (a *ExampleAPI) MethodNoPathParams2(
	string,
	float64,
	int64,
) (*structs.ReturnType, error) {
	a.logger.Info("called method_no_path_params_2")
	return nil, errors.New("method 'method_no_path_params_2' not implemented")
}

// MethodNoPathParams3 example.
func (a *ExampleAPI) MethodNoPathParams3(
	*structs.BodyType,
	string,
	float64,
	int64,
) (*structs.ReturnType, error) {
	a.logger.Info("called method_no_path_params_3")
	return nil, errors.New("method 'method_no_path_params_3' not implemented")
}

// MethodNoPathParams4 example.
func (a *ExampleAPI) MethodNoPathParams4(
	int64,
	float64,
	string,
) (*structs.ReturnType, error) {
	a.logger.Info("called method_no_path_params_4")
	return nil, errors.New("method 'method_no_path_params_4' not implemented")
}

// MethodNoPathParams5 example.
func (a *ExampleAPI) MethodNoPathParams5(
	*structs.BodyType,
	int64,
	float64,
	string,
) (*structs.ReturnType, error) {
	a.logger.Info("called method_no_path_params_5")
	return nil, errors.New("method 'method_no_path_params_5' not implemented")
}

// MethodNoPathParams6 example.
func (a *ExampleAPI) MethodNoPathParams6(
	string,
	float64,
	int64,
	int64,
	float64,
	string,
) (*structs.ReturnType, error) {
	a.logger.Info("called method_no_path_params_6")
	return nil, errors.New("method 'method_no_path_params_6' not implemented")
}

// MethodNoPathParams7 example.
func (a *ExampleAPI) MethodNoPathParams7(
	*structs.BodyType,
	string,
	float64,
	int64,
	int64,
	float64,
	string,
) (*structs.ReturnType, error) {
	a.logger.Info("called method_no_path_params_7")
	return nil, errors.New("method 'method_no_path_params_7' not implemented")
}

// /method/{path_param_0:int}/{path_param_1:string}

// Method0 example.
func (a *ExampleAPI) Method0(
	int64,
	string,
) (*structs.ReturnType, error) {
	a.logger.Info("called method_0")
	return nil, errors.New("method 'method_0' not implemented")
}

// Method1 example.
func (a *ExampleAPI) Method1(
	*structs.BodyType,
	int64,
	string,
) (*structs.ReturnType, error) {
	a.logger.Info("called method_1")
	return nil, errors.New("method 'method_1' not implemented")
}

// Method2 example.
func (a *ExampleAPI) Method2(
	int64,
	string,
	string,
	float64,
	int64,
) (*structs.ReturnType, error) {
	a.logger.Info("called method_2")
	return nil, errors.New("method 'method_2' not implemented")
}

// Method3 example.
func (a *ExampleAPI) Method3(
	*structs.BodyType,
	int64,
	string,
	string,
	float64,
	int64,
) (*structs.ReturnType, error) {
	a.logger.Info("called method_3")
	return nil, errors.New("method 'method_3' not implemented")
}

// Method4 example.
func (a *ExampleAPI) Method4(
	int64,
	string,
	int64,
	float64,
	string,
) (*structs.ReturnType, error) {
	a.logger.Info("called method_4")
	return nil, errors.New("method 'method_4' not implemented")
}

// Method5 example.
func (a *ExampleAPI) Method5(
	*structs.BodyType,
	int64,
	string,
	int64,
	float64,
	string,
) (*structs.ReturnType, error) {
	a.logger.Info("called method_5")
	return nil, errors.New("method 'method_5' not implemented")
}

// Method6 example.
func (a *ExampleAPI) Method6(
	int64,
	string,
	string,
	float64,
	int64,
	int64,
	float64,
	string,
) (*structs.ReturnType, error) {
	a.logger.Info("called method_6")
	return nil, errors.New("method 'method_6' not implemented")
}

// Method7 example.
func (a *ExampleAPI) Method7(
	*structs.BodyType,
	int64,
	string,
	string,
	float64,
	int64,
	int64,
	float64,
	string,
) (*structs.ReturnType, error) {
	a.logger.Info("called method_7")
	return nil, errors.New("method 'method_7' not implemented")
}
