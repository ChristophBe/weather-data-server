package jwt

import (
	"errors"
	"net/http"
	"strings"
)



func GetTokenFormResponse(r *http.Request) (string, error){


	auth := r.Header.Get("Authorization")


	authHeaderParts := strings.Split(auth," ")
	if len(authHeaderParts)<2 {
		return "", errors.New("Token not Found")
	}
	token := authHeaderParts[1]
	return token, nil


}

func GetUserIdBy(request *http.Request)  (int64, error){

	token, err := GetTokenFormResponse(request)
	if err != nil {
		return  0, err
	}

	tp , err := Verify(token)
	if err != nil {
		return  0, err
	}

	return tp.Subject,nil
}
