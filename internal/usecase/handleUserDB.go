package usecase

import (
	"context"
	"quicktables/internal/userDB"
)

func (uc UseCase) DeleteUserDB(username, dbName string) error {
	err := uc.Service.DB.DeleteDB(username, dbName)
	if err != nil {
		return err
	}

	return nil
}

func (uc UseCase) AddUserDB(username, dbName, connStr, vendorName string) error {
	activeDB := uc.userDBs.GetUserDBs(username)
	err := activeDB.RecordConnection(dbName, connStr, vendorName)
	if err != nil {
		return err
	}

	err = uc.Service.DB.AddDB(dbName, connStr, username, vendorName, "")
	if err != nil {
		return err
	}

	return nil
}

func (uc UseCase) HandleUserDB(username, dbName string) error {
	if len(*uc.userDBs) == 0 {
		uc.initUserDBs(username)
	}

	connStr, driver, id := uc.Service.DB.GetDBInfobyName(username, dbName)

	activeDB := uc.userDBs.GetUserDBs(username)
	if id != "" && (activeDB == nil || !activeDB.IsDBCached(dbName)) {
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

	if activeDB == nil || activeDB.Active == nil {
		if err := activeDB.RecordConnection(dbName, connStr, driver); err != nil {
			return err
		}
	}

	if err := activeDB.SetMainDbByName(dbName, connStr, driver); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) initUserDBs(username string) error {
	(*uc.userDBs)[username] = &userDB.ConnStorage{}
	return nil
}
