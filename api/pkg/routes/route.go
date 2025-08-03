// Package routes defines the API routes for the application
package routes

import (
	"github.com/Alfred-Onuada/todo-list-with-cassandra.git/pkg/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// Heatlh check
	router.GET("/", handlers.HeatlhCheck)

	router.GET("/todos", handlers.GetTodos)
	router.GET("/todos/:id", handlers.GetTodo)

	router.POST("/todos", handlers.CreateTodo)

	router.PATCH("/todos/:id", handlers.UpdateTodo)

	router.DELETE("/todos/:id", handlers.DeleteTodo)
}
