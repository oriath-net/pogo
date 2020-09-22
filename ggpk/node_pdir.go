package ggpk

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type DirectoryNode struct {
	nodeCommon
	name         string
	childOffsets []int64
	signature    [32]byte
}

func (n *DirectoryNode) Name() string      { return n.name }
func (n *DirectoryNode) Type() string      { return "PDIR" }
func (n *DirectoryNode) Signature() []byte { return n.signature[:] }

func (g *File) initNodePDIR(offset int64, data []byte) (*DirectoryNode, error) {
	var node physPDIR
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &node); err != nil {
		return nil, fmt.Errorf("unable to read PDIR header at %08x: %w", offset, err)
	}

	if node.Length != uint32(48+g.sizeofName(node.NameLen)+12*node.ChildCount) {
		return nil, fmt.Errorf("PDIR at %08x has unexpected length %d (%+v)", offset, node.Length, node)
	}

	// Read the rest of the PDIR into memory if it isn't already
	if int(node.Length) > len(data) {
		data = make([]byte, node.Length)
		_, err := g.file.ReadAt(data, offset)
		if err != nil {
			return nil, fmt.Errorf("unable to read PDIR data at %08x: %w", offset, err)
		}
	}

	br := bytes.NewReader(data[48:])

	name, err := g.readStringFrom(br)
	if err != nil {
		return nil, fmt.Errorf("unable to read PDIR name at %08x: %w", offset, err)
	}

	offsets := make([]int64, node.ChildCount)
	for i := range offsets {
		var child struct {
			Timestamp int32
			Offset    int64
		}
		if err := binary.Read(br, binary.LittleEndian, &child); err != nil {
			return nil, fmt.Errorf("unable to read PDIR contents at %08x: %w", offset, err)
		}
		offsets[i] = child.Offset
	}

	return &DirectoryNode{
		nodeCommon: nodeCommon{
			src:    g,
			offset: offset,
			length: int64(node.Length),
		},
		name:         name,
		signature:    node.Signature,
		childOffsets: offsets,
	}, nil
}

func (n *DirectoryNode) Children() ([]AnyNode, error) {
	nodes := make([]AnyNode, len(n.childOffsets))
	for i := range nodes {
		n, err := n.childAtIndex(i)
		if err != nil {
			return nil, err
		}
		nodes[i] = n
	}
	return nodes, nil
}

func (n *DirectoryNode) childAtIndex(i int) (AnyNode, error) {
	return n.src.getNode(n.childOffsets[i])
}

func (n *DirectoryNode) childNamed(name string) (AnyNode, error) {
	children, err := n.Children()
	if err != nil {
		return nil, err
	}
	for i := range children {
		if children[i].Name() == name {
			return children[i], nil
		}
	}
	return nil, nil
}
