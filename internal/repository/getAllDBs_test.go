package repository

import (
	"reflect"
	"testing"
)

func TestStorage_GetAllDBs(t *testing.T) {
	cfg := &Config{"sqlite3", "test"}
	storage, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	storage.DB.Exec(queryCreateUserDBs)
	defer storage.DB.Exec(queryDropUserDBs)

	dbs := storage.GetAllDBs("admin")
	if !reflect.DeepEqual(dbs, make([][]string, 0, 5)) {
		t.Fatalf("The slice must be empty: %v", err)
	}

	dbs = storage.GetAllDBs("test")
	if reflect.DeepEqual(dbs, make([][]string, 0, 5)) {
		t.Fatalf("The slice should not be empty: %v", err)
	}
}
