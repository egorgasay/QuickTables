package usecase

import (
	"errors"
)

func (uc *UseCase) CheckPassword(username, password string) bool {
	return uc.service.DB.CheckPassword(username, password)
}

func (uc *UseCase) CheckAndGetDBs(username string) ([][]string, error) {
	if !uc.service.DB.CheckDB(username) {
		return nil, errors.New("user don't have dbs")
	}

	return uc.service.DB.GetAllDBs(username)
}
