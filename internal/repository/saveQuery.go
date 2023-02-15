package repository

import "time"

func (s Storage) SaveQuery(status uint8, query, author, dbName, execTime string) error {
	_, err := s.DB.Exec(`INSERT INTO historyOfQueries (Author, DBName, Query, Status, ExecutionTime, Date)
VALUES (?, ?, ?, ?, ?, ?)`, author, dbName, query, status, execTime, time.Now().Unix())
	return err
}
