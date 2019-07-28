package data

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type MeasuringNode struct {
	Id   int64	  	`json:"id"`
	Name string 	`json:"name"`
	Lat  float64 	`json:"lat"`
	Lng  float64 	`json:"lng"`
	IsPublic bool 	`json:"is_public"`
	IsOutdoors bool `json:"is_outdoors"`
}

const(
	measuringModeQueryString                  = "id(m), m.name, m.lat, m.lng, m.isPublic, m.isOutdoors"
	createStationStatement                    = "MATCH (o:User) WHERE id(o) = $ownerId CREATE (o)-[:OWNER]->(m:MeasuringNode {name: $name, lat: $lat, lng: $lng, isPublic: $isPublic, isOutdoors: $isOutdoors}) RETURN " + measuringModeQueryString
	fetchMeasuringNodesStmt                   = "MATCH (m:MeasuringNode) RETURN " + measuringModeQueryString
	fetchAllPublicMeasuringNodesStmt          = "MATCH (m:MeasuringNode) WHERE m.isPublic = true RETURN " + measuringModeQueryString
	fetchAllVisibleMeasuringNodesByUserIdStmt = "MATCH (m:MeasuringNode) WITH m OPTIONAL MATCH (m)<-[]-(u:User) WITH m, u WHERE m.isPublic OR id(u) = $userId RETURN " + measuringModeQueryString
	fetchAllOwnedMeasuringNodesByUserIdStmt   = "MATCH (m:MeasuringNode)<-[:OWNER]-(u:User) WITH m, u WHERE id(u) = $userId RETURN " + measuringModeQueryString
	fetchMeasuringNodesByIdStmt               = "MATCH (m:MeasuringNode) WHERE id(m) = $nodeId RETURN  "+ measuringModeQueryString
	fetchMeasuringNodesUserRelations          = "MATCH (u:User)-[r]->(n:MeasuringNode) WHERE id(u) = $userId and id(n) = $nodeId return type(r)"
	createAuthorisationRelation 			  = "MATCH (u:User), (n:MeasuringNode) WHERE id(u) = $userId and id(n) = $nodeId CREATE (u)-[r:IS_AUTHORIZED]->(n) RETURN r"
)


type MeasuringNodeRepository struct {}


func (MeasuringNodeRepository) parseMeasuringNodeFromRow(row []interface{}) MeasuringNode {
	node := MeasuringNode{
		Id:			parseRowInt(row[0],0),
		Name:		parseRowString(row[1],""),
		Lat:		parseRowFloat(row[2],0.0),
		Lng:		parseRowFloat(row[3],0.0),
		IsPublic: 	parseRowBool(row[4],false),
		IsOutdoors: parseRowBool(row[5],true),
	}
	return node
}

func (m *MeasuringNodeRepository) handleMeasuringNodeListResultHandler(result neo4j.Result)(interface{},error){
	results := make([]MeasuringNode,0)
	for ; result.Next(); {
		node := m.parseMeasuringNodeFromRow(result.Record().Values())

		results = append(results, node)
	}
	if result.Err() != nil{
		return nil, result.Err()
	}
	return results, nil
}

func (m *MeasuringNodeRepository) handleMeasuringNodeResultHandler(result neo4j.Result)(interface{},error){
	if result.Next() {
		node := m.parseMeasuringNodeFromRow(result.Record().Values())
		return node, nil
	}
	return nil, result.Err()
}

func (m *MeasuringNodeRepository) FetchMeasuringNodeById(nodeId int64 ) (MeasuringNode,error) {
	params := map[string]interface{}{"nodeId": nodeId}

	result, err := doReadTransaction(fetchMeasuringNodesByIdStmt,params,m.handleMeasuringNodeResultHandler)

	if err != nil {
		return MeasuringNode{},err
	}
	return result.(MeasuringNode),nil
}


func (MeasuringNodeRepository) FetchAllMeasuringNodeUserRelations(nodeId int64 , userId int64) ([]string ,error){
	params := map[string]interface{}{
		"userId": userId,
		"nodeId": nodeId,
	}

	result, err := doReadTransaction(fetchMeasuringNodesUserRelations,params, func(result neo4j.Result) (res interface{}, err error) {
		var relations []string

		if result.Err() != nil{
			return relations, result.Err()
		}
		for ; result.Next() ; {
			relation := result.Record().GetByIndex(0).(string)
			relations = append(relations,relation)
		}
		return relations, nil
	})

	if err != nil {
		return make([]string, 0),err
	}
	return result.([]string),nil
}



func (m *MeasuringNodeRepository) FetchAllMeasuringNodes() ([]MeasuringNode,error) {
	params := map[string]interface{}{}

	results, err := doReadTransaction(fetchMeasuringNodesStmt, params,m.handleMeasuringNodeListResultHandler)

	if err != nil{
		return []MeasuringNode{}, err
	}
	return results.([]MeasuringNode),nil
}


func (m *MeasuringNodeRepository) CreateMeasuringNode(node MeasuringNode, userId int64) (MeasuringNode, error) {
	params := map[string]interface{}{
		"name": node.Name,
		"lat":  node.Lat,
		"lng":  node.Lng,
		"isPublic": node.IsPublic,
		"isOutdoors": node.IsOutdoors,
		"ownerId": userId,
	}

	result, err := doWriteTransaction(createStationStatement,params, m.handleMeasuringNodeResultHandler)

	if err != nil {
		return MeasuringNode{},err
	}
	return result.(MeasuringNode), nil
}

func (m *MeasuringNodeRepository) FetchNodesOwnedByUserId(userId int64) ([]MeasuringNode, error) {

	params := map[string]interface{}{
		"userId": userId,
	}

	results, err := doReadTransaction(fetchAllOwnedMeasuringNodesByUserIdStmt, params, m.handleMeasuringNodeListResultHandler)
	if err != nil{
		return []MeasuringNode{}, err
	}
	return results.([]MeasuringNode),nil
}

func (m *MeasuringNodeRepository) FetchAllPublicNodes() ([]MeasuringNode, error) {

	params := map[string]interface{}{}

	results, err := doReadTransaction(fetchAllPublicMeasuringNodesStmt, params, m.handleMeasuringNodeListResultHandler)
	if err != nil{
		return []MeasuringNode{}, err
	}
	return results.([]MeasuringNode),nil
}

func (m *MeasuringNodeRepository) FetchAllVisibleNodesByUserId(userId int64) ([]MeasuringNode, error) {

	params := map[string]interface{}{
		"userId": userId,
	}

	results, err := doReadTransaction(fetchAllVisibleMeasuringNodesByUserIdStmt, params, m.handleMeasuringNodeListResultHandler)
	if err != nil{
		return []MeasuringNode{}, err
	}
	return results.([]MeasuringNode),nil
}

func (m *MeasuringNodeRepository) CreateAuthorisationRelation(node MeasuringNode, user User) (err error){
	params := map[string]interface{}{
		"nodeId": node.Id,
		"userId": user.Id,
	}

	_,err = doWriteTransaction(createAuthorisationRelation,params, func(result neo4j.Result) (res interface{}, err error) {
		return nil, result.Err()
	})
	return
}
