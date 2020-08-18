package services

import (
	"github.com/ChristophBe/weather-data-server/data/database"
	"github.com/ChristophBe/weather-data-server/data/models"
)

type InvitationService interface {
	HandleInvitation(user models.User, invitationId int64) error
}

func GetInvitationService() InvitationService {
	return invitationServiceImpl{
		measuringNodeRepository: database.GetMeasuringNodeRepository(),
	}
}
