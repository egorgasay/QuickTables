package service

import (
	"quicktables/internal/repository"
)

func New(db *repository.Storage) *Service {
	return &Service{DB: db}
}

//go:generate mockgen -source=service.go -destination=mocks/mock.go

// IService work with main DB
type IService interface {
	SaveQuery(status int, query, author, dbName, execTime string) error
	GetQueries(username, dbName string) ([]repository.QueryInfo, error)
	Disconnect() error
	DeleteAccount() error
	DeleteDB(username, dbName string) error
	CreateUser(username, password string) error
	CheckPassword(username, password string) bool
	AddDB(dbName, strCon, owner, driver string) error
	CheckDB(owner string) bool
	GetAllDBs(owner string) [][]string
	GetDB(owner string) (dbName, connStr, driver string)
	GetDBbyName(owner, name string) (connStr, driver string)
	GetUserStats(username string) (*repository.UserStats, error)
	ChangeNick(username, nick string) error
	ChangePassword(username, oldPassword, newPassword string) error
}

type Service struct {
	DB IService
}
