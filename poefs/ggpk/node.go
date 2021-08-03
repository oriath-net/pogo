package ggpk

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type anyNode interface {
	Name() string
}

var (
	errNodeWrongLength = errors.New("node has invalid length")
)

func (g *ggpkFS) getNodeAt(offset int64) (anyNode, error) {
	// long enough to fully read all node headers, hopefully long enough for
	// the filename in a FILE or PDIR node as well
	data := make([]byte, 256) // FIXME: FILE/PDIR need to support longer reads

	n, err := g.file.ReadAt(data, offset)
	if err != nil && (n < 8 || err != io.EOF) {
		return nil, fmt.Errorf("unable to read node at %x: %w", offset, err)
	}

	nodeSize := binary.LittleEndian.Uint32(data[0:4])
	nodeType := string(data[4:8])

	// remove size/type from data
	data = data[8:]
	nodeSize -= 8 // size includes size/type header
	if len(data) > int(nodeSize) {
		data = data[:nodeSize]
	}

	switch nodeType {
	case "GGPK":
		return g.newGgpkNode(data, offset+8, nodeSize)
	case "PDIR":
		return g.newPdirNode(data, offset+8, nodeSize)
	case "FILE":
		return g.newFileNode(data, offset+8, nodeSize)
	default:
		return nil, fmt.Errorf("unknown node type %#v at %x", nodeType, offset)
	}
}
