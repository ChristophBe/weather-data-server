package httpHandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ReadJsonBody(r* http.Request , bodyData interface{})  {
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()
	if err != nil {
		HandleBadRequest("Invalid Body", err)
	}


	if err = json.Unmarshal(body, bodyData); err != nil {
		HandleBadRequest("Invalid Body", err)
	}

	return
}
