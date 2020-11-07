// +build !change

package models

import "testing"

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
