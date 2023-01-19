package exception

// common error
type BaseError struct {
	Code int
	Msg  string
}

func New(code int, msg string) *BaseError {
	return &BaseError{
		Code: code,
		Msg:  msg,
	}
}
