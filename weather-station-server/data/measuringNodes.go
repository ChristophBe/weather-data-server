package data

import (
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"io"
)

type MeasuringNode struct {
	Id   int64	  	`json:"id"`
	Name string 	`json:"name"`
	Lat  float64 	`json:"lat"`
	Lng  float64 	`json:"lng"`
	IsPublic bool 	`json:"is_public"`
}

const(
	createStationStatement       = "CREATE (n:MeasuringNode {name: {name}, lat: {lat}, lng: {lng}, isPublic: {isPublic}})"
	fetchMeasuringNodesStmt      = "MATCH (m:MeasuringNode) RETURN id(m), m.name, m.lat, m.lng, m.isPublic "
	fetchMeasuringNodesByIdStmt      = "MATCH (m:MeasuringNode) WHERE id(m) = {nodeId} RETURN id(m), m.name, m.lat, m.lng, m.isPublic "
	fetchMeasuringNodesUserRelations ="MATCH (u:User)-[r]->(n:MeasuringNode) WHERE id(u) = {userId} and id(n) = {nodeId} return type(r)"
)


func CreateMeasuringNode( con bolt.Conn, node MeasuringNode)  {

	st := prepareStatement(createStationStatement,con)

	//{{Name: {Name}, Lat: {Lat}, lng: {lng}

	result, err := st.ExecNeo(map[string]interface{}{
		"name": node.Name,
		"lat":  node.Lat,
		"lng":  node.Lng,
		"isPublic": node.IsPublic,
	})
	handleError(err)

	numResult, err := result.RowsAffected()
	handleError(err)

	fmt.Printf("CREATED ROWS: %d\n", numResult) // CREATED ROWS: 1

	// Closing the statment will also close the rows
	st.Close()
}

func FetchAllMeasuringNodeById(nodeId int64 ) (MeasuringNode,error) {

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
		measuringNode = parseMeasuringNodeFromRow(row)
	}
	return measuringNode, nil
}


func FetchAllMeasuringNodeUserRelations(nodeId int64 , userId int64) []string {

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



func FetchAllMeasuringNodes(con bolt.Conn) []MeasuringNode {

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
			node := parseMeasuringNodeFromRow(row)
			nodes = append(nodes,node)
		}
	}


	st.Close()

	return nodes
}


func parseMeasuringNodeFromRow(row []interface{}) MeasuringNode {
	node := MeasuringNode{
		Id:row[0].(int64),
		Name:row[1].(string),
		Lat:row[2].(float64),
		Lng:row[3].(float64),
		IsPublic:row[4].(bool),
	}
	return node
}

