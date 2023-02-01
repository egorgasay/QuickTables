package usecase

import (
	"context"
	"fmt"
	"quicktables/internal/userDB"
)

func GetListOfUserTables(ctx context.Context, udbs userDB.ConnStorage, username string) ([]string, error) {
	list, err := udbs.GetAllTables(ctx, username)
	if err != nil {
		return nil, err
	}

	return list, err
}

func (uc UseCase) GetUserTable(ctx context.Context, udbs userDB.ConnStorage, username, tname string) (*Table, error) {
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
