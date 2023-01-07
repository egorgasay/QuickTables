package repository

import "testing"

var queryCreateHistoryOfQueries = `
CREATE TABLE historyOfQueries
(
    ID            INTEGER
        primary key autoincrement,
    Author        TEXT    not null
        references Users (name),
    DBName        TEXT    not null
        references userDBs (dbName),
    Query         TEXT    not null,
    Status        INTEGER not null,
    ExecutionTime REAL,
    Date          INTEGER not null
);

INSERT INTO historyOfQueries (Author, DBName, Query, Status, ExecutionTime, Date) 
VALUES ('test', 'test', 'SELECT 1', 1, 1.0, 1)
`

var queryDropHistoryOfQueries = `DROP TABLE IF EXISTS historyOfQueries`

func TestStorage_GetQueries(t *testing.T) {
	cfg := &Config{"sqlite3", "test"}
	storage, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	storage.DB.Exec(queryCreateHistoryOfQueries)
	defer storage.DB.Exec(queryDropHistoryOfQueries)

	ql, err := storage.GetQueries("admin", "aqwd")
	if err != nil {
		t.Fatal(err.Error())
	} else if len(ql) != 0 {
		t.Fatal("The admin user should not have list of queries")
	}

	_, err = storage.GetQueries("test", "test")
	if err != nil {
		t.Fatal(err.Error())
	} else if len(ql) != 0 {
		t.Fatal("The test user should have list of queries")
	}
}
