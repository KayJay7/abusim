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
	// ... I prepare the network configuration to join the data network...
	fmt.Printf("Begun agent creation: %v\n", containerName)
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
	// ... I prepare the container configuration, passing the image, name, stdin configuration and ports...
	config := &container.Config{
		Image:        image,
		Hostname:     containerName,
		ExposedPorts: exposedPorts,
		AttachStdin:  true,
		OpenStdin:    true,
	}
	// ... I also need to configure CPUs
	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			CPUCount:  1,
			CPUShares: 1024,
		},
	}
	// ... and prepare the attach configuration to connect to stdin
	attachConfig := types.ContainerAttachOptions{
		Stream: false, // We don't need blocking operations
		Stdin:  true,
	}

	// ... and I create the container, with the container and network configuration
	cont, err := d.client.ContainerCreate(
		context.Background(),
		config,
		hostConfig,
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
	fmt.Printf("Created container: %v\n", containerName)
	// I join the control network, since the network configuration only allows for a single network to be specified...
	if err := d.client.NetworkConnect(context.Background(), fmt.Sprintf("%s-control", namespace), cont.ID, nil); err != nil {
		return err
	}
	// ... and start the container...
	err = d.client.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Started container: %v\n", containerName)
	// ... and finally I stream the agent serialization on stdin
	conn, err := d.client.ContainerAttach(
		context.Background(),
		cont.ID,
		attachConfig,
	)
	if err != nil {
		return err
	}
	fmt.Printf("Attached container: %v\n", containerName)
	conn.Conn.Write([]byte(agentSerialization))
	conn.Conn.Write([]byte{'\n'}) // Unix files must end with /n
	conn.CloseWrite()
	conn.Conn.Close()
	log.Printf("Created container \"%s\" with ID %s", containerName, cont.ID)
	return nil
}
