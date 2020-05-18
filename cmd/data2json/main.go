package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/duskwuff/pogo/dat"
	"github.com/duskwuff/pogo/ggpk"

	flag "github.com/spf13/pflag"
)

var formats = flag.StringArrayP("format", "f", []string{}, "path to a Go configuration file containing formats")

func main() {
	flag.Parse()

	if len(*formats) != 1 {
		fmt.Fprintf(os.Stderr, "must specify a format file\n")
		os.Exit(1)
	}

	p := dat.InitParser()
	for _, path := range *formats {
		err := p.LoadFormats(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load formats: %s\n", err)
			os.Exit(1)
		}
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: data2json [options] DataFile.dat [<row IDs>]\n")
		os.Exit(1)
	}

	f, err := openPath(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	basename := strings.TrimSuffix(path.Base(flag.Arg(0)), ".dat")

	rows, err := p.Parse(f, basename)
	if err != nil {
		panic(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)

	wantRowIDs := make([]int, 0)
	for i := 2; i < len(args); i++ {
		id, err := strconv.Atoi(args[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "usage: data2json [options] DataFile.dat [<row IDs>]\n")
			os.Exit(1)
		}
		wantRowIDs = append(wantRowIDs, id)
	}

	if len(wantRowIDs) > 0 {
		for _, i := range wantRowIDs {
			err := enc.Encode(rows[i])
			if err != nil {
				panic(err)
			}
		}
	} else {
		for i := range rows {
			err := enc.Encode(rows[i])
			if err != nil {
				panic(err)
			}
		}
	}
}

func openPath(path string) (io.Reader, error) {
	parts := strings.Split(path, ":")
	switch len(parts) {
	case 1:
		return os.Open(path)

	case 2:
		gf, err := ggpk.Open(parts[0])
		if err != nil {
			return nil, fmt.Errorf("couldn't open GGPK: %w", err)
		}
		return gf.Open(parts[1])

	default:
		return nil, fmt.Errorf("%s: too many colons (use fewer)", path)
	}
}
