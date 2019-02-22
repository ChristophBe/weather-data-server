package data

import (
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"io"
	"time"
)

const(
	createMeasurementStatement   = "MATCH (n:MeasuringNode)  WHERE Id(n) = {stationId} CREATE (n)<-[:MEASUREMENT_FOR]-(m:Measuring {timeStamp: {timeStamp}, pressure: {pressure},temperature: {temperature},humidity: {humidity}})"
	fetchAllMeasuringForNodeStmt = "MATCH (m:Measuring)-[:MEASUREMENT_FOR]->(n:MeasuringNode) WHERE id(n) = {nodeId} RETURN id(m), m.timeStamp, m.temperature, m.humidity, m.pressure"
	fetchMeasuringForNodeAfterTimestampStmt = "MATCH (m:Measuring)-[:MEASUREMENT_FOR]->(n:MeasuringNode) WHERE id(n) = {nodeId} and m.timeStamp > {minTime} RETURN id(m), m.timeStamp, m.temperature, m.humidity, m.pressure ORDER BY m.timeStamp DESC"
)

type Measuring struct {
	Id          int64	  `json:"id"`
	TimeStamp   time.Time `json:"timestamp"`
	Pressure    float64   `json:"pressure"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
}


type Measuerings struct{

}

func CreateMeasurement( con bolt.Conn, stationId int64, measurement Measuring) {

	st := prepareStatement(createMeasurementStatement, con)

	result, err := st.ExecNeo(map[string]interface{}{
		"timeStamp":   measurement.TimeStamp.Unix(),
		"pressure":    measurement.Pressure,
		"temperature": measurement.Temperature,
		"humidity":    measurement.Humidity,
		"stationId":   stationId})
	handleError(err)

	numResult, err := result.RowsAffected()
	handleError(err)

	fmt.Printf("CREATED ROWS: %d\n", numResult);

	st.Close()
}



func FetchAllMeasuringByNodeId(con bolt.Conn, nodeId int64) []Measuring {

	st:= prepareStatement(fetchAllMeasuringForNodeStmt,con)
	rows := queryStatement(st,map[string]interface{}{"nodeId":nodeId})
	//var measurings []Measuring
	measurings := parseMeasuringFormRows(rows, st)

	st.Close()
	return measurings
}


func FetchLastMeasuringsByNodeId(con bolt.Conn, nodeId int64, hours int64) []Measuring {

	minTime := time.Now().Add(time.Duration(-hours)*time.Hour)
	st:= prepareStatement(fetchMeasuringForNodeAfterTimestampStmt,con)
	rows := queryStatement(st,map[string]interface{}{"nodeId":nodeId, "minTime":minTime.Unix()})
	//var measurings []Measuring
	measurings := parseMeasuringFormRows(rows, st)

	st.Close()
	return measurings
}

func parseMeasuringFormRows(rows bolt.Rows, st bolt.Stmt) []Measuring {
	measuring := make([]Measuring, 0)
	var err error
	err = nil
	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			panic(err)
		} else if err != io.EOF {

			// id(m), m.timeStamp, m.temperature, m.humidity, m.pressure
			node := paresMeasuringFormLine(row)
			measuring = append(measuring, node)
		}
	}
	return measuring
}

func paresMeasuringFormLine(row []interface{}) Measuring {
	node := Measuring{Id: row[0].(int64), TimeStamp: time.Unix(row[1].(int64), 0), Temperature: row[2].(float64), Humidity: row[3].(float64), Pressure: row[4].(float64)}
	return node
}
