package docker

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
)

// CreateNetworks creates the data and control network for the simulation
func (d DockerClient) CreateNetworks(namespace string) error {
	// I prepare the two network names...
	dataName := fmt.Sprintf("%s-data", namespace)
	controlName := fmt.Sprintf("%s-control", namespace)
	// ... I create the data network...
	_, err := d.client.NetworkCreate(context.Background(), dataName, types.NetworkCreate{})
	if err != nil {
		return err
	}
	log.Printf("Created network \"%s\"", dataName)
	// ... and I create the control network
	_, err = d.client.NetworkCreate(context.Background(), controlName, types.NetworkCreate{})
	if err != nil {
		return err
	}
	log.Printf("Created network \"%s\"", controlName)
	return nil
}

// RemoveNetworks removes the data and control network for the simulation
func (d DockerClient) RemoveNetworks(namespace string) error {
	// I prepare the two network names...
	dataName := fmt.Sprintf("%s-data", namespace)
	controlName := fmt.Sprintf("%s-control", namespace)
	// ... I remove the data network...
	err := d.client.NetworkRemove(context.Background(), dataName)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "Error: No such network:") {
			return err
		}
	}
	log.Printf("Removed network \"%s\"", dataName)
	// ... and I remove the control network
	err = d.client.NetworkRemove(context.Background(), controlName)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "Error: No such network:") {
			return err
		}
	}
	log.Printf("Removed network \"%s\"", controlName)
	return nil
}
