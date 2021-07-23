package command

import (
	"fmt"
	"log"
	"steel-simulator/config"
)

func Up(conf *config.Config) {
	log.Println("up command")
	for name, agent := range conf.Agents {
		log.Println("==========")
		log.Println(name)
		log.Println(agent.MemoryController)
		log.Println(agent.Memory)
		log.Println(agent.Rules)
		serialization, err := agent.Serialize()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(serialization)
		newAgent := config.Agent{}
		newAgent.Deserialize(serialization)
		if fmt.Sprint(agent) != fmt.Sprint(newAgent) {
			log.Fatalln("Deserialized object is non identical to the original one")
		}
	}
}
