package docker

import (
	"github.com/docker/docker/client"
)

// DockerClient represents a connection to the docker daemon
type DockerClient struct {
	client *client.Client
}

// New connects a new client to the docker daemon
func New() (*DockerClient, error) {
	// I create the docker client...
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	// ... and I return it
	return &DockerClient{
		client: cli,
	}, nil
}
