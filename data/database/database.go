package database

import (
	"errors"
	"fmt"
	"github.com/ChristophBe/weather-data-server/config"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func createDriver() (driver neo4j.Driver, err error) {

	configuration, err := config.GetConfigManager().GetConfig()
	if err != nil {
		return
	}


	configForNeo4j40 := func(conf *neo4j.Config) { conf.Encrypted = false }
	databaseHost := fmt.Sprintf("%s:%d", configuration.Neo4j.Host, configuration.Neo4j.Port)
	driver, err = neo4j.NewDriver(databaseHost, neo4j.BasicAuth(configuration.Neo4j.Username, configuration.Neo4j.Password, ""),configForNeo4j40)
	if err != nil {
		err = errors.Errorf("can not create neo4j driver; cause: %w",err)
	}
	return

}
func createSession(driver neo4j.Driver, mode neo4j.AccessMode) (neo4j.Session, error) {
	session, err := driver.Session(mode)
	if err != nil {
		return nil, errors.Errorf("can not start neo4j session; cause: %w",err)
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
	defer driver.Close()
	if err != nil {
		return nil, errors.Errorf("can not do read transaction; cause: %w",err)
	}


	session, err = createSession(driver, neo4j.AccessModeRead)
	defer session.Close()

	if err != nil {
		 return nil, errors.Errorf("can not do read transaction; cause: %w",err)
	}

	return session.ReadTransaction(func(transaction neo4j.Transaction) (res interface{}, err error) {
		result, err = transaction.Run(statement, params)
		if err != nil {
			return nil,  errors.Errorf("can not run read transaction; cause: %w",err)
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
	defer driver.Close()
	if err != nil {
		return nil, errors.Errorf("can not create driver; cause: %w",err)
	}


	session, err = createSession(driver, neo4j.AccessModeWrite)
	defer session.Close()
	if err != nil {
		return nil, errors.Errorf("can not create session; cause: %w",err)
	}


	return session.WriteTransaction(func(transaction neo4j.Transaction) (res interface{}, err error) {
		result, err = transaction.Run(statement, params)
		if err != nil {
			return nil, errors.Errorf("can not run write transaction; cause: %w",err)
		}
		return resultHandler(result)
	})
}

func saveNode(insertStmt string, updateStmt string, params map[string]interface{}, resultHandler func(result neo4j.Result) (interface{}, error)) (interface{}, error) {
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
			if err != nil {
				err=  errors.Errorf("can not run parse result; cause: %w",err)
			}
			return
		}
		err = result.Err()

		if err == nil {
			err = errors.Errorf("node not found; cause: %w", err)
		}
		return
	}

}

func parseListFromResult(parseSingeRecord func(record neo4j.Record) (res interface{}, err error)) func(result neo4j.Result) (res interface{}, err error) {
	return func(result neo4j.Result) (res interface{}, err error) {
		var list []interface{}
		for result.Next() {
			item, err := parseSingeRecord(result.Record())

			if err != nil {
				return nil,  errors.Errorf("can not run parse list form record; cause: %w",err)
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
