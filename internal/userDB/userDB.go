package userDB

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"log"
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
		return nil, errors.New("у юзера нет бд")
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

func CheckConn(username string) bool {
	_, ok := cstMain[username]
	return ok
}

//func New() map[string]*UserDB {
//	mp := make(map[string]*UserDB)
//	return mp
//}

func RecordConnection(name, connStr, username, driver string) error {
	cn, err := NewConn(connStr, driver)
	if err != nil {
		log.Println(err)
		return err
	}

	udb := &UserDB{
		Name:    name,
		Conn:    cn,
		ConnStr: connStr,
		Driver:  driver,
	}

	if isNil := cst[username]; isNil == nil {
		cst[username] = &Storages{username: udb}
	} else {
		(*cst[username])[udb.Name] = udb
	}
	cstMain[username] = udb

	return nil
}

func GetDbNameAndVendor(username string) (name string, vendor string) {
	return cstMain[username].Name, cstMain[username].Driver
}

// сервис (мапы ,структуры)
// абота с глоб бд, сохранение и тд
// раб с лок бд

func SetMainDbByName(name, username, connStr, driver string) error {
	if dbCached, ok := (*cst[username])[name]; ok {
		cstMain[username] = dbCached
		return nil
	}

	return RecordConnection(name, connStr, username, driver)
}

//func GetDbByName(username, dbName string) (connStr, driver string) {
//
//}

//func (ud UserDB) Remove(id int64) error {
//	_, err := ud.DB.Exec("DELETE FROM userDBs WHERE id = ?", id)
//	return err
//}
