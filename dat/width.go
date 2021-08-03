package dat

import (
	"path"
)

type parserWidth int

const (
	Dat32 parserWidth = iota
	Dat64
	Datl32
	Datl64
)

func widthForFilename(filename string) parserWidth {
	switch path.Ext(filename) {
	case ".dat":
		return Dat32
	case ".dat64":
		return Dat64
	case ".datl":
		return Datl32
	case ".datl64":
		return Datl64
	default:
		return Dat32
	}
}

func (w parserWidth) String() string {
	switch w {
	case Dat32:
		return "32-bit with UTF-16 strings"
	case Dat64:
		return "64-bit with UTF-16 strings"
	case Datl32:
		return "32-bit with UTF-32 strings"
	case Datl64:
		return "64-bit with UTF-32 strings"
	default:
		return "invalid"
	}
}

func (w parserWidth) is64Bit() bool {
	return (w == Dat64) || (w == Datl64)
}

func (w parserWidth) isUTF32() bool {
	return (w == Datl32) || (w == Datl64)
}
