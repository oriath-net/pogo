package cmd

import (
	"io"
	"log"
	"os"

	"github.com/oriath-net/pogo/cmd"
	"github.com/oriath-net/pogo/poefs"

	"github.com/spf13/pflag"
	xunicode "golang.org/x/text/encoding/unicode"
)

func init() {
	flags := pflag.NewFlagSet("cat", pflag.ExitOnError)
	utf16 := flags.Bool("utf16", false, "output UTF-16 text files as UTF-8")

	cmd.AddCommand(&cmd.Command{
		Name:        "cat",
		Description: "Extract a file to standard output",
		Usage:       "pogo cat [options] <Content.ggpk|Steam install>:<file>",

		MinArgs: 1,
		MaxArgs: 1,

		Flags: flags,

		Action: func(args []string) {
			var f io.Reader

			f, err := poefs.OpenFile(args[0])
			if err != nil {
				log.Fatal(err)
			}

			if *utf16 {
				f = xunicode.UTF16(xunicode.LittleEndian, xunicode.UseBOM).NewDecoder().Reader(f)
			}

			// Use a 256K buffer to match Oodle compression block size
			buf := make([]byte, 262144)

			_, err = io.CopyBuffer(os.Stdout, f, buf)
			if err != nil {
				log.Fatal(err)
			}
		},
	})
}
