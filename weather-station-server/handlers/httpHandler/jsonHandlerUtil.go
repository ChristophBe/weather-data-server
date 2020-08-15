package httpHandler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

func ReadJsonBody(r *http.Request, bodyData interface{}) (err error) {
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()
	if err != nil {
		err = BadRequest("Invalid Body", err)
		return
	}

	if err = json.Unmarshal(body, bodyData); err != nil {
		err = BadRequest("Invalid Body", err)
	}

	return
}
func ReadPathVariableInt(r * http.Request, name string)(value int64, err error)  {
	vars := mux.Vars(r)
	value, err = strconv.ParseInt(vars[name], 10, 64)
	return
}