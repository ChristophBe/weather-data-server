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
}

const(
	createStationStatement       = "CREATE (n:MeasuringNode {name: {name}, lat: {lat}, lng: {lng}})"
	fetchMeasuringNodesStmt      = "MATCH (m:MeasuringNode) RETURN id(m), m.name, m.lat, m.lng "
)


func CreateMeasuringNode( con bolt.Conn, node MeasuringNode)  {


	st := prepareStatement(createStationStatement,con)

	//{{Name: {Name}, Lat: {Lat}, lng: {lng}

	result, err := st.ExecNeo(map[string]interface{}{
		"name": node.Name,
		"lat":  node.Lat,
		"lng":  node.Lng})
	handleError(err)

	numResult, err := result.RowsAffected()
	handleError(err)

	fmt.Printf("CREATED ROWS: %d\n", numResult) // CREATED ROWS: 1

	// Closing the statment will also close the rows
	st.Close()
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
			node := MeasuringNode{Id:row[0].(int64), Name:row[1].(string),Lat:row[2].(float64),Lng:row[3].(float64)}
			nodes = append(nodes,node)
		}
	}


	st.Close()

	return nodes
}

