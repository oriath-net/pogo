package ggpk

import (
	"strings"
)

type Iterator struct {
	src    *File
	frames []iteratorFrame
	node   AnyNode
	err    error
	skip   bool
}

type iteratorFrame struct {
	node     *DirectoryNode
	position int
}

func (g *File) Iterator() *Iterator {
	root, err := g.RootNode()
	if err != nil {
		return &Iterator{err: err}
	}
	return root.Iterator()
}

func (n *DirectoryNode) Iterator() *Iterator {
	return &Iterator{
		src: n.src,
		frames: []iteratorFrame{
			{node: n, position: 0},
		},
		err: nil,
	}
}

func (t *Iterator) topFrame() *iteratorFrame {
	if len(t.frames) == 0 {
		return nil
	}
	return &t.frames[len(t.frames)-1]
}

func (t *Iterator) Next() bool {
	if t.err != nil {
		return false
	}

	f := t.topFrame()
	if f == nil {
		return false
	}

	if t.node != nil && t.node.Type() == "PDIR" && !t.skip {
		t.frames = append(t.frames, iteratorFrame{
			node:     t.node.(*DirectoryNode),
			position: 0,
		})
		f = t.topFrame()
	}

	for f != nil && f.position >= len(f.node.childOffsets) {
		t.frames = t.frames[:len(t.frames)-1]
		f = t.topFrame()
	}

	if f == nil {
		return false
	}

	node, err := f.node.childAtIndex(f.position)
	if err != nil {
		t.err = err
		return false
	}

	t.node = node
	t.skip = false
	f.position += 1

	return true
}

func (t *Iterator) Error() error {
	return t.err
}

func (t *Iterator) Path() string {
	parts := make([]string, len(t.frames)) // -1 to skip root frame, +1 for node name
	for i := range t.frames[1:] {
		parts[i] = t.frames[i+1].node.Name()
	}
	parts[len(parts)-1] = t.node.Name()
	return strings.Join(parts, "/")
}

func (t *Iterator) Node() AnyNode {
	return t.node
}

func (t *Iterator) Skip() {
	t.skip = true
}
