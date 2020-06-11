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
	UsageText: "pogo ls [command options] <Content.ggpk:Data/>",

	Flags: []cli.Flag{},

	Action: do_ls,
}

func do_ls(c *cli.Context) error {
	if c.NArg() != 1 {
		cli.ShowCommandHelpAndExit(c, "ls", 1)
	}

	parts := strings.SplitN(c.Args().First(), ":", 2)

	gf, err := ggpk.Open(parts[0])
	if err != nil {
		return err
	}

	var root = ""
	if len(parts) > 1 {
		root = path.Clean("/" + parts[1])
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
