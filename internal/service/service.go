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
	AddDB(strCon, driver, owner string) error
	CheckDB(owner string) bool
	GetDB(owner string) (connStr, driver string)
}

type Service struct {
	DB IService
}
