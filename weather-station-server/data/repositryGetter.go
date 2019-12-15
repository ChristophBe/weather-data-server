package data

func GetMeasurementRepository() MeasuringRepository {
	return measuringRepositoryImpl{}
}
func GetNodeAuthTokenRepository() NodeAuthTokenRepository {
	return nodeAuthTokenRepositoryIml{}
}

func GetMeasuringNodeRepository() MeasuringNodeRepository{
	return measuringNodeRepositoryImpl{}
}
func GetUserRepository() UserRepository{
	return userRepositoryImpl{}
}

