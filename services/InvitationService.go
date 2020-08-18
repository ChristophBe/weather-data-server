package services

import (
	"de.christophb.wetter/data/database"
	"de.christophb.wetter/data/models"
)

type InvitationService interface {
	HandleInvitation(user models.User, invitationId int64) error
}

func GetInvitationService() InvitationService {
	return invitationServiceImpl{
		measuringNodeRepository: database.GetMeasuringNodeRepository(),
	}
}
