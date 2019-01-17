package data

import (
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

// Constants to be used throughout the example
const (
	URI                          = "bolt://neo4j:password123@localhost:7687"
	createMeasurementStatement   = "MATCH (n:MeasuringNode)  WHERE Id(n) = {stationId} CREATE (n)<-[:MEASUREMENT_FOR]-(m:Measuring {timeStamp: {timeStamp}, pressure: {pressure},temperature: {temperature},humidity: {humidity}})"
	createStationStatement       = "CREATE (n:MeasuringNode {name: {name}, lat: {lat}, lng: {lng}})"
	fetchMeasuringNodesStmt      = "MATCH (m:MeasuringNode) RETURN id(m), m.name, m.lat, m.lng "
	fetchAllMeasuringForNodeStmt = "MATCH (m:Measuring)-[:MEASUREMENT_FOR]->(n:MeasuringNode) WHERE id(n) = {nodeId} RETURN id(m), m.timeStamp, m.temperature, m.humidity, m.pressure"

)


func CreateConnection() bolt.Conn {
	driver := bolt.NewDriver()
	con, err := driver.OpenNeo(URI)
	handleError(err)
	return con
}

// Here we prepare a new statement. This gives us the flexibility to
// cancel that statement without any request sent to Neo
func prepareStatement(query string, con bolt.Conn) bolt.Stmt {
	st, err := con.PrepareNeo(query)
	handleError(err)
	return st
}




func queryStatement(st bolt.Stmt,params map[string]interface{}) bolt.Rows {
	// Even once I get the rows, if I do not consume them and close the
	// rows, Neo will discard and not send the data
	rows, err := st.QueryNeo(params)
	handleError(err)
	return rows
}




// Here we create a simple function that will take care of errors, helping with some code clean up
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}