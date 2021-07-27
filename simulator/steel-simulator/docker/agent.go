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

func (d DockerClient) CreateAndRunAgentContainer(namespace, image, containerName, agentSerialization string) error {
	cmd := []string{agentSerialization}

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			fmt.Sprintf("%s-data", namespace): {
				NetworkID: fmt.Sprintf("%s-data", namespace),
				Aliases:   []string{fmt.Sprintf("%s-on-data", containerName)},
			},
		},
	}

	exposedPorts := nat.PortSet{
		"5000/tcp": struct{}{},
		"5001/tcp": struct{}{},
	}

	config := &container.Config{
		Image:        image,
		Hostname:     containerName,
		Cmd:          cmd,
		ExposedPorts: exposedPorts,
	}

	cont, err := d.client.ContainerCreate(
		context.Background(),
		config,
		nil,
		networkConfig,
		nil,
		containerName,
	)

	if err != nil {
		if strings.HasPrefix(err.Error(), "Error response from daemon: Conflict.") {
			log.Printf("Found container \"%s\", recreating", containerName)
			d.RemoveContainer(containerName)
			return d.CreateAndRunAgentContainer(namespace, image, containerName, agentSerialization)
		}
		return err
	}

	if err := d.client.NetworkConnect(context.Background(), fmt.Sprintf("%s-control", namespace), cont.ID, nil); err != nil {
		return err
	}

	err = d.client.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	log.Printf("Created container \"%s\" with ID %s", containerName, cont.ID)

	return nil
}
