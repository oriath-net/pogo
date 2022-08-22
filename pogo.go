package main

import (
	"log"

	"github.com/oriath-net/pogo/cmd"

	_ "github.com/oriath-net/pogo/cmd/analyze"
	_ "github.com/oriath-net/pogo/cmd/cat"
	_ "github.com/oriath-net/pogo/cmd/export"
	_ "github.com/oriath-net/pogo/cmd/extract"
	_ "github.com/oriath-net/pogo/cmd/ls"
	_ "github.com/oriath-net/pogo/cmd/schema2json"
	_ "github.com/oriath-net/pogo/cmd/tidy"
	_ "github.com/oriath-net/pogo/cmd/validate"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("pogo: ")

	cmd.Run()
}
