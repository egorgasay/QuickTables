package userDB

import (
	"context"
)

func (cs *ConnStorage) Begin(ctx context.Context) error {
	tx, err := cs.Active.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	cs.Active.Tx = tx

	return nil
}

func (cs *ConnStorage) Commit() error {
	err := cs.Active.Tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (cs *ConnStorage) Rollback() error {
	err := cs.Active.Tx.Rollback()
	if err != nil {
		return err
	}

	return nil
}
