package dockerdb

import (
	"context"
	"github.com/docker/docker/api/types"
)

// Run start a docker container
func (ddb *DockerDB) Run(ctx context.Context) error {
	if err := ddb.cli.ContainerStart(ctx, ddb.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}
