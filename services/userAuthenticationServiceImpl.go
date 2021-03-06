package services

import (
	"errors"
	"github.com/ChristophBe/weather-data-server/data/models"
	"github.com/ChristophBe/weather-data-server/data/repositories"
	"github.com/ChristophBe/weather-data-server/handlers/httpHandler"
	"golang.org/x/crypto/bcrypt"
)

type userAuthenticationServiceImpl struct {
	authTokenService AuthTokenService
	userRepository   repositories.UserRepository
}

func (ua userAuthenticationServiceImpl) GrandUserAccess(authCredentials AuthCredentials) (user models.User, err error) {

	switch authCredentials.GrantType {
	case "password":
		user, err = ua.passwordGrant(authCredentials)
		return
	case "refresh_token":
		user, err = ua.refreshTokenGrant(authCredentials)
		return
	default:
		err = UnknownGrantTypeError{authCredentials.GrantType}
		return
	}
}

func (ua userAuthenticationServiceImpl) passwordGrant(credentials AuthCredentials) (user models.User, err error) {
	if len(credentials.Password) < 4 || len(credentials.Email) < 4 {
		err = httpHandler.Forbidden("Invalid Credentials", errors.New("password or email is to short"))
		return
	}

	user, e := ua.userRepository.FetchUserByEmail(credentials.Email)
	if e != nil || user.Id == 0 {
		err = httpHandler.Forbidden("Invalid Credentials", e)
	}

	e = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(credentials.Password))
	if e != nil {
		err = httpHandler.Forbidden("Invalid Credentials", e)
	}
	return
}

func (ua userAuthenticationServiceImpl) refreshTokenGrant(credentials AuthCredentials) (user models.User, err error) {
	if len(credentials.RefreshToken) < 0 {
		err = httpHandler.Forbidden("Invalid Credentials", errors.New("password or email is to short"))
		return
	}

	userTokenVerifier := ua.authTokenService.GetTokenVerifier(USER_REFRESH)
	userId, err := userTokenVerifier(credentials.RefreshToken)
	if err != nil {
		err = httpHandler.Forbidden("Invalid Credentials", err)
	}
	user, e := ua.userRepository.FetchUserById(userId)
	if e != nil || user.Id == 0 {
		err = httpHandler.Forbidden("Invalid Credentials", e)
	}
	return
}
