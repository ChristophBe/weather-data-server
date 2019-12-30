package jwt

import (
	"de.christophb.wetter/config"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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
	conf,err := config.GetConfigManager().GetConfig()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(conf.Auth.AuthKey), nil
	})
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId = int64(claims["sub"].(float64))
	}
	return
}
