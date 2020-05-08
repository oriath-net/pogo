package main

import (
	"github.com/spf13/pflag"

	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"strings"
)

func findBBBB(data []byte) int {
	return strings.Index(string(data), "\xbb\xbb\xbb\xbb\xbb\xbb\xbb\xbb")
}

func sniffFile(path string) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("cannot read %s: %w", path, err))
	}

	fmt.Println(path)

	boundary := findBBBB(dat)
	if boundary < 0 {
		fmt.Println("\tcannot find boundary")
		return
	}

	var rowCount32 uint32
	err = binary.Read(bytes.NewReader(dat), binary.LittleEndian, &rowCount32)
	rowCount := int(rowCount32)

	rowBytes := (boundary - 4) / rowCount
	remainderBytes := boundary - rowBytes*rowCount - 4
	if remainderBytes == 0 {
		fmt.Printf("\t%d rows, %d bytes per row\n", rowCount, rowBytes)
	} else {
		fmt.Printf("\t%d rows, %d bytes per row + %d unclaimed (WTF?)\n", rowCount, rowBytes, remainderBytes)
	}
	fmt.Printf("\tboundary: %06x\n", boundary)
}

func main() {
	pflag.Parse()
	for _, path := range pflag.Args() {
		sniffFile(path)
	}
}
