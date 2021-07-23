package command

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"steel-simulator-config/config"
	"steel-simulator/args"
	"steel-simulator/docker"
	"sync"
	"syscall"
)

func Up(args *args.ArgsConfig, conf *config.Config, dcli *docker.DockerClient) {
	log.Println("Bringing up the environment...")

	for name, agent := range conf.Agents {
		containerName := fmt.Sprintf("%s-%s", conf.Namespace, name)
		agentSerialization, err := agent.Serialize()
		if err != nil {
			log.Fatalln(err)
		}
		err = dcli.CreateAndRunAgentContainer(conf.Image, containerName, agentSerialization)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if !args.Detached {
		setupCloseHandler(conf, dcli)

		lines := make(chan string)
		wg := sync.WaitGroup{}
		for name := range conf.Agents {
			containerName := fmt.Sprintf("%s-%s", conf.Namespace, name)
			wg.Add(1)
			go func(containerName, name string, lines chan string) {
				defer wg.Done()
				err := dcli.GetAgentLogsLines(containerName, name, lines)
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
}

func setupCloseHandler(conf *config.Config, dcli *docker.DockerClient) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println()
		Down(conf, dcli)
		os.Exit(0)
	}()
}
