package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func Parse(filename string) (*Config, error) {
	rawConfig := parseRawConfig(filename)

	switch rawConfig.Version {
	case "1.0":
		return parseVersion1dot0(rawConfig)
	default:
		return nil, fmt.Errorf("unknown config file version \"%s\"", rawConfig.Version)
	}
}

func parseVersion1dot0(rawConfig rawConfig) (*Config, error) {
	config := Config{}
	config.Agents = make(map[string]Agent)
	for name, agent := range rawConfig.Agents {
		configAgent := NewAgent()
		if agent.Prototype != "" {
			if proto, ok := rawConfig.Prototypes[agent.Prototype]; ok {
				configAgent.SetMemoryController(proto.MemoryController)
				for _, item := range proto.Memory {
					err := configAgent.AddMemoryItem(item)
					if err != nil {
						return nil, err
					}
				}
				for _, rule := range proto.Rules {
					configAgent.AddRule(rule)
				}
			} else {
				return nil, fmt.Errorf("unknown prototype \"%s\"", agent.Prototype)
			}
		}
		configAgent.SetMemoryController(agent.MemoryController)
		for _, item := range agent.Memory {
			err := configAgent.AddMemoryItem(item)
			if err != nil {
				return nil, err
			}
		}
		for _, rule := range agent.Rules {
			configAgent.AddRule(rule)
		}
		config.Agents[name] = *configAgent
	}
	return &config, nil
}

func parseRawConfig(filename string) rawConfig {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Cannot read config file: %v", err)
	}

	config := rawConfig{}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Cannot parse config file: %v", err)
	}

	return config
}
