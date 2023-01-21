package tool

import (
	"fmt"
	"serverless-dbapi/pkg/exception"
)

type Result[T any] struct {
	Data any
	Err  *exception.BaseError
}

func SuccessResult[T any](data T) Result[T] {
	return Result[T]{
		Data: data,
		Err:  nil,
	}
}

func ErrorResult[T any](err *exception.BaseError, args ...string) Result[T] {
	if len(args) > 0 {
		err.Msg = fmt.Sprintf(err.Msg, args)
	}
	return Result[T]{
		Data: nil,
		Err:  err,
	}
}

func SimpleErrorResult[T any](code int, msg string) Result[T] {
	return Result[T]{
		Data: nil,
		Err:  exception.New(code, msg),
	}
}

func (r Result[T]) IsError() bool {
	return r.Err != nil
}
