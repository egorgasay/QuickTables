package userDB

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

//func (cs *ConnStorage) GetDbNameAndVendor() (name string, vendor string) {
//	return cs.Active.Name, cs.Active.Driver
//}

func (udbs *UserDBs) GetDBName(username string) (string, error) {
	if (*udbs)[username].Active == nil {
		return "", errors.New("no active dbs")
	}

	return (*udbs)[username].Active.Name, nil
}

func (udbs *UserDBs) GetDBVendor(username string) (string, error) {
	userDBs := (*udbs)[username]
	if userDBs == nil || userDBs.Active == nil {
		return "", errors.New("no active dbs")
	}

	return userDBs.Active.Driver, nil
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

func (cs *ConnStorage) GetAllTables(ctx context.Context) ([]string, error) {
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
