package command

import (
	"fmt"
	"log"
	"steel-simulator-config/config"
	"steel-simulator/docker"
	"sync"
)

func LogsFollow(conf *config.Config, dcli *docker.DockerClient) {
	lines := make(chan string)
	wg := sync.WaitGroup{}
	for name := range conf.Agents {
		containerName := fmt.Sprintf("%s-%s", conf.Namespace, name)
		wg.Add(1)
		go func(containerName, name string, lines chan string) {
			defer wg.Done()
			err := dcli.GetAgentLogsLines(containerName, name, lines, true)
			if err != nil {
				log.Fatalln(err)
			}
		}(containerName, name, lines)
	}
	wg.Add(1)
	go func() {
		for {
			line := <-lines
			fmt.Println(line)
		}
	}()
	wg.Wait()
}

func Logs(conf *config.Config, dcli *docker.DockerClient) {
	lines := make(chan string)
	close := make(chan struct{}, len(conf.Agents))
	for name := range conf.Agents {
		containerName := fmt.Sprintf("%s-%s", conf.Namespace, name)
		go func(containerName, name string, lines chan string) {
			err := dcli.GetAgentLogsLines(containerName, name, lines, false)
			if err != nil {
				log.Fatalln(err)
			}
			close <- struct{}{}
		}(containerName, name, lines)
	}
	cnt := 0
	for {
		if cnt == len(conf.Agents) {
			break
		}
		select {
		case line := <-lines:
			fmt.Println(line)
		case <-close:
			cnt++
		}
	}
}
