package repository

import (
	"bufio"
	"database/sql"
	"errors"
	"os"
	"strings"
)

func InitDB(cfg *Config) (*sql.DB, error) {
	if cfg == nil {
		return nil, errors.New("invalid cfg")
	}

	simpleTest := "SELECT Name FROM Users"
	db, err := sql.Open(cfg.DriverName, cfg.DataSourceName)
	if err != nil {
		return nil, err
	}
	var tmp string
	if err = db.QueryRow(simpleTest).Scan(&tmp); err == nil {
		return db, nil
	}

	f, err := os.Open("schema.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var scanner = bufio.NewScanner(f)
	var str strings.Builder

	for scanner.Scan() {
		str.WriteString(scanner.Text())
	}

	queries := strings.Split(str.String(), "EOQUERY")
	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
