package services

import (
	"github.com/ChristophBe/weather-data-server/data/models"
)

type TokenType int

const (
	USER_AUTH TokenType = iota
	USER_REFRESH
	USER_INVITATION
	USER_ENABLE
	NODE_AUTH
)

func (t TokenType) toString() string {

	tokenTypeToString := map[TokenType]string{
		USER_AUTH:       "USER_AUTH",
		USER_REFRESH:    "USER_REFRESH",
		USER_INVITATION: "USER_INVITATION",
		USER_ENABLE:     "USER_ENABLE",
		NODE_AUTH:       "NODE_AUTH",
	}
	return tokenTypeToString[t]
}
func tokenTypeByString(typeString string) (TokenType, error) {
	stringToTokenType := map[string]TokenType{
		"USER_AUTH":       USER_AUTH,
		"USER_REFRESH":    USER_REFRESH,
		"USER_INVITATION": USER_INVITATION,
		"USER_ENABLE":     USER_ENABLE,
		"NODE_AUTH":       NODE_AUTH,
	}
	return stringToTokenType[typeString], nil
}


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
	GetTokenVerifier(tokenTyp TokenType) func(token string) (int64, error)
}

func GetAuthTokenService() AuthTokenService {
	return authTokenServiceImpl{}
}
