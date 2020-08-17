package handlers

import (
	"de.christophb.wetter/data/database"
	"de.christophb.wetter/data/models"
	"de.christophb.wetter/data/repositories"
	"de.christophb.wetter/handlers/httpHandler"
	"de.christophb.wetter/services"
	"fmt"
	"net/http"
)

type nodeHandlersImpl struct {
	authTokenService     services.AuthTokenService
	shareNodeService     services.ShareNodeService
	nodeRepository       repositories.MeasuringNodeRepository
	userRepository       repositories.UserRepository
	invitationRepository repositories.InvitationRepository
}

func GetNodeHandlers() NodeHandlers {
	return nodeHandlersImpl{
		authTokenService:     services.GetAuthTokenService(),
		shareNodeService:     services.GetShareNodeService(),
		nodeRepository:       database.GetMeasuringNodeRepository(),
		userRepository:       database.GetUserRepository(),
		invitationRepository: database.GetInvitationRepository(),
	}
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

func (n nodeHandlersImpl) GetShareNodeHandler() http.Handler {
	return httpHandler.AuthorizedAppHandler(n.authTokenService.VerifyUserAccessToken, n.shareNodeHandler)
}

func (n nodeHandlersImpl) isValidMeasuringNodes(node models.MeasuringNode) bool {
	return node.Lat >= -90 && node.Lat <= 90 && node.Lng >= 180 && node.Lng <= -180 && len(node.Name) > 0
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

	userReqId, err := httpHandler.ReadPathVariableInt(r, "userId")
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
	if err != nil {
		err = httpHandler.InternalError(err)
		return
	}

	//TODO: Implement Node Update

	response.Data = node
	response.Status = http.StatusAccepted
	return
}

func (n nodeHandlersImpl) shareNodeHandler(userId int64, request *http.Request) (response httpHandler.HandlerResponse, err error) {

	nodeId, err := httpHandler.ReadPathVariableInt(request, "nodeId")
	if err != nil {
		message := fmt.Sprintf(httpHandler.ErrorMessageParameterf, "nodeId")
		err = httpHandler.BadRequest(message, err)
		return
	}

	node, err := n.nodeRepository.FetchMeasuringNodeById(nodeId)
	if err != nil {
		message := fmt.Sprintf(httpHandler.ErrorMessageNotFoundf, "node")
		err = httpHandler.NotFound(message, err)
		return
	}

	owner, err := n.userRepository.FetchOwnerByMeasuringNode(node.Id)
	if err != nil || userId != owner.Id {
		err = httpHandler.Forbidden(httpHandler.ErrorMessageNotAuthorized, err)
		return
	}

	var requestBody struct {
		Email string `json:"email"`
	}

	err = readBody(request, &requestBody)
	if err != nil {
		err = httpHandler.InternalError(err)
		return
	}

	err = n.shareNodeService.ShareNode(node, owner, requestBody.Email)
	if err != nil {
		err = httpHandler.InternalError(err)
		return
	}

	response.Data = struct {
		Msg string `json:"message"`
	}{Msg: "the node was successfully shared"}
	response.Status = http.StatusOK
	return
}
