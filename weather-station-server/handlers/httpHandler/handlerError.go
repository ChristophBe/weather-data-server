package httpHandler

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
func HandleForbidden(message string, cause error){
	panic(Forbidden(message,cause))
}

func NotFound(message string, cause error) HandlerError{
	return HandlerError{
		Message:    message,
		StatusCode: http.StatusNotFound,
		Cause:      cause,
	}
}
func HandleNotFound(message string, cause error){
	panic(NotFound(message,cause))
}

func InternalError(err error)  HandlerError{
	return HandlerError{
		Message:    "unexpected Error",
		StatusCode: http.StatusInternalServerError,
		Cause:      err,
	}
}
func HandleInternalError(cause error){
	panic(InternalError(cause))
}

func BadRequest(message string, cause error) HandlerError{
	return HandlerError{
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Cause:      cause,
	}
}
func HandleBadRequest(message string, cause error){
	panic(BadRequest(message,cause))
}



