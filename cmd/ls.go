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

	Flags: []cli.Flag{},

	Action: do_ls,
}

func do_ls(c *cli.Context) error {
	if c.NArg() != 1 {
		cli.ShowCommandHelpAndExit(c, "ls", 1)
	}

	arg0 := c.Args().First()
	colonIdx := strings.LastIndex(arg0, ":")

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

	switch n := node.(type) {
	case *ggpk.FileNode:
		listNodeVerbosely(node, root)
	case *ggpk.DirectoryNode:
		children, err := n.Children()
		if err != nil {
			return err
		}
		for i := range children {
			listNodeVerbosely(children[i], path.Join(root, children[i].Name()))
		}
	default:
		log.Fatalf("unexpected %T node", n)
	}

	return nil
}

func listNodeVerbosely(n ggpk.AnyNode, path string) {
	switch n := n.(type) {
	case *ggpk.FileNode:
		fmt.Printf("F %12d %s\n", n.Size(), path)
	case *ggpk.DirectoryNode:
		fmt.Printf("D %12s %s/\n", "", path)
	default:
		log.Fatalf("unexpected %T while iterating", n)
	}
}
