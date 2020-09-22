package ggpk

// These types represent the physical layout of nodes within a GGPK archive.

type physHeader struct {
	Length uint32
	Kind   [4]byte
}

type physGGPK struct {
	physHeader
	Version    int32
	RootOffset int64
	FreeOffset int64
}

type physFREE struct {
	physHeader
	NextOffset int64
	// followed by arbitrary amounts of free space
}

type physFILE struct {
	physHeader
	NameLen   int32
	Signature [32]byte
	// followed by name and file content
}

type physPDIR struct {
	physHeader
	NameLen    int32
	ChildCount int32
	Signature  [32]byte
	// followed by name and child specs
}
