package handlers

import (
	"github.com/ChristophBe/weather-data-server/data/database"
	"github.com/ChristophBe/weather-data-server/services"
	"net/http"
)

type UserHandlers interface {
	GetCreateUserHandler() http.Handler
	GetUserMeHandler() http.Handler
	GetUserEnableHandler() http.Handler
}

func GetUserHandlers() UserHandlers {
	return userHandlersImpl{
		mailService:          services.GetMailService(),
		tokenService:         services.GetAuthTokenService(),
		invitationService:    services.GetInvitationService(),
		userRepository:       database.GetUserRepository(),
		invitationRepository: database.GetInvitationRepository(),
	}
}
