package data

import (
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"io"
	"time"
)

type Measuring struct {
	Id          int64	  `json:"id"`
	TimeStamp   time.Time `json:"timestamp"`
	Pressure    float64   `json:"pressure"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
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

	fmt.Printf("CREATED ROWS: %d\n", numResult) // CREATED ROWS: 1

	// Closing the statment will also close the rows
	st.Close()
}

func FetchAllMeasuringByNodeId(con bolt.Conn, nodeId int64) []Measuring {

	st:= prepareStatement(fetchAllMeasuringForNodeStmt,con)
	rows := queryStatement(st,map[string]interface{}{"nodeId":nodeId})
	//var measuring []Measuring
	measuring := make([]Measuring,0)
	var err error
	err = nil

	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			panic(err)
		} else if err != io.EOF {

			// id(m), m.timeStamp, m.temperature, m.humidity, m.pressure
			node := Measuring{Id:row[0].(int64), TimeStamp:time.Unix(row[1].(int64),0),Temperature:row[2].(float64),Humidity:row[3].(float64),Pressure:row[3].(float64)}
			measuring = append(measuring,node)
		}
	}


	st.Close()

	return measuring
}
