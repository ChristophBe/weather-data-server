package handlers

import (
	"../data"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func MeasurementsByNodeHandler(w http.ResponseWriter, r *http.Request) {
	con := data.CreateConnection()
	defer con.Close()

	vars := mux.Vars(r)

	nodeId, err := strconv.ParseInt(vars["nodeId"], 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case  http.MethodGet:



		measurements := data.FetchAllMeasuringByNodeId(con, nodeId)

		jsonObj, err := json.Marshal(measurements)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonObj)
		return

	case http.MethodPost:
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Unmarshal
		var measuring data.Measuring
		err = json.Unmarshal(b, &measuring)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		measuring.TimeStamp = time.Now()

		data.CreateMeasurement(con,nodeId, measuring)
		return
	}

	http.Error(w,"Method not allowed",http.StatusMethodNotAllowed)
}
