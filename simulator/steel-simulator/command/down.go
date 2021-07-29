package command

import (
	"fmt"
	"log"
	"steel-simulator/config"
	"steel-simulator/docker"
)

// Down tears down the simulation environment
func Down(conf *config.Config, dcli *docker.DockerClient) {
	log.Println("Tearing down the environment...")
	// I range over the agents...
	for name := range conf.Agents {
		// ... and I remove the containers
		containerName := fmt.Sprintf("%s-%s", conf.Namespace, name)
		err := dcli.RemoveContainer(containerName)
		if err != nil {
			log.Println(err)
		}
	}
	// I remove the coordinator container...
	if err := dcli.RemoveContainer(fmt.Sprintf("%s-coordinator", conf.Namespace)); err != nil {
		log.Println(err)
	}
	// ... and I remove the networks
	if err := dcli.RemoveNetworks(conf.Namespace); err != nil {
		log.Println(err)
	}
}
