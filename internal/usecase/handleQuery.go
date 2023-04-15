package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"
)

func (uc *UseCase) HandleQuery(query, username string) (QueryResponse, error) {
	cleanQuery := strings.Trim(query, "\r\n")
	garbage := "\r\n"

	// cleaning a query
	for strings.Contains(garbage, string(cleanQuery[len(cleanQuery)-1])) ||
		strings.Contains(garbage, string(cleanQuery[0])) {
		cleanQuery = strings.Trim(query, garbage)
		log.Println(cleanQuery)
	}

	if cleanQuery[len(cleanQuery)-1] != ';' {
		cleanQuery += ";"
	}

	query = cleanQuery

	queries := strings.Split(query, ";")
	var rows *sql.Rows
	var err error
	var isSelect bool

	ctx := context.Background()
	usDB := uc.userDBs.GetUserDBs(username)

	err = usDB.Begin(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrConnDone) {
			err = usDB.RefreshConn()
			if err == nil {
				return uc.HandleQuery(query, username)
			}
		}
		return QueryResponse{}, err
	}

	defer func() {
		err := usDB.Rollback()
		if err != nil {
			log.Println(err)
		}
	}()

	for _, query := range queries {
		if query == "" {
			continue
		}
		ctx := context.Background()
		query = strings.Trim(query, garbage)
		if !strings.Contains(strings.ToLower(query), "select") {
			_, err = usDB.Exec(ctx, query)
			if err != nil {
				return QueryResponse{}, err
			}
			continue
		}

		rows, err = usDB.Query(ctx, query)
		if err != nil {
			return QueryResponse{}, err
		}

		isSelect = true
	}

	if !isSelect {
		err = usDB.Commit()
		if err != nil {
			return QueryResponse{}, err
		}
		return QueryResponse{IsSelect: isSelect}, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return QueryResponse{}, err
	}

	rowsArr := doTableFromData(cols, rows)
	if len(rowsArr) > 1000 {
		return QueryResponse{Table: Table{HTMLTable: doLargeTable(cols, rowsArr)}, IsSelect: isSelect}, nil
	}

	err = usDB.Commit()
	if err != nil {
		return QueryResponse{}, err
	}

	return QueryResponse{Table: Table{Rows: rowsArr, Cols: cols}, IsSelect: isSelect}, nil
}

func (uc *UseCase) HandleUserQueries(query, username, currentDB string) (QueryResponse, error) {
	start := time.Now()
	qh, err := uc.HandleQuery(query, username)
	if err != nil {
		go uc.SaveQuery(2, query, username, currentDB, "0")
		return qh, err
	}

	go uc.SaveQuery(1, query, username, currentDB, time.Since(start).String())
	return qh, nil
}

func (uc *UseCase) SaveQuery(status uint8, query, author, dbName, execTime string) {
	err := uc.service.DB.SaveQuery(status, query, author, dbName, execTime)
	if err != nil {
		log.Println(err)
	}
}
