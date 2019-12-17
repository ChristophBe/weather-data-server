package database

import (
	"de.christophb.wetter/data/repositories"
)

func GetMeasurementRepository() repositories.MeasuringRepository {
	return measuringRepositoryImpl{}
}
func GetNodeAuthTokenRepository() repositories.NodeAuthTokenRepository {
	return nodeAuthTokenRepositoryIml{}
}

func GetMeasuringNodeRepository() repositories.MeasuringNodeRepository {
	return measuringNodeRepositoryImpl{}
}
func GetUserRepository() repositories.UserRepository {
	return userRepositoryImpl{}
}

