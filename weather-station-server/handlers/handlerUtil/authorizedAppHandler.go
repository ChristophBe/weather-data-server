package handlerUtil

import (
	"errors"
	"net/http"
	"strings"
)

type AuthorizedAppHandler struct {
	Handler  func(sub int64, r * http.Request)(resp interface{},statusCode int, err error)
	VerifyToken func(token string)(sub int64, err error)
}


func (a AuthorizedAppHandler) serveHTTP(w http.ResponseWriter,r *http.Request){
	defer catchErrors(w,r)
	tokenString, err := a.getTokenFormResponse(r)
	if err!=nil {
		panic( Forbidden("not authorized",err))
	}

	sub,err := a.VerifyToken(tokenString)

	if err!=nil {
		panic( Forbidden("not authorized",err))
	}

	AppHandler(func(r *http.Request) (interface{}, int, error) {
		return a.Handler(sub,r)
	}).ServeHTTP(w,r)
}


func (a AuthorizedAppHandler) getTokenFormResponse(r *http.Request) (string, error){
	auth := r.Header.Get("Authorization")

	authHeaderParts := strings.Split(auth," ")
	if len(authHeaderParts)<2 {
		return "", errors.New("token not found")
	}
	token := authHeaderParts[1]
	return token, nil
}

