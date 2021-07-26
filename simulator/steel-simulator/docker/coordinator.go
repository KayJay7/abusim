package docker

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

func (d DockerClient) CreateAndRunCoordinatorContainer(namespace string) error {
	containerName := fmt.Sprintf("%s-coordinator", namespace)

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"4000/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "4000",
				},
			},
		},
	}

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			fmt.Sprintf("%s-control", namespace): {
				NetworkID: fmt.Sprintf("%s-control", namespace),
				Aliases:   []string{"steel-coordinator"},
			},
		},
	}

	exposedPorts := nat.PortSet{
		"4000/tcp": struct{}{},
		"5001/tcp": struct{}{},
	}

	config := &container.Config{
		Image:        "steel-coordinator",
		Hostname:     containerName,
		ExposedPorts: exposedPorts,
	}

	cont, err := d.client.ContainerCreate(
		context.Background(),
		config,
		hostConfig,
		networkConfig,
		nil,
		containerName,
	)

	if err != nil {
		if strings.HasPrefix(err.Error(), "Error response from daemon: Conflict.") {
			log.Printf("Found container \"%s\", recreating", containerName)
			d.RemoveContainer(containerName)
			return d.CreateAndRunCoordinatorContainer(namespace)
		}
		return err
	}

	err = d.client.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	log.Printf("Created coordinator container \"%s\" with ID %s", containerName, cont.ID)

	return nil
}
