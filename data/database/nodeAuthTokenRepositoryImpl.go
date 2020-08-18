package database

import (
	"de.christophb.wetter/data/models"
	"errors"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"time"
)

type nodeAuthTokenRepositoryIml struct{}

func (nodeAuthTokenRepositoryIml) tokenResultHandler(record neo4j.Record) (result interface{}, err error) {

	rawNode, ok := record.Get("t")
	if !ok {
		err = errors.New("can not parse node auth token")
		return
	}
	node := rawNode.(neo4j.Node)
	props := node.Props()
	nodeAuthToken := models.NodeAuthToken{
		Id:           node.Id(),
		CreationTime: parseTimeProp(props["creationTime"], time.Unix(0, 0)),
		TokenHash:    parseByteArrayProp(props["tokenHash"], []byte("")),
	}
	return nodeAuthToken, nil
}

func (r nodeAuthTokenRepositoryIml) InsertNodeAuthToken(nodeId int64, token models.NodeAuthToken) (savedToken models.NodeAuthToken, err error) {
	params := map[string]interface{}{
		"creationTime": token.CreationTime.Unix(),
		"tokenHash":    string(token.TokenHash),
		"nodeId":       nodeId,
	}

	stmt := "MATCH (n:MeasuringNode)  WHERE Id(n) = {nodeId} CREATE (n)<-[:AUTH_TOKEN_FOR]-(m:NodeAuthToken {creationTime: {creationTime}, tokenHash: {tokenHash}}) RETURN t"

	res, err := doWriteTransaction(stmt, params, parseSingleItemFromResult(r.tokenResultHandler))

	if err != nil {
		return
	}
	savedToken = res.(models.NodeAuthToken)
	return
}

func (r nodeAuthTokenRepositoryIml) FetchAuthTokenByNodeId(nodeId int64) (token models.NodeAuthToken, err error) {
	params := map[string]interface{}{"nodeId": nodeId}

	stmt := "MATCH (t:NodeAuthToken)-[:AUTH_TOKEN_FOR]->(n:MeasuringNode) WHERE id(n) = {nodeId} RETURN t"

	res, err := doWriteTransaction(stmt, params, parseSingleItemFromResult(r.tokenResultHandler))

	if err != nil {
		return
	}
	token = res.(models.NodeAuthToken)
	return
}
