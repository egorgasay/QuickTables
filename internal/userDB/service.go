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

type ConnStorage struct {
	DBs    map[string]*UserDB
	Active *UserDB
}

type UserDBs map[string]*ConnStorage

func New() *UserDBs {
	dbs := make(UserDBs)
	return &dbs
}

func (udbs *UserDBs) GetUserDBs(username string) *ConnStorage {
	if (*udbs)[username] == nil {
		(*udbs)[username] = &ConnStorage{}
	}

	return (*udbs)[username]
}
