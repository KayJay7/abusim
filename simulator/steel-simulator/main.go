package main

import (
	"log"
	"steel-simulator/args"
	"steel-simulator/command"
	"steel-simulator/config"
	"steel-simulator/docker"
)

func main() {
	argsConfig := args.ParseArgs()

	conf, err := config.Parse(argsConfig.ConfigFile)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	dcli, err := docker.New()
	if err != nil {
		log.Fatalf("Error creating docker client: %v", err)
	}

	switch argsConfig.SubCommand {
	case args.SUBCOMMAND_UP:
		command.Up(conf, dcli)
	case args.SUBCOMMAND_DOWN:
		command.Down(conf, dcli)
	}
}
