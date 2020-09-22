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
	version    int32
}

func (n *HeaderNode) Name() string      { return "" }
func (n *HeaderNode) Type() string      { return "GGPK" }
func (n *HeaderNode) Signature() []byte { return nil }

func (g *File) initNodeGGPK(offset int64, data []byte) (*HeaderNode, error) {
	var node physGGPK
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &node); err != nil {
		return nil, fmt.Errorf("unable to read GGPK header at %08x: %w", offset, err)
	}

	return &HeaderNode{
		nodeCommon: nodeCommon{
			src:    g,
			offset: offset,
			length: int64(node.Length),
		},
		version:    node.Version,
		rootOffset: int64(node.RootOffset),
		freeOffset: int64(node.FreeOffset),
	}, nil
}
