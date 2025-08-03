package db

import (
	"fmt"
	"os"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

type DBConnection struct {
	session    *gocql.Session
	connection *gocql.ClusterConfig
}

var db *DBConnection

func Connect() {
	fmt.Println("Connecting to the Cassandra DB")

	db = &DBConnection{}
	db.connection = gocql.NewCluster(os.Getenv("CASSANDRA_HOST"))
	db.connection.Keyspace = "todo_db"
	db.connection.Consistency = gocql.Quorum

	// connect to the db
	var err error
	db.session, err = db.connection.CreateSession()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the Cassandra DB")
}

func GetDB() *gocql.Session {
	if db == nil {
		Connect()
	}

	return db.session
}
