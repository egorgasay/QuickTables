package userDB

import (
	"database/sql"
	"github.com/docker/docker/client"
)

type UserDB struct {
	Name      string
	Conn      *sql.Conn
	Driver    string
	ConnStr   string
	DockerCli *client.Client
	Tx        *sql.Tx
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
	Vendor   string
}

// Add mutex here
type Storages map[string]*UserDB
type ConnStorage map[string]*Storages
type ConnStorageMain map[string]*UserDB

var st = Storages{"": nil}
var cst = ConnStorage{"": &st}
var cstMain = make(ConnStorageMain)
