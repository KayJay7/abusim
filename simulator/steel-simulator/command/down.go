package command

import (
	"fmt"
	"log"
	"steel-simulator-config/config"
	"steel-simulator/docker"
)

func Down(conf *config.Config, dcli *docker.DockerClient) {
	log.Println("Tearing down the environment...")

	for name := range conf.Agents {
		containerName := fmt.Sprintf("%s-%s", conf.Namespace, name)
		err := dcli.RemoveContainer(containerName)
		if err != nil {
			log.Println(err)
		}
	}
}
