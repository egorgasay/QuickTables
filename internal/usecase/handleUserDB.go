package usecase

import (
	"context"
)

func (uc UseCase) DeleteUserDB(username, dbName string) error {
	err := uc.Service.DB.DeleteDB(username, dbName)
	if err != nil {
		return err
	}

	return nil
}

func (uc UseCase) HandleUserDB(username, dbName string) error {
	connStr, driver, id := uc.Service.DB.GetDBInfobyName(username, dbName)
	activeDB := uc.userDBs.GetActiveDB(username)
	if id != "" && !activeDB.IsDBCached(dbName) {
		ctx := context.Background()

		err := uc.runDBFromDocker(ctx, id)
		if err != nil {
			return err
		}
	}

	if err := activeDB.SetMainDbByName(dbName, connStr, driver); err != nil {
		return err
	}

	return nil
}
