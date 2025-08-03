package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
)

type Priority uint8

var (
	PriorityImportant Priority = 0
	PriorityCasual    Priority = 1
)

type Todo struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Due       time.Time `json:"due"`
	Priority  Priority  `json:"priority"`
	Completed bool      `json:"completed"`
}

type UpdateTodoType struct {
	Content   *string    `json:"content"`
	Due       *time.Time `json:"due"`
	Priority  *Priority  `json:"priority"`
	Completed *bool      `json:"completed"`
}

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

func CreateTodo(todo Todo) error {
	// insert the todo into the db
	return db.session.Query(`
		INSERT INTO todos (id, content, due, priority, completed)
		VALUES (?, ?, ?, ?, ?)
	`, todo.ID, todo.Content, todo.Due, todo.Priority, todo.Completed).Exec()
}

func DeleteTodo(id string) error {
	// delete the todo with the given id
	return db.session.Query(`
		DELETE FROM todos WHERE id = ?
	`, id).Exec()
}

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
