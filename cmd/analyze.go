package cmd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"strings"

	cli "github.com/urfave/cli/v2"
)

var Analyze = cli.Command{
	Name:      "analyze",
	Usage:     "Analyze a .dat file",
	UsageText: "pogo analyze [command options] [<Content.ggpk>:]<Data/File.dat>",

	Flags: []cli.Flag{},

	Action: do_analyze,
}

func findBBBB(data []byte) int {
	return strings.Index(string(data), "\xbb\xbb\xbb\xbb\xbb\xbb\xbb\xbb")
}

func do_analyze(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Must specify a data file")
	}

	for _, path := range c.Args().Slice() {
		f, err := openGgpkPath(path)
		if err != nil {
			return fmt.Errorf("%s: %s", path, err)
		}

		dat, err := ioutil.ReadAll(f)
		if err != nil {
			return fmt.Errorf("%s: %s", path, err)
		}

		var rowCount32 uint32
		err = binary.Read(bytes.NewReader(dat), binary.LittleEndian, &rowCount32)
		if err != nil {
			return fmt.Errorf("%s: unable to get row count: %s", path, err)
		}
		rowCount := int(rowCount32)

		boundary := findBBBB(dat)
		varBytes := len(dat) - boundary

		rowBytes := (boundary - 4) / rowCount
		remainderBytes := boundary - rowBytes*rowCount - 4

		fmt.Println(path)
		if remainderBytes == 0 {
			fmt.Printf("  %d rows, %d bytes per row\n", rowCount, rowBytes)
		} else {
			fmt.Printf("  %d rows, %d bytes per row + %d unclaimed (WTF?)\n", rowCount, rowBytes, remainderBytes)
		}
		fmt.Printf("  boundary at 0x%06x\n", boundary)
		fmt.Printf("  variable data: %d bytes (%.1f per row)\n", varBytes, float64(varBytes)/float64(rowCount))
	}

	return nil
}
