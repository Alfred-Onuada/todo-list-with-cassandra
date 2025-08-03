// Package handlers defines all API handlers
package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Priority string

var (
	CasualPriority    Priority = "casual"
	ImportantPriority Priority = "important"
)

type Todos struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Due       time.Time `json:"due"`
	Priority  Priority  `json:"priority"`
	Completed bool      `json:"completed"`
}

// holds my current todos
var todos []Todos = []Todos{}

func HeatlhCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf(
			"Welcome to my todos API, the goal of this is to play with Gin and Cassandra DB, the current time is %s",
			time.Now().String(),
		),
	})
}

func GetTodos(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Todos retrieved successfully",
		"data":    todos,
	})
}

func GetTodo(c *gin.Context) {
	id := c.Param("id")

	todo := Todos{}
	for _, t := range todos {
		if t.ID == id {
			todo = t
			break
		}
	}

	if todo.ID == "" {
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

func CreateTodo(c *gin.Context) {
	var todo Todos
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	todo.ID = uuid.NewString()
	todos = append(todos, todo)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Todo created successfully",
		"data":    todo,
	})
}

func DeleteTodo(c *gin.Context) {
	id := c.Param("id")

	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "Todo deleted successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "Todo not found",
	})
}

func UpdateTodo(c *gin.Context) {
	id := c.Param("id")

	var todo Todos
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	for i, t := range todos {
		if t.ID == id {
			todos[i] = todo
			c.JSON(http.StatusOK, gin.H{
				"message": "Todo updated successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "Todo not found",
	})
}
