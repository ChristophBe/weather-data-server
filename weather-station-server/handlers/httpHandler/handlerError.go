package httpHandler

import (
	"fmt"
	"net/http"
)

const (
	ErrorMessageNotAuthorized   = "invalid or insufficient authorization"
	ErrorMessageNotFound        = "not found"
	ErrorMessageNotFoundf       = "%v not found"
	ErrorMessageUnexpectedError = "unexpected error"
	ErrorMessageInvalidBody     = "invalid body"
	ErrorMessageParameterf      = "invalid value for parameter %v"
)

type HandlerError struct {
	Message    string
	StatusCode int
	Cause      error
}

func (e HandlerError) Error() string {
	return fmt.Sprintf("Message: %s StatusCode: %d", e.Message, e.StatusCode)
}
func (e HandlerError) Unwrap() error {
	return e.Cause
}

func Forbidden(message string, cause error) HandlerError {
	return HandlerError{
		Message:    message,
		StatusCode: http.StatusForbidden,
		Cause:      cause,
	}
}

func NotFound(message string, cause error) HandlerError {
	return HandlerError{
		Message:    message,
		StatusCode: http.StatusNotFound,
		Cause:      cause,
	}
}

func InternalError(err error) HandlerError {
	return HandlerError{
		Message:    ErrorMessageUnexpectedError,
		StatusCode: http.StatusInternalServerError,
		Cause:      err,
	}
}

func BadRequest(message string, cause error) HandlerError {
	return HandlerError{
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Cause:      cause,
	}
}
