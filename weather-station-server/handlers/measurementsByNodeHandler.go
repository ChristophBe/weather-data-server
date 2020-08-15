package handlers

import (
	"de.christophb.wetter/data/database"
	"de.christophb.wetter/data/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)


func getNodeIDFormRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	nodeId, err := strconv.ParseInt(vars["nodeId"], 10, 64)
	return nodeId, err
}
func getNodeFormRequest(r *http.Request) (node models.MeasuringNode, err error) {
	nodeId, err := getNodeIDFormRequest(r)
	if err != nil {
		return
	}
	nodeRepo := database.GetMeasuringNodeRepository()

	node, err = nodeRepo.FetchMeasuringNodeById(nodeId)
	return
}
