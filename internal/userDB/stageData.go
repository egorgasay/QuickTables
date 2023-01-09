package userDB

import (
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/sijms/go-ora/v2"
	"log"
	"time"
)

func CheckConn(username string) bool {
	_, ok := cstMain[username]
	return ok
}

func RecordConnPostgres(conf *CustomDB) (string, error) {
	connStr := fmt.Sprintf(
		"host=localhost user=%s password='%s' dbname=%s port=%s sslmode=disable",
		conf.DB.User, conf.DB.Password, conf.DB.Name, conf.Port)
	time.Sleep(8 * time.Second)

	return connStr, RecordConnection(conf.DB.Name, connStr, conf.Username, "postgres")
}

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

func SetMainDbByName(name, username, connStr, driver string) error {
	if !CheckConn(username) {
		return errors.New("Authentication failed")
	}

	if dbCached, ok := (*cst[username])[name]; ok {
		cstMain[username] = dbCached
		return nil
	}

	return RecordConnection(name, connStr, username, driver)
}

//func (ud UserDB) Remove(id int64) error {
//	_, err := ud.DB.Exec("DELETE FROM userDBs WHERE id = ?", id)
//	return err
//}
