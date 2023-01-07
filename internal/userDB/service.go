package userDB

import (
	"context"
	"database/sql"
	"errors"
)

type UserDB struct {
	Name    string
	Conn    *sql.Conn
	Driver  string
	ConnStr string
}

type Storages map[string]*UserDB
type ConnStorage map[string]*Storages
type ConnStorageMain map[string]*UserDB

var st = Storages{"": nil}
var cst = ConnStorage{"": &st}
var cstMain = make(ConnStorageMain)

func Query(ctx context.Context, username, query string) (*sql.Rows, error) {
	if !CheckConn(username) {
		return nil, errors.New("Authentication failed")
	}

	return cstMain[username].Conn.QueryContext(ctx, query)
}

func NewConn(cred, driver string) (*sql.Conn, error) {
	db, err := sql.Open(driver, cred)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return db.Conn(ctx)
}
