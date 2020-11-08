// +build !change

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// min coverage: app,client,models,utils 90%

func mockTodo() *Todo {
	return &Todo{
		ID:       1,
		Title:    "t",
		Content:  "content",
		Finished: false,
	}
}

func TestMarkFinished(t *testing.T) {
	todo := mockTodo()

	todo.Finished = false
	todo.MarkFinished()
	assert.Equal(
		t,
		true,
		todo.Finished,
		"MarkFinished didn't save correct result to todo val",
	)
}

func TestMarkUnfinished(t *testing.T) {
	todo := mockTodo()

	todo.Finished = true
	todo.MarkUnfinished()
	assert.Equal(
		t,
		false,
		todo.Finished,
		"MarkUnfinished didn't save correct result to todo val",
	)
}
