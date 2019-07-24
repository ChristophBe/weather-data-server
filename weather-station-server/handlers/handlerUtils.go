package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"
)


const(
	InvalidBody = "invalid Body"
	NotAuthorized = "user not Authorized"
)


type handlerError struct {
	Err          error     `json:"-"`
	ErrorMessage string    `json:"error"`
	Timestamp    time.Time `json:"timestamp"`
	Status    	 int	   `json:"-"`
}

func writeJsonResponse( data interface{}, w http.ResponseWriter) error{

	jsonObject,err:= json.Marshal(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Write(jsonObject)
	w.Header().Set("Content-Type", "application/json")

	return nil
}

func handleError(w http.ResponseWriter, error handlerError, httpStatus int){

	error.Timestamp= time.Now()
	log.Println(error)


	jsonObj ,err  := json.Marshal(error)


	if err != nil {
		log.Fatal("can not marshal error")
	}

	http.Error(w, string(jsonObj), httpStatus)
	w.Header().Set("Content-Type", "application/json")
	//panic(error.Err)
}

func readBody(r *http.Request , item interface{}) error {


	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, item)
	if err != nil {
		return err
	}
	return nil
}


func recoverHandlerErrors(w http.ResponseWriter) {
	if r := recover(); r != nil {
		if reflect.TypeOf(r).String() == "handlerError" {
			handleError(w, r.(handlerError), r.(handlerError).Status)
		} else {
			log.Print(r)
			handleError(w, handlerError{Err: nil, ErrorMessage: "unexpected Failure"}, http.StatusBadRequest)
		}

	}

}

func panicIfErrorNonNil(err error, errorMessage string, status int){
	if err != nil {
		panic(handlerError{Err:err,ErrorMessage:errorMessage,Status:status})
	}
}