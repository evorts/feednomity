package cmd

import (
	"fmt"
	"github.com/evorts/feednomity/pkg/cli"
)

var Blaster = &cli.Command{
	Description: "Mail blaster command line, cron job style",
	Run: func(cmd *cli.Command, args []string) {
		fmt.Println("mail blaster executed")
	},
}
