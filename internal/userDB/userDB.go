package userDB

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type UserDB struct {
	Username string
	Conn     *sql.Conn
	Driver   string
	ConnStr  string
}

type ConnStorage map[string]*UserDB

var cst = make(ConnStorage)

func Query(ctx context.Context, username, query string) (*sql.Rows, error) {
	if !CheckConn(username) {
		return nil, errors.New("у юзера нет бд")
	}

	return cst[username].Conn.QueryContext(ctx, query)
}

func NewConn(cred, driver string) (*sql.Conn, error) {
	db, err := sql.Open(driver, cred)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return db.Conn(ctx)
}

func CheckConn(username string) bool {
	_, ok := cst[username]
	return ok
}

func New() map[string]*UserDB {
	mp := make(map[string]*UserDB)
	return mp
}

func RecordConnection(connStr string, username string, driver string) error {
	cn, err := NewConn(connStr, driver)
	if err != nil {
		return err
	}

	cst[username] = &UserDB{
		Username: username,
		Conn:     cn,
		ConnStr:  connStr,
		Driver:   driver,
	}
	return nil
}

//func (ud UserDB) Remove(id int64) error {
//	_, err := ud.DB.Exec("DELETE FROM userDBs WHERE id = ?", id)
//	return err
//}
