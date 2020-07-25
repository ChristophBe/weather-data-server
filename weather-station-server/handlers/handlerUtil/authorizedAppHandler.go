package handlerUtil

import (
	"errors"
	"net/http"
	"strings"
)

func AuthorizedAppHandler(
	tokenVerifier func(token string)(sub int64, err error),
	handler  func(sub int64, r * http.Request)(resp interface{},statusCode int, err error),
)AppHandler{
	return func(r *http.Request) (interface{}, int, error) {
		tokenString, err := readTokenFormRequest(r)

		if err!=nil {
			panic( Forbidden("not authorized",err))
		}

		sub,err := tokenVerifier(tokenString)

		if err!=nil {
			panic( Forbidden("not authorized",err))
		}

		return handler(sub,r)
	}
}



func readTokenFormRequest(r *http.Request) (string, error){
	auth := r.Header.Get("Authorization")

	authHeaderParts := strings.Split(auth," ")
	if len(authHeaderParts)<2 {
		return "", errors.New("token not found")
	}
	token := authHeaderParts[1]
	return token, nil
}

