package database

import (
	"de.christophb.wetter/config"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/pkg/errors"
	"log"
)

func createDriver() (driver neo4j.Driver, err error) {

	configuration ,err:= config.GetConfigManager().GetConfig()
	if err != nil{
		return
	}

	databaseHost := fmt.Sprintf("%s:%d",configuration.Neo4j.Host,configuration.Neo4j.Port)
	driver, err = neo4j.NewDriver(databaseHost, neo4j.BasicAuth(configuration.Neo4j.Username, configuration.Neo4j.Password, ""))
	return

}
func createSession(driver neo4j.Driver, mode neo4j.AccessMode) (neo4j.Session, error) {
	session, err := driver.Session(mode)
	if err != nil {
		return nil, err
	}

	return session, err
}

func doReadTransaction(statement string, params map[string]interface{}, resultHandler func(result neo4j.Result) (res interface{}, err error)) (res interface{}, err error) {

	var (
		driver  neo4j.Driver
		session neo4j.Session
		result  neo4j.Result
	)

	driver, err = createDriver()
	if err != nil {
		return nil, err
	}
	defer driver.Close()

	session, err = createSession(driver, neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return session.ReadTransaction(func(transaction neo4j.Transaction) (res interface{}, err error) {
		result, err = transaction.Run(statement, params)
		if err != nil {
			return nil, err
		}
		return resultHandler(result)
	})
}

func doWriteTransaction(statement string, params map[string]interface{}, resultHandler func(result neo4j.Result) (res interface{}, err error)) (res interface{}, err error) {

	var (
		driver  neo4j.Driver
		session neo4j.Session
		result  neo4j.Result
	)

	driver, err = createDriver()
	if err != nil {
		return nil, err
	}
	defer driver.Close()

	session, err = createSession(driver, neo4j.AccessModeWrite)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return session.WriteTransaction(func(transaction neo4j.Transaction) (res interface{}, err error) {
		result, err = transaction.Run(statement, params)
		if err != nil {
			return nil, err
		}
		return resultHandler(result)
	})
}

func saveNode(insertStmt string, updateStmt string, params map[string]interface{}, resultHandler func(result neo4j.Result) (interface{}, error)) (interface{}, error) {
	log.Print(params)
	if params["id"] != nil {
		return doWriteTransaction(updateStmt, params, resultHandler)
	} else {
		return doWriteTransaction(insertStmt, params, resultHandler)
	}
}

func parseSingleItemFromResult(parseSingeRecord func(record neo4j.Record) (res interface{}, err error)) func(result neo4j.Result) (res interface{}, err error) {
	return func(result neo4j.Result) (res interface{}, err error) {
		if result.Next() {
			res, err = parseSingeRecord(result.Record())
			return
		}
		err = result.Err()

		if err == nil {
			err = errors.New("node not found")
		}
		return
	}

}

func parseListFromResult(parseSingeRecord func(record neo4j.Record) (res interface{}, err error)) func(result neo4j.Result) (res interface{}, err error) {
	return func(result neo4j.Result) (res interface{}, err error) {
		var list []interface{}
		for ; result.Next(); {
			item, err := parseSingeRecord(result.Record())

			if err != nil {
				return nil, err
			}
			list = append(list, item)
		}
		if result.Err() != nil {
			return
		}
		res = list
		return
	}
}
