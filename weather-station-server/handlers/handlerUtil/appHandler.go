package handlerUtil

import (
	"encoding/json"
	"log"
	"net/http"
)

type AppHandler func( r * http.Request) (interface{}, int,error)



func catchErrors(w http.ResponseWriter, r* http.Request)  {
	recovered := recover()
	if recovered == nil{
		return
	}

	handlerError, ok := recovered.(HandlerError)
	if ok {
		http.Error(w, handlerError.Message, handlerError.StatusCode)
		log.Printf("Request failed: {url: %s,methonde:%s, message:%s,status: %d, cause:%v", r.URL,r.Method,handlerError.Message,handlerError.StatusCode,handlerError.Cause)
		return
	}

	http.Error(w,"unexpected error",http.StatusInternalServerError)
	log.Printf("Request failed with unknown cause: {url: %s,methonde:%s, cause:%v", r.URL,r.Method,recovered)

	return

}

func (fn AppHandler)ServeHTTP(w http.ResponseWriter,r* http.Request)  {
	defer catchErrors(w,r)

	resp, statusCode,err := fn(r)
	if err!= nil{
		panic(err)
	}
	writeJsonResponse(resp,statusCode, w)
}



func writeJsonResponse(resp interface{},statusCode int, w http.ResponseWriter)  {
	jsonResponse ,err:= json.Marshal(resp)
	if	err != nil {
		panic(err)
	}

	if _, err = w.Write(jsonResponse);err !=nil{
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(statusCode)


}
