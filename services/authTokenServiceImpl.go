package services

import (
	"errors"
	"fmt"
	"github.com/ChristophBe/weather-data-server/config"
	"github.com/ChristophBe/weather-data-server/data/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)


type authTokenContext struct {
	Type     TokenType
	Expiring time.Duration
	Sub      int64
}

type authTokenServiceImpl struct {}

func (a authTokenServiceImpl) GetTokenVerifier(tokenTyp TokenType) func(token string) (int64, error) {
	return func(token string) (int64, error) {
		return a.verifyToken(token, tokenTyp)
	}
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
		Expiring: time.Hour * time.Duration(24*30), //30 days
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

func (a authTokenServiceImpl) GenerateUserInvitationToken(node models.Invitation) (signedToken string, err error) {
	tokenContext := authTokenContext{
		Type:     USER_INVITATION,
		Expiring: time.Hour * time.Duration(24*30),
		Sub:      node.Id,
	}
	signedToken, err = a.generateToken(tokenContext)
	return
}

func (a authTokenServiceImpl) GenerateUserEnableToken(user models.User) (signedToken string, err error) {
	tokenContext := authTokenContext{
		Type:     USER_ENABLE,
		Expiring: time.Hour * time.Duration(24*3),
		Sub:      user.Id,
	}
	signedToken, err = a.generateToken(tokenContext)
	return
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

func (a authTokenServiceImpl) verifyToken(tokenString string, expectedType TokenType) (sub int64, err error) {
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
		var tokenType TokenType

		//Check if Token Type is Valid
		tokenType, err = tokenTypeByString(claims["type"].(string))
		if err != nil {
			return
		}
		if tokenType != expectedType {
			err = errors.New("unexpected token typ")
			return
		}
		if tokenType != NODE_AUTH {
			exp := int64(claims["exp"].(float64))
			if time.Unix(exp, 0).Before(time.Now()) {
				err = TokenExpiredError{}
			}
		}


		sub = int64(claims["sub"].(float64))
	}
	return
}
