package handlers

import (
	"net/http"
)

type NodeHandlers interface {
	GetFetchNodesHandler() http.Handler
	GetFetchNodesByOwnerHandler() http.Handler
	GetSaveNodeHandler() http.Handler
	GetShareNodeHandler() http.Handler
}
