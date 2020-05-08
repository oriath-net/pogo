package main

import (
	"github.com/duskwuff/pogo"
	"github.com/spf13/pflag"

	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var formats = pflag.StringArrayP("formats", "f", []string{}, "path to a Go configuration file containing formats")

func main() {
	pflag.Parse()

	if len(*formats) != 1 {
		fmt.Fprintf(os.Stderr, "must specify a format file\n")
		os.Exit(1)
	}

	p := pogo.InitParser()
	for _, path := range *formats {
		err := p.LoadFormats(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load formats: %s\n", err)
			os.Exit(1)
		}
	}

	args := pflag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: data2json [options] DataFile.dat [<row IDs>]\n")
		os.Exit(1)
	}

	rows, err := p.ParseFile(pflag.Arg(0), nil)
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
