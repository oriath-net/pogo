package analyze

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/oriath-net/pogo/cmd"
	"github.com/oriath-net/pogo/poefs"

	"github.com/spf13/pflag"
)

func init() {
	flags := pflag.NewFlagSet("analyze", pflag.ExitOnError)
	flagShort := flags.Bool("short", false, "Display overview only")

	cmd.AddCommand(&cmd.Command{
		Name:        "analyze",
		Description: "Analyze a .dat file",
		Usage:       "pogo analyze [options] [<Content.ggpk>:]<Data/File.dat>",

		MinArgs: 1,
		MaxArgs: -1,

		Flags: flags,

		Action: func(args []string) {
			for _, path := range args {
				f, err := poefs.OpenFile(path)
				if err != nil {
					log.Fatal(err)
				}

				dat, err := ioutil.ReadAll(f)
				if err != nil {
					log.Fatal(err)
				}
				f.Close()

				analyzeData(dat, path, *flagShort)
			}
		},
	})
}

func analyzeData(dat []byte, name string, shortMode bool) {
	var rowCount32 uint32
	err := binary.Read(bytes.NewReader(dat), binary.LittleEndian, &rowCount32)
	if err != nil {
		log.Fatalf("%s: unable to get row count: %s\n\n", name, err)
	}
	rowCount := int(rowCount32)

	fmt.Printf("%s:", name)
	if rowCount == 0 {
		log.Printf("  %d rows, row size indeterminate\n\n", rowCount)
		return
	}

	boundary := findBBBB(dat)
	varBytes := len(dat) - boundary
	varData := dat[boundary:]

	rowBytes := (boundary - 4) / rowCount
	remainderBytes := boundary - rowBytes*rowCount - 4

	if remainderBytes == 0 {
		fmt.Printf("  %d rows, %d bytes per row\n", rowCount, rowBytes)
	} else {
		fmt.Printf("  %d rows, %d bytes per row + %d unclaimed (WTF?)\n", rowCount, rowBytes, remainderBytes)
	}
	fmt.Printf("  boundary at 0x%06x\n", boundary)
	fmt.Printf("  variable data: %d bytes (%.1f per row)\n", varBytes, float64(varBytes)/float64(rowCount))

	fmt.Println("")

	if shortMode {
		return
	}

	hdr := make([]byte, 3*rowBytes)
	sep := strings.Repeat("-- ", rowBytes)

	for i := 0; i < rowBytes; i++ {
		tmp := fmt.Sprintf("%02x", i%256)
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
	analyzeStrings(varData)

	fmt.Println("")
}

func findBBBB(dat []byte) int {
	return strings.Index(string(dat), "\xbb\xbb\xbb\xbb\xbb\xbb\xbb\xbb")
}

func hexdump(data []byte) string {
	if len(data) == 0 {
		return "(empty)"
	}
	s := ""
	h := hex.EncodeToString(data)
	for i := range data {
		s += h[i*2:i*2+2] + " "
	}
	return s[:len(s)-1]
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
		dumper := hex.Dumper(os.Stdout)
		dumper.Write(data)
		dumper.Close()
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
