package repository

import "testing"

func TestStorage_DeleteDB(t *testing.T) {
	cfg := &Config{"sqlite3", "test"}
	storage, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	storage.DB.Exec(queryCreateUserDBs)
	defer storage.DB.Exec(queryDropUserDBs)

	err = storage.DeleteDB("test", "test")
	if err != nil {
		t.Fatal(err)
	}
}
