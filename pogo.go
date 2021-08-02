package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/oriath-net/pogo/cmd"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("pogo: ")

	app := &cli.App{
		Name:  "pogo",
		Usage: "Go tools for Path of Exile",

		Flags: []cli.Flag{},

		Commands: []*cli.Command{
			&cmd.Analyze,
			&cmd.Cat,
			&cmd.Data2json,
			&cmd.Extract,
			&cmd.Ls,
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
