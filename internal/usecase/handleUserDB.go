package usecase

import (
	"context"
	"quicktables/internal/service"
	"quicktables/internal/userDB"
)

func (uc UseCase) DeleteUserDB(appDB service.IService, username, dbName string) error {
	err := appDB.DeleteDB(username, dbName)
	if err != nil {
		return err
	}

	return nil
}

func (uc UseCase) HandleUserDB(appDB service.IService, username, dbName string, udbs userDB.ConnStorage) (userDB.ConnStorage, error) {
	connStr, driver, id := appDB.GetDBInfobyName(username, dbName)
	if id != "" && !udbs.IsDBCached(dbName) {
		ctx := context.Background()

		err := uc.runDBFromDocker(ctx, id)
		if err != nil {
			return userDB.ConnStorage{}, err
		}
	}

	if err := udbs.SetMainDbByName(dbName, connStr, driver); err != nil {
		return userDB.ConnStorage{}, err
	}

	return udbs, nil
}
