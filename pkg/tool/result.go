package tool

import "serverless-dbapi/pkg/exception"

type Result[T any] struct {
	Data any
	Err  *exception.BaseError
}

func SuccessResult(T any) Result[any] {
	return Result[any]{
		Data: T,
		Err:  nil,
	}
}

func ErrorResult(err *exception.BaseError) Result[any] {
	return Result[any]{
		Data: nil,
		Err:  err,
	}
}

func SimpleErrorResult(code int, msg string) Result[any] {
	return Result[any]{
		Data: nil,
		Err:  exception.New(code, msg),
	}
}

func (r Result[T]) IsError() bool {
	return r.Err != nil
}
