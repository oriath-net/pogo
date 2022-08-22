package cmd

import (
	"log"

	"github.com/spf13/pflag"
)

func init() {
	AddCommand(&Command{
		Name:        "help",
		Description: "display help",
		Usage:       "pogo help",

		Action:  help,
		MaxArgs: 1,
	})
}

func help(args []string) {
	log.SetPrefix("")

	// generic help
	if len(args) == 0 {
		log.Print("Usage:")
		log.Print("\tpogo <command> [options and arguments...]")

		log.Print("")
		log.Print("Commands:")
		for _, c := range commands {
			log.Printf("\t%-15s %s", c.Name, c.Description)
		}

		log.Print()
		log.Print("Global options:")
		pflag.PrintDefaults()
	} else {
		c := find(args[0])
		if c == nil {
			log.Printf("Unknown command '%s'. Run 'pogo help' for a list of commands.", args[0])
			return
		}

		log.Print("Usage:")
		if c.Usage != "" {
			log.Printf("\t%s", c.Usage)
		} else {
			log.Printf("\tpogo %s [options...]", c.Name)
		}

		log.Print()
		log.Print("Global options:")
		GlobalFlags.PrintDefaults()

		if c.Flags != nil {
			log.Print("")
			log.Print("Options:")
			c.Flags.PrintDefaults()
		}
	}
}
