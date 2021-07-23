package command

import (
	"log"
	"steel-simulator/config"
)

func Up(config *config.Config) {
	log.Println("up command")
	for name, agent := range config.Agents {
		log.Println("==========")
		log.Println(name)
		log.Println(agent.MemoryController)
		log.Println(agent.Memory)
		log.Println(agent.Rules)
	}
}
