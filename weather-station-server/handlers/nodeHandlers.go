package handlers

import (
	"de.christophb.wetter/data/database"
	"de.christophb.wetter/services"
	"net/http"
)

type NodeHandlers interface {
	GetFetchNodesHandler() http.Handler
	GetFetchNodesByOwnerHandler() http.Handler
	GetSaveNodeHandler() http.Handler
}

func GetNodeHandlers() NodeHandlers {
	return nodeHandlersImpl{
		authTokenService: services.GetAuthTokenService(),
		nodeRepository:   database.GetMeasuringNodeRepository(),
	}
}
