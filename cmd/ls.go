package cmd

import (
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/oriath-net/pogo/poefs"
	"github.com/oriath-net/pogo/util"

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
			Name:    "longer",
			Aliases: []string{"ll"},
			Usage:   "Include way too much information in the listing",
		},

		&cli.BoolFlag{
			Name:  "json",
			Usage: "JSON output for nerds",
		},

		&cli.BoolFlag{
			Name:    "recurse",
			Aliases: []string{"R"},
			Usage:   "Recurse into subdirectories",
		},
	},

	Action: do_ls,
}

type jsonOutput struct {
	Name   string `json:"name"`
	IsDir  bool   `json:"isdir"`
	Size   int64  `json:"size"`
	Source string `json:"source"`
	Sha256 string `json:"sha256,omitempty"`
}

func do_ls(c *cli.Context) error {
	if c.NArg() != 1 {
		cli.ShowCommandHelpAndExit(c, "ls", 1)
	}

	srcPath, localPath := poefs.SplitPath(c.Args().First())

	var srcFs fs.FS
	if srcPath == "" {
		if filepath.IsAbs(localPath) {
			srcFs = os.DirFS("/")
			localPath = localPath[1:]
		} else {
			srcFs = os.DirFS(".")
		}
	} else {
		var err error
		srcFs, err = poefs.Open(srcPath)
		if err != nil {
			log.Fatalf("Unable to open %s: %s", srcPath, err)
		}
	}

	doJson := c.Bool("json")
	longer := c.Bool("longer")
	long := longer || c.Bool("long")
	recurse := c.Bool("recurse")

	err := fs.WalkDir(srcFs, localPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d == nil {
			log.Fatalf("%s doesn't exist", path)
		}

		di, err := d.Info()
		if err != nil {
			return err
		}

		if doJson {
			jout := jsonOutput{
				Name:  path,
				IsDir: di.IsDir(),
				Size:  di.Size(),
			}
			if se, ok := di.(poefs.StatExtensions); ok {
				jout.Source = se.Provenance()
				if sig := se.Signature(); sig != nil {
					jout.Sha256 = hex.EncodeToString(sig)
				}
			}
			util.WriteJson(os.Stdout, &jout, false)
		} else if long {
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
			if longer {
				if se, ok := di.(poefs.StatExtensions); ok {
					fmt.Printf("%14s %s\n", "src =", se.Provenance())
					if sig := se.Signature(); sig != nil {
						fmt.Printf("%14s %s\n", "sha256 =", hex.EncodeToString(sig))
					}
					fmt.Println("")
				}
			}
		} else {
			if di.IsDir() {
				fmt.Println(path + "/")
			} else {
				fmt.Println(path)
			}
		}

		if !recurse && di.IsDir() && path != localPath {
			return fs.SkipDir
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
