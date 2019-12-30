package handlers

import (
	"de.christophb.wetter/data/models"
	"net/http"
)
import "de.christophb.wetter/data/database"
import "de.christophb.wetter/jwt"

func isValidMeasuringNodes(node models.MeasuringNode)bool{
	return node.Lat >= -90 && node.Lat <= 90 && node.Lng >= 180 && node.Lng <= -180 && len(node.Name)>0
}

func AddNodeHandler(w http.ResponseWriter, r *http.Request)  {

	defer recoverHandlerErrors(w)

	nodeRepo := database.GetMeasuringNodeRepository()

	userId , err := jwt.GetUserIdByRequest(r)

	panicIfErrorNonNil(err, NotAuthorized, http.StatusForbidden)

	var node models.MeasuringNode
	err = readBody(r,&node)
	panicIfErrorNonNil(err, InvalidBody, http.StatusBadRequest)

	if !isValidMeasuringNodes(node){
		handleError(w,handlerError{ErrorMessage:InvalidBody},http.StatusBadRequest)
	}

	node, err = nodeRepo.CreateMeasuringNode(node, userId)
	panicIfErrorNonNil(err, "failed to save node", http.StatusInternalServerError)

	writeJsonResponse(node,w)
}