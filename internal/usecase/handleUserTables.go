package usecase

import (
	"context"
	"fmt"
)

func (uc UseCase) GetListOfUserTables(ctx context.Context, username string) ([]string, error) {
	activeDB := uc.userDBs.GetUserDBs(username)
	return activeDB.GetAllTables(ctx)
}

func (uc UseCase) GetUserTable(ctx context.Context, username, tname string) (*Table, error) {
	udbs := uc.userDBs.GetUserDBs(username)
	err := udbs.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer udbs.Rollback()

	var query string
	if udbs.Active.Driver == "postgres" {
		query = fmt.Sprintf(`SELECT * FROM "%s"`, tname)
	} else {
		query = fmt.Sprintf(`SELECT * FROM %s`, tname)
	}

	rows, err := udbs.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	cols, _ := rows.Columns()

	rowsArr := doTableFromData(cols, rows)

	if len(rowsArr) > 1000 {
		return &Table{HTMLTable: doLargeTable(cols, rowsArr)}, nil
	}

	return &Table{Rows: rowsArr, Cols: cols}, nil
}
