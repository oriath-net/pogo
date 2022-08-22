package export

import (
	"log"
	"os"
	"path"
	"strconv"

	"github.com/oriath-net/pogo/cmd"
	"github.com/oriath-net/pogo/dat"
	"github.com/oriath-net/pogo/poefs"
	"github.com/oriath-net/pogo/util"

	"github.com/spf13/pflag"
)

func init() {
	flags := pflag.NewFlagSet("export", pflag.ExitOnError)
	fDebug := flags.CountP("debug", "d", "Display warnings and debugging messages while parsing data (use twice for more)")
	fStrict := flags.Bool("strict", false, "Fail on warnings")
	fPretty := flags.Bool("pretty", false, "Pretty-print output")

	cmd.AddCommand(&cmd.Command{
		Name:        "data2json",
		Description: "Convert .dat / .dat64 files to JSON",
		Usage:       "pogo data2json [options] [<Content.ggpk>:]<Data/File.dat> [<row id...>]",

		MinArgs: 1,
		MaxArgs: -1,

		Flags: flags,

		Action: func(args []string) {
			vers, _ := cmd.GlobalFlags.GetString("version")
			p := dat.InitParser(vers)

			fmtDir, _ := cmd.GlobalFlags.GetString("fmt")
			if fmtDir != "" {
				p.SetFormatDir(fmtDir)
			}

			p.SetDebug(*fDebug)

			if *fStrict {
				p.SetStrict(1)
			}

			f, err := poefs.OpenFile(args[0])
			if err != nil {
				log.Fatal(err)
			}

			_, filename := poefs.SplitPath(args[0])
			rows, err := p.Parse(f, path.Base(filename))
			if err != nil {
				log.Fatal(err)
			}

			wantRowIDs := make([]int, 0)
			for _, arg := range args[1:] {
				id, err := strconv.Atoi(arg)
				if err != nil {
					log.Fatalf("Invalid row ID '%s'", arg)
				}
				wantRowIDs = append(wantRowIDs, id)
			}

			if len(wantRowIDs) > 0 {
				for _, i := range wantRowIDs {
					err := util.WriteJson(os.Stdout, rows[i], *fPretty)
					if err != nil {
						log.Fatal(err)
					}
				}
			} else {
				for i := range rows {
					err := util.WriteJson(os.Stdout, rows[i], *fPretty)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		},
	})
}
