package usecase

import (
	"context"
	"quicktables/internal/service"
	"quicktables/internal/userDB"
)

func DeleteUserDB(appDB service.IService, username, dbName string) error {
	err := appDB.DeleteDB(username, dbName)
	if err != nil {
		return err
	}

	return nil
}

func HandleUserDB(appDB service.IService, username, dbName string, udbs *userDB.ConnStorage) error {
	connStr, driver, id := appDB.GetDBInfobyName(username, dbName)
	if id != "" && !udbs.IsDBCached(dbName) {
		ctx := context.Background()

		err := runDBFromDocker(ctx, id)
		if err != nil {
			return err
		}
	}

	if err := udbs.SetMainDbByName(dbName, connStr, driver); err != nil {
		return err
	}

	return nil
}
