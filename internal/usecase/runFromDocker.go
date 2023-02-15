package usecase

import (
	"context"
	"github.com/egorgasay/dockerdb"
)

func (uc *UseCase) runDBFromDocker(ctx context.Context, id string) error {
	err := dockerdb.Run(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
