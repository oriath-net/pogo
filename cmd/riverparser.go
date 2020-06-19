package cmd

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/oriath-net/pogo/river"
	cli "github.com/urfave/cli/v2"
)

var RiverParser = cli.Command{
	Name: "riverparser",

	// This is more of a test than a practical example
	Hidden: true,
	//Usage
	//UsageText

	Flags: []cli.Flag{},

	Action: do_riverparser,
}

func do_riverparser(c *cli.Context) error {
	var err error

	filename := c.Args().Get(0)

	var r io.Reader

	r, err = os.Open(filename)
	if err != nil {
		return err
	}

	if strings.HasSuffix(filename, ".gz") {
		rz, err := gzip.NewReader(r)
		if err != nil {
			return err
		}
		r = rz
	}

	d := json.NewDecoder(r)

	river := river.RiverOutput{}
	err = d.Decode(&river)
	if err != nil {
		return err
	}

	fmt.Printf("%#v\n", river)

	return nil
}
