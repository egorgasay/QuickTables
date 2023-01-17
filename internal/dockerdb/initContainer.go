package dockerdb

import (
	"context"
	"errors"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"quicktables/internal/userDB"
)

// TODO: IMPLEMENT DOWNLOAD ALL IMAGES

func InitContainer(conf *userDB.CustomDB) (*client.Client, string, error) {
	if conf == nil {
		return nil, "", errors.New("conf must be not nil")
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv,
		client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, "", err
	}

	var env []string
	var portDocker nat.Port

	switch conf.Vendor {
	case "postgres":
		portDocker = "5432/tcp"
		env = []string{"POSTGRES_DB=" + conf.DB.Name, "POSTGRES_USER=" + conf.DB.User, "POSTGRES_PASSWORD=" + conf.DB.Password}
	case "mysql":
		portDocker = "3306/tcp"
		env = []string{"MYSQL_DATABASE=" + conf.DB.Name, "MYSQL_USER=" + conf.DB.User, "MYSQL_ROOT_PASSWORD=" + conf.DB.Password,
			"MYSQL_PASSWORD=" + conf.DB.Password}
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			portDocker: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: conf.Port,
				},
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: conf.Vendor,
		Env:   env,
	}, hostConfig, nil, nil, conf.DB.User+"_"+conf.DB.Name)
	if err != nil {
		return nil, "", err
	}

	return cli, resp.ID, nil
}
