// +build !change

package models

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// min coverage: app,client,models,utils 90%

func TestNewInMemoryStorage(t *testing.T) {
	db := NewInMemoryStorage()
	if db == nil {
		t.Fatalf("NewInMemoryStorage returns nil db")
	}

	if db.todos == nil {
		t.Fatalf("NewInMemoryStorage returns db with nil todos field")
	}
}

func TestAddTodo(t *testing.T) {
	db := NewInMemoryStorage()

	title, content := "title", "content"

	todo, err := db.AddTodo(title, content)

	assert.NoError(t, err, "AddTodo returns no nil err")

	assert.Equal(t, 1, len(db.todos), "AddTodo adds too more fields to store")

	savedTodo, ok := db.todos[todo.ID]
	assert.Equal(
		t,
		true,
		ok,
		`AddTodo returns todo with ID that does not exist in store`,
	)

	assert.Equal(
		t,
		todo,
		savedTodo,
		`AddTodo returns incorrect data. savedTodo != todo`,
	)

	assert.Equal(
		t,
		todo.Content,
		content,
		`AddTodo adds incorrect content to store`,
	)

	assert.Equal(
		t,
		todo.Title,
		title,
		`AddTodo adds incorrect title to store`,
	)

	assert.Equal(
		t,
		todo.Finished,
		false,
		`AddTodo adds todo with Finished status to store`,
	)
}

func TestGetTodo(t *testing.T) {
	db := NewInMemoryStorage()
	incorrectID := -1 // after create store is empty

	_, err := db.GetTodo(ID(incorrectID))
	assert.Error(t, err, "GetTodo returns not nil result with empty state")

	todo, _ := db.AddTodo("title", "content")
	savedTodo, err := db.GetTodo(todo.ID)
	assert.NoError(t, err, "GetTodo returns not nil result with empty state")

	assert.Equal(
		t,
		todo,
		savedTodo,
		"GetTodo returns val != saved returns from AddTodo val",
	)
}

func TestGetAll(t *testing.T) {
	const N = 5

	db := NewInMemoryStorage()

	localStore := make(map[*Todo]bool)

	for i := 0; i < N; i++ {
		todo, _ := db.AddTodo(
			"title"+strconv.Itoa(i),
			"content"+strconv.Itoa(i),
		)
		localStore[todo] = true
	}

	todos, err := db.GetAll()
	assert.NoError(t, err, "GetAll returns error")

	assert.NotNil(t, todos, "todos is nil")

	assert.Equal(t, len(localStore), len(todos), "GetAll returns error")

	for _, todo := range todos {
		_, ok := localStore[todo]
		assert.Equal(t, ok, true, "return val does not exist in expected store")
	}
}
