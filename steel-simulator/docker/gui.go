package docker

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

// CreateAndRunGUIContainer creates and runs a container for the GUI
func (d DockerClient) CreateAndRunGUIContainer(namespace, image string, port int) error {
	// I prepare the container name...
	containerName := fmt.Sprintf("%s-gui", namespace)
	// ... I prepare the host configuration to open a port...
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"80/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprintf("%d", port),
				},
			},
		},
	}
	// ... I prepare the exposed ports...
	exposedPorts := nat.PortSet{
		"80/tcp": struct{}{},
	}
	//  ... I prepare the container configuration, passing the image, name and ports...
	config := &container.Config{
		Image:        image,
		Hostname:     containerName,
		ExposedPorts: exposedPorts,
	}
	// ... and I create the container, with the container, host and network configuration
	cont, err := d.client.ContainerCreate(
		context.Background(),
		config,
		hostConfig,
		nil,
		nil,
		containerName,
	)
	// If I get an error...
	if err != nil {
		// ... and this error is due to a duplicate container...
		if strings.HasPrefix(err.Error(), "Error response from daemon: Conflict.") {
			// ... I remove the duplicate...
			log.Printf("Found container \"%s\", recreating", containerName)
			d.RemoveContainer(containerName)
			// ... and I restart the creation
			return d.CreateAndRunGUIContainer(namespace, image, port)
		}
		return err
	}
	// Finally, I start the container
	err = d.client.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	log.Printf("Created GUI container \"%s\" with ID %s", containerName, cont.ID)
	return nil
}
