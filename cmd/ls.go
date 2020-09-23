package cmd

import (
	"fmt"
	"log"
	"path"
	"strings"

	cli "github.com/urfave/cli/v2"

	"github.com/oriath-net/pogo/ggpk"
)

var Ls = cli.Command{
	Name:      "ls",
	Usage:     "List contents within a single directory in a GGPK",
	UsageText: "pogo ls [options] <Content.ggpk>[:<path>]",

	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "long",
			Aliases: []string{"l"},
			Usage:   "Include more information in the listing",
		},

		&cli.BoolFlag{
			Name:    "longer",
			Aliases: []string{"ll"},
			Usage:   "Whoa, too much information",
		},

		&cli.BoolFlag{
			Name:    "directory",
			Aliases: []string{"d"},
			Usage:   "Don't list the contents of a directory argument",
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

	arg0 := c.Args().First()
	colonIdx := strings.LastIndex(arg0, ":")

	recurse := c.Bool("recurse")
	directory := c.Bool("directory")

	var verbosity int
	if c.Bool("longer") {
		verbosity = 2
	} else if c.Bool("long") {
		verbosity = 1
	} else {
		verbosity = 0
	}

	var ggpkPath string
	var root string
	if colonIdx < 0 {
		ggpkPath = arg0
		root = "/"
	} else {
		ggpkPath = arg0[:colonIdx]
		root = path.Clean("/" + string(arg0[colonIdx+1:]))
	}

	gf, err := ggpk.Open(ggpkPath)
	if err != nil {
		return err
	}

	node, err := gf.NodeAtPath(root)
	if err != nil {
		return err
	}

	if dn, isdir := node.(*ggpk.DirectoryNode); isdir && !recurse && !directory {
		children, err := dn.Children()
		if err != nil {
			log.Fatalf("Unable to list children of %s: %w", root, err)
		}
		for i := range children {
			listNode(children[i], children[i].Name(), verbosity, recurse)
		}
	} else {
		listNode(node, root, verbosity, recurse)
	}

	return nil
}

func listNode(n ggpk.AnyNode, name string, verbosity int, recurse bool) {
	switch verbosity {
	case 0:
		switch n := n.(type) {
		case *ggpk.FileNode:
			fmt.Printf("%s\n", name)
		case *ggpk.DirectoryNode:
			fmt.Printf("%s/\n", name)
		default:
			log.Fatalf("unexpected %T while iterating", n)
		}

	case 1:
		switch n := n.(type) {
		case *ggpk.FileNode:
			fmt.Printf(
				"%s %12d %s\n",
				n.Type(),
				n.Size(),
				name,
			)
		case *ggpk.DirectoryNode:
			fmt.Printf(
				"%s %12d %s/\n",
				n.Type(),
				n.ChildCount(),
				name,
			)
		default:
			log.Fatalf("unexpected %T while iterating", n)
		}

	case 2:
		fmt.Printf(
			"%s %016x:%-12x %s\n",
			n.Type(),
			n.Offset(),
			n.Length(),
			name,
		)
		fmt.Printf("     sha256:%64x\n", n.Signature())
	}

	if dn, isdir := n.(*ggpk.DirectoryNode); isdir && recurse {
		children, err := dn.Children()
		if err != nil {
			log.Fatalf("Unable to list children of %s: %w", name, err)
		}
		for i := range children {
			listNode(
				children[i],
				path.Join(name, children[i].Name()),
				verbosity,
				recurse,
			)
		}
	}
}
