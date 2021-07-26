package docker

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
)

func (d DockerClient) CreateNetworks(namespace string) error {
	dataName := fmt.Sprintf("%s-data", namespace)
	controlName := fmt.Sprintf("%s-control", namespace)
	_, err := d.client.NetworkCreate(context.Background(), dataName, types.NetworkCreate{})
	if err != nil {
		return err
	}
	log.Printf("Created network \"%s\"", dataName)
	_, err = d.client.NetworkCreate(context.Background(), controlName, types.NetworkCreate{})
	if err != nil {
		return err
	}
	log.Printf("Created network \"%s\"", controlName)
	return nil
}

func (d DockerClient) RemoveNetworks(namespace string) error {
	dataName := fmt.Sprintf("%s-data", namespace)
	controlName := fmt.Sprintf("%s-control", namespace)
	err := d.client.NetworkRemove(context.Background(), dataName)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "Error: No such network:") {
			return err
		}
	}
	log.Printf("Removed network \"%s\"", dataName)
	err = d.client.NetworkRemove(context.Background(), controlName)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "Error: No such network:") {
			return err
		}
	}
	log.Printf("Removed network \"%s\"", controlName)
	return nil
}
