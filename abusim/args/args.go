package args

import (
	"flag"
	"log"
)

type SubCommand int

const (
	SUBCOMMAND_UP   SubCommand = iota
	SUBCOMMAND_DOWN SubCommand = iota
	SUBCOMMAND_LOGS SubCommand = iota
)

// ArgsConfig represents the configuration given through the command line
type ArgsConfig struct {
	ConfigFile string
	GUI        bool
	GUIPort    int
	GUIImage   string
	SubCommand SubCommand
	Detached   bool
	FollowLogs bool
}

// ParseArgs parses the command line and returns the configuration
func ParseArgs() *ArgsConfig {
	// I create an empty configuration...
	config := ArgsConfig{}
	// ... I set up the global flags...
	flag.StringVar(&config.ConfigFile, "c", "abusim.yml", "configuration file")
	flag.BoolVar(&config.GUI, "g", false, "spawn GUI with simulator")
	flag.IntVar(&config.GUIPort, "gui-port", 8080, "GUI docker port")
	flag.StringVar(&config.GUIImage, "gui-image", "steel-gui", "GUI docker image")
	// ... I set up the "up" subcommand with its flags...
	upCmd := flag.NewFlagSet("up", flag.ExitOnError)
	upCmd.BoolVar(&config.Detached, "d", false, "detached")
	// ... I set up the "down" subcommand...
	downCmd := flag.NewFlagSet("down", flag.ExitOnError)
	// ... and I set up the "logs" subcommand with its flags
	logsCmd := flag.NewFlagSet("logs", flag.ExitOnError)
	logsCmd.BoolVar(&config.FollowLogs, "f", false, "follow")
	// I parse the command line for the global flags...
	flag.Parse()
	// ... I search for a subcommand...
	args := flag.Args()
	if len(args) == 0 {
		log.Fatalln("A command is needed: up, down, logs")
	}
	// ... and I parse the corresponding flags
	switch args[0] {
	case "up":
		upCmd.Parse(args[1:])
		config.SubCommand = SUBCOMMAND_UP
	case "down":
		downCmd.Parse(args[1:])
		config.SubCommand = SUBCOMMAND_DOWN
	case "logs":
		logsCmd.Parse(args[1:])
		config.SubCommand = SUBCOMMAND_LOGS
	default:
		log.Fatalf("unknown subcommand \"%s\", see help for more details", args[0])
	}
	// Finally, I return the configuration
	return &config
}
