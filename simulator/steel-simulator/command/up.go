package command

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"steel-simulator-config/config"
	"steel-simulator/args"
	"steel-simulator/docker"
	"syscall"
)

func Up(args *args.ArgsConfig, conf *config.Config, dcli *docker.DockerClient) {
	log.Println("Bringing up the environment...")

	if err := dcli.CreateNetworks(conf.Namespace); err != nil {
		log.Fatalln(err)
	}

	if err := dcli.CreateAndRunCoordinatorContainer(conf.Namespace); err != nil {
		log.Fatalln(err)
	}

	for name, agent := range conf.Agents {
		containerName := fmt.Sprintf("%s-%s", conf.Namespace, name)
		agentSerialization, err := agent.Serialize()
		if err != nil {
			log.Fatalln(err)
		}
		err = dcli.CreateAndRunAgentContainer(conf.Namespace, conf.Image, containerName, agentSerialization)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if !args.Detached {
		setupCloseHandler(conf, dcli)
		LogsFollow(conf, dcli)
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
