package data

import (
	"fmt"
	"io"
	"time"

	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

const(
	createNodeAuthTokenStatement  = "MATCH (n:MeasuringNode)  WHERE Id(n) = {nodeId} CREATE (n)<-[:AUTH_TOKEN_FOR]-(m:NodeAuthToken {creationTime: {creationTime}, tokenHash: {tokenHash}})"
	fetchNodeAuthTokenForNodeStmt = "MATCH (t:NodeAuthToken)-[:AUTH_TOKEN_FOR]->(n:MeasuringNode) WHERE id(n) = {nodeId} RETURN id(t), t.creationTime, t.tokenHash"
)

type NodeAuthToken struct {
	Id           int64
	TokenHash    []byte
	CreationTime time.Time
}



func CreateNodeAuthToken( con bolt.Conn, nodeId int64, token NodeAuthToken) {

	st := prepareStatement(createNodeAuthTokenStatement, con)

	result, err := st.ExecNeo(map[string]interface{}{
		"creationTime":   token.CreationTime.Unix(),
		"tokenHash":    string(token.TokenHash),
		"nodeId": nodeId})
	handleError(err)

	numResult, err := result.RowsAffected()
	handleError(err)

	fmt.Printf("CREATED ROWS: %d\n", numResult) // CREATED ROWS: 1

	// Closing the statment will also close the rows
	st.Close()
}


func FetchAuthTokenByNodeId(con bolt.Conn, nodeId int64) (NodeAuthToken,error){

	st:= prepareStatement(fetchNodeAuthTokenForNodeStmt,con)
	rows := queryStatement(st,map[string]interface{}{"nodeId":nodeId})
	//var measuring []Measuring
	var authToken NodeAuthToken

	row, _, err := rows.NextNeo()
	if err != nil && err != io.EOF {
		panic(err)
	} else if err != io.EOF {

		// id(m), m.timeStamp, m.temperature, m.humidity, m.pressure

		fmt.Println(row)

		


		authToken = NodeAuthToken{Id:row[0].(int64), CreationTime:time.Unix(row[1].(int64),0),TokenHash:[]byte(row[2].(string))}
	}


	st.Close()

	return authToken, err
}
