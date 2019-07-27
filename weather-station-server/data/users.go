package data

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"io"
	"log"
	"time"
)

const(
	userReturnString                   = "id(u),u.lastLogin, u.creationTime, u.email, u.username, u.passwordHash, u.isEnabled, u.enableSecretHash"
	updateUserStatement                = "Match (u:User) WHERE id(u) = {userId} SET u.lastLogin = {lastLogin}, u.isEnabled = {isEnabled},u.username = {username}, u.passwordHash = {passwordHash} RETURN " + userReturnString
	insertUserStatement                = "CREATE (u:User {lastLogin: {lastLogin}, creationTime: {creationTime},email: {email},username: {username},isEnabled: {isEnabled},enableSecretHash: {enableSecretHash}, passwordHash: {passwordHash}}) RETURN " + userReturnString
	fetchUserByUSerIdStatement         = "MATCH (u:User) where id(u)={userId} return " + userReturnString
	fetchUserByEmailStatement          = "MATCH (u:User) where u.email={email} return " + userReturnString
	fetchOwnerByMeasuringNodeStatement = "MATCH (u:User)-[:OWNER]->(n:MeasuringNode) WHERE id(n) = {nodeId} RETURN " + userReturnString
	fetchUserByUsernameStatement       = "MATCH (u:User) where u.username={username} return " + userReturnString
)

type User struct {
	Id           int64     `json:"id"`
	LastLogin    time.Time `json:"last-login"`
	CreationTime time.Time `json:"creation-time"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	IsEnabled    bool      `json:"enabled"`
	EnableSecretHash []byte`json:"-"`
	PasswordHash []byte    `json:"-"`

}

func userResultHandler(result neo4j.Result) (interface{},error){
	if result.Next() {
		user := paresUserFormLine(result.Record().Values())
		return user, nil
	}
	return nil, result.Err()
}


func paresUserFormLine(row []interface{}) User {
	user := User{
		Id:               parseRowInt(row[0],0),
		LastLogin:        parseRowTime(row[1],time.Unix(0,0)),
		CreationTime:     parseRowTime(row[2],time.Unix(0,0)),
		Email:            parseRowString(row[3],""),
		Username:         parseRowString(row[4],""),
		PasswordHash:     parseRowBytes(row[5],[]byte("")),
		IsEnabled:        parseRowBool(row[6],false),
		EnableSecretHash: parseRowBytes(row[7],[]byte(""))}

	return user
}

func UpsertUser( user User)(User,error) {
	if user.Id == 0 {
		return insertUser(user)
	}else {
		return updateUser(user)
	}
}

/*func updateUser( user User) {
	con := CreateConnection()
	defer con.Close()



	st := prepareStatement(updateUserStatement, con)
	defer st.Close()

	result, err := st.ExecNeo(map[string]interface{}{
		"userId":   	user.Id,
		"username": 	user.Username,
		"lastLogin":   	user.LastLogin.Unix(),
		"isEnabled": 	user.IsEnabled,
		"enableSecretHash": string(user.EnableSecretHash),
		"passwordHash": string(user.PasswordHash)})
	handleError(err)

	_, err = result.RowsAffected()
	handleError(err)
}*/

func updateUser( user User)(insertedUser User,err error){
	params := map[string]interface{}{
		"userId":   	user.Id,
		"username": 	user.Username,
		"lastLogin":   	user.LastLogin.Unix(),
		"isEnabled": 	user.IsEnabled,
		"enableSecretHash": string(user.EnableSecretHash),
		"passwordHash": string(user.PasswordHash)}

	res, err:= doWriteTransaction(updateUserStatement, params , userResultHandler)

	log.Print(res, err)
	if err != nil{
		return User{}, err
	}

	return res.(User) ,nil
}

func insertUser( user User) (insertedUser User,err error) {
	params := map[string]interface{}{
		"lastLogin":   	user.LastLogin.Unix(),
		"creationTime":	user.CreationTime.Unix(),
		"email":    	user.Email,
		"username": 	user.Username,
		"isEnabled": 	user.IsEnabled,
		"enableSecretHash": string(user.EnableSecretHash),
		"passwordHash": string(user.PasswordHash)}

	res, err:=doWriteTransaction(insertUserStatement, params , userResultHandler)

	log.Print(res, err)
	if err != nil{
		return User{}, err
	}

	return res.(User) ,nil
}
/*
func insertUser( user User) (User,error){


	con := CreateConnection()
	defer con.Close()



	st := prepareStatement(insertUserStatement, con)
	defer st.Close()

	result, err := st.ExecNeo(map[string]interface{}{
		"lastLogin":   	user.LastLogin.Unix(),
		"creationTime":	user.CreationTime.Unix(),
		"email":    	user.Email,
		"username": 	user.Username,
		"isEnabled": 	user.IsEnabled,
		"enableSecretHash": string(user.EnableSecretHash),
		"passwordHash": string(user.PasswordHash)})
	handleError(err)

	_, err = result.RowsAffected()
	handleError(err)
}

 */

func fetchUser(statement string, params map[string]interface{}) (User,error){

	con := CreateConnection()
	defer con.Close()

	st:= prepareStatement(statement,con)
	defer st.Close()


	rows := queryStatement(st,params)

	var user User

	row, _, err := rows.NextNeo()
	if err != nil && err != io.EOF {
		return user, err

	} else if err != io.EOF {
		user = paresUserFormLine(row)
	}
	return user , nil
}
func FetchUserById( userId int64) (User,error){
	return fetchUser(fetchUserByUSerIdStatement, map[string]interface{}{"userId": userId })
}
func FetchOwnerByMesuringNode(nodeId int64) (User,error){
	return fetchUser(fetchOwnerByMeasuringNodeStatement, map[string]interface{}{"nodeId": nodeId })
}
func FetchUserByEmail(email string) (User,error){

	con := CreateConnection()
	defer con.Close()

	st:= prepareStatement(fetchUserByEmailStatement,con)
	defer st.Close()


	rows := queryStatement(st,map[string]interface{}{"email":email})

	var user User

	row, _, err := rows.NextNeo()
	if err != nil && err != io.EOF {
		return user, err

	} else if err != io.EOF {
		user = paresUserFormLine(row)
	}
	return user , nil
}

func FetchUserByUsername(username string) (User,error){

	con := CreateConnection()
	defer con.Close()

	st:= prepareStatement(fetchUserByUsernameStatement,con)
	defer st.Close()


	rows := queryStatement(st,map[string]interface{}{"username": username})

	var user User

	row, _, err := rows.NextNeo()
	if err != nil && err != io.EOF {
		return user, err

	} else if err != io.EOF {
		user = paresUserFormLine(row)
	}
	return user , nil
}
