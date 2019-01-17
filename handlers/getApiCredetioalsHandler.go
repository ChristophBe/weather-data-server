package handlers

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
)

func RandStringRunes(n int) string {

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type credentialsExternal struct {
	Secret string
	ClientId string
}
func GenerateApiCredentialsHandler(w http.ResponseWriter, request *http.Request) {

	secret:= RandStringRunes(32)

	hash,err := bcrypt.GenerateFromPassword([]byte(secret),bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	clientId := RandStringRunes(16)
	//inter := credentialsInternal{Hash: hash,ClientId:clientId}
	ext := credentialsExternal{Secret: secret,ClientId:clientId}

	//TODO: Save the credentials for for the give Station to the db


	extJSON,err:= json.Marshal(ext)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(extJSON)
	return

}