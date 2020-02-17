package httpErrorHandling

import (
	"de.christophb.wetter/handlers"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HttpHandlerError struct {
	Message string
	Cause   error
	Status  int
	Request *http.Request
}

type ErrorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Message string 		`json:"message"`
}
func (e HttpHandlerError) Error()string {
	return fmt.Sprintf("HttpHandlerError: %s, %s, %d, %s, %v",e.Request.URL,e.Request.Method,e.Status,e.Message,e.Cause)
}

// The error
func ErrorHandlingMiddleWare(handler http.Handler) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer errorHandler(writer, request)
		handler.ServeHTTP(writer,request)
	})
}

func errorHandler(writer http.ResponseWriter, request *http.Request) {
	recovered := recover()

	if recovered != nil{
		httpHandlerError,ok := recovered.(HttpHandlerError)

		if !ok {
			err,ok := recovered.(error)
			if ok {
				httpHandlerError = HttpHandlerError{
					Message: "unexpected failure",
					Cause:   err,
					Status:  http.StatusInternalServerError,
					Request: request,
				}
			}
		}
		if ok {
			log.Println(httpHandlerError)
			writerHttpErrorHandler(writer,httpHandlerError)
			return
		}
	}
	panic(recovered)
}

func writerHttpErrorHandler(writer http.ResponseWriter, handlerError HttpHandlerError) {
	resp := ErrorResponse{
		Timestamp: time.Now(),
		Message:   handlerError.Message,
	}
	err := handlers.WriteJsonResponse(resp,writer)
	panic(err)
}
