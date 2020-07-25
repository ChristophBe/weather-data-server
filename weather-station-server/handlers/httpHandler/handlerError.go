package handlerUtil

import (
	"fmt"
	"net/http"
)

type HandlerError struct {
	Message string
	StatusCode int
	Cause error
}
func (e HandlerError) Error() string {
	return fmt.Sprintf("Message: %s StatusCode: %d",e.Message,e.StatusCode)
}

func Forbidden(message string, cause error) HandlerError{
	return HandlerError{
		Message:    message,
		StatusCode: http.StatusForbidden,
		Cause:      cause,
	}
}

func NotFound(message string, cause error) HandlerError{
	return HandlerError{
		Message:    message,
		StatusCode: http.StatusNotFound,
		Cause:      cause,
	}
}

func InternalError(err error)  HandlerError{
	return HandlerError{
		Message:    "unexpected Error",
		StatusCode: http.StatusInternalServerError,
		Cause:      err,
	}
}

func BadRequest(message string, cause error) HandlerError{
	return HandlerError{
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Cause:      cause,
	}
}




