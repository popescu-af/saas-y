package logic

import "foo-service/pkg/structs"

// API defines the operations supported by the foo-service service.
type API interface {
	// /method_no_path_params
	MethodNoPathParams0() (*structs.ReturnType, error)

	MethodNoPathParams1(
		*structs.BodyType,
	) (*structs.ReturnType, error)

	MethodNoPathParams2(
		string,
		float64,
		int64,
	) (*structs.ReturnType, error)

	MethodNoPathParams3(
		*structs.BodyType,
		string,
		float64,
		int64,
	) (*structs.ReturnType, error)

	MethodNoPathParams4(
		int64,
		float64,
		string,
	) (*structs.ReturnType, error)

	MethodNoPathParams5(
		*structs.BodyType,
		int64,
		float64,
		string,
	) (*structs.ReturnType, error)

	MethodNoPathParams6(
		string,
		float64,
		int64,
		int64,
		float64,
		string,
	) (*structs.ReturnType, error)

	MethodNoPathParams7(
		*structs.BodyType,
		string,
		float64,
		int64,
		int64,
		float64,
		string,
	) (*structs.ReturnType, error)

	// /method/{path_param_0:int}/{path_param_1:string}
	Method0(
		int64,
		string,
	) (*structs.ReturnType, error)

	Method1(
		*structs.BodyType,
		int64,
		string,
	) (*structs.ReturnType, error)

	Method2(
		int64,
		string,
		string,
		float64,
		int64,
	) (*structs.ReturnType, error)

	Method3(
		*structs.BodyType,
		int64,
		string,
		string,
		float64,
		int64,
	) (*structs.ReturnType, error)

	Method4(
		int64,
		string,
		int64,
		float64,
		string,
	) (*structs.ReturnType, error)

	Method5(
		*structs.BodyType,
		int64,
		string,
		int64,
		float64,
		string,
	) (*structs.ReturnType, error)

	Method6(
		int64,
		string,
		string,
		float64,
		int64,
		int64,
		float64,
		string,
	) (*structs.ReturnType, error)

	Method7(
		*structs.BodyType,
		int64,
		string,
		string,
		float64,
		int64,
		int64,
		float64,
		string,
	) (*structs.ReturnType, error)
}