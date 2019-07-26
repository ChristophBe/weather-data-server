package data

import (
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"io"
	"log"
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
	measuringModeQueryString 	 = "id(m), m.name, m.lat, m.lng, m.isPublic, m.isOutdoors"
	createStationStatement       = "MATCH (o:User) WHERE id(o) = $ownerId CREATE (o)-[:OWNER]->(m:MeasuringNode {name: $name, lat: $lat, lng: $lng, isPublic: $isPublic, isOutdoors: $isOutdoors}) RETURN " + measuringModeQueryString
	fetchMeasuringNodesStmt      = "MATCH (m:MeasuringNode) RETURN " + measuringModeQueryString
	fetchAllPublicMeasuringNodesStmt      = "MATCH (m:MeasuringNode) WHERE m.isPublic = true RETURN " + measuringModeQueryString
	fetchAllVisibaleMeasuringNodesByUserIdStmt      = "MATCH (m:MeasuringNode) WITH m OPTIONAL MATCH (m)<-[]-(u:User) WITH m, u WHERE m.isPublic OR id(u) = $userId RETURN " + measuringModeQueryString
	fetchMeasuringNodesByIdStmt      = "MATCH (m:MeasuringNode) WHERE id(m) = {nodeId} RETURN  "+ measuringModeQueryString
	fetchMeasuringNodesUserRelations =" MATCH (u:User)-[r]->(n:MeasuringNode) WHERE id(u) = {userId} and id(n) = {nodeId} return type(r)"
)


type MeasuringNodeRepository struct {

}
func (m *MeasuringNodeRepository) FetchAllMeasuringNodeById(nodeId int64 ) (MeasuringNode,error) {

	con := CreateConnection()
	defer con.Close()

	st:= prepareStatement(fetchMeasuringNodesByIdStmt,con)
	defer st.Close()


	rows := queryStatement(st,map[string]interface{}{"nodeId": nodeId})

	var measuringNode MeasuringNode

	row, _, err := rows.NextNeo()
	if err != nil && err != io.EOF {
		return measuringNode, err

	} else if err != io.EOF {
		measuringNode = m.parseMeasuringNodeFromRow(row)
	}
	return measuringNode, nil
}


func (MeasuringNodeRepository) FetchAllMeasuringNodeUserRelations(nodeId int64 , userId int64) []string {

	con := CreateConnection()
	defer con.Close()

	st:= prepareStatement(fetchMeasuringNodesUserRelations,con)
	defer st.Close()

	rows := queryStatement(st , map[string]interface{}{
		"userId": userId,
		"nodeId": nodeId,
	})
	var relations []string
	var err error
	err = nil

	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			panic(err)
		} else if err != io.EOF {
			relation := row[0].(string)
			relations = append(relations,relation)
		}
	}
	st.Close()

	return relations
}



func (m *MeasuringNodeRepository) FetchAllMeasuringNodes(con bolt.Conn) []MeasuringNode {

	st:= prepareStatement(fetchMeasuringNodesStmt,con)
	rows := queryStatement(st ,nil)
	var nodes []MeasuringNode
	var err error
	err = nil

	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			panic(err)
		} else if err != io.EOF {
			node := m.parseMeasuringNodeFromRow(row)
			nodes = append(nodes,node)
		}
	}


	st.Close()

	return nodes
}


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


func (m *MeasuringNodeRepository) CreateMeasuringNode(node MeasuringNode, userId int64) (MeasuringNode, error) {
	var (
		err      error
		driver   neo4j.Driver
		session  neo4j.Session
		result   neo4j.Result
		resultNode 	 MeasuringNode
	)

	driver, err = createDriver()
	if err != nil {
		return resultNode, err
	}
	defer driver.Close()

	session, err = createSession(driver,neo4j.AccessModeWrite)
	if err != nil {
		return resultNode, err
	}
	defer session.Close()

	log.Print(node)

	params := map[string]interface{}{
		"name": node.Name,
		"lat":  node.Lat,
		"lng":  node.Lng,
		"isPublic": node.IsPublic,
		"isOutdoors": node.IsOutdoors,
		"ownerId": userId,
	}

	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err = transaction.Run( createStationStatement, params)

		if err != nil {
			return nil, err
		}

		if result.Next() {
			resultNode = m.parseMeasuringNodeFromRow(result.Record().Values())
			return resultNode, nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return resultNode, err
	}

	return resultNode, nil
}

func (m *MeasuringNodeRepository) FetchAllPublicNodes() ([]MeasuringNode, error) {
	var (
		err      error
		driver   neo4j.Driver
		session  neo4j.Session
		result   neo4j.Result
		results	 []MeasuringNode
	)

	driver, err = createDriver()
	if err != nil {
		return results, err
	}
	defer driver.Close()

	session, err = createSession(driver,neo4j.AccessModeRead)
	if err != nil {
		return results, err
	}
	defer session.Close()



	_, err = session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err = transaction.Run(fetchAllPublicMeasuringNodesStmt, map[string]interface{}{})

		if err != nil {
			return nil, err
		}

		for ; result.Next(); {
			node := m.parseMeasuringNodeFromRow(result.Record().Values())

			results = append(results, node)

		}
		if result.Err() != nil{
			return nil, result.Err()
		}
		return results, nil

	})
	if err != nil {
		return results, err
	}

	return results, nil
}

func (m *MeasuringNodeRepository) FetchAllVisibleNodesByUserId(userId int64) ([]MeasuringNode, error) {
	var (
		err      error
		driver   neo4j.Driver
		session  neo4j.Session
		result   neo4j.Result
		results	 []MeasuringNode
	)

	driver, err = createDriver()
	if err != nil {
		return results, err
	}
	defer driver.Close()

	session, err = createSession(driver,neo4j.AccessModeRead)
	if err != nil {
		return results, err
	}
	defer session.Close()



	_, err = session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err = transaction.Run(fetchAllVisibaleMeasuringNodesByUserIdStmt, map[string]interface{}{
			"userId": userId,
		})

		if err != nil {
			return nil, err
		}

		for ; result.Next(); {
			node := m.parseMeasuringNodeFromRow(result.Record().Values())

			results = append(results, node)

		}
		if result.Err() != nil{
			return nil, result.Err()
		}
		return results, nil

	})
	if err != nil {
		return results, err
	}

	return results, nil
}