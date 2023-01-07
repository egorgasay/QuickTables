package repository

import "testing"

func TestStorage_ChangeNick(t *testing.T) {
	cfg := &Config{"sqlite3", "test"}
	storage, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	storage.DB.Exec(queryCreateUsers)
	defer storage.DB.Exec(queryDropUsers)

	err = storage.ChangeNick("admin", "test")
	if err != nil {
		t.Fatal(err)
	}
}
