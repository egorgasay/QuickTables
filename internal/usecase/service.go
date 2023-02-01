package usecase

import (
	"database/sql"
	"quicktables/internal/service"
	"quicktables/internal/userDB"
)

type UseCase struct {
	Service *service.Service
	userDBs *userDB.UserDBs
}

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

func New(service *service.Service, userDBs *userDB.UserDBs) UseCase {
	if service == nil {
		panic("storage is nil")
	}

	if userDBs == nil {
		panic("userDBs is nil")
	}

	return UseCase{
		Service: service,
		userDBs: userDBs,
	}
}
