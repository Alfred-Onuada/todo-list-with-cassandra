// Package handlers defines all API handlers
package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Alfred-Onuada/todo-list-with-cassandra.git/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HealthCheck is a simple handler to check if the API is running.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf(
			"Welcome to my todos API, the goal of this is to play with Gin and Cassandra DB, the current time is %s",
			time.Now().String(),
		),
	})
}

// GetTodos retrieves all todo items from the database.
// It returns a JSON response with the list of todos or an error message if something goes wrong
func GetTodos(c *gin.Context) {
	// fetch the todos
	todos, err := db.GetTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Todos retrieved successfully",
		"data":    todos,
	})
}

// GetTodo retrieves a single todo item by its ID.
// It returns a JSON response with the todo item or an error message if something goes wrong
func GetTodo(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID",
			"error":   "ID param can not be empty",
		})
	}

	todo, err := db.GetTodoByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
		return
	}

	if todo == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Todo retrieved successfully",
		"data":    todo,
	})
}

// CreateTodo creates a new todo item in the database.
// It expects a JSON request body with the todo details.
// It returns a JSON response with the created todo item or an error message if something goes wrong
func CreateTodo(c *gin.Context) {
	var todo db.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}
	// generate a new UUID for the todo
	todo.ID = uuid.NewString()
	todo.Completed = false // default to not completed

	if err := db.CreateTodo(todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Todo created successfully",
		"data":    todo,
	})
}

// DeleteTodo deletes a todo item by its ID.
// It returns a JSON response indicating success or failure.
// If the ID is invalid or the todo item does not exist, it returns an appropriate error message.
func DeleteTodo(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID",
			"error":   "ID param can not be empty",
		})
		return
	}

	if err := db.DeleteTodo(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Todo deleted successfully",
	})
}

// UpdateTodo updates an existing todo item in the database.
// It expects a JSON request body with the fields to update.
// It returns a JSON response indicating success or failure.
func UpdateTodo(c *gin.Context) {
	id := c.Param("id")

	var todo db.UpdateTodoType
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if err := db.UpdateTodo(id, todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Todo updated successfully",
	})
}
