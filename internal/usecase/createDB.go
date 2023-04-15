package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/egorgasay/dockerdb/v2"
	"log"
	"os"
	"time"
)

func (uc *UseCase) CreateSqlite(username, dbName string) error {
	path := fmt.Sprintf("users/%s/", username)
	err := os.MkdirAll(path, 0777)
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

func strConnBuilder(conf *dockerdb.CustomDB) (connStr string) {
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

func (uc *UseCase) HandleDocker(username string, conf dockerdb.CustomDB) error {
	ctx := context.TODO()
	ddb, err := dockerdb.New(ctx, conf)
	if err != nil {
		return fmt.Errorf("can't create a dockerdb instance: %w", err)
	}

	connStr := strConnBuilder(&conf)
	if err := uc.checkConnDocker(connStr, conf.Vendor); err != nil {
		log.Println("checkConnDocker: ", err)
	}

	usDBs := uc.userDBs.GetUserDBs(username)

	err = usDBs.SetMainDbByName(conf.DB.Name, connStr, conf.Vendor)
	if err != nil {
		log.Println("SetMainDbByName: ", err)
	}

	err = uc.service.DB.AddDB(conf.DB.Name, connStr, username, conf.Vendor, ddb.ID)
	if err != nil {
		return fmt.Errorf("can't add a db: %w", err)
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
