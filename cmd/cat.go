package cmd

import (
	"io"
	"os"

	cli "github.com/urfave/cli/v2"
	xunicode "golang.org/x/text/encoding/unicode"
)

var Cat = cli.Command{
	Name:      "cat",
	Usage:     "Extract a file from a GGPK to standard output",
	UsageText: "pogo cat [options] <Content.ggpk>:<Data/File.dat>",

	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "utf16",
			Usage: "output UTF-16 text files as UTF-8",
		},
	},

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

	if c.Bool("utf16") {
		f = xunicode.UTF16(xunicode.LittleEndian, xunicode.UseBOM).NewDecoder().Reader(f)
	}

	_, err = io.Copy(os.Stdout, f)
	return err
}
