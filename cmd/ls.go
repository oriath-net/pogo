package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/oriath-net/pogo/poefs"

	cli "github.com/urfave/cli/v2"
)

var Ls = cli.Command{
	Name:      "ls",
	Usage:     "List files present in an install",
	UsageText: "pogo ls [options] <Content.ggpk|Steam install>[:<path>]",

	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "long",
			Aliases: []string{"l"},
			Usage:   "Include more information in the listing",
		},

		&cli.BoolFlag{
			Name:    "recurse",
			Aliases: []string{"R"},
			Usage:   "Recurse into subdirectories",
		},
	},

	Action: do_ls,
}

func do_ls(c *cli.Context) error {
	if c.NArg() != 1 {
		cli.ShowCommandHelpAndExit(c, "ls", 1)
	}

	srcPath, localPath := poefs.SplitPath(c.Args().First())

	var srcFs fs.FS
	if srcPath == "" {
		srcFs = os.DirFS(".")
	} else {
		var err error
		srcFs, err = poefs.Open(srcPath)
		if err != nil {
			log.Fatalf("Unable to open %s: %s", srcPath, err)
		}
	}

	long := c.Bool("long")
	recurse := c.Bool("recurse")

	fs.WalkDir(srcFs, localPath, func(path string, d fs.DirEntry, err error) error {
		if d == nil {
			log.Fatalf("%s doesn't exist", path)
		}

		di, err := d.Info()
		if err != nil {
			return err
		}

		if !long {
			if di.IsDir() {
				fmt.Println(path + "/")
			} else {
				fmt.Println(path)
			}
		} else {
			if di.IsDir() {
				fmt.Printf(
					"%s %12s %s\n",
					"D",
					"",
					path+"/",
				)
			} else {
				fmt.Printf(
					"%s %12d %s\n",
					"F",
					di.Size(),
					path,
				)
			}
		}

		if !recurse && di.IsDir() && path != localPath {
			return fs.SkipDir
		}
		return nil
	})

	return nil
}
