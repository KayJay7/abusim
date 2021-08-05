package command

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"steel-simulator/args"
	"steel-simulator/config"
	"steel-simulator/docker"
	"syscall"
)

// Up brings up the simulation environment
func Up(args *args.ArgsConfig, conf *config.Config, dcli *docker.DockerClient) {
	log.Println("Bringing up the environment...")
	// I create the networks...
	if err := dcli.CreateNetworks(conf.Namespace); err != nil {
		log.Fatalln(err)
	}
	// ... and I run the coordinator container
	if err := dcli.CreateAndRunCoordinatorContainer(conf.Namespace, conf.CoordinatorImage); err != nil {
		log.Fatalln(err)
	}
	// I create a list of the already created agents...
	endpoints := []string{}
	// ... and I range over all the agents to be created
	for name, agent := range conf.Agents {
		// I set the endpoints (all the containers created so far)...
		agent.Endpoints = endpoints
		// ... I set the container name...
		containerName := fmt.Sprintf("%s-%s", conf.Namespace, name)
		// ... and I set the configuration serialization
		agentSerialization, err := agent.Serialize()
		if err != nil {
			log.Fatalln(err)
		}
		// I run the agent container...
		err = dcli.CreateAndRunAgentContainer(conf.Namespace, conf.Image, containerName, agentSerialization)
		if err != nil {
			log.Fatalln(err)
		}
		// ... and I append it to the created agent endpoints
		endpoints = append(endpoints, fmt.Sprintf("%s-on-data:5000", containerName))
	}
	// If the command it was not invoked as detached...
	if !args.Detached {
		// ... I create an handler to tear down the environment at exit...
		setupCloseHandler(conf, dcli)
		// ... and I show the logs
		LogsFollow(conf, dcli)
	}
}

// setupCloseHandler waits for a SIGTERM and then tears down the environment
func setupCloseHandler(conf *config.Config, dcli *docker.DockerClient) {
	// I register for the SIGTERMs...
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// ... and I run a goroutine to handle their arrival
	go func() {
		// I block until a SIGTERM...
		<-c
		fmt.Println()
		// ... I tear down the environment...
		Down(conf, dcli)
		// ... and I exit
		os.Exit(0)
	}()
}
