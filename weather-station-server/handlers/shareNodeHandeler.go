package handlers

import (
	"de.christophb.wetter/configs"
	"de.christophb.wetter/data"
	"de.christophb.wetter/email"
	"de.christophb.wetter/jwt"
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

func ShareNodeHandler(w http.ResponseWriter, request *http.Request) {


	defer recoverHandlerErrors(w)
	nodeId, err := getNodeIDFormRequest(request)
	panicIfErrorNonNil(err,"missing nodeId", http.StatusNotFound)

	userId, err := jwt.GetUserIdBy(request)
	panicIfErrorNonNil(err, "can not authenticate user", http.StatusForbidden)

	owner ,err :=  data.GetUserRepository().FetchOwnerByMeasuringNode(nodeId)
	if err != nil || userId != owner.Id {
		panic(handlerError{Err:err, ErrorMessage:"user is not owner",Status: http.StatusForbidden})
	}

	var shareNodeDTO ShareNodeDTO
	err = readBody(request,&shareNodeDTO)
	panicIfErrorNonNil(err, InvalidBody, http.StatusBadRequest)



	nodeRepo := data.GetMeasuringNodeRepository()
	node, err := nodeRepo.FetchMeasuringNodeById(nodeId)
	panicIfErrorNonNil(err,"unexpected error",http.StatusInternalServerError)

	user, err := data.GetUserRepository().FetchUserByEmail(shareNodeDTO.Email)
	panicIfErrorNonNil(err,"unexpected error",http.StatusInternalServerError)

	isNewUser := user.Id == 0

	emailParams := shareMailParams{
		Username: owner.Username,
		NodeName: node.Name,
		IsNewUser: isNewUser,
		NodeUrl: configs.FRONTEND_BASE_URL + "/nodes/" +  strconv.Itoa(int(nodeId)),
	}

	if isNewUser {
		enableHash, enableToken ,err  := generateEnableToken(shareNodeDTO.Email)
		panicIfErrorNonNil(err,"unexpected error",http.StatusInternalServerError)

		emailParams.ActivationLink = configs.FRONTEND_BASE_URL + "/users/create/" + enableToken

		user.Email = shareNodeDTO.Email
		user.IsEnabled = false
		user.CreationTime = time.Now()
		user.EnableSecretHash = enableHash

		user,err = data.GetUserRepository().SaveUser(user)

		panicIfErrorNonNil(err,"unexpected error",http.StatusInternalServerError)
	}


	err = nodeRepo.CreateAuthorisationRelation(node,user)
	panicIfErrorNonNil(err,"unexpected error",http.StatusInternalServerError)

	go sendShareMail(user.Email, emailParams)

	respones := struct {
		Msg string `json:"message"`
	}{Msg:"the node was successfully shared"}

	err = writeJsonResponse(respones,w)
	panicIfErrorNonNil(err,"unexpected error",http.StatusInternalServerError)
}

func sendShareMail( recipient string, params shareMailParams)  {

	err:=email.SendHtmlMail(recipient,"Die Wetterstation \"" + params.NodeName + "\" wurde mit dir geteilt.","shareNodeMailTemplate.html",params)

	if err != nil {
		log.Fatal(err)
	}
	log.Print(err)
}