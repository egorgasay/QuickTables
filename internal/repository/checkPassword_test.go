package repository

import "testing"

var queryCreateUsers = `
	CREATE TABLE Users 
	(
		Name TEXT not null unique,
		ID INTEGER not null
		primary key autoincrement
		unique,
		Role     TEXT,
		Password TEXT,
		Nickname TEXT
	);

	INSERT INTO Users (Name, Role, Password)
	VALUES ('admin', 'admin', 'admin')
`

var queryDropUsers = `DROP TABLE Users`

func TestStorage_CheckPassword(t *testing.T) {
	cfg := &Config{"sqlite3", "test"}
	storage, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	storage.DB.Exec(queryCreateUsers)
	defer storage.DB.Exec(queryDropUsers)

	ok := storage.CheckPassword("test", "qwe")
	if ok {
		t.Fatalf("Failed to check password: %v", err)
	}

	ok = storage.CheckPassword("admin", "admin")
	if !ok {
		t.Fatalf("Failed to check password: %v", err)
	}
}
