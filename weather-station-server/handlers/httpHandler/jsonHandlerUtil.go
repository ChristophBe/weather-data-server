package httpHandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ReadJsonBody(r* http.Request , bodyData interface{})  (err error){
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
