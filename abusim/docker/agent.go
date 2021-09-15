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

// CreateAndRunAgentContainer creates and runs a container for an agent
func (d DockerClient) CreateAndRunAgentContainer(namespace, image, containerName, agentSerialization string) error {
	// I prepare the command line for the agent...
	cmd := []string{agentSerialization}
	// ... I prepare the network configuration to join the data network...
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			fmt.Sprintf("%s-data", namespace): {
				NetworkID: fmt.Sprintf("%s-data", namespace),
				Aliases:   []string{fmt.Sprintf("%s-on-data", containerName)},
			},
		},
	}
	// ... I prepare the exposed ports...
	exposedPorts := nat.PortSet{
		"5000/tcp": struct{}{},
		"5001/tcp": struct{}{},
	}
	//  ... I prepare the container configuration, passing the image, name, command and ports...
	config := &container.Config{
		Image:        image,
		Hostname:     containerName,
		Cmd:          cmd,
		ExposedPorts: exposedPorts,
	}
	// ... and I create the container, with the container and network configuration
	cont, err := d.client.ContainerCreate(
		context.Background(),
		config,
		nil,
		networkConfig,
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
			return d.CreateAndRunAgentContainer(namespace, image, containerName, agentSerialization)
		}
		return err
	}
	// I join the control network, since the network configuration only allows for a single network to be specified...
	if err := d.client.NetworkConnect(context.Background(), fmt.Sprintf("%s-control", namespace), cont.ID, nil); err != nil {
		return err
	}
	// ... and I start the container
	err = d.client.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	log.Printf("Created container \"%s\" with ID %s", containerName, cont.ID)
	return nil
}
