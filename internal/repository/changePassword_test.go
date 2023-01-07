package repository

import "testing"

func TestStorage_ChangePassword(t *testing.T) {
	cfg := &Config{"sqlite3", "test"}
	storage, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	storage.DB.Exec(queryCreateUsers)
	defer storage.DB.Exec(queryDropUsers)

	err = storage.ChangePassword("admin", "admin", "admin2")
	if err != nil {
		t.Fatal(err)
	}

	err = storage.ChangePassword("qw", "admin", "admin2")
	if err == nil {
		t.Fatal("There must be an error when changing the password of a non-existent user")
	}
}
