package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	cli "github.com/urfave/cli/v2"

	"github.com/oriath-net/pogo/dat"
	"github.com/oriath-net/pogo/poefs"
)

var Data2json = cli.Command{
	Name:      "data2json",
	Usage:     "Convert .dat files to JSON",
	UsageText: "pogo data2json [options] [<Content.ggpk>:]<Data/File.dat> [<row id...>]",

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "fmt",
			Usage: "path to a Go configuration file containing formats",
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Display warnings while parsing data",
		},
		&cli.BoolFlag{
			Name:  "verbose-debug",
			Usage: "Display insanely verbose debugging messages",
		},
		&cli.StringFlag{
			Name:        "version",
			Usage:       "Path of Exile version to assume for formats",
			DefaultText: "9.99",
		},
	},

	Action: do_data2json,
}

func do_data2json(c *cli.Context) error {
	vers := c.String("version")
	if vers == "" {
		vers = "9.99"
	}

	p := dat.InitParser(vers)
	if c.Bool("verbose-debug") {
		p.SetDebug(2)
	} else if c.Bool("debug") {
		p.SetDebug(1)
	}

	fmtDir := c.String("fmt")
	if fmtDir != "" {
		p.SetFormatDir(fmtDir)
	}

	if !c.Args().Present() {
		return fmt.Errorf("Must specify a data file")
	}

	dat_path := c.Args().First()
	f, err := poefs.OpenFile(dat_path)
	if err != nil {
		return err
	}

	_, filename := poefs.SplitPath(dat_path)
	rows, err := p.Parse(f, path.Base(filename))
	if err != nil {
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)

	wantRowIDs := make([]int, 0)
	for _, arg := range c.Args().Tail() {
		id, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("Invalid row ID '%s'", arg)
		}
		wantRowIDs = append(wantRowIDs, id)
	}

	if len(wantRowIDs) > 0 {
		for _, i := range wantRowIDs {
			err := enc.Encode(rows[i])
			if err != nil {
				return err
			}
		}
	} else {
		for i := range rows {
			err := enc.Encode(rows[i])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
