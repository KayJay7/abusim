package main

import (
	"log"
	"steel-simulator/args"
	"steel-simulator/command"
	"steel-simulator/config"
)

func main() {
	argsConfig := args.ParseArgs()

	config, err := config.Parse(argsConfig.ConfigFile)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	switch argsConfig.SubCommand {
	case args.SUBCOMMAND_UP:
		command.Up(config)
	case args.SUBCOMMAND_DOWN:
		command.Down()
	}
}
