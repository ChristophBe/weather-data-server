package services

import (
	"de.christophb.wetter/data/database"
	"de.christophb.wetter/data/models"
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
