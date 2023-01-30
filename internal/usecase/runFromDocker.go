package usecase

import (
	"context"
	"quicktables/internal/dockerdb"
)

func runDBFromDocker(ctx context.Context, id string) error {
	ddb, err := dockerdb.New(nil)
	if err != nil {
		return err
	}

	ddb.ID = id

	err = ddb.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}
