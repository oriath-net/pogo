package cmd

import (
	"io"
	"os"

	"github.com/oriath-net/pogo/poefs"

	cli "github.com/urfave/cli/v2"
	xunicode "golang.org/x/text/encoding/unicode"
)

var Cat = cli.Command{
	Name:      "cat",
	Usage:     "Extract a file from a GGPK to standard output",
	UsageText: "pogo cat [options] <Content.ggpk|Steam install>:<file>",

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

	var f io.Reader

	f, err := poefs.OpenFile(c.Args().First())
	if err != nil {
		return err
	}

	if c.Bool("utf16") {
		f = xunicode.UTF16(xunicode.LittleEndian, xunicode.UseBOM).NewDecoder().Reader(f)
	}

	// Use a 256K buffer to match Oodle compression block size
	buf := make([]byte, 262144)

	_, err = io.CopyBuffer(os.Stdout, f, buf)
	return err
}
