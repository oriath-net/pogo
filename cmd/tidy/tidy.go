package tidy

import (
	"bytes"
	"log"
	"os"
	"path/filepath"

	"github.com/oriath-net/pogo/cmd"
	"github.com/oriath-net/pogo/dat"
	"github.com/oriath-net/pogo/util"
)

func init() {
	cmd.AddCommand(&cmd.Command{
		Name:        "tidy",
		Description: "Clean up and validate one or more schema JSON files",
		Usage:       "pogo tidy dat/formats/*.json",

		MinArgs: 1,
		MaxArgs: -1,

		Action: func(args []string) {
			for _, f := range args {
				tidyFile(f)
			}
		},
	})
}

func tidyFile(f string) {
	if filepath.Ext(f) != ".json" {
		log.Printf("skipping non-json file %s", f)
		return
	}

	jfmt := dat.JsonFormat{}
	err := util.ReadJsonFromFile(f, &jfmt)
	if err != nil {
		log.Fatalf("failed to load %s: %w", f, err)
	}

	buf := bytes.Buffer{}
	err = util.WriteJson(&buf, &jfmt, true)
	if err != nil {
		log.Fatalf("failed to reserialize %s: %w", f, err)
	}

	os.Rename(f, f+"~")

	err = os.WriteFile(f, buf.Bytes(), 0o666)
	if err != nil {
		log.Fatalf("failed to reserialize %s: %w", f, err)
	}
}
