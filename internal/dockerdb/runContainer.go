package dockerdb

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func RunContainer(ctx context.Context, cli *client.Client, id string) error {
	if err := cli.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}
