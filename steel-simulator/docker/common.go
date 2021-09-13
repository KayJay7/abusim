package docker

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
)

// RemoveContainer removes a container
func (d DockerClient) RemoveContainer(containerName string) error {
	// I remove the specified container, forcefully...
	err := d.client.ContainerRemove(
		context.Background(),
		containerName,
		types.ContainerRemoveOptions{
			Force: true,
		},
	)
	// If I get an error...
	if err != nil {
		// ... and this error is due to a missing container, I ignore it
		if !strings.HasPrefix(err.Error(), "Error: No such container:") {
			return err
		}
	}
	log.Printf("Removed container \"%s\"", containerName)
	return nil
}

// getLogs returns a log reader for a container
func (d DockerClient) getLogs(containerName string, follow bool) (io.ReadCloser, error) {
	// I get the logs for the specified container, eventually following them...
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
	// ... and I return the reader
	return reader, nil
}

// GetContainerLogsLines returns a channel of log lines from a container
func (d DockerClient) GetContainerLogsLines(containerName, name string, lines chan string, follow bool) error {
	// I get the logs reader for the specified container...
	logs, err := d.getLogs(containerName, follow)
	if err != nil {
		log.Fatalln(err)
	}
	// ... I prepare a channer for exiting the lines goroutine...
	close := make(chan struct{}, 1)
	// ... I prepare a bell to exit the function after a second...
	exitBell := time.AfterFunc(time.Second, func() { close <- struct{}{} })
	// ... and I wrap the reader
	reader := bufio.NewReader(logs)
	// I spawn a goroutine to get the log lines...
	go func() {
		for {
			// ... I get the 8 bytes of message header...
			header := make([]byte, 8)
			io.ReadFull(reader, header)
			// ... I read the message length (the last 4 bytes)...
			bufLen := binary.BigEndian.Uint32(header[4:])
			// ... and I read the message
			buf := make([]byte, bufLen)
			io.ReadFull(reader, buf)
			bufLines := string(buf)
			// I split the lines in the message...
			for _, line := range strings.Split(bufLines, "\n") {
				if len(line) > 0 {
					// ... I send it to the channel...
					lines <- fmt.Sprintf("%s: %s", name, line)
					// ... and I reset the bell
					exitBell.Reset(time.Second)
				}
			}
		}
	}()
	// If the logs are not following, I exit at the bell
	if !follow {
		<-close
	}
	return nil
}
