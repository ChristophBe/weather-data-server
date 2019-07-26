package handlers

import (
	"../data"
	"../jwt"
	"net/http"
)


func FetchNodesHandler(w http.ResponseWriter , r * http.Request){

	defer recoverHandlerErrors(w)




	userId,err := jwt.GetUserIdBy(r)

	var nodes []data.MeasuringNode
	nodeRepo := data.MeasuringNodeRepository{}

	if err != nil {
		//user is not authorized
		nodes ,err = nodeRepo.FetchAllPublicNodes()
		panicIfErrorNonNil(err, "can not fetch Nodes",http.StatusInternalServerError)

	} else {
		nodes ,err = nodeRepo.FetchAllVisibleNodesByUserId(userId)
		panicIfErrorNonNil(err, "can not fetch Nodes",http.StatusInternalServerError)
	}

	err = writeJsonResponse(nodes, w)
	panicIfErrorNonNil(err, "unexpected Error", http.StatusInternalServerError)
}