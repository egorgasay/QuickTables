package repository

import (
	"testing"
)

func TestStorage_SaveQuery(t *testing.T) {
	cfg := &Config{"sqlite3", "test"}
	storage, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	storage.DB.Exec(queryCreateHistoryOfQueries)
	defer storage.DB.Exec(queryDropHistoryOfQueries)

	err = storage.SaveQuery(1, "test", "test", "test", "1")
	if err != nil {
		t.Fatalf(err.Error())
	}
}
