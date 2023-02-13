package userDB

import (
	"database/sql"
	"github.com/docker/docker/client"
	"sync"
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
	Mu     sync.RWMutex
}

type UserDBs struct {
	DBs map[string]*ConnStorage
	Mu  sync.RWMutex
}

func New() *UserDBs {
	dbs := make(map[string]*ConnStorage)
	return &UserDBs{DBs: dbs}
}
