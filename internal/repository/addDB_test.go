package repository

import "testing"

func TestStorage_AddDB(t *testing.T) {
	cfg := &Config{"sqlite3", "test"}
	storage, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	queryCreate := `
	CREATE TABLE userDBs 
	( ID INT, connStr TEXT, owner TEXT, driver TEXT, dbName TEXT);
`
	queryDrop := `DROP TABLE userDBs`

	storage.DB.Exec(queryCreate)
	defer storage.DB.Exec(queryDrop)

	err = storage.AddDB("test", "test", "test", "test")
	if err != nil {
		t.Fatalf("Failed to add database: %v", err)
	}
}
