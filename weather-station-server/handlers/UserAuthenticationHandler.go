package handlers

import (
	"de.christophb.wetter/data/database"
	"de.christophb.wetter/data/models"
	"de.christophb.wetter/handlers/handlerUtil"
	"de.christophb.wetter/services"
	"encoding/json"
	"errors"
	"io/ioutil"
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

func UserAuthenticationHandler(r *http.Request) (resp interface{}, statusCode int, err error) {

	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()
	if err != nil {
		err = handlerUtil.BadRequest("Invalid Request Body", err)
		return
	}

	var authCredentials services.AuthCredentials

	if err = json.Unmarshal(body, &authCredentials); err != nil {
		err = handlerUtil.BadRequest("Invalid Request Body", err)
		return
	}

	user,err := services.GetUserAuthenticationService().GrandUserAccess(authCredentials)
	if err != nil {
		var UnknownGrantTypeError *services.UnknownGrantTypeError
		if errors.As(err,&UnknownGrantTypeError) {
			err = handlerUtil.Forbidden("unknown grant_type",err)
			return
		}
		var TokenExpired  *services.TokenExpiredError
		if errors.As(err,&TokenExpired) {
			err = handlerUtil.Forbidden("token expired",err)
			return
		}

		err = handlerUtil.Forbidden("invalid credentials",err)
		return
	}

	if !user.IsEnabled {
		err = handlerUtil.Forbidden("User is not enabled", err)
		return
	}

	go updateLastLogin(user)

	if resp, err = generateAuthTokenResponse(user); err != nil {
		err = handlerUtil.InternalError(err)
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