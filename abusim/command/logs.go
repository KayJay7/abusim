package command

import (
	"abusim/config"
	"abusim/docker"
	"fmt"
	"log"
	"sync"
)

// LogsFollow prints the logs of the containers continuously
func LogsFollow(conf *config.Config, dcli *docker.DockerClient) {
	// I prepare a channel for the log lines...
	lines := make(chan string)
	// ... and I prepare a wait group for syncronization
	wg := sync.WaitGroup{}
	// I spawn a goroutine to get the log lines of the coordinator
	wg.Add(1)
	go func(lines chan string) {
		defer wg.Done()
		err := dcli.GetContainerLogsLines(fmt.Sprintf("%s-coordinator", conf.Namespace), "coordinator", lines, true)
		if err != nil {
			log.Fatalln(err)
		}
	}(lines)
	// I spawn a goroutine for each agent to get their log lines
	for name := range conf.Agents {
		containerName := fmt.Sprintf("%s-%s", conf.Namespace, name)
		wg.Add(1)
		go func(containerName, name string, lines chan string) {
			defer wg.Done()
			err := dcli.GetContainerLogsLines(containerName, name, lines, true)
			if err != nil {
				log.Fatalln(err)
			}
		}(containerName, name, lines)
	}
	// I spawn a goroutine to print the lines
	wg.Add(1)
	go func() {
		for {
			// I wait for a line...
			line := <-lines
			// ... and I print it
			fmt.Println(line)
		}
	}()
	wg.Wait()
}

// LogsFollow prints the logs of the containers
func Logs(conf *config.Config, dcli *docker.DockerClient) {
	// I prepare a channel for the log lines...
	lines := make(chan string)
	// ... and I prepare a channel for exiting the function
	close := make(chan struct{}, len(conf.Agents))
	// I spawn a goroutine to get the log lines of the coordinator
	go func(lines chan string) {
		err := dcli.GetContainerLogsLines(fmt.Sprintf("%s-coordinator", conf.Namespace), "coordinator", lines, true)
		if err != nil {
			log.Fatalln(err)
		}
		close <- struct{}{}
	}(lines)
	// I spawn a goroutine for each agent to get their log lines
	for name := range conf.Agents {
		containerName := fmt.Sprintf("%s-%s", conf.Namespace, name)
		go func(containerName, name string, lines chan string) {
			err := dcli.GetContainerLogsLines(containerName, name, lines, false)
			if err != nil {
				log.Fatalln(err)
			}
			close <- struct{}{}
		}(containerName, name, lines)
	}
	// I prepare a counter for the exited goroutines
	cnt := 0
	for {
		// If all the goroutines are exited I also exit...
		if cnt == len(conf.Agents)+1 {
			break
		}
		// ... otherwise I wait on the channels
		select {
		// If a line arrives I print it...
		case line := <-lines:
			fmt.Println(line)
		// ... otherwise, if a closing signal arrives, I increment the counter
		case <-close:
			cnt++
		}
	}
}
