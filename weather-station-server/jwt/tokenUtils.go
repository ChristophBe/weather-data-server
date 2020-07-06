package jwt

import (
	"de.christophb.wetter/services"
	"errors"
	"net/http"
	"strings"
)


func GetTokenFormResponse(r *http.Request) (string, error){
	auth := r.Header.Get("Authorization")

	authHeaderParts := strings.Split(auth," ")
	if len(authHeaderParts)<2 {
		return "", errors.New("token not found")
	}
	token := authHeaderParts[1]
	return token, nil
}

func GetUserIdByRequest(request *http.Request)  (userId int64, err error){
	tokenString, err := GetTokenFormResponse(request)

	userId,err = services.GetAuthTokenService().VerifyNodeAccessToken(tokenString)
	return
}
