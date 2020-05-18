package ggpk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

// File represents a GGPK archive.
type File struct {
	file *os.File
	root HeaderNode
}

type AnyNode interface {
	Name() string
	Type() string
}

type nodeCommon struct {
	src    *File
	offset int64
	length int64
}

// Open() opens a GGPK archive at a specified path and returns a handle to the archive.
func Open(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	ggpk := &File{
		file: f,
	}

	node, err := ggpk.getNode(0)
	if err != nil {
		return nil, fmt.Errorf("unable to read root node: %w", err)
	}

	typedNode, ok := node.(*HeaderNode)
	if !ok {
		return nil, fmt.Errorf("not a GGPK file")
	}
	ggpk.root = *typedNode

	return ggpk, nil
}

func (g *File) getNode(offset int64) (AnyNode, error) {
	// long enough to fully read all node headers, hopefully long enough for
	// the filename in a FILE or PDIR node as well
	data := make([]byte, 128)
	n, err := g.file.ReadAt(data, offset)
	if err != nil && (n < 8 || err != io.EOF) {
		return nil, fmt.Errorf("unable to read node at %08x: %w", offset, err)
	}

	var nodeHeader physHeader
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &nodeHeader); err != nil {
		return nil, fmt.Errorf("unable to read node at %08x: %w", offset, err)
	}

	switch string(nodeHeader.Kind[:]) {
	case "GGPK":
		return g.initNodeGGPK(offset, data)
	case "PDIR":
		return g.initNodePDIR(offset, data)
	case "FILE":
		return g.initNodeFILE(offset, data)
	default:
		return nil, fmt.Errorf("unknown node type %#v at %08x", nodeHeader.Kind, offset)
	}
}

func (g *File) RootNode() (*DirectoryNode, error) {
	node, err := g.getNode(g.root.rootOffset)
	if err != nil {
		return nil, fmt.Errorf("unable to get root PDIR: %w", err)
	}
	typedNode, ok := node.(*DirectoryNode)
	if !ok {
		return nil, fmt.Errorf("root node is not a directory")
	}
	return typedNode, nil
}

func (g *File) NodeAtPath(path string) (AnyNode, error) {
	node, err := g.RootNode()
	if err != nil {
		return nil, err
	}

	if path == "." {
		return node, nil
	}

	parts := strings.Split(path, "/")
	for i := range parts {
		child, err := node.childNamed(parts[i])
		if err != nil {
			return nil, err
		}
		if child == nil {
			return nil, fmt.Errorf("%s: file not found", path)
		}
		if i == len(parts)-1 {
			return child, nil
		}
		dirChild, ok := child.(*DirectoryNode)
		if !ok {
			return nil, fmt.Errorf("%s: not a directory", path)
		}
		node = dirChild
	}
	panic("unreachable")
}

func (g *File) Open(path string) (io.Reader, error) {
	node, err := g.NodeAtPath(path)
	if err != nil {
		return nil, err
	}
	fileNode, ok := node.(*FileNode)
	if !ok {
		return nil, fmt.Errorf("%s: not a file", path)
	}
	return fileNode.Reader(), nil
}
