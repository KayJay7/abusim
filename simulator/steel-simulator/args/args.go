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

type ArgsConfig struct {
	ConfigFile string
	SubCommand SubCommand
	Detached   bool
	FollowLogs bool
}

func ParseArgs() *ArgsConfig {
	config := ArgsConfig{}

	flag.StringVar(&config.ConfigFile, "c", "steel-simulator.yml", "configuration file")

	upCmd := flag.NewFlagSet("up", flag.ExitOnError)
	upCmd.BoolVar(&config.Detached, "d", false, "detached")

	downCmd := flag.NewFlagSet("down", flag.ExitOnError)

	logsCmd := flag.NewFlagSet("logs", flag.ExitOnError)
	logsCmd.BoolVar(&config.FollowLogs, "f", false, "follow")

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatalln("A command is needed: up, down, logs")
	}
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

	return &config
}
