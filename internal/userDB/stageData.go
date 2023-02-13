package userDB

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/sijms/go-ora/v2"
	"log"
	"time"
)

func (cs *ConnStorage) CheckConn() bool {
	if cs.Active != nil {
		return true
	}

	return false
}

func (cs *ConnStorage) IsDBCached(dbname string) bool {
	_, ok := cs.DBs[dbname]
	return ok
}

func (cs *ConnStorage) NewConn(cred, driver string) (*sql.Conn, error) {
	db, err := sql.Open(driver, cred)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	cwe := make(chan connWithErr)
	var userDB connWithErr

	go cs.createConn(ctx, db, cwe)

	select {
	case <-time.After(7 * time.Second):
		return nil, errors.New("timed out")
	case userDB = <-cwe:
	}

	return userDB.conn, userDB.err
}

type connWithErr struct {
	conn *sql.Conn
	err  error
}

func (cs *ConnStorage) createConn(ctx context.Context, db *sql.DB, cwe chan connWithErr) {
	conn, err := db.Conn(ctx)
	cwe <- connWithErr{conn, err}
}

func (cs *ConnStorage) RecordConnection(name, connStr, driver string) error {
	cn, err := cs.NewConn(connStr, driver)
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

	if cs == nil {
		cs = &ConnStorage{}
	}

	if cs.DBs == nil {
		cs.DBs = make(map[string]*UserDB)
	}

	cs.DBs[udb.Name] = udb
	cs.Active = udb

	return nil
}

func (cs *ConnStorage) SetMainDbByName(name, connStr, driver string) error {
	if cs.IsDBCached(name) {
		cs.Active = cs.DBs[name]
		return nil
	}

	return cs.RecordConnection(name, connStr, driver)
}

//func AddDockerCli(cli *client.Client, conf *CustomDB) error {
//	if !CheckConn(conf.Username) {
//		return errors.New("Authentication failed")
//	}
//
//	(*cst[conf.Username])[conf.DB.Name].DockerCli = cli
//
//	return nil
//}

//func (ud UserDB) Remove(id int64) error {
//	_, err := ud.DB.Exec("DELETE FROM userDBs WHERE id = ?", id)
//	return err
//}
