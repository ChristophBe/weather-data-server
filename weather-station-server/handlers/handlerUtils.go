package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type handlerError struct {
	Err          error     `json:"-"`
	ErrorMessage string    `json:"error"`
	Timestamp    time.Time `json:"timestamp"`
	status    	 int	   `json:"-"`
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