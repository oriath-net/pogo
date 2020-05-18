package ggpk

import (
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf16"
)

func readStringFrom(rr io.Reader) (string, error) {
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
