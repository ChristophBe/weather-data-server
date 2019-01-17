package handlers

import (
	"../data"
	"encoding/json"
	"net/http"
)

func NodesHandler(w http.ResponseWriter, r *http.Request) {

	con := data.CreateConnection()
	defer con.Close()


	switch r.Method {
	case http.MethodGet:
		nodes := data.FetchAllMeasuringNodes(con)


		jsonObj, err := json.Marshal(nodes)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonObj)
		return

	}

	http.Error(w,"Method not allowed",http.StatusMethodNotAllowed)




}