package data

import (
	"io"
	"time"
)

const(
	createUserStatement                = "CREATE (m:User {lastLogin: {lastLogin}, creationTime: {creationTime},email: {email},username: {username},passwordHash: {passwordHash}})"
	fetchUserByUSerIdStatement         = "MATCH (u:User) where id(u)={userId} return id(u),u.lastLogin, u.creationTime, u.email, u.username, u.passwordHash"
	fetchUserByEmailStatement          = "MATCH (u:User) where u.email={email} return id(u),u.lastLogin, u.creationTime, u.email, u.username, u.passwordHash"
	fetchOwnerByMeasuringNodeStatement = "MATCH (u:User)-[:OWNER]->(n:MeasuringNode) WHERE id(n) = {nodeId} RETURN u"
	fetchUserByUsernameStatement       = "MATCH (u:User) where u.username={username} return id(u),u.lastLogin, u.creationTime, u.email, u.username, u.passwordHash"
)

type User struct {
	Id           int64     `json:"id"`
	LastLogin    time.Time `json:"timestamp"`
	CreationTime time.Time `json:"timestamp"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash []byte    `json:"-"`

}



func paresUserFormLine(row []interface{}) User {
	node := User{
		Id: row[0].(int64),
		LastLogin: time.Unix(row[1].(int64), 0),
		CreationTime: time.Unix(row[2].(int64),0),
		Email: row[3].(string),  Username: row[4].(string),
		PasswordHash: []byte(row[5].(string))}
	return node
}

func CreateUser( user User) {
	con := CreateConnection()
	defer con.Close()

	st := prepareStatement(createUserStatement, con)
	defer st.Close()

	result, err := st.ExecNeo(map[string]interface{}{
		"lastLogin":   	user.LastLogin.Unix(),
		"creationTime":	user.CreationTime.Unix(),
		"email":    	user.Email,
		"username": 	user.Username,
		"passwordHash": string(user.PasswordHash)})
	handleError(err)

	_, err = result.RowsAffected()
	handleError(err)
}

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
