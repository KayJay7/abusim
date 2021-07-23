package docker

import (
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func (d DockerClient) CreateAndRunAgentContainer(image string, containerName string, agentSerialization string) error {
	cmd := []string{agentSerialization}
	config := &container.Config{
		Image:    image,
		Hostname: containerName,
		Cmd:      cmd,
	}

	cont, err := d.client.ContainerCreate(
		context.Background(),
		config,
		nil,
		nil,
		nil,
		containerName,
	)

	if err != nil {
		if strings.HasPrefix(err.Error(), "Error response from daemon: Conflict.") {
			log.Printf("Found container \"%s\", recreating", containerName)
			d.RemoveAgentContainer(containerName)
			return d.CreateAndRunAgentContainer(image, containerName, agentSerialization)
		}
		return err
	}

	d.client.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	log.Printf("Created container \"%s\" with ID %s", containerName, cont.ID)

	return nil
}

func (d DockerClient) RemoveAgentContainer(containerName string) error {
	err := d.client.ContainerRemove(
		context.Background(),
		containerName,
		types.ContainerRemoveOptions{
			Force: true,
		},
	)

	if err != nil {
		if !strings.HasPrefix(err.Error(), "Error: No such container:") {
			return err
		}
	}
	log.Printf("Removed container \"%s\"", containerName)

	return nil
}
