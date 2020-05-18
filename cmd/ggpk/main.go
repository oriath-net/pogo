package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/duskwuff/pogo/ggpk"
	flag "github.com/spf13/pflag"
)

var excludePaths = flag.StringSlice("exclude", []string{}, "exclude these paths")
var ggpkPath = flag.StringP("ggpk", "f", "", "source GGPK file (required)")
var intoPath = flag.String("into", "", "extract into this directory")
var toStdout = flag.Bool("stdout", false, "extract files to standard output")
var verbose = flag.BoolP("verbose", "v", false, "display verbose listings and extraction progress")
var noRecurse = flag.BoolP("no-recurse", "n", false, "don't recurse into directories")

type fileVisitor func(path string, node ggpk.AnyNode) error

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
	}

	if *ggpkPath == "" {
		fmt.Fprintf(os.Stderr, "Must specify a GGPK\n")
		os.Exit(1)
	}

	f, err := ggpk.Open(*ggpkPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open GGPK: %s\n", err)
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case "list", "l":
		iterate(f, doList, flag.Args()[1:])

	case "extract", "x":
		if *intoPath == "" && !*toStdout {
			fmt.Fprintf(os.Stderr, "--into is required if --stdout is not in use\n")
			os.Exit(1)
		}
		iterate(f, doExtract, flag.Args()[1:])

	default:
		fmt.Fprintf(os.Stderr, "Unknown command. %s --help for details.\n", os.Args[0])
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(
		os.Stderr,
		"Usage:\n"+
			"    ggpk [options] list [<paths...>]\n"+
			"    ggpk [options] extract [<paths...>]\n"+
			"\n"+
			"Options are:\n",
	)
	flag.PrintDefaults()
	os.Exit(1)
}

func excluded(path string) bool {
	for _, xpath := range *excludePaths {
		if path == xpath {
			return true
		}
	}
	return false
}

func iterate(f *ggpk.File, fn fileVisitor, paths []string) {
	if len(paths) == 0 {
		paths = []string{""}
	}
	for i := range paths {
		rootPath := path.Clean(paths[i])
		rootNode, err := f.NodeAtPath(rootPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", rootPath, err)
			os.Exit(1)
		}

		switch n := rootNode.(type) {
		case *ggpk.FileNode:
			err := fn(rootPath, rootNode)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", rootPath, err)
				os.Exit(1)
			}

		case *ggpk.DirectoryNode:
			iter := n.Iterator()
			for iter.Next() {
				fullPath := path.Join(rootPath, iter.Path())
				if excluded(fullPath) {
					iter.Skip()
					continue
				}
				err := fn(fullPath, iter.Node())
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s: %s\n", fullPath, err)
					os.Exit(1)
				}
				if *noRecurse {
					iter.Skip()
				}
			}
			if err := iter.Error(); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}

		default:
			panic("unexpected node type")
		}
	}
}

func doList(path string, node ggpk.AnyNode) error {
	if *verbose {
		switch n := node.(type) {
		case *ggpk.FileNode:
			fmt.Printf("F %12d %s\n", n.Size(), path)
		case *ggpk.DirectoryNode:
			fmt.Printf("D %12s %s/\n", "", path)
		default:
			panic(fmt.Errorf("unexpected %T while iterating", n))
		}
	} else {
		switch n := node.(type) {
		case *ggpk.FileNode:
			fmt.Printf("%s\n", path)
		case *ggpk.DirectoryNode:
			fmt.Printf("%s/\n", path)
		default:
			panic(fmt.Errorf("unexpected %T while iterating", n))
		}
	}
	return nil
}

func doExtract(srcpath string, node ggpk.AnyNode) error {
	var tgt io.WriteCloser

	fileNode, ok := node.(*ggpk.FileNode)
	if !ok {
		return nil
	}

	if *verbose {
		fmt.Fprintf(os.Stderr, "%s\n", srcpath)
	}

	if *toStdout {
		tgt = os.Stdout

	} else {
		truepath := path.Join(*intoPath, srcpath)

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

	if !*toStdout {
		err := tgt.Close()
		return err
	}

	return nil
}
