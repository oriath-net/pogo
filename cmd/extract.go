package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/oriath-net/pogo/poefs"

	cli "github.com/urfave/cli/v2"
)

var Extract = cli.Command{
	Name:      "extract",
	Usage:     "Extract files from an install, including bundled files",
	UsageText: "pogo extract [options] <Content.ggpk|Steam install>[:<dir>]",

	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "print filenames as they're extracted",
		},

		&cli.StringSliceFlag{
			Name:  "include",
			Usage: "only extract files matching this glob",
		},

		&cli.StringSliceFlag{
			Name:  "exclude",
			Usage: "don't extract files matching this glob",
		},

		&cli.BoolFlag{
			Name:  "no-recurse",
			Usage: "don't recurse into subdirectories",
		},

		&cli.StringFlag{
			Name:    "into",
			Aliases: []string{"o"},
			Usage:   "extract into this directory",
			Value:   "out",
		},
	},

	Action: do_extract,
}

func do_extract(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowCommandHelpAndExit(c, "extract", 1)
	}

	srcPath, localPath := poefs.SplitPath(c.Args().First())

	srcFs, err := poefs.Open(srcPath)
	if err != nil {
		log.Fatalf("Unable to open %s: %s", srcPath, err)
	}

	into := c.String("into")
	verbose := c.Bool("verbose")
	noRecurse := c.Bool("no-recurse")
	include := c.StringSlice("include")
	exclude := c.StringSlice("exclude")

	err = fs.WalkDir(srcFs, localPath, func(p string, d fs.DirEntry, err error) error {
		if d == nil {
			log.Fatalf("%s doesn't exist", p)
		}

		if err != nil {
			return err
		}

		di, err := d.Info()
		if err != nil {
			return err
		}

		// check excludes first
		if filterPath(p, exclude, false) {
			if di.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if di.IsDir() {
			if noRecurse && p != localPath {
				fmt.Printf("skipping dir %s\n", p)
				return fs.SkipDir
			}

		} else if filterPath(p, include, true) {
			if verbose {
				fmt.Println(p)
			}

			tgtPath := path.Join(into, p)

			err := os.MkdirAll(path.Dir(tgtPath), 0777)
			if err != nil {
				return err
			}

			tgt, err := os.OpenFile(tgtPath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				return err
			}

			src, err := srcFs.Open(p)
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

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	/*
		f.IterateBundledFiles(prefix, func(p string) error {
			if includePath(c, p) {
				if c.Bool("verbose") {
					fmt.Println(p)
				}

				tgtPath := path.Join(c.String("into"), p)

				err := os.MkdirAll(path.Dir(tgtPath), 0777)
				if err != nil {
					return err
				}

				tgt, err := os.OpenFile(tgtPath, os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					return err
				}

				src, err := f.Open(p)
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
			return nil
		})
	*/

	return nil
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
