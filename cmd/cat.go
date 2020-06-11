package cmd

import (
	"io"
	"os"

	cli "github.com/urfave/cli/v2"
)

var Cat = cli.Command{
	Name:      "cat",
	Usage:     "Extract a file from a GGPK to standard output",
	UsageText: "pogo cat [command options] <Content.ggpk:Data/File.dat>",

	Flags: []cli.Flag{},

	Action: do_cat,
}

func do_cat(c *cli.Context) error {
	if c.NArg() != 1 {
		cli.ShowCommandHelpAndExit(c, "cat", 1)
	}

	f, err := openGgpkPath(c.Args().First())
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, f)
	return err
}
