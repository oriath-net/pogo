package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	cli "github.com/urfave/cli/v2"

	"github.com/oriath-net/pogo/dat"
)

var Data2json = cli.Command{
	Name:      "data2json",
	Usage:     "Convert .dat files to JSON",
	UsageText: "pogo data2json [command options] [<Content.ggpk>:]<Data/File.dat> [<row id...>]",

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "fmt",
			Usage:    "path to a Go configuration file containing formats",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable format debugging log messages",
		},
	},

	Action: do_data2json,
}

func do_data2json(c *cli.Context) error {
	p := dat.InitParser()

	if c.Bool("debug") {
		p.EnableDebug()
	}

	err := p.LoadFormats(c.String("fmt"))
	if err != nil {
		return fmt.Errorf("Failed to load formats: %s", err)
	}

	if !c.Args().Present() {
		return fmt.Errorf("Must specify a data file")
	}

	dat_path := c.Args().First()
	f, err := openGgpkPath(dat_path)
	if err != nil {
		return err
	}

	basename := strings.TrimSuffix(path.Base(dat_path), ".dat")
	rows, err := p.Parse(f, basename)
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
