package database

import (
	"github.com/ChristophBe/weather-data-server/data/repositories"
)

func GetMeasurementRepository() repositories.MeasuringRepository {
	return measuringRepositoryImpl{}
}
func GetMeasuringNodeRepository() repositories.MeasuringNodeRepository {
	return measuringNodeRepositoryImpl{}
}
func GetUserRepository() repositories.UserRepository {
	return userRepositoryImpl{}
}
func GetInvitationRepository() repositories.InvitationRepository {
	return invitationRepositoryImpl{}
}
