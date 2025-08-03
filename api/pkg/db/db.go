package db

import (
	"fmt"
	"log"
	"os"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cassandra" // registers Cassandra driver
	_ "github.com/golang-migrate/migrate/v4/source/file"        // registers file source
)

type DBConnection struct {
	session    *gocql.Session
	connection *gocql.ClusterConfig
}

var db *DBConnection

// Connect initializes the connection to the Cassandra database.
// It ensures the keyspace exists and runs any necessary migrations.
// It should be called once at the start of the application. It will panic if it fails to connect or run migrations.
// Example usage:
//
//	db.Connect()
func Connect() {
	keyspace := "todos_db"
	dbAddress := os.Getenv("CASSANDRA_HOST")

	fmt.Println("Running DB Migrations")
	ensureKeySpaceExists(dbAddress, keyspace)
	runMigrations(dbAddress, keyspace)

	fmt.Println("Connecting to the Cassandra DB at", dbAddress)

	// connect to DB
	db = &DBConnection{}
	db.connection = gocql.NewCluster(dbAddress)
	db.connection.Keyspace = keyspace
	db.connection.Consistency = gocql.Quorum

	var err error
	db.session, err = db.connection.CreateSession()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the Cassandra DB")
}

// ensureKeySpaceExists checks if the keyspace exists and creates it if it does not.
func ensureKeySpaceExists(dbAddress string, keyspace string) {
	connection := gocql.NewCluster(dbAddress)
	connection.Consistency = gocql.Quorum

	session, err := connection.CreateSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.Query(fmt.Sprintf(`
		CREATE KEYSPACE IF NOT EXISTS %s 
    WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'}
	`, keyspace)).Exec()
	if err != nil {
		panic(err)
	}
}

// runMigrations applies the database migrations using the migrate package.
// It connects to the Cassandra database and applies all migrations found in the specified source URL.
// The source URL points to the directory containing migration files.
// It will log an error and exit if it fails to create the migrate instance or apply migrations
func runMigrations(dbAddress string, keyspace string) {
	// file holding the migrations
	sourceURL := "file://./pkg/db/migrations"
	databaseURL := fmt.Sprintf("cassandra://%s:9042/%s", dbAddress, keyspace)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Apply all up migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migrations applied successfully")
}
