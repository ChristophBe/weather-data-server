package handlers

import (
	"de.christophb.wetter/data/database"
	"de.christophb.wetter/data/models"
	"de.christophb.wetter/jwt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)


func FetchNodesHandler(w http.ResponseWriter , r * http.Request){

	defer recoverHandlerErrors(w)

	userId,err := jwt.GetUserIdByRequest(r)

	var nodes []models.MeasuringNode
	nodeRepo := database.GetMeasuringNodeRepository()

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


func FetchNodesByOwnerHandler(w http.ResponseWriter , r * http.Request){

	defer recoverHandlerErrors(w)

	vars := mux.Vars(r)
	userReqId, err := strconv.ParseInt(vars["userId"], 10, 64)
	panicIfErrorNonNil(err, "not found",http.StatusNotFound)

	userId,err := jwt.GetUserIdByRequest(r)
	panicIfErrorNonNil(err, NotAuthorized,http.StatusForbidden)

	if userId != userReqId {
		panic(handlerError{ErrorMessage:NotAuthorized,Status: http.StatusForbidden})
	}

	var nodes []models.MeasuringNode
	nodeRepo := database.GetMeasuringNodeRepository()


	nodes ,err = nodeRepo.FetchNodesOwnedByUserId(userId)
	panicIfErrorNonNil(err, "can not fetch Nodes",http.StatusInternalServerError)


	err = writeJsonResponse(nodes, w)
	panicIfErrorNonNil(err, "unexpected Error", http.StatusInternalServerError)
}

