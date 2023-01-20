package tool

import (
	"fmt"
	"serverless-dbapi/pkg/exception"
)

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

func ErrorResult(err *exception.BaseError, args ...string) Result[any] {
	if len(args) > 0 {
		err.Msg = fmt.Sprintf(err.Msg, args)
	}
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
