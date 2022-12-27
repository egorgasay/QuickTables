package service

import (
	"quicktables/internal/repository"
)

func New(db *repository.Storage) *Service {
	return &Service{DB: db}
}

// IService work with main DB
type IService interface {
	SaveQuery(status int, query, author, dbName, execTime string) error
	GetQueries(username, dbName string) ([]repository.QueryInfo, error)
	Disconnect() error
	DeleteAccount() error
	CreateUser(username, password string) error
	CheckPassword(username, password string) bool
	AddDB(dbName, strCon, owner, driver string) error
	CheckDB(owner string) bool
	GetAllDBs(owner string) [][]string
	GetDB(owner string) (dbName, connStr, driver string)
	GetDBbyName(owner, name string) (connStr, driver string)
}

type Service struct {
	DB IService
}
