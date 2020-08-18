package services

import (
	"fmt"
	"github.com/ChristophBe/weather-data-server/config"
	"github.com/ChristophBe/weather-data-server/data/models"
	"github.com/ChristophBe/weather-data-server/data/repositories"
	"github.com/ChristophBe/weather-data-server/email"
	"github.com/ChristophBe/weather-data-server/handlers/httpHandler"
	"log"
	"strconv"
	"time"
)

type shareMailParams struct {
	Username, NodeName, ActivationLink, NodeUrl string
	IsNewUser                                   bool
}

type shareNodeServiceImpl struct {
	authTokenService     AuthTokenService
	userRepository       repositories.UserRepository
	nodeRepository       repositories.MeasuringNodeRepository
	invitationRepository repositories.InvitationRepository
}

func (s shareNodeServiceImpl) ShareNode(node models.MeasuringNode, nodeOwner models.User, email string) (err error) {

	user, err := s.userRepository.FetchUserByEmail(email)

	isNewUser := err != nil || user.Id == 0

	conf, err := config.GetConfigManager().GetConfig()
	if err != nil {
		return err
	}

	emailParams := shareMailParams{
		Username:  nodeOwner.Username,
		NodeName:  node.Name,
		IsNewUser: isNewUser,
		NodeUrl:   conf.FrontendBaseUrl + "/nodes/" + strconv.Itoa(int(node.Id)),
	}

	if isNewUser {

		invitationToken, err := s.handleInvitationForNewUser(node, email)
		if err != nil {
			return err
		}

		emailParams.ActivationLink = conf.FrontendBaseUrl + "/users/create/" + invitationToken

	} else {
		err = s.nodeRepository.CreateAuthorisationRelation(node, user)
		if err != nil {
			err = httpHandler.InternalError(err)
			return
		}
	}

	recipient := user.Email
	if isNewUser {
		recipient = email
	}
	go s.sendShareMail(recipient, emailParams)
	return
}

func (s shareNodeServiceImpl) handleInvitationForNewUser(node models.MeasuringNode, email string) (invitationToken string, err error) {
	var invitation models.Invitation
	invitation, err = s.invitationRepository.FetchInvitationByEmail(email)

	if err != nil {
		invitation = models.Invitation{
			Email:        email,
			CreationTime: time.Now(),
		}

		invitation, err = s.invitationRepository.SaveInvitation(invitation)
		if err != nil {
			err = fmt.Errorf("can not save invitation caused by: %w", err)
			return
		}
	}

	err = s.invitationRepository.AddNodeToInvitation(invitation, node)
	if err != nil {
		err = fmt.Errorf("unabled add node(%d) to invitation(%d): %w", node.Id, invitation.Id, err)
		return
	}

	invitationToken, err = s.authTokenService.GenerateUserInvitationToken(invitation)
	if err != nil {
		err = httpHandler.InternalError(err)
		return
	}

	return
}

func (s shareNodeServiceImpl) sendShareMail(recipient string, params shareMailParams) {

	err := email.SendHtmlMail(recipient, "Die Wetterstation \""+params.NodeName+"\" wurde mit dir geteilt.", "shareNodeMailTemplate.html", params)

	if err != nil {
		log.Fatal(err)
	}

}
