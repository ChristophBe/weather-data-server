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

func UserAuthenticationHandler(r *http.Request) (resp interface{}, statusCode int) {

	authService := services.GetUserAuthenticationService()
	var authCredentials services.AuthCredentials
	httpHandler.ReadJsonBody(r,&authCredentials)

	user,err := authService.GrandUserAccess(authCredentials)
	if err != nil {
		var UnknownGrantTypeError *services.UnknownGrantTypeError
		if errors.As(err,&UnknownGrantTypeError) {
			httpHandler.HandleForbidden("unknown grant_type",err)
		}
		var TokenExpired  *services.TokenExpiredError
		if errors.As(err,&TokenExpired) {
			httpHandler.HandleForbidden("token expired",err)

		}
		httpHandler.HandleForbidden("invalid credentials",err)
		return
	}

	if !user.IsEnabled {
		httpHandler.HandleForbidden("User is not enabled", err)
	}

	go updateLastLogin(user)

	if resp, err = generateAuthTokenResponse(user); err != nil {
		httpHandler.HandleInternalError(err)
	}
	statusCode = http.StatusOK
	return
}
func generateAuthTokenResponse(user models.User) (token authTokenResponse, err error) {

	tokenService := services.GetAuthTokenService()
	token.Type = "Bearer"
	token.Token, err = tokenService.GenerateUserAccessToken(user)
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