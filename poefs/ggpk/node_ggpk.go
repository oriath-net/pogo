package ggpk

import (
	"encoding/binary"
	"fmt"
)

type ggpkNode struct {
	src        *ggpkFS
	version    int32
	rootOffset int64
	freeOffset int64
}

func (n ggpkNode) Name() string {
	return "<GGPK header node>"
}

func (g *ggpkFS) newGgpkNode(data []byte, offset int64, length uint32) (*ggpkNode, error) {
	if len(data) != 20 {
		return nil, errNodeTooShort
	}

	return &ggpkNode{
		src:        g,
		version:    int32(binary.LittleEndian.Uint32(data[0:])),
		rootOffset: int64(binary.LittleEndian.Uint64(data[4:])),
		freeOffset: int64(binary.LittleEndian.Uint64(data[12:])),
	}, nil
}

func (n ggpkNode) getRoot() (*pdirNode, error) {
	nn, err := n.src.getNodeAt(n.rootOffset)
	if err != nil {
		return nil, fmt.Errorf("getRoot: %w", err)
	}
	if pn, ok := nn.(*pdirNode); ok {
		return pn, nil
	} else {
		return nil, errStructure
	}
}
