package handlers

import (
	"../data"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type EnableTokenDTO struct {
	Token string `json:"token"`
}


type EnableTokenResp struct {
	Msg string `json:"message"`
}
func EnableUserHandler(w http.ResponseWriter, request *http.Request)  {

	defer recoverHandlerErrors(w)

	var enableToken EnableTokenDTO

	err := readBody(request, &enableToken)
	if err != nil {
		panic(handlerError{Err: err, ErrorMessage: "invalid body or invalid token"})
	}

	email, secret := parseToken(enableToken)

	user, err := data.FetchUserByEmail(email)

	if err != nil {
		panic(handlerError{Err: err, ErrorMessage: "invalid body or invalid token"})
	}

	if user.IsEnabled {
		err = writeJsonResponse(EnableTokenResp{Msg: "Your Account was enabled successfully"}, w)
		if err != nil {
			panic(handlerError{Err: err, ErrorMessage: "invalid body or invalid token"})
		}
		return
	}

	err = bcrypt.CompareHashAndPassword(user.EnableSecretHash, []byte(secret))

	if err != nil {
		panic(handlerError{Err: err, ErrorMessage: "invalid body or invalid token"})
	}

	enableUser(user)

	err = writeJsonResponse(EnableTokenResp{Msg: "Your Account was enabled successfully"}, w)
	if err != nil {
		panic(handlerError{Err: err, ErrorMessage: "invalid body or invalid token"})
	}

}



func enableUser(user data.User) {
	user.IsEnabled = true
	user.EnableSecretHash = []byte("")
	data.UpsertUser(user)
}

func parseToken(enableToken EnableTokenDTO) (string, string) {
	decodedToken, err := base64.StdEncoding.DecodeString(enableToken.Token)
	if err != nil {
		panic(handlerError{Err: err, ErrorMessage: "invalid body or invalid token"})
	}
	tokenParts := strings.Split(string(decodedToken), ":")
	email, secret := tokenParts[0], tokenParts[1]
	return email, secret
}
