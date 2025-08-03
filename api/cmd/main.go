package main

import (
	"fmt"
	"os"

	"github.com/Alfred-Onuada/todo-list-with-cassandra.git/pkg/db"
	"github.com/Alfred-Onuada/todo-list-with-cassandra.git/pkg/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// connect to the db
	db.Connect()

	// define routes
	routes.RegisterRoutes(router)

	// setup the address
	port := os.Getenv("API_PORT")
	address := fmt.Sprintf(":%s", port)

	router.Run(address)
}
