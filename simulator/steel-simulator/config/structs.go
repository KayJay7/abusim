package config

type rawConfig struct {
	Version    string                  `yaml:"version"`
	Image      string                  `yaml:"image"`
	Namespace  string                  `yaml:"namespace"`
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
