package repository

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	DriverName     string
	DataSourceName string
}

type IStorage interface {
	Ping() error
	Close() error
	Query(string, ...any) (*sql.Rows, error)
	Exec(string, ...any) (sql.Result, error)
	QueryRow(string, ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
}

type Storage struct {
	DB IStorage
}

func New(cfg *Config) (*Storage, error) {
	if cfg == nil {
		return nil, errors.New("invalid cfg")
	}

	db, err := InitDB(cfg)
	if err != nil {
		return nil, err
	}
	return &Storage{DB: db}, nil
}

func (s Storage) Disconnect() error {
	return nil
}

func (s Storage) DeleteAccount() error {
	return nil
}
