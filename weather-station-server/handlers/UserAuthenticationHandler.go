package handlers

import (
	"de.christophb.wetter/data/database"
	"de.christophb.wetter/data/models"
	"de.christophb.wetter/handlers/httpHandler"
	"de.christophb.wetter/services"
	"errors"
	"log"
	"net/http"
	"time"
)

type authCredentialBody struct {
	GrantType    string `json:"grant_type"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

type authTokenResponse struct {
	Type    string `json:"token_type"`
	Token   string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

func UserAuthenticationHandler(r *http.Request) (response httpHandler.HandlerResponse, err error) {

	authService := services.GetUserAuthenticationService()
	var authCredentials services.AuthCredentials
	err = httpHandler.ReadJsonBody(r,&authCredentials)
	if err != nil {
		return
	}

	user,err := authService.GrandUserAccess(authCredentials)
	if err != nil {
		err =  annotateGrandUserAccessError(err)
		return
	}

	if !user.IsEnabled {
		err = httpHandler.Forbidden("User is not enabled", nil)
		return
	}

	go updateLastLogin(user)


	if response.Data, err = generateAuthTokenResponse(user); err != nil {
		err = httpHandler.InternalError(err)
	}
	response.Status = http.StatusOK
	return
}

func annotateGrandUserAccessError(err error) error {
	var UnknownGrantTypeError *services.UnknownGrantTypeError
	if errors.As(err, &UnknownGrantTypeError) {
		return httpHandler.Forbidden("unknown grant_type", err)

	}
	var TokenExpired *services.TokenExpiredError
	if errors.As(err, &TokenExpired) {

	}
	return httpHandler.Forbidden("invalid credentials", err)

}
func generateAuthTokenResponse(user models.User) (token authTokenResponse, err error) {

	tokenService := services.GetAuthTokenService()
	token.Type = "Bearer"
	token.Token, err = tokenService.GenerateUserAccessToken(user)
	if err != nil {
		return
	}
	token.Refresh, err = tokenService.GenerateUserRefreshToken(user)
	return
}


func updateLastLogin(user models.User) {

	user.LastLogin = time.Now()

	_, err := database.GetUserRepository().SaveUser(user)

	if err != nil {
		log.Print(err)
	}
}