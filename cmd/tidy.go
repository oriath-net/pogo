package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/oriath-net/pogo/dat"
	"github.com/oriath-net/pogo/util"

	cli "github.com/urfave/cli/v2"
)

var Tidy = cli.Command{
	Name:      "tidy",
	Usage:     "Clean up and validate one or more schema JSON files",
	UsageText: "pogo tidy dat/formats/*.json",

	Action: do_tidy,
}

func do_tidy(c *cli.Context) error {
	if !c.Args().Present() {
		return errNotEnoughArguments
	}

	for i := 0; i < c.NArg(); i++ {
		arg := c.Args().Get(i)
		if filepath.Ext(arg) != ".json" {
			log.Printf("skipping non-json file %s", arg)
			continue
		}

		jfmt := dat.JsonFormat{}
		err := util.ReadJsonFromFile(arg, &jfmt)
		if err != nil {
			return fmt.Errorf("failed to load %s: %w", arg, err)
		}

		buf := bytes.Buffer{}
		err = util.WriteJson(&buf, &jfmt, true)
		if err != nil {
			return fmt.Errorf("failed to reserialize %s: %w", arg, err)
		}

		os.Rename(arg, arg+"~")

		err = os.WriteFile(arg, buf.Bytes(), 0666)
		if err != nil {
			return fmt.Errorf("failed to reserialize %s: %w", arg, err)
		}
	}

	return nil
}
