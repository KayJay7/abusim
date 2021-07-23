package docker

import (
	"bufio"
	"context"
	"fmt"
	"io"
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

func (d DockerClient) getAgentFollowLogs(containerName string) (io.ReadCloser, error) {
	reader, err := d.client.ContainerLogs(
		context.Background(),
		containerName,
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
		},
	)

	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (d DockerClient) GetAgentLogsLines(containerName, name string, lines chan string) error {
	logs, err := d.getAgentFollowLogs(containerName)
	if err != nil {
		log.Fatalln(err)
	}
	reader := bufio.NewReader(logs)
	for {
		header := make([]byte, 8)
		io.ReadFull(reader, header)
		var bufLen uint32
		bufLen |= uint32(header[4]) << 24
		bufLen |= uint32(header[5]) << 16
		bufLen |= uint32(header[6]) << 8
		bufLen |= uint32(header[7])
		buf := make([]byte, bufLen)
		io.ReadFull(reader, buf)
		bufLines := string(buf)
		for _, line := range strings.Split(bufLines, "\n") {
			if len(line) > 0 {
				lines <- fmt.Sprintf("%s: %s", name, line)
			}
		}
	}
	return nil
}
