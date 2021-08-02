package ggpk

import (
	"errors"
	"io"
	"io/fs"
	"strings"
)

var (
	errSignature = errors.New("not a GGPK file")
	errVersion   = errors.New("unsupported GGPK version")
	errStructure = errors.New("unexpected node type")
	errNotDir    = errors.New("non-directory in path")
	errNotFile   = errors.New("that's a directory")
)

type ggpkFS struct {
	file          io.ReaderAt
	root          ggpkNode
	useBundles    bool
	useUTF32Names bool
}

func NewLoader(src io.ReaderAt) (*ggpkFS, error) {
	g := &ggpkFS{file: src}
	node, err := g.getNodeAt(0)
	if err != nil {
		return nil, err
	}

	root, ok := node.(*ggpkNode)
	if !ok {
		return nil, errSignature
	}
	g.root = *root

	switch g.root.version {
	case 2:
		break
	case 3:
		g.useBundles = true
	case 4:
		g.useBundles = true
		g.useUTF32Names = true
	default:
		return nil, errVersion
	}

	return g, nil
}

func (g *ggpkFS) UseBundles() bool {
	return g.useBundles
}

func (g *ggpkFS) useUTF32() bool {
	return g.useUTF32Names
}

func (g *ggpkFS) Open(name string) (fs.File, error) {
	var node anyNode

	node, err := g.root.getRoot()
	if err != nil {
		return nil, err
	}

	parts := strings.Split(name, "/")
	for i := range parts {
		if parts[i] == "" || parts[i] == "." {
			continue
		}

		pdirNode, ok := node.(*pdirNode)
		if !ok {
			return nil, &fs.PathError{
				Op:   "open",
				Path: name,
				Err:  errNotDir,
			}
		}

		cn, err := pdirNode.ChildNamed(parts[i])
		if err != nil {
			return nil, &fs.PathError{
				Op:   "open",
				Path: name,
				Err:  err,
			}
		}

		node = cn
	}

	switch n := node.(type) {
	case *fileNode:
		return n.Reader()
	case *pdirNode:
		return n.Reader()
	default:
		return nil, &fs.PathError{
			Op:   "open",
			Path: name,
			Err:  errNotFile,
		}
	}
}
