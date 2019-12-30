package handlers

import (
	"de.christophb.wetter/data/database"
	"de.christophb.wetter/data/models"
	"de.christophb.wetter/jwt"
	"de.christophb.wetter/utils"
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
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
func IsValidEmail(email string) bool {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return len(email) < 254 && rxEmail.MatchString(email)
}

func isValidUserDTO(userDTO UserDTO) bool {
	return len(userDTO.Username) > 4 && len(userDTO.Password) > 4 && IsValidEmail(userDTO.Email)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	userRepo := database.GetUserRepository()

	invalidBodyError := handlerError{Err: err, ErrorMessage: "Invalid Request Body"}
	if err != nil {
		handleError(w, invalidBodyError, http.StatusBadRequest)
		return
	}

	// Unmarshal
	var userDTO UserDTO
	err = json.Unmarshal(b, &userDTO)
	if err != nil {
		handleError(w, invalidBodyError, http.StatusBadRequest)
		return
	}


	if isValidUserDTO(userDTO){
		handleError(w,invalidBodyError,http.StatusBadRequest)
		return
	}

	//Check if User with given Email is Existing
	user, err := userRepo.FetchUserByEmail(userDTO.Email)
	if err != nil || user.Id != 0 {
		handleError(w, invalidBodyError, http.StatusBadRequest)
		return
	}

	//Check if User with given Email is Existing
	user, err = userRepo.FetchUserByUsername(userDTO.Username)
	if err != nil || user.Id != 0 {
		handleError(w, invalidBodyError, http.StatusBadRequest)
		return
	}

	enableHash, enableToken, err := generateEnableToken(userDTO.Email)

	if err != nil {
		handleError(w, invalidBodyError, http.StatusBadRequest)
		return
	}

	log.Print("enableToken: " + enableToken)

	//Create user object
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)

	newUser := models.User{CreationTime: time.Now(), Email: userDTO.Email, Username: userDTO.Username, PasswordHash: passwordHash, IsEnabled: false, EnableSecretHash: enableHash}

	//Save user to DB
	user, err = userRepo.SaveUser(newUser)
	if err != nil {
		handleError(w, handlerError{Err: err, ErrorMessage: "unexpected error"}, http.StatusInternalServerError)
		return
	}

	err = writeJsonResponse(user, w)

	if err != nil {
		handleError(w, handlerError{Err: err, ErrorMessage: "unexpected error"}, http.StatusInternalServerError)
		return
	}
}

func generateEnableToken(identifier string) ([]byte, string, error) {

	enabledSecret := utils.RandStringRunes(32)
	enabledSecretHash, err := bcrypt.GenerateFromPassword([]byte(enabledSecret), bcrypt.DefaultCost)

	if err != nil {
		return []byte(""), "", err
	}
	enableToken := base64.StdEncoding.EncodeToString([]byte( identifier + ":" + enabledSecret))

	return enabledSecretHash, enableToken, nil
}

func UsersMe(w http.ResponseWriter, r *http.Request) {
	userId, err := jwt.GetUserIdByRequest(r)
	if err != nil {
		handleError(w, handlerError{Err: err, ErrorMessage: "not authenticated"}, http.StatusForbidden)
		return
	}

	user, err := database.GetUserRepository().FetchUserById(userId)

	if err != nil {
		handleError(w, handlerError{Err: err, ErrorMessage: "not authenticated"}, http.StatusForbidden)
		return
	}

	writeJsonResponse(user, w)
}
