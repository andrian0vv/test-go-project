package errors

type codifiedError struct {
	code string
	err  error
}

func WithCode(err error, code string) *codifiedError {
	return &codifiedError{err: err, code: code}
}

func (e *codifiedError) Error() string {
	return e.err.Error()
}

func (e *codifiedError) Code() string {
	return e.code
}

func (e *codifiedError) Unwrap() error {
	return e.err
}
