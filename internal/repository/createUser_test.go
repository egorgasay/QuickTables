package repository

import "testing"

func TestStorage_CreateUser(t *testing.T) {
	cfg := &Config{"sqlite3", "test"}
	storage, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	queryCreate := `
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
	queryDrop := `DROP TABLE Users`

	storage.DB.Exec(queryCreate)
	defer storage.DB.Exec(queryDrop)

	err = storage.CreateUser("test", "qwe")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	err = storage.CreateUser("admin", "admin")
	if err == nil {
		t.Fatalf("Unique constraint failed: %v", err)
	}
}
