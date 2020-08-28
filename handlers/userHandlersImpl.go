package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ChristophBe/weather-data-server/config"
	"github.com/ChristophBe/weather-data-server/data/models"
	"github.com/ChristophBe/weather-data-server/data/repositories"
	"github.com/ChristophBe/weather-data-server/data/transitory"
	"github.com/ChristophBe/weather-data-server/handlers/httpHandler"
	"github.com/ChristophBe/weather-data-server/services"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

type userHandlersImpl struct {
	mailService          services.MailService
	tokenService         services.AuthTokenService
	invitationService    services.InvitationService
	userRepository       repositories.UserRepository
	invitationRepository repositories.InvitationRepository
}

func (u userHandlersImpl) GetCreateUserHandler() http.Handler {
	return httpHandler.JsonHandler(u.createUser)
}

func (u userHandlersImpl) GetUserMeHandler() http.Handler {
	return httpHandler.AuthorizedAppHandler(u.tokenService.GetTokenVerifier(services.USER_AUTH), u.usersMe)
}
func (u userHandlersImpl) GetUserEnableHandler() http.Handler {
	return httpHandler.JsonHandler(u.enableUser)
}

func (u userHandlersImpl) createUser(r *http.Request) (response httpHandler.HandlerResponse, err error) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		err = httpHandler.BadRequest(httpHandler.ErrorMessageInvalidBody, err)
		return
	}

	// Unmarshal
	var body transitory.UserCreateBody
	err = json.Unmarshal(b, &body)
	if err != nil {
		err = httpHandler.BadRequest(httpHandler.ErrorMessageInvalidBody, err)
		return
	}

	if !body.IsValid() {
		err = httpHandler.BadRequest(httpHandler.ErrorMessageInvalidBody, err)
		return
	}

	invitationId, err := u.tokenService.GetTokenVerifier(services.USER_INVITATION)(body.InvitationToken)
	if err != nil {
		err = httpHandler.BadRequest("invalid invitation_token", err)
		return
	}
	invitation, err := u.invitationRepository.FetchInvitationById(invitationId)
	if err != nil {
		err = httpHandler.BadRequest("invalid invitation_token", err)
		return
	}

	if u.userRepository.HasUserWithUsername(body.Username) || u.userRepository.HasUserWithEmail(body.Email) {
		err = httpHandler.BadRequest(httpHandler.ErrorMessageInvalidBody, nil)
		return
	}

	//Create user object
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		err = httpHandler.InternalError(err)
		return
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
		err = httpHandler.InternalError(err)
		return
	}

	if !user.IsEnabled {
		go u.sendEnableToken(user)
	}

	go func() {
		err := u.invitationService.HandleInvitation(user, invitationId)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to handle invitation cause:%w", err))
		}
	}()
	response.Status = http.StatusOK
	response.Data = user
	return
}

func (u userHandlersImpl) enableUser(r *http.Request) (response httpHandler.HandlerResponse, err error) {

	var body struct {
		Token string `json:"token"`
	}

	err = httpHandler.ReadJsonBody(r, &body)

	if err != nil {
		err = httpHandler.BadRequest(httpHandler.ErrorMessageInvalidBody, err)
		return
	}

	userId, err := u.tokenService.GetTokenVerifier(services.USER_ENABLE)(body.Token)

	if err != nil {
		err = httpHandler.BadRequest("invalid token", err)
		return
	}

	user, err := u.userRepository.FetchUserById(userId)
	if err != nil {
		err = httpHandler.BadRequest("invalid token", err)
		return
	}

	user.IsEnabled = true

	user, err = u.userRepository.SaveUser(user)
	if err != nil {
		err = httpHandler.InternalError(err)
		return
	}
	response.Status = http.StatusOK
	response.Data = user
	return
}

func (u userHandlersImpl) usersMe(userId int64, _ *http.Request) (response httpHandler.HandlerResponse, err error) {
	user, err := u.userRepository.FetchUserById(userId)

	if err != nil {
		err = httpHandler.Forbidden("not authorized", err)
	}
	response.Status = http.StatusOK
	response.Data = user
	return

}

func (u userHandlersImpl) sendEnableToken(user models.User) {

	conf, err := config.GetConfigManager().GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	enableToken, err := services.GetAuthTokenService().GenerateUserEnableToken(user)
	if err != nil {
		log.Fatal(err)
	}
	params := struct {
		Username       string
		ActivationLink string
	}{
		Username:       user.Username,
		ActivationLink: fmt.Sprintf("%s/users/enable?token=%s ", conf.FrontendBaseUrl, enableToken),
	}

	subject := "Bestätige deine E-Mail Adresse"
	err = u.mailService.SendHtmlMail(mail.Address{Address: user.Email}, subject, "static/enableMailTemplate.html", params)

	if err != nil {
		log.Fatal(err)
	}
	log.Print(err)

}
