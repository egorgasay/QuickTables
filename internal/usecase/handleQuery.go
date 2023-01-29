package usecase

import (
	"context"
	"database/sql"
	"log"
	"quicktables/internal/userDB"
	"strings"
)

type QueryResponse struct {
	Status    uint8
	Rows      [][]sql.NullString
	Cols      []string
	HTMLTable string
	IsSelect  bool
}

func HandleQuery(udbs *userDB.ConnStorage, query string) (QueryResponse, error) {
	cleanQuery := strings.Trim(query, "\r\n")
	garbage := "\r\n "

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

	lines := strings.Split(query, "\n")
	queries := make([]string, 0, len(lines))
	var rows *sql.Rows
	var err error
	var isSelect bool

	ctx := context.Background()

	err = udbs.Begin(ctx)
	if err != nil {
		return QueryResponse{}, err
	}

	defer func() {
		err := udbs.Rollback()
		if err != nil {
			log.Println(err)
		}
	}()

	for _, line := range lines {
		line = strings.Trim(line, " \r")
		if !strings.HasSuffix(line, ";") {
			queries = append(queries, line)
			continue
		}
		ctx := context.Background()
		shortQuery := strings.Join(queries, "\n") + line

		if !strings.HasPrefix(strings.ToLower(shortQuery), "select") {
			queries = make([]string, 0, len(lines))

			_, err = udbs.Exec(ctx, shortQuery)
			if err != nil {
				return QueryResponse{}, err
			}
			continue
		}

		rows, err = udbs.Query(ctx, shortQuery)
		if err != nil {
			return QueryResponse{}, err
		}

		isSelect = true
		queries = make([]string, 0, len(lines))
	}

	if !isSelect {
		err = udbs.Commit()
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
		return QueryResponse{HTMLTable: doLargeTable(cols, rowsArr), IsSelect: isSelect}, nil
	}

	err = udbs.Commit()
	if err != nil {
		return QueryResponse{}, err
	}

	return QueryResponse{Rows: rowsArr, Cols: cols, IsSelect: isSelect}, nil
}
