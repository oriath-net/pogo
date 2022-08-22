package cmd

import (
	"log"
	"os"

	"github.com/spf13/pflag"
)

type Command struct {
	Name               string
	Aliases            []string
	Description, Usage string

	MinArgs, MaxArgs int
	Flags            *pflag.FlagSet

	Action func(args []string)
}

var commands = []Command{}

func AddCommand(c *Command) {
	commands = append(commands, *c)
}

func find(name string) *Command {
	for _, c := range commands {
		if c.Name == name {
			return &c
		}
		for _, alias := range c.Aliases {
			if alias == name {
				return &c
			}
		}
	}
	return nil
}

func Run() {
	if len(os.Args) < 2 {
		help([]string{})
		os.Exit(1)
	}

	name := os.Args[1]
	c := find(name)
	if c == nil {
		log.Printf("Unknown command '%s'. Run 'pogo help' for a list of commands.", name)
		os.Exit(1)
	}

	pflag.CommandLine.AddFlagSet(GlobalFlags)
	if c.Flags != nil {
		pflag.CommandLine.AddFlagSet(c.Flags)
	}
	pflag.SetInterspersed(true)
	pflag.Parse()

	args := pflag.Args()[1:] // skip command name which will be in the first slot
	if len(args) < c.MinArgs {
		log.Printf("Not enough arguments. Run 'pogo help %s' for help.", c.Name)
		os.Exit(1)
	}
	if len(args) > c.MaxArgs && c.MaxArgs >= 0 { // allow MaxArgs=-1 for unlimited
		log.Printf("Too many arguments. Run 'pogo help %s' for help.", c.Name)
		os.Exit(1)
	}

	c.Action(args)
}
