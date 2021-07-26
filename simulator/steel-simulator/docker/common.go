package docker

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
)

func (d DockerClient) RemoveContainer(containerName string) error {
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

func (d DockerClient) getLogs(containerName string, follow bool) (io.ReadCloser, error) {
	reader, err := d.client.ContainerLogs(
		context.Background(),
		containerName,
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     follow,
		},
	)

	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (d DockerClient) GetContainerLogsLines(containerName, name string, lines chan string, follow bool) error {
	logs, err := d.getLogs(containerName, follow)
	if err != nil {
		log.Fatalln(err)
	}
	close := make(chan struct{}, 1)
	exitBell := time.AfterFunc(time.Second, func() { close <- struct{}{} })
	reader := bufio.NewReader(logs)
	go func() {
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
					exitBell.Reset(time.Second)
				}
			}
		}
	}()
	if !follow {
		<-close
	}
	return nil
}
