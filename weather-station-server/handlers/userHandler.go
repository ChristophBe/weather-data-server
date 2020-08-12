package handlers

import (
	"de.christophb.wetter/data/database"
	"de.christophb.wetter/services"
	"net/http"
)

type UserHandlers interface {
	GetCreateUserHandler() http.Handler
	GetUserMeHandler() http.Handler
	GetUserEnableHandler() http.Handler
}


func GetUserHandlers() UserHandlers {
	return userHandlersImpl{
		tokenService:   services.GetAuthTokenService(),
		invitationService:   services.GetInvitationService(),
		userRepository: database.GetUserRepository(),
		invitationRepository: database.GetInvitationRepository(),
	}
}
