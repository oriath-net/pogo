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
	file          io.ReaderAt
	root          HeaderNode
	useBundles    bool
	useUTF32Names bool
}

type AnyNode interface {
	Name() string
	Type() string
	Signature() []byte
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

	switch ggpk.root.version {
	case 2:
		break
	case 3:
		ggpk.useBundles = true
		break
	case 4:
		ggpk.useBundles = true
		ggpk.useUTF32Names = true
		break
	default:
		return nil, fmt.Errorf("unknown GGPK version %d", ggpk.root.version)
	}

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
	var node AnyNode

	node, err := g.RootNode()
	if err != nil {
		return nil, err
	}

	parts := strings.Split(path, "/")
	for i := range parts {
		if parts[i] == "" || parts[i] == "." {
			continue
		}
		dirNode, ok := node.(*DirectoryNode)
		if !ok {
			return nil, fmt.Errorf("%s: not a directory", path)
		}
		node, err = dirNode.childNamed(parts[i])
		if err != nil {
			return nil, err
		}
		if node == nil {
			return nil, fmt.Errorf("%s: file not found", path)
		}
	}
	return node, nil
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
