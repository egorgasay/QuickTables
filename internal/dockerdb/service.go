package dockerdb

import (
	"github.com/docker/docker/client"
	"quicktables/internal/userDB"
)

type DockerDB struct {
	ID   string
	cli  *client.Client
	conf *userDB.CustomDB
}

func New(conf *userDB.CustomDB) (*DockerDB, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv,
		client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &DockerDB{cli: cli, conf: conf}, nil
}
