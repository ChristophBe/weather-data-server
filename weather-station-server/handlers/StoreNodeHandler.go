package handlers

import "net/http"
import "de.christophb.wetter/data"
import "de.christophb.wetter/jwt"

func AddNodeHandler(w http.ResponseWriter, r *http.Request)  {

	defer recoverHandlerErrors(w)

	nodeRepo := data.GetMeasuringNodeRepository()

	userId , err := jwt.GetUserIdBy(r)

	panicIfErrorNonNil(err, NotAuthorized, http.StatusForbidden)

	var node data.MeasuringNode
	err = readBody(r,&node)
	panicIfErrorNonNil(err, InvalidBody, http.StatusBadRequest)
	//TODO Validate input Node

	node, err = nodeRepo.CreateMeasuringNode(node, userId)
	panicIfErrorNonNil(err, "failed to save node", http.StatusInternalServerError)

	writeJsonResponse(node,w)
}