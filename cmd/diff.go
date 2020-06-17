package cmd

import (
	"bytes"
	"fmt"
	"path"

	cli "github.com/urfave/cli/v2"

	"github.com/oriath-net/pogo/ggpk"
)

var Diff = cli.Command{
	Name:      "diff",
	Usage:     "Show changes between two GGPKs",
	UsageText: "pogo diff [command options] <Content1.ggpk> <Content2.ggpk>",

	Flags: []cli.Flag{},

	Action: do_diff,
}

func do_diff(c *cli.Context) error {
	if c.NArg() != 2 {
		cli.ShowCommandHelpAndExit(c, "diff", 1)
	}

	gf1, err := ggpk.Open(c.Args().Get(0))
	if err != nil {
		return err
	}
	r1, err := gf1.RootNode()
	if err != nil {
		return err
	}

	gf2, err := ggpk.Open(c.Args().Get(1))
	if err != nil {
		return err
	}
	r2, err := gf2.RootNode()
	if err != nil {
		return err
	}

	diff_at("/", r1, r2)

	return nil
}

func diff_at(root string, n1, n2 *ggpk.DirectoryNode) error {
	if bytes.Equal(n1.Signature(), n2.Signature()) {
		return nil
	}

	ca1, err := n1.Children()
	if err != nil {
		return err
	}

	ca2, err := n2.Children()
	if err != nil {
		return err
	}

	c1map := make(map[string]ggpk.AnyNode)
	for _, c := range ca1 {
		c1map[c.Name()] = c
	}

	for _, c2 := range ca2 {
		c1, exists := c1map[c2.Name()]
		if !exists {
			show_added(root, c2)
		} else {
			delete(c1map, c2.Name())
			if !bytes.Equal(c1.Signature(), c2.Signature()) {
				c1_dir, c1_isdir := c1.(*ggpk.DirectoryNode)
				c2_dir, c2_isdir := c2.(*ggpk.DirectoryNode)
				if c1_isdir && c2_isdir {
					diff_at(path.Join(root, c1.Name()), c1_dir, c2_dir)
				} else if !c1_isdir && !c2_isdir {
					show_changed_file(root, c1.(*ggpk.FileNode), c2.(*ggpk.FileNode))
				} else {
					// weird case that I'm not sure ever comes up -- file changed to directory or vice versa
					show_removed(root, c1)
					show_added(root, c2)
				}
			}
		}
	}

	for _, c1 := range c1map {
		show_removed(root, c1)
	}

	return nil
}

func show_added(root string, n ggpk.AnyNode) error {
	p := path.Join(root, n.Name())

	switch n := n.(type) {
	case *ggpk.FileNode:
		fmt.Printf("+ F %20d %s\n", n.Size(), p)
	case *ggpk.DirectoryNode:
		fmt.Printf("+ D %20s %s/\n", "", p)
		children, err := n.Children()
		if err != nil {
			return err
		}
		for _, c := range children {
			show_added(path.Join(root, n.Name()), c)
		}
	}
	return nil
}

func show_removed(root string, n ggpk.AnyNode) {
	p := path.Join(root, n.Name())
	switch n := n.(type) {
	case *ggpk.FileNode:
		fmt.Printf("- F %20d %s\n", n.Size(), p)
	case *ggpk.DirectoryNode:
		fmt.Printf("- D %20s %s\n", "", p)
	}
}

func show_changed_file(root string, n1, n2 *ggpk.FileNode) {
	p := path.Join(root, n1.Name())
	fmt.Printf("± F %8d -> %8d %s\n", n1.Size(), n2.Size(), p)
}
