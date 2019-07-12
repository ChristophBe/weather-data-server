package handlers

import (
	"../data"
	"../jwt"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type UserDTO struct {
	Email string 	`json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}


func isValidUserDTO(userDTO UserDTO) bool{
	mailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return len(userDTO.Username) > 4 && len(userDTO.Password) > 4 && mailRegex.MatchString(userDTO.Email)
}

func CreateUserHandler(w http.ResponseWriter, r * http.Request) {


	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	invalidBodyError := handlerError{Err:err, ErrorMessage:"Invalid Request Body"}
	if err != nil {
		handleError(w,invalidBodyError,http.StatusBadRequest)
		return
	}

	// Unmarshal
	var userDTO UserDTO
	err = json.Unmarshal(b, &userDTO)
	if err != nil {
		handleError(w,invalidBodyError, http.StatusBadRequest)
		return
	}

	if isValidUserDTO(userDTO){
		handleError(w,invalidBodyError,http.StatusBadRequest)
	}

	//Check if User with given Email is Existing
	user , err := data.FetchUserByEmail(userDTO.Email)
	if err != nil || user.Id != 0{
		handleError(w,invalidBodyError,http.StatusBadRequest)
		return
	}

	//Check if User with given Email is Existing
	user , err = data.FetchUserByUsername(userDTO.Username)
	if err != nil || user.Id != 0{
		handleError(w,invalidBodyError,http.StatusBadRequest)
		return
	}


	//Create user object
	passwordHash,err := bcrypt.GenerateFromPassword([]byte(userDTO.Password),bcrypt.DefaultCost)
	newUser := data.User{CreationTime:time.Now(), Email:userDTO.Email,Username:userDTO.Username,PasswordHash:passwordHash}

	//Save user to DB
	data.CreateUser(newUser)

	//Read created user From DB
	user, err = data.FetchUserByEmail(newUser.Email)

	writeJsonResponse(user,w)
}

//func auth

type authCredentialsDTO struct {
	 Email string `json:"email"`
	 Password string `json:"password"`
}

type authTokenDTO struct {
	Token string `json:"token"`
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



	tokenString, err := jwt.GenerateToken(user)
	if err != nil || user.Id == 0 {
		handleError(w,handlerError{Err:err,ErrorMessage:"Failed to Create Token"},http.StatusBadRequest)
		return
	}
	
	token := authTokenDTO{Token:tokenString}
	err = writeJsonResponse(token,w)

	if err != nil || user.Id == 0 {
		handleError(w,handlerError{Err:err,ErrorMessage:""},http.StatusInternalServerError)
		return
	}

}







func CheckTokenHandler(w http.ResponseWriter,r *http.Request)  {
	auth := r.Header.Get("Authorization")
	authHeaderParts := strings.Split(auth," ")

	token := authHeaderParts[1]
	payload,err:= jwt.Verify(token)

	if err != nil {
		handleError(w,handlerError{Err:err,ErrorMessage:"invalid Token"}, http.StatusForbidden)
		return
	}

	writeJsonResponse(payload,w)
}