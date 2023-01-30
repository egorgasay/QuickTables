package usecase

import "database/sql"

type Table struct {
	HTMLTable string
	Rows      [][]sql.NullString
	Cols      []string
}

type QueryResponse struct {
	Status   uint8
	IsSelect bool
	Table    Table
}
