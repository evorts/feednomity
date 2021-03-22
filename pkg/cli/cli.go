package cli

import (
	"flag"
	"github.com/pkg/errors"
)

type Command struct {
	Label string
	Description string
	Run  func(cmd *Command, args []string)
}

type manager struct {
	commands map[string]*Command
}

type IManager interface {
	AddCommand(key string, command *Command)
	Listen() error
}

func NewCli() IManager {
	return &manager{
		commands: make(map[string]*Command, 0),
	}
}

func (c *manager) AddCommand(key string, command *Command) {
	c.commands[key] = command
}

func (c *manager) GetCommand(key string) *Command {
	if cc, ok := c.commands[key]; ok {
		return cc
	}
	return nil
}

func (c *manager) commandExist(key string) bool {
	if _, ok := c.commands[key]; ok {
		return true
	}
	return false
}

func (c *manager) Listen() error {
	// get arguments
	flag.Parse()
	a := flag.Args()
	if len(a) < 1 {
		return errors.New("please provide main command key")
	}
	key := a[0]
	args := make([]string, 0)
	if len(a) > 1 {
		args = a[1:]
	}
	if !c.commandExist(key) {
		return errors.New("command not exist")
	}
	cmd := c.GetCommand(key)
	cmd.Run(cmd, args)
	return nil
}
