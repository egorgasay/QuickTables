package userDB

import (
	"context"
	"database/sql"
)

type UserDB struct {
	Name    string
	Conn    *sql.Conn
	Driver  string
	ConnStr string
}

type DB struct {
	Name     string
	User     string
	Password string
}

type CustomDB struct {
	DB       DB
	Username string
	Port     string
}

// Add mutex here
type Storages map[string]*UserDB
type ConnStorage map[string]*Storages
type ConnStorageMain map[string]*UserDB

var st = Storages{"": nil}
var cst = ConnStorage{"": &st}
var cstMain = make(ConnStorageMain)

func NewConn(cred, driver string) (*sql.Conn, error) {
	db, err := sql.Open(driver, cred)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return db.Conn(ctx)
}
