package userDB

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

func (cs *ConnStorage) GetDbNameAndVendor() (name string, vendor string) {
	return cs.Active.Name, cs.Active.Driver
}

func (cs *ConnStorage) CheckConnDocker(strConn, driver string) error {
	for attempt := 0; attempt < 25; attempt++ {
		db, err := sql.Open(driver, strConn)
		if err == nil && db.Ping() == nil {
			db.Close()
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return errors.New("can't connect")
}

func (cs *ConnStorage) Query(ctx context.Context, query string) (*sql.Rows, error) {
	if !cs.CheckConn() {
		return nil, errors.New("Authentication failed")
	}

	return cs.Active.Tx.QueryContext(ctx, query)
}

func (cs *ConnStorage) Exec(ctx context.Context, query string) (sql.Result, error) {
	if !cs.CheckConn() {
		return nil, errors.New("Authentication failed")
	}

	return cs.Active.Tx.ExecContext(ctx, query)
}

func (ud *UserDBs) GetDbName(username string) (string, error) {
	if (*ud)[username] == nil {
		return "", errors.New("no active dbs")
	}

	return (*ud)[username].Active.Name, nil
}

func (cs *ConnStorage) GetAllTables(ctx context.Context, username string) ([]string, error) {
	var query string
	if !cs.CheckConn() {
		return nil, errors.New("Authentication failed")
	}

	driver := cs.GetDbDriver()

	switch driver {
	case "mysql":
		query = fmt.Sprintf(`SELECT table_name
				FROM information_schema.tables
				WHERE table_type='BASE TABLE'
      			AND table_schema = '%s'`, cs.GetSysDbName())
	case "postgres":
		query = "SELECT table_name FROM information_schema.tables WHERE table_schema='public'"
	case "mssql":
		query = `SELECT name
		FROM sys.objects
		WHERE type_desc = 'USER_TABLE'`

	case "sqlite3":
		query = `SELECT name FROM sqlite_master 
		WHERE type IN ('table','view') 
		AND name NOT LIKE 'sqlite_%'
		ORDER BY 1;`
	}

	rows, err := cs.Active.Conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, 10)

	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}

	if len(names) == 0 {
		return nil, errors.New("no tables in default schema")
	}

	return names, nil
}

func (cs *ConnStorage) GetDbDriver() string {
	return cs.Active.Driver
}

func (cs *ConnStorage) GetSysDbName() string {
	connStr := strings.Split(cs.Active.ConnStr, "/")
	return connStr[len(connStr)-1]
}
