package httpHandler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type HandlerResponse struct {
	Data   interface{}
	Status int
}

type JsonHandler func( r * http.Request) (HandlerResponse, error)

func (fn JsonHandler) catchErrors(w http.ResponseWriter, r* http.Request)  {

	recovered := recover()
	if recovered == nil{
		return
	}

	handlerError, ok := recovered.(HandlerError)
	if ok {
		fn.Error(w, handlerError.Message, handlerError.StatusCode)
		log.Printf("Request failed: {url: %s,methonde:%s, message:%s,status: %d, cause:%v", r.URL,r.Method,handlerError.Message,handlerError.StatusCode,handlerError.Cause)
		return
	}

	fn.Error(w,"unexpected Error",http.StatusInternalServerError)
	log.Printf("Request failed with unknown cause: {url: %s,methonde:%s, cause:%v", r.URL,r.Method,recovered)

	return
}

func (fn JsonHandler)ServeHTTP(w http.ResponseWriter,r* http.Request)  {
	defer fn.catchErrors(w,r)

	response, err := fn(r)
	if err != nil {
		panic(err)
	}

	status := response.Status
	if status == 0{
		status = http.StatusOK
	}

	fn.writeJsonResponse(response.Data,status, w)
}



func (fn JsonHandler) Error(w http.ResponseWriter,message string, code int )  {
	response := struct {
		Message string `json:"message"`
		Timestamp time.Time `json:"timestamp"`
	}{
		Message: message,
		Timestamp: time.Now(),
	}
	fn.writeJsonResponse(response,code,w)
}



func (fn JsonHandler) writeJsonResponse(resp interface{},statusCode int, w http.ResponseWriter)  {
	jsonResponse ,err:= json.Marshal(resp)

	if	err != nil {
		fn.Error(w,"unexpected error", http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	_,err = w.Write(jsonResponse)
	if	err != nil {
		fn.Error(w,"unexpected error", http.StatusInternalServerError)
		return
	}
}
