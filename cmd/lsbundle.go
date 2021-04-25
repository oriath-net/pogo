package cmd

import (
	"fmt"

	cli "github.com/urfave/cli/v2"

	"github.com/oriath-net/pogo/ggpk"
)

var LsBundle = cli.Command{
	Name:      "lsbundle",
	Usage:     "List the contents of all bundles (in no particular order)",
	UsageText: "pogo lsbundle <Content.ggpk>",

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "prefix",
			Usage: "Only list files starting with this prefix",
		},
	},

	Action: do_lsbundle,
}

func do_lsbundle(c *cli.Context) error {
	if c.NArg() != 1 {
		cli.ShowCommandHelpAndExit(c, "lsbundle", 1)
	}

	prefix := c.String("prefix")

	f, err := ggpk.Open(c.Args().First())
	if err != nil {
		return fmt.Errorf("failed to open GGPK: %s", err)
	}

	f.DumpBundleIndex(prefix)

	return nil
}
