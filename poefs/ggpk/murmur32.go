package ggpk

import (
	"strings"
	"unicode/utf16"
)

func murmur32_utf16(name string) uint32 {
	codepoints := utf16.Encode([]rune(strings.ToLower(name)))

	m := uint32(0x5bd1e995)
	h := uint32(len(codepoints) * 2)
	r := 24

	for i := 0; i+1 < len(codepoints); i += 2 {
		k := uint32(codepoints[i]) + uint32(codepoints[i+1])<<16
		k *= m
		k ^= k >> r
		k *= m
		h *= m
		h ^= k
	}

	if len(codepoints)%2 == 1 {
		h ^= uint32(codepoints[len(codepoints)-1])
		h *= m
	}

	h ^= h >> 13
	h *= m
	h ^= h >> 15

	return h
}
