package models

import (
	"errors"
	"quicktables/internal/userDB"
)

type User struct {
	UserID   int64
	UserName string
	Password string
	DBs      []*userDB.UserDB
}

type IUser interface {
	AddDB(connectionString string) error
	RemoveDB(dbID int64) error
}

func New(userName string, password string) *User {
	// userID := SELECT count(*) FROM users
	return &User{}
}

func (u *User) AddDB(connectionString string) error {
	var newDB *userDB.UserDB
	db, err := newDB.New(connectionString)
	if err != nil {
		return err
	}
	u.DBs = append(u.DBs, db)
	return nil
}

func (u *User) RemoveDB(dbID int64) error {
	for i, db := range u.DBs {
		if db.Id == dbID {
			err := db.Remove(dbID)
			if err != nil {
				return err
			}
			u.DBs = append(u.DBs[:i], u.DBs[i+1:]...)
			return nil
		}
	}
	return errors.New("указанной бд не существует")
}
