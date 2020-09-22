package ggpk

import (
	"encoding/binary"
	"fmt"
	"io"

	"unicode/utf16"
)

func (g *File) readStringFrom(rr io.Reader) (string, error) {
	if g.useUTF32Names {
		runes := make([]rune, 0, 64)
		for {
			var ch rune
			err := binary.Read(rr, binary.LittleEndian, &ch)
			if err != nil {
				return "", fmt.Errorf("failed reading string: %w", err)
			}
			if ch == 0 {
				break
			}
			runes = append(runes, ch)
		}
		return string(runes), nil
	} else {
		str := make([]uint16, 0, 64)
		for {
			var ch uint16
			err := binary.Read(rr, binary.LittleEndian, &ch)
			if err != nil {
				return "", fmt.Errorf("failed reading string: %w", err)
			}
			if ch == 0 {
				break
			}
			str = append(str, ch)
		}
		return string(utf16.Decode(str)), nil
	}
}

func (g *File) sizeofName(n int32) int32 {
	if g.useUTF32Names {
		return 4 * n
	} else {
		return 2 * n
	}
}
