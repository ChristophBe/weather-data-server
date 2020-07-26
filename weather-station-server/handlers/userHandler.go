package handlers

import (
	"de.christophb.wetter/config"
	"de.christophb.wetter/data/database"
	"de.christophb.wetter/data/models"
	"de.christophb.wetter/data/repositories"
	"de.christophb.wetter/data/transitory"
	"de.christophb.wetter/email"
	"de.christophb.wetter/handlers/httpHandler"
	"de.christophb.wetter/services"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type UserHandlers interface {
	GetCreateUserHandler() http.Handler
	GetUserMeHandler() http.Handler
	GetUserEnableHandler() http.Handler
}
type userHandlersImpl struct {
	tokenService         services.AuthTokenService
	userRepository       repositories.UserRepository
	invitationRepository repositories.InvitationRepository
}

func (u userHandlersImpl) GetCreateUserHandler() http.Handler {
	return httpHandler.JsonHandler(u.createUser)
}

func (u userHandlersImpl) GetUserMeHandler() http.Handler {
	return httpHandler.AuthorizedAppHandler(u.tokenService.VerifyUserAccessToken, u.usersMe)
}
func (u userHandlersImpl) GetUserEnableHandler() http.Handler {
	return httpHandler.JsonHandler(u.enableUser)
}

func GetUserHandlers() UserHandlers {
	return userHandlersImpl{
		tokenService:   services.GetAuthTokenService(),
		userRepository: database.GetUserRepository(),
		invitationRepository: database.GetInvitationRepository(),
	}
}

func (u userHandlersImpl) createUser(r *http.Request) (response interface{}, statusCode int) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		httpHandler.HandleBadRequest(InvalidBody, err)
	}

	// Unmarshal
	var body transitory.UserCreateBody
	err = json.Unmarshal(b, &body)
	if err != nil {
		httpHandler.HandleBadRequest(InvalidBody, err)
	}

	if !body.IsValid() {
		httpHandler.HandleBadRequest(InvalidBody, err)
	}

	invitationId, err := u.tokenService.VerifyUserInvitationToken(body.InvitationToken)
	if err != nil {
		httpHandler.HandleBadRequest("invalid invitation_token", err)
	}
	invitation, err := u.invitationRepository.FetchInvitationById(invitationId)
	if err != nil {
		httpHandler.HandleBadRequest("invalid invitation_token", err)
	}

	if u.userRepository.HasUserWithEmail(body.Email) {
		httpHandler.HandleBadRequest(InvalidBody, nil)
	}

	if u.userRepository.HasUserWithUsername(body.Username) {
		httpHandler.HandleBadRequest(InvalidBody, nil)
	}

	//Create user object
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		httpHandler.HandleInternalError(err)
	}

	newUser := models.User{
		CreationTime: time.Now(),
		Email:        strings.ToLower(body.Email),
		Username:     body.Username,
		PasswordHash: passwordHash,
		IsEnabled:    strings.ToLower(invitation.Email) == strings.ToLower(body.Email),
	}

	//Save user to DB
	user, err := u.userRepository.SaveUser(newUser)
	if err != nil {
		httpHandler.HandleInternalError(err)
	}

	if !user.IsEnabled {
		go sendEnableToken(user)
	}

	go handleInvitation(user, invitationId)
	statusCode = http.StatusOK
	response = user
	return
}


func (u userHandlersImpl) enableUser(r* http.Request)  (response interface{}, statusCode int)  {

	var body struct{
		Token string `json:"token"`
	}

	httpHandler.ReadJsonBody(r,&body)

	userId, err:= u.tokenService.VerifyUserEnableToken(body.Token)

	if err != nil{
		httpHandler.HandleBadRequest("invalid token",err)
	}

	user , err := u.userRepository.FetchUserById(userId)
	if err != nil{
		httpHandler.HandleBadRequest("invalid token",err)
	}

	user.IsEnabled = true

	user ,err  = u.userRepository.SaveUser(user)
	if err != nil{
		httpHandler.HandleInternalError(err)
	}
	statusCode = http.StatusOK
	response = user
	return
}

func (u userHandlersImpl) usersMe(userId int64, _ *http.Request)(response interface{},statusCode int){
	response, err := u.userRepository.FetchUserById(userId)

	if err != nil {
		httpHandler.HandleForbidden("not authorized", err)
	}
	statusCode = http.StatusOK
	return

}


func sendEnableToken(user models.User) {

	conf ,err := config.GetConfigManager().GetConfig()
	if err != nil{
		log.Fatal(err)
	}

	enableToken ,err := services.GetAuthTokenService().GenerateUserEnableToken(user)
	if err != nil{
		log.Fatal(err)
	}
	params := struct {
		Username string
		ActivationLink string
	}{
		Username: user.Username,
		ActivationLink: fmt.Sprintf("%s/users/enable?token=%s ",conf.FrontendBaseUrl,enableToken),
	}

	err=email.SendHtmlMail(user.Email,"Bestätige deine E-Mail Adresse","enableMailTemplate.html",params)

	if err != nil {
		log.Fatal(err)
	}
	log.Print(err)
}

func handleInvitation(user models.User, invitationId int64) {
	nodeRepo := database.GetMeasuringNodeRepository()

	nodes, err:= nodeRepo.FetchAllNodesByInvitationId(invitationId)

	if err!= nil{
		log.Fatalf("Failed to fetch nodes by invitation cause: %v\n",err)
		return
	}
	for _,node := range nodes{
		err:= nodeRepo.CreateAuthorisationRelation(node,user)

		if err!= nil{
			log.Panicln(fmt.Sprintf("Failed to add Node-Auth-Relation {nodeID: %d ,cause: %v}",node.Id,err))
		}
	}

}