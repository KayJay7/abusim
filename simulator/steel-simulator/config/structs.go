package config

type Config struct {
	Agents map[string]Agent
}

type Agent struct {
	MemoryController string
	Memory           map[string]map[string][]string
	Rules            []string
}

type rawConfig struct {
	Version    string                  `yaml:"version"`
	Agents     map[string]rawAgent     `yaml:"agents"`
	Prototypes map[string]rawPrototype `yaml:"prototypes"`
}

type rawAgent struct {
	Prototype    string `yaml:"prototype"`
	rawPrototype `yaml:",inline"`
}

type rawPrototype struct {
	MemoryController string   `yaml:"memorycontroller"`
	Memory           []string `yaml:"memory"`
	Rules            []string `yaml:"rules"`
}
