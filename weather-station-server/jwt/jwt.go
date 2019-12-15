package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"de.christophb.wetter/data"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"
)

type TokenHeader struct {
	Algorithm string `json:"alg"`
	Typ string `json:"typ"`
}

type TokenPayload struct {
	Expiration time.Time `json:"exp"`
	Subject    int64  `json:"sub"`
	Name       string `json:"name"`
}

type Token struct {
	Header TokenHeader
	TokenPayload TokenPayload
	secret string
}

func GenerateToken(user data.User)( signedToken string,err error ){
	header :=TokenHeader{Algorithm:"HS256", Typ:"JWT"}

	encodedHeader, err := encodeTokenPayload(header)
	if err != nil {
		return
	}

	expirationTime := time.Now().Add(time.Hour * time.Duration(24))

	tokenPayload := TokenPayload{Expiration:expirationTime,Subject:user.Id, Name: user.Username}

	encodedData ,err := encodeTokenPayload(tokenPayload)


	secret := base64.StdEncoding.EncodeToString(signToken( encodedHeader +"." + encodedData))

	token :=  encodedHeader +"." + encodedData + "." + secret

	// store token
	return token, err
}

func encodeTokenPayload(payload interface{})( encodedPayload string,err error ){

	headerJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encodedPayload = base64.StdEncoding.EncodeToString(headerJson)
	return encodedPayload, nil
}


var key = []byte("QMMFGQCtVpWUHYEQcHQnkD9u6fZMMh44uZUaMRtMMbpjdyjRxRA3Pw57sUpBAhRv")

func signToken(token string) []byte{
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(token))
	return mac.Sum(nil)
}

func validMAC(expacted string, messageMAC []byte) bool {
	expectedMAC := signToken(expacted)
	return hmac.Equal(messageMAC, expectedMAC)
}

func decodeTokenPart(encoded string, destination interface{}) error{

	headerBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}
	err  = json.Unmarshal(headerBytes, destination)
	if err != nil {
		return err
	}
	return nil
}


func Verify(token string)(TokenPayload,error) {
	tokenParts := strings.Split(token,".")


	if len(tokenParts) < 3 {
		return TokenPayload{}, errors.New("token not found")
	}
	encodedHeader := tokenParts[0]
	encodedPayload := tokenParts[1]
	encodedSecret := tokenParts[2]

	// Decode token-header
	var header TokenHeader
	err := decodeTokenPart(encodedHeader,&header)
	if err != nil {
		return TokenPayload{}, err
	}

	if header.Typ != "JWT" || header.Algorithm != "HS256" {
		return TokenPayload{}, errors.New("wrong token-type or signing-algorithm")
	}


	//decode token-secret

	secret, err := base64.StdEncoding.DecodeString(encodedSecret)
	if err != nil {
		return TokenPayload{}, err
	}

	isValid := validMAC(encodedHeader + "." + encodedPayload, secret)

	if !isValid {
		return TokenPayload{}, errors.New("invalid token")
	}





	// Decode token-payload
	var payload TokenPayload
	err = decodeTokenPart(encodedPayload,&payload)
	if err != nil {
		return TokenPayload{}, err
	}

	now := time.Now()
	log.Print(now)
	log.Print(payload.Expiration)


	if now.After(payload.Expiration) {
		return TokenPayload{}, errors.New("token expired")
	}

	return payload, nil

}