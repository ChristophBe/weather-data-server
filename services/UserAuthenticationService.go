package services

import (
	"fmt"
	"github.com/ChristophBe/weather-data-server/data/database"
	"github.com/ChristophBe/weather-data-server/data/models"
)

type InvalidCredentialsError struct{}

func (ue InvalidCredentialsError) Error() string {
	return "invalid credentials"
}

type UnknownGrantTypeError struct {
	grantType string
}

func (ue UnknownGrantTypeError) Error() string {
	return fmt.Sprintf("unknowen grant_type of type \"%s\"", ue.grantType)
}

type AuthCredentials struct {
	GrantType    string `json:"grant_type"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

type UserAuthenticationService interface {
	GrandUserAccess(AuthCredentials) (models.User, error)
}

func GetUserAuthenticationService() UserAuthenticationService {
	return userAuthenticationServiceImpl{
		authTokenService: GetAuthTokenService(),
		userRepository:   database.GetUserRepository(),
	}

}
