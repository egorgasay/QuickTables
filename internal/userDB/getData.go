package userDB

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

func GetDbNameAndVendor(username string) (name string, vendor string) {
	return cstMain[username].Name, cstMain[username].Driver
}

func GetDbName(username string) string {
	return cstMain[username].Name
}

func GetAllTables(ctx context.Context, username string) ([]string, error) {
	var query string
	if !CheckConn(username) {
		return nil, errors.New("Authentication failed")
	}

	driver := GetDbDriver(username)

	switch driver {
	case "mysql":
		query = fmt.Sprintf(`SELECT table_name
				FROM information_schema.tables
				WHERE table_type='BASE TABLE'
      			AND table_schema = '%s'`, GetSysDbName(username))
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

	rows, err := cstMain[username].Conn.QueryContext(ctx, query)
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

func GetDbDriver(username string) string {
	return cstMain[username].Driver
}

func GetSysDbName(username string) string {
	connStr := strings.Split(cstMain[username].ConnStr, "/")
	return connStr[len(connStr)-1]
}
