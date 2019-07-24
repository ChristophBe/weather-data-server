package handlers

import (
	"../data"
	"net/http"
)


func FetchNodesHandler(w http.ResponseWriter , r * http.Request){
	con := data.CreateConnection()
	defer con.Close()

	nodeRepo := data.MeasuringNodeRepository{}
	nodes := nodeRepo.FetchAllMeasuringNodes(con)
	writeJsonResponse(nodes, w)
}