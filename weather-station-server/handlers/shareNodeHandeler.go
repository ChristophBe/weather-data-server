package handlers

import (
	"de.christophb.wetter/config"
	"de.christophb.wetter/data/database"
	"de.christophb.wetter/data/models"
	"de.christophb.wetter/email"
	"de.christophb.wetter/handlers/handlerUtil"
	"de.christophb.wetter/services"
	"log"
	"net/http"
	"strconv"
	"time"
)


type ShareNodeDTO struct {
	Email string `json:"email"`
}

type shareMailParams struct {
	Username, NodeName, ActivationLink, NodeUrl string
	IsNewUser bool
}

func ShareNodeHandler(userId int64, request *http.Request)(response interface{},statusCode int ,err error) {

	userRepo := database.GetUserRepository()
	invitationRepo := database.GetInvitationRepository()
	nodeRepo := database.GetMeasuringNodeRepository()

	node, err := getNodeFormRequest(request)
	if err!= nil{
		err = handlerUtil.NotFound("node not found",err)
	}

	owner ,err :=  userRepo.FetchOwnerByMeasuringNode(node.Id)
	if err != nil || userId != owner.Id {
		err = handlerUtil.Forbidden("user is not owner",err)
	}

	var shareNodeDTO ShareNodeDTO
	err = readBody(request,&shareNodeDTO)
	if err != nil{
		err = handlerUtil.InternalError(err)
		return
	}

	user, _ := userRepo.FetchUserByEmail(shareNodeDTO.Email)

	isNewUser := user.Id == 0

	conf, err := config.GetConfigManager().GetConfig()
	if err != nil {
		err = handlerUtil.InternalError(err)
		return
	}

	emailParams := shareMailParams{
		Username:  owner.Username,
		NodeName:  node.Name,
		IsNewUser: isNewUser,
		NodeUrl:   conf.FrontendBaseUrl + "/nodes/" +  strconv.Itoa(int(node.Id)),
	}

	if isNewUser {

		var invitation models.Invitation
		invitation, err = invitationRepo.FetchInvitationByEmail(shareNodeDTO.Email)
		if err != nil {
			invitation = models.Invitation{
				Email: shareNodeDTO.Email,
				CreationTime: time.Now(),
			}
			log.Print(invitation)
			invitation ,err= invitationRepo.SaveInvitation(invitation)
			if err != nil {
				return
			}
		}

		err= invitationRepo.AddNodeToInvitation(invitation,node)
		if err != nil {
			return
		}

		var invitationToken string
		invitationToken, err = services.GetAuthTokenService().GenerateUserInvitationToken(invitation)
		if err != nil {
			return
		}

		emailParams.ActivationLink = conf.FrontendBaseUrl + "/users/create/" + invitationToken

	}	else {
		err = nodeRepo.CreateAuthorisationRelation(node,user)
		if err != nil {
			return
		}
	}


	recipient :=user.Email
	if isNewUser {
		recipient = shareNodeDTO.Email
	}
	go sendShareMail(recipient, emailParams)

	response = struct {
		Msg string `json:"message"`
	}{Msg:"the node was successfully shared"}
	statusCode = http.StatusOK
	return
}

func sendShareMail( recipient string, params shareMailParams)  {

	err:=email.SendHtmlMail(recipient,"Die Wetterstation \"" + params.NodeName + "\" wurde mit dir geteilt.","shareNodeMailTemplate.html",params)

	if err != nil {
		log.Fatal(err)
	}
	log.Print(err)
}