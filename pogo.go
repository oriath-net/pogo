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
			&cmd.Ggpk,
			&cmd.Data2json,
			&cmd.Analyze,
			&cmd.Cat,
			&cmd.Ls,
			&cmd.Diff,
			&cmd.JsonSchema,
			&cmd.RiverParser,
			&cmd.RiverDl,
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
