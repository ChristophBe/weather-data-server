package handlers

import (
	"../data"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func RandStringRunes(n int) string {

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

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

	if err == nil {

		handleError(w, handlerError{Err:err,ErrorMessage:"can not create token"}, http.StatusBadRequest)
		return
	}

	//TODO: Check if current user is owner of this measuring Node
	secret:= RandStringRunes(32)

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

