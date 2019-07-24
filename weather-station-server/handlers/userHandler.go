package handlers

import (
	"../data"
	"../jwt"
	"../utils"
	"encoding/base64"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

type UserDTO struct {
	Email string 	`json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}


func isValidUserDTO(userDTO UserDTO) bool{

	//TODO: fix regex to accept all valid email Addresses.
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


	//TODO add Input Validation
	/*if isValidUserDTO(userDTO){
		handleError(w,invalidBodyError,http.StatusBadRequest)
		return
	}*/

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


	enableHash, enableToken ,err  := generateEnableToken(userDTO.Email)

	if err != nil{
		handleError(w,invalidBodyError,http.StatusBadRequest)
		return
	}


	log.Print("enableToken: " + enableToken)

	//Create user object
	passwordHash,err := bcrypt.GenerateFromPassword([]byte(userDTO.Password),bcrypt.DefaultCost)


	newUser := data.User{CreationTime:time.Now(), Email:userDTO.Email,Username:userDTO.Username,PasswordHash:passwordHash, IsEnabled:false,EnableSecretHash: enableHash}

	//Save user to DB
	data.UpsertUser(newUser)

	//Read created user From DB
	user, err = data.FetchUserByEmail(newUser.Email)

	writeJsonResponse(user,w)
}

func generateEnableToken(identifier string) ([]byte, string, error) {

	enabledSecret := utils.RandStringRunes(32)
	enabledSecretHash, err :=  bcrypt.GenerateFromPassword([]byte(enabledSecret),bcrypt.DefaultCost)

	if err != nil {
		return  []byte(""),"",err
	}
	enableToken := base64.StdEncoding.EncodeToString([]byte( identifier + ":" + enabledSecret))

	return enabledSecretHash, enableToken, nil
}


func UsersMe(w http.ResponseWriter, r * http.Request){
	userId, err :=  jwt.GetUserIdBy(r)
	if err != nil{
		handleError(w,handlerError{Err: err,ErrorMessage:"not authenticated"}, http.StatusForbidden)
		return
	}

	user, err := data.FetchUserById(userId)

	if err != nil{
		handleError(w,handlerError{Err: err,ErrorMessage:"not authenticated"}, http.StatusForbidden)
		return
	}

	writeJsonResponse(user,w)
}
