package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"quicktables/internal/dockerdb"
	"quicktables/internal/userDB"
	"time"
)

func (uc *UseCase) CreateSqlite(username, dbName string) error {
	path := fmt.Sprintf("users/%s/", username)
	err := os.MkdirAll(path, 777)
	if err != nil {
		return fmt.Errorf("create sqlite: %w", err)
	}

	err = uc.service.DB.AddDB(dbName, path+dbName, username, "sqlite3", "")
	if err != nil {
		return fmt.Errorf("create sqlite: %w", err)
	}

	usDBs := uc.userDBs.GetUserDBs(username)

	err = usDBs.RecordConnection(dbName, path+dbName, "sqlite3")
	if err != nil {
		return fmt.Errorf("create sqlite: %w", err)
	}

	return nil
}

func strConnBuilder(conf *userDB.CustomDB) (connStr string) {
	if conf.Vendor == "postgres" {
		connStr = fmt.Sprintf(
			"host=localhost user=%s password='%s' dbname=%s port=%s sslmode=disable",
			conf.DB.User, conf.DB.Password, conf.DB.Name, conf.Port)
	} else if conf.Vendor == "mysql" {
		connStr = fmt.Sprintf(
			"%s:%s@tcp(127.0.0.1:%s)/%s",
			conf.DB.User, conf.DB.Password, conf.Port, conf.DB.Name)
	}

	return connStr
}

func (uc *UseCase) HandleDocker(username string, ddb *dockerdb.DockerDB, conf *userDB.CustomDB) error {
	if ddb == nil {
		return errors.New("docker db is nil")
	}

	connStr := strConnBuilder(conf)
	if err := uc.checkConnDocker(connStr, conf.Vendor); err != nil {
		return err
	}

	usDBs := uc.userDBs.GetUserDBs(username)

	err := usDBs.SetMainDbByName(conf.DB.Name, connStr, conf.Vendor)
	if err != nil {
		return err
	}

	err = uc.service.DB.AddDB(conf.DB.Name, connStr, username, conf.Vendor, ddb.ID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) checkConnDocker(strConn, driver string) error {
	for attempt := 0; attempt < 25; attempt++ {
		db, err := sql.Open(driver, strConn)
		if err == nil && db.Ping() == nil {
			db.Close()
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return errors.New("can't connect")
}
