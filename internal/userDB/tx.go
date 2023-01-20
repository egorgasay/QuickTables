package userDB

import (
	"context"
	"errors"
)

func Begin(ctx context.Context, username string) error {
	db := cstMain[username]
	if db == nil {
		return errors.New("no active dbs")
	}

	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	cstMain[username].Tx = tx

	return nil
}

func Commit(username string) error {
	db := cstMain[username]
	if db == nil {
		return errors.New("no active dbs")
	}

	err := db.Tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func Rollback(username string) error {
	db := cstMain[username]
	if db == nil {
		return errors.New("no active dbs")
	}

	err := db.Tx.Rollback()
	if err != nil {
		return err
	}

	return nil
}
