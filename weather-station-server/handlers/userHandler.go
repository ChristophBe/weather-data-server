package handlers

import (
	"de.christophb.wetter/config"
	"de.christophb.wetter/data/database"
	"de.christophb.wetter/data/models"
	"de.christophb.wetter/data/repositories"
	"de.christophb.wetter/data/transitory"
	"de.christophb.wetter/email"
	"de.christophb.wetter/handlers/handlerUtil"
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
	return handlerUtil.AppHandler(u.createUser)
}

func (u userHandlersImpl) GetUserMeHandler() http.Handler {
	return handlerUtil.AuthorizedAppHandler(u.tokenService.VerifyUserAccessToken, u.usersMe)
}
func (u userHandlersImpl) GetUserEnableHandler() http.Handler {
	return handlerUtil.AppHandler(u.enableUser)
}

func GetUserHandlers() UserHandlers {
	return userHandlersImpl{
		tokenService:   services.GetAuthTokenService(),
		userRepository: database.GetUserRepository(),
		invitationRepository: database.GetInvitationRepository(),
	}
}

func (u userHandlersImpl) createUser(r *http.Request) (response interface{}, statusCode int, err error) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		err = handlerUtil.BadRequest(InvalidBody, err)
		return
	}

	// Unmarshal
	var body transitory.UserCreateBody
	err = json.Unmarshal(b, &body)
	if err != nil {
		err = handlerUtil.BadRequest(InvalidBody, err)
		return
	}

	if !body.IsValid() {
		err = handlerUtil.BadRequest(InvalidBody, err)
		return
	}

	invitationId, err := u.tokenService.VerifyUserInvitationToken(body.InvitationToken)
	if err != nil {
		err = handlerUtil.BadRequest("invalid invitation_token", err)
		return
	}
	invitation, err := u.invitationRepository.FetchInvitationById(invitationId)
	if err != nil {
		err = handlerUtil.BadRequest("invalid invitation_token", err)
		return
	}

	if u.userRepository.HasUserWithEmail(body.Email) {
		err = handlerUtil.BadRequest(InvalidBody, nil)
		return
	}

	if u.userRepository.HasUserWithUsername(body.Username) {
		err = handlerUtil.BadRequest(InvalidBody, nil)
		return
	}

	//Create user object
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)


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
		err = handlerUtil.InternalError(err)
		return
	}

	if !user.IsEnabled {
		go sendEnableToken(user)
	}

	go handleInvitation(user, invitationId)
	statusCode = http.StatusOK
	response = user
	return
}


func (u userHandlersImpl) enableUser(r* http.Request)  (response interface{}, statusCode int, err error)  {

	var body struct{
		Token string `json:"token"`
	}

	if err = readBody(r, &body);err != nil{
		err = handlerUtil.BadRequest("invalid body",err)
		return
	}

	userId, err:= u.tokenService.VerifyUserEnableToken(body.Token)

	if err != nil{
		err = handlerUtil.BadRequest("invalid token",err)
		return
	}

	user , err := u.userRepository.FetchUserById(userId)
	if err != nil{
		err = handlerUtil.BadRequest("invalid token",err)
		return
	}

	user.IsEnabled = true

	user ,err  = u.userRepository.SaveUser(user)
	if err != nil{
		err = handlerUtil.InternalError(err)
		return
	}
	statusCode = http.StatusOK
	response = user
	return
}

func (u userHandlersImpl) usersMe(userId int64, _ *http.Request)(response interface{},statusCode int,err error){
	response, err = u.userRepository.FetchUserById(userId)

	if err != nil {
		err = handlerUtil.Forbidden("not authorized", err)
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