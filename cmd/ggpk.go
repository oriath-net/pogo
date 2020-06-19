package cmd

import (
	"fmt"
	cli "github.com/urfave/cli/v2"
	"io"
	"log"
	"os"
	"path"

	"github.com/oriath-net/pogo/ggpk"
)

var Ggpk = cli.Command{
	Name:      "ggpk",
	Usage:     "List or extract files in a GGPK archive",
	UsageText: "ggpk <subcommand> [options] <Content.ggpk> [<path...>]",

	Subcommands: []*cli.Command{
		{
			Name:      "list",
			Aliases:   []string{"l"},
			Usage:     "List files within the archive",
			UsageText: "ggpk list [options...] <Content.ggpk> [<path...>]",

			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "verbose",
					Aliases: []string{"v"},
					Usage:   "display extra details for each file",
				},

				&cli.BoolFlag{
					Name:    "no-recurse",
					Aliases: []string{"n"},
					Usage:   "don't recurse into subdirectories",
				},

				&cli.StringSliceFlag{
					Name:  "exclude",
					Usage: "exclude these paths",
				},
			},
			Action: func(c *cli.Context) error { return runGgpkCommon(c, do_list) },
		},

		{
			Name:      "extract",
			Aliases:   []string{"x"},
			Usage:     "extract files from an archive",
			UsageText: "ggpk extract [options] <Content.ggpk> [<path...>]",

			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "verbose",
					Aliases: []string{"v"},
					Usage:   "display extra details for each file",
				},

				&cli.BoolFlag{
					Name:    "no-recurse",
					Aliases: []string{"n"},
					Usage:   "don't recurse into subdirectories",
				},

				&cli.StringSliceFlag{
					Name:  "exclude",
					Usage: "exclude these paths",
				},

				&cli.StringFlag{
					Name:    "into",
					Aliases: []string{"o"},
					Usage:   "extract into this directory",
					Value:   "out",
				},

				&cli.BoolFlag{
					Name:  "stdout",
					Usage: "extract files to standard output",
				},
			},
			Action: func(c *cli.Context) error { return runGgpkCommon(c, do_extract) },
		},
	},
}

type ggpkRuntimeContext struct {
	ggpk *ggpk.File

	excludePaths []string
	intoPath     string
	noRecurse    bool
	toStdout     bool
	verbose      bool
}

type fileVisitor func(rc *ggpkRuntimeContext, path string, node ggpk.AnyNode) error

func runGgpkCommon(c *cli.Context, fn fileVisitor) error {
	if c.NArg() < 1 {
		log.Println("Need a GGPK filename")
		cli.ShowCommandHelpAndExit(c, "", 1)
	}
	f, err := ggpk.Open(c.Args().First())
	if err != nil {
		return fmt.Errorf("failed to open GGPK: %s", err)
	}
	rc := ggpkRuntimeContext{
		ggpk:         f,
		excludePaths: c.StringSlice("exclude"),
		intoPath:     c.String("into"),
		noRecurse:    c.Bool("no-recurse"),
		toStdout:     c.Bool("stdout"),
		verbose:      c.Bool("verbose"),
	}
	var paths = c.Args().Tail()
	if len(paths) == 0 {
		paths = []string{""}
	}
	iterate(&rc, fn, paths)
	return nil
}

func excluded(rc *ggpkRuntimeContext, path string) bool {
	for _, xpath := range rc.excludePaths {
		if path == xpath {
			return true
		}
	}
	return false
}

func iterate(rc *ggpkRuntimeContext, fn fileVisitor, paths []string) {
	for i := range paths {
		rootPath := path.Clean("/" + paths[i])
		rootNode, err := rc.ggpk.NodeAtPath(rootPath)
		if err != nil {
			log.Fatal(err)
		}

		switch n := rootNode.(type) {
		case *ggpk.FileNode:
			err := fn(rc, rootPath, rootNode)
			if err != nil {
				log.Fatalf("%s: %s", rootPath, err)
			}

		case *ggpk.DirectoryNode:
			iter := n.Iterator()
			for iter.Next() {
				fullPath := path.Join(rootPath, iter.Path())
				if excluded(rc, fullPath) {
					iter.Skip()
					continue
				}
				err := fn(rc, fullPath, iter.Node())
				if err != nil {
					log.Fatalf("%s: %s", fullPath, err)
				}
				if rc.noRecurse {
					iter.Skip()
				}
			}
			if err := iter.Error(); err != nil {
				log.Fatal(err)
			}

		default:
			log.Fatalf("unexpected node type %T", n)
		}
	}
}

func do_list(rc *ggpkRuntimeContext, path string, node ggpk.AnyNode) error {
	if rc.verbose {
		listNodeVerbosely(node, path)
	} else {
		switch n := node.(type) {
		case *ggpk.FileNode:
			fmt.Printf("%s\n", path)
		case *ggpk.DirectoryNode:
			fmt.Printf("%s/\n", path)
		default:
			log.Fatalf("unexpected %T while iterating", n)
		}
	}
	return nil
}

func do_extract(rc *ggpkRuntimeContext, srcpath string, node ggpk.AnyNode) error {
	var tgt io.WriteCloser

	fileNode, ok := node.(*ggpk.FileNode)
	if !ok {
		return nil
	}

	if rc.verbose {
		fmt.Fprintf(os.Stderr, "%s\n", srcpath)
	}

	if rc.toStdout {
		tgt = os.Stdout

	} else {
		truepath := path.Join(rc.intoPath, srcpath)

		err := os.MkdirAll(path.Dir(truepath), 0777)
		if err != nil {
			return err
		}

		tgt, err = os.OpenFile(truepath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
	}

	_, err := io.Copy(tgt, fileNode.Reader())
	if err != nil {
		return err
	}

	if !rc.toStdout {
		err := tgt.Close()
		return err
	}

	return nil
}
