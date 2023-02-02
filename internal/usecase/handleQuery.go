package usecase

import (
	"context"
	"database/sql"
	"log"
	"strings"
)

func (uc UseCase) HandleQuery(query, username string) (QueryResponse, error) {
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
	usDB := uc.userDBs.GetActiveDB(username)

	err = usDB.Begin(ctx)
	if err != nil {
		return QueryResponse{}, err
	}

	defer func() {
		err := usDB.Rollback()
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

			_, err = usDB.Exec(ctx, shortQuery)
			if err != nil {
				return QueryResponse{}, err
			}
			continue
		}

		rows, err = usDB.Query(ctx, shortQuery)
		if err != nil {
			return QueryResponse{}, err
		}

		isSelect = true
		queries = make([]string, 0, len(lines))
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
