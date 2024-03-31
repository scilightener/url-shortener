package api

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "Ok"
	StatusError = "Error"

	validateRequired = "required"
	validateURL      = "url"
)

func OkResponse() Response {
	return Response{Status: StatusOk}
}

func ErrResponse(err string) Response {
	return Response{Status: StatusError, Error: err}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case validateRequired:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case validateURL:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
