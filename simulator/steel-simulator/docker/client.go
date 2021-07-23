package docker

import (
	"github.com/docker/docker/client"
)

type DockerClient struct {
	client *client.Client
}

func New() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	return &DockerClient{
		client: cli,
	}, nil
}
