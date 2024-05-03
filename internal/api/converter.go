package api

import ie "github.com/andrian0vv/test-go-project/internal/errors"

func ToError(err error) ErrorResponse {
	code, ok := ie.Code(err)
	if !ok {
		code = codeInternalError
	}

	return ErrorResponse{
		Error: &struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    code,
			Message: err.Error(),
		},
	}
}
