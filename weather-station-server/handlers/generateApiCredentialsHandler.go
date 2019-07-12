package handlers

import (
	"../data"
	"../jwt"
	"../utils"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)


type credentialsExternal struct {
	Secret string
	CreationTime time.Time
}
func GenerateApiCredentialsHandler(w http.ResponseWriter, request *http.Request) {
	con := data.CreateConnection()
	defer con.Close()
	vars := mux.Vars(request)
	nodeId, err := strconv.ParseInt(vars["nodeId"], 10, 64)

	_, err = data.FetchAuthTokenByNodeId(con, nodeId)


	userId, err := jwt.GetUserIdBy(request)
	if err != nil{

		handleError(w,handlerError{Err:err, ErrorMessage:"can not authenticate user"}, http.StatusForbidden)
		return
	}

	owner ,err :=  data.FetchOwnerByMesuringNode(nodeId)
	if err != nil || userId != owner.Id {
		handleError(w,handlerError{Err:err, ErrorMessage:"user is not owner"}, http.StatusForbidden)
	}


	secret:= utils.RandStringRunes(32)

	hash,err := bcrypt.GenerateFromPassword([]byte(secret),bcrypt.DefaultCost)
	if err != nil {
		handleError(w,handlerError{Err:err, ErrorMessage:"something went wrong"}, http.StatusInternalServerError)
		return
	}

	nodeToken := data.NodeAuthToken{TokenHash:hash,CreationTime:time.Now()}


	//inter := credentialsInternal{TokenHash: hash,ClientId:clientId}
	ext := credentialsExternal{Secret: secret,CreationTime:nodeToken.CreationTime}

	data.CreateNodeAuthToken(con,nodeId,nodeToken)

	writeJsonResponse(ext,w)
}

