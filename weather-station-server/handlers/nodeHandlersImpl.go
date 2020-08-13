package handlers

import (
	"de.christophb.wetter/data/models"
	"de.christophb.wetter/data/repositories"
	"de.christophb.wetter/handlers/httpHandler"
	"de.christophb.wetter/services"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type nodeHandlersImpl struct {
	authTokenService services.AuthTokenService
	nodeRepository   repositories.MeasuringNodeRepository
}

func (n nodeHandlersImpl) GetFetchNodesHandler() http.Handler {
	return httpHandler.AuthorizedAppHandlerWithUnauthorisedFallback(n.authTokenService.VerifyUserAccessToken, n.fetchNodesHandlerAuthorized, n.fetchNodesHandlerUnauthorized)
}

func (n nodeHandlersImpl) GetFetchNodesByOwnerHandler() http.Handler {
	return httpHandler.AuthorizedAppHandler(n.authTokenService.VerifyUserAccessToken, n.fetchNodesByOwnerHandler)
}

func (n nodeHandlersImpl) GetSaveNodeHandler() http.Handler {
	return httpHandler.AuthorizedAppHandler(n.authTokenService.VerifyUserAccessToken, n.saveNodeHandler)
}

func (n nodeHandlersImpl) fetchNodesHandlerAuthorized(userId int64, r *http.Request) (response httpHandler.HandlerResponse, err error) {
	nodes, err := n.nodeRepository.FetchAllVisibleNodesByUserId(userId)
	if err != nil {
		err = httpHandler.InternalError(err)
		return
	}

	response.Data = nodes
	response.Status = http.StatusOK
	return
}
func (n nodeHandlersImpl) fetchNodesHandlerUnauthorized(r *http.Request) (response httpHandler.HandlerResponse, err error) {
	nodes, err := n.nodeRepository.FetchAllPublicNodes()
	if err != nil {
		err = httpHandler.InternalError(err)
		return
	}

	response.Data = nodes
	return
}

func (n nodeHandlersImpl) fetchNodesByOwnerHandler(userId int64, r *http.Request) (response httpHandler.HandlerResponse, err error) {

	vars := mux.Vars(r)
	userReqId, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		msg := fmt.Sprintf(httpHandler.ErrorMessageParameterf, "userId")
		err = httpHandler.BadRequest(msg, err)
	}

	if userId != userReqId {
		err = httpHandler.Forbidden(httpHandler.ErrorMessageNotAuthorized, nil)
	}

	nodes, err := n.nodeRepository.FetchNodesOwnedByUserId(userId)
	if err != nil {
		err = httpHandler.InternalError(err)
		return
	}

	response.Data = nodes
	return
}

func (n nodeHandlersImpl) isValidMeasuringNodes(node models.MeasuringNode) bool {
	return node.Lat >= -90 && node.Lat <= 90 && node.Lng >= 180 && node.Lng <= -180 && len(node.Name) > 0
}

func (n nodeHandlersImpl) saveNodeHandler(userId int64, r *http.Request) (response httpHandler.HandlerResponse, err error) {

	var node models.MeasuringNode
	err = httpHandler.ReadJsonBody(r, &node)
	if err != nil {
		err = httpHandler.BadRequest(httpHandler.ErrorMessageInvalidBody, err)
		return
	}

	if !n.isValidMeasuringNodes(node) {
		err = httpHandler.BadRequest(httpHandler.ErrorMessageInvalidBody, nil)
		return
	}

	node, err = n.nodeRepository.CreateMeasuringNode(node, userId)
	panicIfErrorNonNil(err, "failed to save node", http.StatusInternalServerError)

	//TODO: Implement Node Update

	response.Data = node
	response.Status = http.StatusAccepted
	return
}
