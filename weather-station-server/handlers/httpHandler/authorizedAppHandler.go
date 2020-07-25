package httpHandler

import (
	"errors"
	"net/http"
	"strings"
)

func AuthorizedAppHandler(
	tokenVerifier func(token string)(sub int64, err error),
	handler  func(sub int64, r * http.Request)(resp interface{},statusCode int, err error),
) JsonHandler {
	return func(r *http.Request) (interface{}, int, error) {
		_, tokenString, err := readTokenFormRequest(r)

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

func readTokenFormRequest(r *http.Request) (tokenType string, tokenString string, err error){
	auth := r.Header.Get("Authorization")

	authHeaderParts := strings.Split(auth," ")
	if len(authHeaderParts)<2 {
		err = errors.New("token not found")
		return
	}
	tokenType = authHeaderParts[0]
	tokenString = authHeaderParts[1]
	return
}

