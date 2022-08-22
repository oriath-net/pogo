package ggpk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"time"
)

type fileNode struct {
	src        *ggpkFS
	name       string
	signature  []byte
	dataOffset int64
	dataLength int64
}

func (n *fileNode) Name() string {
	return n.name
}

func (g *ggpkFS) newFileNode(data []byte, offset int64, length uint32) (*fileNode, error) {
	if len(data) < 36 {
		return nil, errNodeTooShort
	}

	nameLength := int(binary.LittleEndian.Uint32(data[0:]))
	signature := data[4:36]

	headerLength := 36 + g.sizeofChars(nameLength)
	if headerLength > int(length) {
		return nil, errNodeTooShort
	}
	if len(data) < headerLength {
		data = make([]byte, headerLength)
		_, err := g.file.ReadAt(data, offset)
		if err != nil {
			return nil, err
		}
	}

	br := bytes.NewReader(data[36:])
	name, err := g.readStringFrom(br)
	if err != nil {
		return nil, fmt.Errorf("failed to read FILE name: %w", err)
	}

	n := &fileNode{
		src:        g,
		name:       name,
		signature:  signature,
		dataOffset: offset + int64(headerLength),
		dataLength: int64(length) - int64(headerLength),
	}

	return n, nil
}

func (n *fileNode) Reader() (fs.File, error) {
	return &fsFileNode{
		n,
		io.NewSectionReader(
			n.src.file,
			n.dataOffset,
			n.dataLength,
		),
	}, nil
}

// fsFileNode

type fsFileNode struct {
	src *fileNode
	*io.SectionReader
}

func (ffn *fsFileNode) Close() error {
	return nil
}

func (ffn *fsFileNode) Stat() (fs.FileInfo, error) {
	return &fsFileNodeStat{ffn}, nil
}

// fsFileNodeStat

type fsFileNodeStat struct {
	*fsFileNode
}

func (ffs *fsFileNodeStat) Name() string {
	return ffs.src.name
}

func (ffs *fsFileNodeStat) Size() int64 {
	return ffs.src.dataLength
}

func (ffs *fsFileNodeStat) Mode() fs.FileMode {
	return 0o444
}

func (ffi *fsFileNodeStat) ModTime() time.Time {
	return time.Unix(0, 0)
}

func (ffi *fsFileNodeStat) IsDir() bool {
	return false
}

func (ffi *fsFileNodeStat) Sys() any {
	return nil
}

func (ffi *fsFileNodeStat) Provenance() string {
	return "GGPK"
}

func (ffi *fsFileNodeStat) Signature() []byte {
	return ffi.src.signature
}
