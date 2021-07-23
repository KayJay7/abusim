package command

import (
	"fmt"
	"log"
	"steel-simulator-config/config"
	"steel-simulator/docker"
)

func Up(conf *config.Config, dcli *docker.DockerClient) {
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
}
