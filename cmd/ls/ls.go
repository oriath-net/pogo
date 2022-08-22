package ls

import (
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/oriath-net/pogo/cmd"
	"github.com/oriath-net/pogo/poefs"
	"github.com/oriath-net/pogo/util"

	"github.com/spf13/pflag"
)

type listingHandler func(string, fs.FileInfo)

func init() {
	flags := pflag.NewFlagSet("ls", pflag.ExitOnError)
	fLong := flags.CountP("long", "l", "Include more information in the listing")
	fJson := flags.Bool("json", false, "JSON output for nerds")
	fRecurse := flags.BoolP("recurse", "R", false, "Recurse into subdirectories")

	cmd.AddCommand(&cmd.Command{
		Name:        "ls",
		Description: "List files in a GGPK or Steam install",
		Usage:       "pogo ls [options] <Content.ggpk|Steam install>[:path]",

		MinArgs: 1,
		MaxArgs: 1,

		Flags: flags,

		Action: func(args []string) {
			srcPath, localPath := poefs.SplitPath(args[0])
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

			var fileHandler listingHandler
			if *fJson {
				fileHandler = jsonHandler()
			} else if *fLong > 0 {
				fileHandler = longHandler(*fLong)
			} else {
				fileHandler = shortHandler()
			}

			err := fs.WalkDir(srcFs, localPath, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d == nil {
					log.Fatalf("%s doesn't exist", path)
				}

				if path == "." || path == ".." {
					return nil
				}

				di, err := d.Info()
				if err != nil {
					return err
				}

				fileHandler(path, di)

				if !*fRecurse && di.IsDir() && path != localPath {
					return fs.SkipDir
				}
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		},
	})
}

type jsonOutput struct {
	Name   string `json:"name"`
	IsDir  bool   `json:"isdir"`
	Size   int64  `json:"size"`
	Source string `json:"source"`
	Sha256 string `json:"sha256,omitempty"`
}

func jsonHandler() listingHandler {
	return func(path string, di fs.FileInfo) {
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
	}
}

func longHandler(verbosity int) listingHandler {
	return func(path string, di fs.FileInfo) {
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

		if verbosity > 1 {
			if se, ok := di.(poefs.StatExtensions); ok {
				fmt.Printf("%14s %s\n", "src =", se.Provenance())
				if sig := se.Signature(); sig != nil {
					fmt.Printf("%14s %s\n", "sha256 =", hex.EncodeToString(sig))
				}
				fmt.Println("")
			}
		}
	}
}

func shortHandler() listingHandler {
	return func(path string, di fs.FileInfo) {
		if di.IsDir() {
			fmt.Println(path + "/")
		} else {
			fmt.Println(path)
		}
	}
}
