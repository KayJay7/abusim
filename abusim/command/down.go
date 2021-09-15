package command

import (
	"fmt"
	"log"

	"github.com/abu-lang/abusim/abusim/args"
	"github.com/abu-lang/abusim/abusim/config"
	"github.com/abu-lang/abusim/abusim/docker"
)

// Down tears down the simulation environment
func Down(args *args.ArgsConfig, conf *config.Config, dcli *docker.DockerClient) {
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
	// I eventually remove the GUI container...
	if args.GUI {
		if err := dcli.RemoveContainer(fmt.Sprintf("%s-gui", conf.Namespace)); err != nil {
			log.Println(err)
		}
	}
	// ... I remove the coordinator container...
	if err := dcli.RemoveContainer(fmt.Sprintf("%s-coordinator", conf.Namespace)); err != nil {
		log.Println(err)
	}
	// ... and I remove the networks
	if err := dcli.RemoveNetworks(conf.Namespace); err != nil {
		log.Println(err)
	}
}
