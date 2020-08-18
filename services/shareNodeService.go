package services

import (
	"github.com/ChristophBe/weather-data-server/data/database"
	"github.com/ChristophBe/weather-data-server/data/models"
)

type ShareNodeService interface {
	ShareNode(node models.MeasuringNode, nodeOwner models.User, email string) error
}

func GetShareNodeService() ShareNodeService {
	return shareNodeServiceImpl{
		authTokenService:     GetAuthTokenService(),
		userRepository:       database.GetUserRepository(),
		nodeRepository:       database.GetMeasuringNodeRepository(),
		invitationRepository: database.GetInvitationRepository(),
	}
}
