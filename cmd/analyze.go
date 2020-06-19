package cmd

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	cli "github.com/urfave/cli/v2"
)

var Analyze = cli.Command{
	Name:      "analyze",
	Usage:     "Analyze a .dat file",
	UsageText: "pogo analyze [options] [<Content.ggpk>:]<Data/File.dat>",

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
		varData := dat[boundary:]

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

		fmt.Println("")

		hdr := make([]byte, 3*rowBytes)
		sep := strings.Repeat("-- ", rowBytes)

		for i := 0; i < rowBytes; i++ {
			tmp := fmt.Sprintf("%02d", i%100)
			hdr[i*3+0] = tmp[0]
			hdr[i*3+1] = tmp[1]
			hdr[i*3+2] = ' '
		}

		colMax := make([]byte, rowBytes)

		fmt.Printf("%6s  %s\n", "", hdr)
		fmt.Printf("%6s  %s\n", "", sep)

		for i := 0; i < rowCount; i++ {
			row := dat[4+i*rowBytes : 4+(i+1)*rowBytes]
			fmt.Printf("%6d: %s\n", i, hexdump(row))
			for j := 0; j < rowBytes; j++ {
				if row[j] > colMax[j] {
					colMax[j] = row[j]
				}
			}
		}

		fmt.Printf("%6s  %s\n", "", sep)
		fmt.Printf("%6s  %s\n", "", hdr)
		fmt.Printf("%6s: %s\n", "MAX", hexdump(colMax))

		fmt.Printf("\n\nVariable data:\n")
		if len(varData) > 8 {
			analyzeStrings(varData)
		} else {
			dumper := hex.Dumper(os.Stdout)
			dumper.Write(varData)
			dumper.Close()
		}

		fmt.Println("")
	}

	return nil
}

func hexdump(data []byte) string {
	if len(data) == 0 {
		return "(empty)"
	}
	s := make([]byte, 3*len(data))
	h := hex.EncodeToString(data)
	for i := range data {
		s[i*3] = h[i*2]
		s[i*3+1] = h[i*2+1]
		s[i*3+2] = ' '
	}
	return string(s[:len(s)-1])
}

func analyzeStrings(data []byte) {
	p := 0
	lastStringEnd := -1

	for p < len(data) {
		ok, str := isStringAt(data, p)
		if ok {
			if p > lastStringEnd && lastStringEnd != -1 {
				fmt.Printf("%06x: %s\n", lastStringEnd, hexdump(data[lastStringEnd:p]))
			}
			fmt.Printf("%06x: %q\n", p, str)
			p += len(str)*2 + 4
			lastStringEnd = p
		} else {
			p += 1
		}
	}

	if lastStringEnd < 0 {
		fmt.Println("No strings found")
	} else if p > lastStringEnd {
		fmt.Printf("%06x: %s\n", lastStringEnd, hexdump(data[lastStringEnd:]))
	}
}

func isStringAt(data []byte, offset int) (bool, string) {
	p := offset
	for {
		if p+3 >= len(data) {
			return false, ""
		}
		c, c2, c3, c4 := data[p], data[p+1], data[p+2], data[p+3]
		if c == 0 && c2 == 0 && c3 == 0 && c4 == 0 {
			// end of string
			break
		}
		if c2 != 0 {
			// high byte not zero
			return false, ""
		}
		if (c < 0x20 && c != '\n' && c != '\r' && c != '\t') || c > 0x7e {
			// non-ASCII low byte
			return false, ""
		}
		p += 2
	}
	length := (p - offset) / 2
	if length < 2 {
		// too short (length=1 makes small integers show up as strings)
		return false, ""
	}
	sa := make([]byte, length)
	for i := range sa {
		sa[i] = data[offset+2*i]
	}
	return true, string(sa)
}
