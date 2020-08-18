package services

import (
	"github.com/ChristophBe/weather-data-server/data/models"
)

type TokenExpiredError struct{}

func (TokenExpiredError) Error() string {
	return "token is expired"
}

type AuthTokenService interface {
	GenerateUserAccessToken(user models.User) (string, error)
	GenerateUserRefreshToken(user models.User) (string, error)
	GenerateNodeAccessToken(node models.MeasuringNode) (string, error)

	GenerateUserInvitationToken(invitation models.Invitation) (string, error)
	GenerateUserEnableToken(user models.User) (string, error)
	VerifyUserAccessToken(token string) (int64, error)
	VerifyUserRefreshToken(token string) (int64, error)
	VerifyNodeAccessToken(token string) (int64, error)
	VerifyUserInvitationToken(token string) (int64, error)
	VerifyUserEnableToken(token string) (int64, error)
}

func GetAuthTokenService() AuthTokenService {
	return authTokenServiceImpl{}
}
