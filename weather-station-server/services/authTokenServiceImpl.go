package services

import (
	"de.christophb.wetter/config"
	"de.christophb.wetter/data/models"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type tokenType int

const (
	USER_AUTH tokenType = iota
	USER_REFRESH
	NODE_AUTH
)


func (t tokenType) toString() string {

	tokenTypeToString := map[tokenType]string{
		USER_AUTH:    "USER_AUTH",
		USER_REFRESH: "USER_REFRESH",
		NODE_AUTH:    "NODE_AUTH",
	}
	return tokenTypeToString[t]
}
func tokenTypeByString(typeString string) (tokenType, error) {
	stringToTokenType := map[string]tokenType{
		"USER_AUTH":    USER_AUTH,
		"USER_REFRESH": USER_REFRESH,
		"NODE_AUTH":    NODE_AUTH,
	}
	return stringToTokenType[typeString], nil
}

type authTokenContext struct {
	Type     tokenType
	Expiring time.Duration
	Sub      int64
}

type authTokenServiceImpl struct {
}

func (a authTokenServiceImpl) GenerateUserAccessToken(user models.User) (signedToken string, err error) {
	tokenContext := authTokenContext{
		Type:     USER_AUTH,
		Expiring: time.Hour * time.Duration(24),
		Sub:      user.Id,
	}
	signedToken, err = a.generateToken(tokenContext)
	return
}

func (a authTokenServiceImpl) GenerateUserRefreshToken(user models.User) (signedToken string, err error) {
	tokenContext := authTokenContext{
		Type:     USER_REFRESH,
		Expiring: time.Hour *  time.Duration(24* 30),//30 days
		Sub:      user.Id,
	}
	signedToken, err = a.generateToken(tokenContext)
	return
}

func (a authTokenServiceImpl) GenerateNodeAccessToken(node models.MeasuringNode) (signedToken string, err error) {
	tokenContext := authTokenContext{
		Type:     NODE_AUTH,
		Expiring: time.Hour * time.Duration(24),
		Sub:      node.Id,
	}
	signedToken, err = a.generateToken(tokenContext)
	return
}

func (a authTokenServiceImpl) VerifyUserAccessToken(token string) (int64, error) {
	return a.verifyToken(token, USER_AUTH)
}

func (a authTokenServiceImpl) VerifyUserRefreshToken(token string) (int64, error) {
	return a.verifyToken(token, USER_REFRESH)
}

func (a authTokenServiceImpl) VerifyNodeAccessToken(token string) (int64, error) {
	return a.verifyToken(token, NODE_AUTH)
}

func (a authTokenServiceImpl) generateToken(context authTokenContext) (signedToken string, err error) {
	expirationTime := time.Now().Add(context.Expiring)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  context.Sub,
		"type": context.Type.toString(),
		"exp":  expirationTime.Unix(),
	})

	conf, err := config.GetConfigManager().GetConfig()
	if err != nil {
		return
	}
	signedToken, err = token.SignedString([]byte(conf.Auth.AuthKey))
	return
}

func (a authTokenServiceImpl) verifyToken(tokenString string, expectedType tokenType) (sub int64, err error) {
	conf, err := config.GetConfigManager().GetConfig()

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
		var tokenType tokenType

		//Check if Token Type is Valid
		tokenType, err = tokenTypeByString(claims["type"].(string))
		if err != nil {
			return
		}
		if tokenType != expectedType {
			err = errors.New("unexpected token typ")
			return
		}
		exp := int64(claims["exp"].(float64))
		if time.Unix(exp,0).Before(time.Now()) {
			err=TokenExpiredError{}
		}


		sub = int64(claims["sub"].(float64))
	}
	return
}
