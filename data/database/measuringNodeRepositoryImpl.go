package database

import (
	"errors"
	"github.com/ChristophBe/weather-data-server/data/models"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type measuringNodeRepositoryImpl struct{}

func (measuringNodeRepositoryImpl) parseMeasuringNodeFromRecord(record neo4j.Record) (interface{}, error) {

	nodeData, ok := record.Get("m")
	if !ok {
		err := errors.New("can not parse measuring form record")
		return nil, err
	}
	node := nodeData.(neo4j.Node)
	props := node.Props()
	measuringNode := models.MeasuringNode{
		Id:         node.Id(),
		Name:       parseStringProp(props["name"], ""),
		Lat:        parseFloatProp(props["lat"], 0.0),
		Lng:        parseFloatProp(props["lng"], 0.0),
		IsPublic:   parseBoolProp(props["isPublic"], false),
		IsOutdoors: parseBoolProp(props["isOutdoors"], true),
	}
	return measuringNode, nil
}

func (measuringNodeRepositoryImpl) castListOfMeasuringNodes(input interface{}) (nodes []models.MeasuringNode) {
	items := input.([]interface{})
	nodes = make([]models.MeasuringNode, len(items))
	for key, x := range items {
		nodes[key] = x.(models.MeasuringNode)
	}
	return
}

func (m measuringNodeRepositoryImpl) FetchMeasuringNodeById(nodeId int64) (measuringNode models.MeasuringNode, err error) {
	params := map[string]interface{}{"nodeId": nodeId}

	stmt := "MATCH (m:MeasuringNode) WHERE id(m) = $nodeId RETURN  m"
	result, err := doReadTransaction(stmt, params, parseSingleItemFromResult(m.parseMeasuringNodeFromRecord))

	if err != nil {
		return
	}
	return result.(models.MeasuringNode), nil
}

func (measuringNodeRepositoryImpl) FetchAllMeasuringNodeUserRelations(nodeId int64, userId int64) ([]string, error) {
	params := map[string]interface{}{
		"userId": userId,
		"nodeId": nodeId,
	}

	stmt := "MATCH (u:User)-[r]->(n:MeasuringNode) WHERE id(u) = $userId and id(n) = $nodeId return type(r)"

	result, err := doReadTransaction(stmt, params, func(result neo4j.Result) (res interface{}, err error) {
		var relations []string

		if result.Err() != nil {
			return relations, result.Err()
		}
		for result.Next() {
			relation := result.Record().GetByIndex(0).(string)
			relations = append(relations, relation)
		}
		return relations, nil
	})

	if err != nil {
		return make([]string, 0), err
	}
	return result.([]string), nil
}

func (m measuringNodeRepositoryImpl) FetchAllMeasuringNodes() ([]models.MeasuringNode, error) {
	params := map[string]interface{}{}

	stmt := "MATCH (m:MeasuringNode) RETURN m"
	results, err := doReadTransaction(stmt, params, parseListFromResult(m.parseMeasuringNodeFromRecord))

	if err != nil {
		return []models.MeasuringNode{}, err
	}
	return results.([]models.MeasuringNode), nil
}

func (m measuringNodeRepositoryImpl) CreateMeasuringNode(node models.MeasuringNode, userId int64) (models.MeasuringNode, error) {
	params := map[string]interface{}{
		"name":       node.Name,
		"lat":        node.Lat,
		"lng":        node.Lng,
		"isPublic":   node.IsPublic,
		"isOutdoors": node.IsOutdoors,
		"ownerId":    userId,
	}

	stmt := "MATCH (o:User) WHERE id(o) = $ownerId CREATE (o)-[:OWNER]->(m:MeasuringNode {name: $name, lat: $lat, lng: $lng, isPublic: $isPublic, isOutdoors: $isOutdoors}) RETURN m"
	result, err := doWriteTransaction(stmt, params, parseSingleItemFromResult(m.parseMeasuringNodeFromRecord))

	if err != nil {
		return models.MeasuringNode{}, err
	}
	return result.(models.MeasuringNode), nil
}

func (m measuringNodeRepositoryImpl) FetchNodesOwnedByUserId(userId int64) ([]models.MeasuringNode, error) {

	params := map[string]interface{}{
		"userId": userId,
	}

	stmt := "MATCH (m:MeasuringNode)<-[:OWNER]-(u:User) WITH m, u WHERE id(u) = $userId RETURN m"
	results, err := doReadTransaction(stmt, params, parseListFromResult(m.parseMeasuringNodeFromRecord))
	if err != nil {
		return []models.MeasuringNode{}, err
	}
	return m.castListOfMeasuringNodes(results), nil
}

func (m measuringNodeRepositoryImpl) FetchAllPublicNodes() ([]models.MeasuringNode, error) {

	params := map[string]interface{}{}

	stmt := "MATCH (m:MeasuringNode) WHERE m.isPublic = true RETURN m"
	results, err := doReadTransaction(stmt, params, parseListFromResult(m.parseMeasuringNodeFromRecord))
	if err != nil {
		return []models.MeasuringNode{}, err
	}
	return m.castListOfMeasuringNodes(results), nil
}

func (m measuringNodeRepositoryImpl) FetchAllVisibleNodesByUserId(userId int64) ([]models.MeasuringNode, error) {

	params := map[string]interface{}{
		"userId": userId,
	}

	stmt := "MATCH (m:MeasuringNode) WITH m OPTIONAL MATCH (m)<-[]-(u:User) WITH m, u WHERE m.isPublic OR id(u) = $userId RETURN m"
	results, err := doReadTransaction(stmt, params, parseListFromResult(m.parseMeasuringNodeFromRecord))
	if err != nil {
		return []models.MeasuringNode{}, err
	}
	return m.castListOfMeasuringNodes(results), nil
}
func (m measuringNodeRepositoryImpl) FetchAllNodesByInvitationId(invitationId int64) ([]models.MeasuringNode, error) {

	params := map[string]interface{}{
		"invitationId": invitationId,
	}

	stmt := "MATCH (m:MeasuringNode)<-[:INVITATION_FOR]-(i:Invitation) WHERE id(i) = $invitationId RETURN m"
	results, err := doReadTransaction(stmt, params, parseListFromResult(m.parseMeasuringNodeFromRecord))
	if err != nil {
		return []models.MeasuringNode{}, err
	}
	return m.castListOfMeasuringNodes(results), nil
}

func (m measuringNodeRepositoryImpl) CreateAuthorisationRelation(node models.MeasuringNode, user models.User) (err error) {
	params := map[string]interface{}{
		"nodeId": node.Id,
		"userId": user.Id,
	}

	stmt := "MATCH (u:User), (n:MeasuringNode) WHERE id(u) = $userId and id(n) = $nodeId CREATE (u)-[r:IS_AUTHORIZED]->(n) RETURN r"
	_, err = doWriteTransaction(stmt, params, func(result neo4j.Result) (res interface{}, err error) {
		return nil, result.Err()
	})
	return
}
