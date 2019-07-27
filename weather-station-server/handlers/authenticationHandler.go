package handlers

import (
	"../data"
	"../jwt"
	"encoding/json"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/log"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"time"
)



type authCredentialsDTO struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type authTokenDTO struct {
	Type string `json:"token_type"`
	Token string `json:"access_token"`
}


func AuthenticationHandler(w http.ResponseWriter,r *http.Request)  {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		handleError(w,handlerError{Err:err, ErrorMessage:"Invalid Request Body"},http.StatusBadRequest)
		return
	}



	// Unmarshal
	var authCredentials authCredentialsDTO
	err = json.Unmarshal(b, &authCredentials)
	if err != nil {
		handleError(w,handlerError{Err:err,ErrorMessage:"Invalid Request Body"}, http.StatusBadRequest)
		return
	}

	user, err :=data.FetchUserByEmail(authCredentials.Email)


	if err != nil || user.Id == 0 {
		handleError(w,handlerError{Err:err,ErrorMessage:"Invalid Credentials"},http.StatusForbidden)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash,[]byte(authCredentials.Password))

	if err != nil{
		handleError(w,handlerError{Err:err,ErrorMessage:"Invalid Credentials"},http.StatusForbidden)
		return
	}

	if !user.IsEnabled {
		handleError(w,handlerError{Err:nil, ErrorMessage:"User is not enabled"},http.StatusForbidden)
		return
	}



	tokenString, err := jwt.GenerateToken(user)
	if err != nil || user.Id == 0 {
		handleError(w,handlerError{Err:err,ErrorMessage:"Failed to Create Token"},http.StatusBadRequest)
		return
	}

	go  updateLastLogin(user)

	token := authTokenDTO{Token:tokenString, Type: "Bearer"}
	err = writeJsonResponse(token,w)

	if err != nil || user.Id == 0 {
		handleError(w,handlerError{Err:err,ErrorMessage:""},http.StatusInternalServerError)
		return
	}

}

func updateLastLogin(user data.User){

	user.LastLogin = time.Now()

	_, err := data.UpsertUser(user)

	if err != nil {
		log.Error(err)
	}
}
