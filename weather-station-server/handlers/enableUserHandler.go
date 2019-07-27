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
	Username string `json:"username"`
	Password string `json:"password"`
}


func EnableUserHandler(w http.ResponseWriter, request *http.Request)  {

	defer recoverHandlerErrors(w)

	var enableUserDTO EnableTokenDTO

	err := readBody(request, &enableUserDTO)
	panicIfErrorNonNil(err, "invalid body or invalid token",http.StatusBadRequest)

	email, secret := parseToken(enableUserDTO)

	user, err := data.FetchUserByEmail(email)
	panicIfErrorNonNil(err, "invalid body or invalid token",http.StatusBadRequest)

	if user.IsEnabled {
		err = writeJsonResponse(user, w)
		panicIfErrorNonNil(err, "invalid body or invalid token",http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.EnableSecretHash, []byte(secret))
	panicIfErrorNonNil(err, "invalid body or invalid token",http.StatusBadRequest)


	userNameTest, err := data.FetchUserByUsername(enableUserDTO.Username)
	if err != nil || userNameTest.Id!= 0 {
		panic(handlerError{Err:err, ErrorMessage:"invalid body or invalid token",Status:http.StatusBadRequest})
	}


	user.EnableSecretHash = []byte("")
	user.IsEnabled = true
	user.Username = enableUserDTO.Username

	
	passwordHash,err := bcrypt.GenerateFromPassword([]byte(enableUserDTO.Password),bcrypt.DefaultCost)
	panicIfErrorNonNil(err, "unexpected error", http.StatusInternalServerError)
	user.PasswordHash = passwordHash
	
	
	user,err = data.UpsertUser(user)
	panicIfErrorNonNil(err, "unexpected error", http.StatusInternalServerError)

	err = writeJsonResponse(user, w)
	panicIfErrorNonNil(err, "unexpected error", http.StatusInternalServerError)


}



func enableUser(user data.User)(data.User,error) {
	user.IsEnabled = true
	user.EnableSecretHash = []byte("")
	return data.UpsertUser(user)
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
