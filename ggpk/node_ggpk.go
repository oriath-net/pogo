package ggpk

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type HeaderNode struct {
	nodeCommon
	rootOffset int64
	freeOffset int64
}

func (n *HeaderNode) Name() string { return "" }
func (n *HeaderNode) Type() string { return "GGPK" }

func (g *File) initNodeGGPK(offset int64, data []byte) (*HeaderNode, error) {
	var node physGGPK
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &node); err != nil {
		return nil, fmt.Errorf("unable to read GGPK header at %08x: %w", offset, err)
	}

	if node.NodeCount != 2 {
		return nil, fmt.Errorf("count is %d in GGPK header at %08x, wtf?", node.NodeCount, offset)
	}

	return &HeaderNode{
		nodeCommon: nodeCommon{
			src:    g,
			offset: offset,
			length: int64(node.Length),
		},
		rootOffset: int64(node.RootOffset),
		freeOffset: int64(node.FreeOffset),
	}, nil
}
