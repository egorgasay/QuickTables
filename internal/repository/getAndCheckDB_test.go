package repository

import (
	"testing"
)

var queryCreateUserDBs = `
	CREATE TABLE userDBs 
	( ID INT, connStr TEXT, owner TEXT, driver TEXT, dbName TEXT);
INSERT INTO userDBs (ID, connStr, owner, driver, dbName) 
VALUES (1, 'test', 'test', 'test', 'test')
`
var queryDropUserDBs = `DROP TABLE IF EXISTS userDBs`

func TestStorage_CheckDB(t *testing.T) {
	cfg := &Config{"sqlite3", "test"}
	storage, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	storage.DB.Exec(queryCreateUserDBs)
	defer storage.DB.Exec(queryDropUserDBs)

	status := storage.CheckDB("admin")
	if status {
		t.Fatalf("The admin does not have a database")
	}

	status = storage.CheckDB("test")
	if !status {
		t.Fatalf("The user has a database")
	}
}

func TestStorage_GetDB(t *testing.T) {
	cfg := &Config{"sqlite3", "test"}
	storage, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	storage.DB.Exec(queryCreateUserDBs)
	defer storage.DB.Exec(queryDropUserDBs)

	connStr, driver, dbName := storage.GetDB("admin")
	if !(connStr == "" && driver == "" && dbName == "") {
		t.Fatalf("The admin user should not have a database")
	}

	connStr, driver, dbName = storage.GetDB("test")
	if connStr == "" && driver == "" && dbName == "" {
		t.Fatalf("The test user should have a database")
	}
}

func TestStorage_GetDBbyName(t *testing.T) {
	cfg := &Config{"sqlite3", "test"}
	storage, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	storage.DB.Exec(queryCreateUserDBs)
	defer storage.DB.Exec(queryDropUserDBs)

	connStr, driver, _ := storage.GetDBInfobyName("admin", "adqwd")
	if !(connStr == "" && driver == "") {
		t.Fatalf("The admin user should not have a database")
	}

	connStr, driver, _ = storage.GetDBInfobyName("test", "test")
	if connStr == "" && driver == "" {
		t.Fatalf("The test user should have a database")
	}
}
