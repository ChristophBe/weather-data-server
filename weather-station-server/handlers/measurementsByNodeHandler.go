package handlers

import (
	"../data"
	"../jwt"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)



func CheckNodePermissionForUser(r *http.Request) bool {

	nodeId, err := getNodeIDFormRequest(r)
	if err != nil{
		return false
	}


	nodeRepo := data.MeasuringNodeRepository{}
	node, err := nodeRepo.FetchAllMeasuringNodeById(nodeId)
	if err != nil{
		return false
	}

	if node.IsPublic {
		return true
	}

	userId, err := jwt.GetUserIdBy(r)
	if err != nil{
		return false
	}

	relations := nodeRepo.FetchAllMeasuringNodeUserRelations(nodeId,userId)
	return len(relations) > 0
}

func PostMeasurementForNodeHandler(w http.ResponseWriter, r *http.Request){
	con := data.CreateConnection()
	defer con.Close()

	nodeId, err := getNodeIDFormRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = NodeAuthorisationHandler(nodeId ,r)
	if err != nil {
		handleError(w,handlerError{ ErrorMessage: "no valid Credentials", Err:err},http.StatusForbidden)
		//http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		handleError(w,handlerError{Err:err, ErrorMessage:"Invalid Request Body"},http.StatusBadRequest)
		return
	}

	// Unmarshal
	var measuring data.Measuring
	err = json.Unmarshal(b, &measuring)
	if err != nil {
		handleError(w,handlerError{Err:err,ErrorMessage:"Invalid Request Body"}, http.StatusBadRequest)
		return
	}
	measuring.TimeStamp = time.Now()

	data.CreateMeasurement(con, nodeId, measuring)
}

func GetAllMeasurementsByNodeHandler(w http.ResponseWriter, r *http.Request){
	con := data.CreateConnection()
	defer con.Close()

	isNotAllowedToAccess := !CheckNodePermissionForUser(r)
	if isNotAllowedToAccess {
		handleError(w, handlerError{Err:nil,ErrorMessage:"Access Forbidden"}, http.StatusForbidden)
		return

	}


	nodeId, err := getNodeIDFormRequest(r)

	if err != nil {
		handleError(w, handlerError{Err:err,ErrorMessage:"missing NodeId"}, http.StatusBadRequest)
		return
	}

	measurements := data.FetchAllMeasuringByNodeId(con, nodeId)

	writeJsonResponse(measurements, w)
}


func GetLastMeasurementsByNodeHandler(w http.ResponseWriter, r *http.Request){
	con := data.CreateConnection()
	defer con.Close()

	nodeId, err := getNodeIDFormRequest(r)
	if err != nil {
		handleError(w, handlerError{Err:err,ErrorMessage:"missing NodeId"}, http.StatusBadRequest)
		return
	}


	limit,err := strconv.ParseInt(r.FormValue("limit"), 10, 64)

	if err != nil {
		handleError(w, handlerError{Err:err,ErrorMessage:"invalid value for param limit"}, http.StatusBadRequest)
		return
	}


	measurements := data.FetchLastMeasuringsByNodeId(con, nodeId, limit)

	writeJsonResponse(measurements,w)
}


func getNodeIDFormRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	nodeId, err := strconv.ParseInt(vars["nodeId"], 10, 64)
	return nodeId, err
}



func NodeAuthorisationHandler(nodeId int64, r *http.Request) error{

	con := data.CreateConnection()
	defer con.Close()

	reqToken := r.Header.Get("Authorization")
	splited := strings.Split(reqToken, " ")

	token := splited[len(splited)-1]

	authToken, err := data.FetchAuthTokenByNodeId(con, nodeId)

	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(authToken.TokenHash, []byte(token))

	if err != nil {
		return err
	}
	return nil

}

