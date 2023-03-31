package usecase

import (
	"context"
)

func (uc *UseCase) DeleteUserDB(username, dbName string) error {
	err := uc.service.DB.DeleteDB(username, dbName)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) AddUserDB(username, dbName, connStr, vendorName string) error {
	DBs := uc.userDBs.GetUserDBs(username)
	err := DBs.RecordConnection(dbName, connStr, vendorName)
	if err != nil {
		return err
	}

	err = uc.service.DB.AddDB(dbName, connStr, username, vendorName, "")
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) HandleUserDB(username, dbName string) error {
	connStr, driver, id := uc.service.DB.GetDBInfobyName(username, dbName)

	udbs := uc.userDBs.GetUserDBs(username)
	if udbs.Active == nil {
		if err := udbs.RecordConnection(dbName, connStr, driver); err != nil {
			return err
		}
	}

	if id != "" && !udbs.IsDBCached(dbName) {
		ctx := context.Background()

		err := uc.runDBFromDocker(ctx, id)
		if err != nil {
			return err
		}

		err = uc.checkConnDocker(connStr, driver)
		if err != nil {
			return err
		}
	}

	if err := udbs.SetMainDbByName(dbName, connStr, driver); err != nil {
		return err
	}

	return nil
}
