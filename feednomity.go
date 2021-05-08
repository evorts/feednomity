package main

import (
	"github.com/evorts/feednomity/cmd"
	"github.com/evorts/feednomity/pkg/cli"
	"log"
)

func main() {
	commands := cli.NewCli()
	commands.AddCommand("api", cmd.Api)
	commands.AddCommand("web", cmd.Web)
	commands.AddCommand("blaster", cmd.Blaster)
	if err := commands.Listen(); err != nil {
		log.Fatal(err)
	}
}
