package extract

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/oriath-net/pogo/cmd"
	"github.com/oriath-net/pogo/poefs"

	"github.com/spf13/pflag"
)

func init() {
	flags := pflag.NewFlagSet("extract", pflag.ExitOnError)
	fVerbose := cmd.GlobalVerbose
	fIncludes := flags.StringArray("include", []string{}, "Only extract files matching one of these path globs")
	fExcludes := flags.StringArray("exclude", []string{}, "Don't extract files matching any of these path globs")
	fNoRecurse := flags.Bool("no-recurse", false, "Don't recurse into subdirectories")
	fInto := flags.StringP("into", "o", "out", "Extract into this directory")

	cmd.AddCommand(&cmd.Command{
		Name:        "extract",
		Description: "Extract files from an install, including bundled files",
		Usage:       "pogo extract [options] <Content.ggpk|Steam install>[:<dir>]",

		Flags:   flags,
		MinArgs: 1,
		MaxArgs: 1,

		Action: func(args []string) {
			srcPath, localPath := poefs.SplitPath(args[0])

			srcFs, err := poefs.Open(srcPath)
			if err != nil {
				log.Fatalf("Unable to open %s: %s", srcPath, err)
			}

			err = fs.WalkDir(srcFs, localPath, func(p string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d == nil {
					log.Fatalf("%s doesn't exist", p)
				}

				if p == "." || p == ".." {
					return nil
				}

				di, err := d.Info()
				if err != nil {
					return err
				}

				// check excludes first
				if filterPath(p, *fExcludes, false) {
					if *fVerbose > 1 {
						log.Printf("skipping excluded %s", p)
					}
					if di.IsDir() {
						return fs.SkipDir
					}
					return nil
				} else if di.IsDir() {
					if *fNoRecurse && p != localPath {
						if *fVerbose > 1 {
							log.Printf("not recursing into directory %s", p)
						}
						return fs.SkipDir
					}
					return nil
				} else if !filterPath(p, *fIncludes, true) {
					if *fVerbose > 1 {
						log.Printf("skipping non-included %s", p)
					}
					return nil
				} else {
					if *fVerbose > 0 {
						fmt.Println(p)
					}
					src, err := srcFs.Open(p)
					if err != nil {
						return err
					}
					return extractFile(src, path.Join(*fInto, p))
				}
			})
			if err != nil {
				log.Fatal(err)
			}
		},
	})
}

func filterPath(p string, filters []string, fallback bool) bool {
	if len(filters) == 0 {
		return fallback
	}

	for _, f := range filters {
		matched, err := path.Match(f, p)
		if err != nil {
			log.Fatal(err)
		}
		if matched {
			return true
		}
	}

	return false
}

func extractFile(src io.Reader, tgtPath string) error {
	err := os.MkdirAll(path.Dir(tgtPath), 0o777)
	if err != nil {
		return err
	}

	tgt, err := os.OpenFile(tgtPath, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		return err
	}

	buf := make([]byte, 262144)
	_, err = io.CopyBuffer(tgt, src, buf)
	if err != nil {
		return err
	}

	err = tgt.Close()
	return err
}
