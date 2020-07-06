package services

import (
	"de.christophb.wetter/data/models"
)

type TokenExpiredError struct {}
func (TokenExpiredError) Error() string{
	return "token is expired"
}


type AuthTokenService interface {
	GenerateUserAccessToken(user models.User)(string,error)
	GenerateUserRefreshToken(user models.User)(string,error)
	GenerateNodeAccessToken(node models.MeasuringNode)(string,error)
	VerifyUserAccessToken(token string)(int64,error)
	VerifyUserRefreshToken(token string)(int64,error)
	VerifyNodeAccessToken(token string)(int64,error)
}

func GetAuthTokenService() AuthTokenService {
	return authTokenServiceImpl{}
}
