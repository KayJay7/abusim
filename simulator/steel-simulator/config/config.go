package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"steel-simulator-common/config"

	"gopkg.in/yaml.v2"
)

// Config represents a configuration
type Config struct {
	Image            string
	CoordinatorImage string
	Namespace        string
	Agents           map[string]config.Agent
}

// rawConfig represents a YAML configuration file
type rawConfig struct {
	Version          string                  `yaml:"version"`
	Image            string                  `yaml:"image"`
	CoordinatorImage string                  `yaml:"coordinatorimage"`
	Namespace        string                  `yaml:"namespace"`
	Includes         []string                `yaml:"includes"`
	Agents           map[string]rawAgent     `yaml:"agents"`
	Prototypes       map[string]rawPrototype `yaml:"prototypes"`
}

type rawAgent struct {
	Prototype    string `yaml:"prototype"`
	rawPrototype `yaml:",inline"`
}

type rawPrototype struct {
	MemoryController string   `yaml:"memorycontroller"`
	Memory           []string `yaml:"memory"`
	Rules            []string `yaml:"rules"`
	Tick             string   `yaml:"tick"`
}

// Parse parses a configuration file given its path
func Parse(filename string) (*Config, error) {
	// I get the raw configuration data...
	rawConfig := parseRawConfig(filename)
	// ... and I parse it accordingly to its version
	switch rawConfig.Version {
	case "1.0":
		return parseVersion1dot0(rawConfig, path.Dir(filename))
	default:
		return nil, fmt.Errorf("unknown config file version \"%s\"", rawConfig.Version)
	}
}

// parseVersion1dot0 returns a valid configuration from a raw one
func parseVersion1dot0(rawConfig rawConfig, base string) (*Config, error) {
	// I create the configuration and I save the agent and coordinator image name and the namespace...
	conf := Config{}
	conf.Image = rawConfig.Image
	conf.CoordinatorImage = rawConfig.CoordinatorImage
	conf.Namespace = rawConfig.Namespace
	// ... and I create a map for the agents
	conf.Agents = make(map[string]config.Agent)
	// I loop over the raw configuration agents...
	for name, agent := range rawConfig.Agents {
		// I create the new agent structure
		configAgent := config.NewAgent(name)
		// If it has a prototype...
		if agent.Prototype != "" {
			// ... I check if the prototype exists...
			if proto, ok := rawConfig.Prototypes[agent.Prototype]; ok {
				// ... I set its memory controller and tick fields...
				configAgent.SetMemoryController(proto.MemoryController)
				err := configAgent.SetTick(proto.Tick)
				if err != nil {
					return nil, err
				}
				// ... I set the memory...
				for _, item := range proto.Memory {
					err := configAgent.AddMemoryItem(item)
					if err != nil {
						return nil, err
					}
				}
				// ... and I set the rules
				configAgent.Rules = append(configAgent.Rules, proto.Rules...)
			} else {
				return nil, fmt.Errorf("unknown prototype \"%s\"", agent.Prototype)
			}
		}
		// Moreover I assign the agent-specific fields, overwriting the prototype ones
		// I set the memory controller and tick fields...
		configAgent.SetMemoryController(agent.MemoryController)
		err := configAgent.SetTick(agent.Tick)
		if err != nil {
			return nil, err
		}
		// ... I set the memory...
		for _, item := range agent.Memory {
			err := configAgent.AddMemoryItem(item)
			if err != nil {
				return nil, err
			}
		}
		// ... and I set the rules
		configAgent.Rules = append(configAgent.Rules, agent.Rules...)
		// Finally, I add the agent to the agents list
		conf.Agents[name] = *configAgent
	}
	for _, include := range rawConfig.Includes {
		filenameInclude := path.Join(base, include)
		rawConfigInclude := parseRawConfig(filenameInclude)
		confInclude, err := parseVersion1dot0(rawConfigInclude, path.Dir(filenameInclude))
		if err != nil {
			return nil, err
		}
		for k, v := range confInclude.Agents {
			if _, ok := conf.Agents[k]; !ok {
				conf.Agents[k] = v
			} else {
				return nil, fmt.Errorf("duplicate agent \"%s\"", k)
			}
		}
	}
	return &conf, nil
}

// parseRawConfig parses a YAML file into a raw configuration
func parseRawConfig(filename string) rawConfig {
	// I read the file content...
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Cannot read config file: %v", err)
	}
	// ... and I unmarshal the YAML into the raw structure
	rawconfig := rawConfig{}
	err = yaml.Unmarshal(yamlFile, &rawconfig)
	if err != nil {
		log.Fatalf("Cannot parse config file: %v", err)
	}
	// Finally, I return the raw structure
	return rawconfig
}
