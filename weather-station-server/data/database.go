package data

import (
	"../configs"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// Constants to be used throughout the example
const (
	URI = "bolt://neo4j:" + configs.NEO4J_PASSWORD + "@" + configs.NEO4J_HOST
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


func createDriver()(neo4j.Driver, error){

	driver, err := neo4j.NewDriver("bolt://" + configs.NEO4J_HOST, neo4j.BasicAuth("neo4j", configs.NEO4J_PASSWORD, ""))
	if err != nil {
		return nil, err
	}

	return driver,err


}
func createSession( driver neo4j.Driver ,mode neo4j.AccessMode)(neo4j.Session, error){
	session, err := driver.Session(mode)
	if err != nil {
		return nil, err
	}

	return session,err
}