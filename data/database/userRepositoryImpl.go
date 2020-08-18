package database

import (
	"errors"
	"github.com/ChristophBe/weather-data-server/data/models"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"time"
)

type userRepositoryImpl struct{}

func (u userRepositoryImpl) parseUserFormRecord(record neo4j.Record) (res interface{}, err error) {
	nodeData, ok := record.Get("u")
	if !ok {
		err = errors.New("can not parse user-node form record")
		return
	}

	node := nodeData.(neo4j.Node)
	props := node.Props()
	user := models.User{
		Id:               node.Id(),
		LastLogin:        parseTimeProp(props["lastLogin"], time.Unix(0, 0)),
		CreationTime:     parseTimeProp(props["creationTime"], time.Unix(0, 0)),
		Email:            parseStringProp(props["email"], ""),
		Username:         parseStringProp(props["username"], ""),
		PasswordHash:     parseByteArrayProp(props["passwordHash"], []byte("")),
		IsEnabled:        parseBoolProp(props["isEnabled"], false),
		EnableSecretHash: parseByteArrayProp(props["enableSecretHash"], []byte("")),
	}

	res = user
	return
}

func (u userRepositoryImpl) SaveUser(user models.User) (savedUser models.User, err error) {

	params := map[string]interface{}{
		"username":         user.Username,
		"email":            user.Email,
		"lastLogin":        user.LastLogin.Unix(),
		"creationTime":     user.CreationTime.Unix(),
		"isEnabled":        user.IsEnabled,
		"enableSecretHash": string(user.EnableSecretHash),
		"passwordHash":     string(user.PasswordHash)}
	if user.Id != 0 {
		params["id"] = user.Id
	}

	insertUserStatement := "CREATE (u:User {lastLogin: {lastLogin}, creationTime: {creationTime},email: {email},username: {username},isEnabled: {isEnabled},enableSecretHash: {enableSecretHash}, passwordHash: {passwordHash}}) RETURN u"
	updateUserStatement := "Match (u:User) WHERE id(u) = {id} SET u.lastLogin = {lastLogin}, u.isEnabled = {isEnabled},u.username = {username}, u.passwordHash = {passwordHash} RETURN u"

	result, err := saveNode(insertUserStatement, updateUserStatement, params, parseSingleItemFromResult(u.parseUserFormRecord))

	if err != nil {
		return
	}
	savedUser = result.(models.User)
	return
}

func (u userRepositoryImpl) FetchUserById(userId int64) (user models.User, err error) {
	params := map[string]interface{}{"userId": userId}
	stmt := "MATCH (u:User) where id(u)={userId} return u"
	res, err := doReadTransaction(stmt, params, parseSingleItemFromResult(u.parseUserFormRecord))
	if err != nil {
		return
	}
	user = res.(models.User)
	return
}

func (u userRepositoryImpl) FetchOwnerByMeasuringNode(nodeId int64) (user models.User, err error) {
	params := map[string]interface{}{"nodeId": nodeId}
	stmt := "MATCH (u:User)-[:OWNER]->(n:MeasuringNode) WHERE id(n) = {nodeId} RETURN u"
	res, err := doReadTransaction(stmt, params, parseSingleItemFromResult(u.parseUserFormRecord))
	if err != nil {
		return
	}
	user = res.(models.User)
	return
}

func (u userRepositoryImpl) FetchUserByEmail(email string) (user models.User, err error) {
	params := map[string]interface{}{"email": email}
	stmt := "MATCH (u:User) WHERE u.email = {email} RETURN u"
	res, err := doReadTransaction(stmt, params, parseSingleItemFromResult(u.parseUserFormRecord))
	if err != nil {
		return
	}
	user = res.(models.User)
	return
}

func (u userRepositoryImpl) FetchUserByUsername(username string) (user models.User, err error) {
	params := map[string]interface{}{"username": username}
	stmt := "MATCH (u:User) WHERE u.username = {username} RETURN u"
	res, err := doReadTransaction(stmt, params, parseSingleItemFromResult(u.parseUserFormRecord))
	if err != nil {
		return
	}
	user = res.(models.User)
	return
}

func (u userRepositoryImpl) HasUserWithEmail(email string) bool {
	user, err := u.FetchUserByEmail(email)
	return err == nil && user.Id != 0
}

func (u userRepositoryImpl) HasUserWithUsername(username string) bool {
	user, err := u.FetchUserByUsername(username)
	return err == nil && user.Id != 0
}
