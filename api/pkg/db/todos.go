package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
)

// Priority represents the priority of a todo item.
// It can be either Important (0) or Casual (1).
// It is defined as a uint8 type for efficient storage.
type Priority uint8

var (
	PriorityImportant Priority = 0
	PriorityCasual    Priority = 1
)

// Todo represents a todo item in the database.
// It includes fields for ID, content, due date, priority, and completion status.
type Todo struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Due       time.Time `json:"due"`
	Priority  Priority  `json:"priority"`
	Completed bool      `json:"completed"`
}

// UpdateTodoType is used for updating a todo item.
type UpdateTodoType struct {
	Content   *string    `json:"content"`
	Due       *time.Time `json:"due"`
	Priority  *Priority  `json:"priority"`
	Completed *bool      `json:"completed"`
}

// GetTodos retrieves all todo items from the database.
// It returns a slice of Todo items or an error if the query fails.
func GetTodos() (*[]Todo, error) {
	// create a slice to hold the
	result := []Todo{}

	// a variable to hold each row
	var todo Todo

	// query the db and get the itr
	itr := db.session.Query("SELECT id, content, due, priority, completed FROM todos").Iter()

	for itr.Scan(&todo.ID, &todo.Content, &todo.Due, &todo.Priority, &todo.Completed) {
		result = append(result, todo)
	}

	// check if there was an error
	if err := itr.Close(); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTodoByID retrieves a todo item by its ID from the database.
// It returns the Todo item if found, or nil if not found, along with an error if any occurs.
func GetTodoByID(id string) (*Todo, error) {
	var todo Todo

	err := db.session.Query(`
		SELECT id, content, due, priority, completed 
		FROM todos
		WHERE id = ?
	`, id).Scan(&todo.ID, &todo.Content, &todo.Due, &todo.Priority, &todo.Completed)
	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil // return nil if todo not found
		}
		return nil, err // return error if any other error occurs
	}

	return &todo, nil
}

// CreateTodo inserts a new todo item into the database.
// It takes a Todo struct as input and returns an error if the insertion fails.
func CreateTodo(todo Todo) error {
	// insert the todo into the db
	return db.session.Query(`
		INSERT INTO todos (id, content, due, priority, completed)
		VALUES (?, ?, ?, ?, ?)
	`, todo.ID, todo.Content, todo.Due, todo.Priority, todo.Completed).Exec()
}

// DeleteTodo removes a todo item from the database by its ID.
// It returns an error if the deletion fails.
func DeleteTodo(id string) error {
	// delete the todo with the given id
	return db.session.Query(`
		DELETE FROM todos WHERE id = ?
	`, id).Exec()
}

// UpdateTodo updates an existing todo item in the database.
// It takes an ID and an UpdateTodoType struct containing the fields to update.
// It returns an error if the update fails or if no fields are provided to update.
func UpdateTodo(id string, todo UpdateTodoType) error {
	assignments := []string{}
	values := []any{}

	// check which fields are provided and build the query accordingly
	if todo.Content != nil {
		assignments = append(assignments, "content = ?")
		values = append(values, *todo.Content)
	}
	if todo.Due != nil {
		assignments = append(assignments, "due = ?")
		values = append(values, *todo.Due)
	}
	if todo.Priority != nil {
		assignments = append(assignments, "priority = ?")
		values = append(values, *todo.Priority)
	}
	if todo.Completed != nil {
		assignments = append(assignments, "completed = ?")
		values = append(values, *todo.Completed)
	}

	if len(assignments) == 0 {
		return fmt.Errorf("no fields to update")
	}

	values = append(values, id) // add the id to the end of the values slice

	query := fmt.Sprintf(`
		UPDATE todos 
		SET %s
		WHERE id = ?
	`, strings.Join(assignments, ", "))

	// update the todo with the given id
	return db.session.Query(query, values...).Exec()
}
