package usecase

import (
	"database/sql"
	"github.com/jedib0t/go-pretty/table"
)

func doTableFromData(cols []string, rows *sql.Rows) [][]sql.NullString {
	readCols := make([]interface{}, len(cols))
	writeCols := make([]sql.NullString, len(cols))

	rowsArr := make([][]sql.NullString, 0, 1000)
	for i := 0; rows.Next(); i++ {

		for i := range writeCols {
			readCols[i] = &writeCols[i]
		}

		err := rows.Scan(readCols...)
		if err != nil {
			panic(err)
		}
		rowsArr = append(rowsArr, make([]sql.NullString, len(cols)))
		copy(rowsArr[i], writeCols)
	}

	return rowsArr
}

func doLargeTable(cols []string, rowsArr [][]sql.NullString) string {
	t := table.NewWriter()

	colsForTable := make(table.Row, 0, 10)
	for _, el := range cols {
		colsForTable = append(colsForTable, el)
	}

	t.AppendHeader(colsForTable)

	rowsForTable := make([]table.Row, 0, 2000)
	for _, el := range rowsArr {
		rowForTable := make(table.Row, 0, 10)

		for _, el := range el {
			rowForTable = append(rowForTable, el.String)
		}

		rowsForTable = append(rowsForTable, rowForTable)
	}

	t.AppendRows(rowsForTable)

	return t.RenderHTML()
}
