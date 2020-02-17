package handlers

import (
	"de.christophb.wetter/config"
	"de.christophb.wetter/data/database"
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

	userId, err := jwt.GetUserIdByRequest(request)
	panicIfErrorNonNil(err, "can not authenticate user", http.StatusForbidden)

	owner ,err :=  database.GetUserRepository().FetchOwnerByMeasuringNode(nodeId)
	if err != nil || userId != owner.Id {
		panic(handlerError{Err:err, ErrorMessage:"user is not owner",Status: http.StatusForbidden})
	}

	var shareNodeDTO ShareNodeDTO
	err = readBody(request,&shareNodeDTO)
	panicIfErrorNonNil(err, InvalidBody, http.StatusBadRequest)



	nodeRepo := database.GetMeasuringNodeRepository()
	node, err := nodeRepo.FetchMeasuringNodeById(nodeId)
	panicIfErrorNonNil(err,"unexpected error",http.StatusInternalServerError)

	user, err := database.GetUserRepository().FetchUserByEmail(shareNodeDTO.Email)
	panicIfErrorNonNil(err,"unexpected error",http.StatusInternalServerError)

	isNewUser := user.Id == 0


	conf, err := config.GetConfigManager().GetConfig()
	panicIfErrorNonNil(err,"unexpected error",http.StatusInternalServerError)

	emailParams := shareMailParams{
		Username:  owner.Username,
		NodeName:  node.Name,
		IsNewUser: isNewUser,
		NodeUrl:   conf.FrontendBaseUrl + "/nodes/" +  strconv.Itoa(int(nodeId)),
	}

	if isNewUser {
		enableHash, enableToken ,err  := generateEnableToken(shareNodeDTO.Email)
		panicIfErrorNonNil(err,"unexpected error",http.StatusInternalServerError)

		emailParams.ActivationLink = conf.FrontendBaseUrl + "/users/create/" + enableToken

		user.Email = shareNodeDTO.Email
		user.IsEnabled = false
		user.CreationTime = time.Now()
		user.EnableSecretHash = enableHash

		user,err = database.GetUserRepository().SaveUser(user)

		panicIfErrorNonNil(err,"unexpected error",http.StatusInternalServerError)
	}


	err = nodeRepo.CreateAuthorisationRelation(node,user)
	panicIfErrorNonNil(err,"unexpected error",http.StatusInternalServerError)

	go sendShareMail(user.Email, emailParams)

	respones := struct {
		Msg string `json:"message"`
	}{Msg:"the node was successfully shared"}

	err = WriteJsonResponse(respones,w)
	panicIfErrorNonNil(err,"unexpected error",http.StatusInternalServerError)
}

func sendShareMail( recipient string, params shareMailParams)  {

	err:=email.SendHtmlMail(recipient,"Die Wetterstation \"" + params.NodeName + "\" wurde mit dir geteilt.","shareNodeMailTemplate.html",params)

	if err != nil {
		log.Fatal(err)
	}
	log.Print(err)
}