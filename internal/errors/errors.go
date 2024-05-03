package errors

import "errors"

type code interface {
	Code() string
}

func Code(err error) (string, bool) {
	var e code
	if errors.As(err, &e) {
		return e.Code(), true
	}

	return "", false
}
