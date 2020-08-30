package database

import (
	"errors"
	"github.com/ChristophBe/weather-data-server/data/models"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"time"
)

type measuringRepositoryImpl struct{}

func (measuringRepositoryImpl) measuringResultHandler(record neo4j.Record) (interface{}, error) {

	nodeData, ok := record.Get("m")
	if !ok {
		err := errors.New("can not parse measuring form record")
		return nil, err
	}
	node := nodeData.(neo4j.Node)
	props := node.Props()
	measurement := models.Measurement{
		Id:          node.Id(),
		TimeStamp:   parseTimeProp(props["timeStamp"], time.Unix(0, 0)),
		Temperature: parseFloatProp(props["pressure"], .0),
		Humidity:    parseFloatProp(props["temperature"], .0),
		Pressure:    parseFloatProp(props["humidity"], .0),
	}

	return measurement, nil
}

func (measuringRepositoryImpl) castResultList(input interface{}) (measurements []models.Measurement) {
	if input == nil {
		measurements = make([]models.Measurement, 0)
		return
	}
	items := input.([]interface{})
	measurements = make([]models.Measurement, len(items))

	for key, item := range items {
		measurements[key] = item.(models.Measurement)
	}
	return
}

func (r measuringRepositoryImpl) CreateMeasurement(stationId int64, measurement models.Measurement) (savedMeasurement models.Measurement, err error) {
	params := map[string]interface{}{
		"timeStamp":   measurement.TimeStamp.Unix(),
		"pressure":    measurement.Pressure,
		"temperature": measurement.Temperature,
		"humidity":    measurement.Humidity,
		"stationId":   stationId,
	}

	stmt := "MATCH (n:MeasuringNode)  WHERE Id(n) = $stationId CREATE (n)<-[:MEASUREMENT_FOR]-(m:Measurement {timeStamp: $timeStamp, pressure: $pressure,temperature: $temperature,humidity: $humidity}) RETURN m"
	res, err := doWriteTransaction(stmt, params, parseSingleItemFromResult(r.measuringResultHandler))

	if err != nil {
		return
	}

	savedMeasurement = res.(models.Measurement)
	return
}

func (r measuringRepositoryImpl) FetchAllMeasuringsByNodeId(nodeId int64) (measurements []models.Measurement, err error) {

	params := map[string]interface{}{"nodeId": nodeId}

	stmt:= "MATCH (m:Measurement)-[:MEASUREMENT_FOR]->(n:MeasuringNode) WHERE id(n) = $nodeId RETURN m ORDER BY m.timeStamp DESC"
	results, err := doReadTransaction(stmt, params, parseListFromResult(r.measuringResultHandler))

	measurements = r.castResultList(results)
	return
}

func (r measuringRepositoryImpl) FetchLastMeasuringsByNodeId(nodeId int64, hours int64) (measurements []models.Measurement, err error) {

	minTime := time.Now().Add(time.Duration(-hours) * time.Hour)
	params := map[string]interface{}{"nodeId": nodeId, "minTime": minTime.Unix()}

	stmt:= "MATCH (m:Measurement)-[:MEASUREMENT_FOR]->(n:MeasuringNode) WHERE id(n) = $nodeId and m.timeStamp > $minTime RETURN m ORDER BY m.timeStamp DESC"

	results, err := doReadTransaction(stmt, params, parseListFromResult(r.measuringResultHandler))

	measurements = r.castResultList(results)
	return
}
