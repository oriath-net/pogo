package ggpk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"io"
	"strings"
)

type BundleIndex struct {
	bundles []bundle_index_info
	files   map[uint64]bundle_file_info
}

type bundle_index_info struct {
	Name             string
	UncompressedSize uint32
}

type bundle_file_info struct {
	path     string
	bundleId uint32
	offset   uint32
	size     uint32
}

type bundle_pathrep struct {
	offset        uint32
	size          uint32
	recursiveSize uint32
}

func loadBundleIndex(g *File) (*BundleIndex, error) {
	index_file, err := g.Open("Bundles2/_.index.bin")
	if err != nil {
		return nil, fmt.Errorf("unable to load index bundle: %w", err)
	}

	index_bundle, err := LoadBundle(index_file)
	if err != nil {
		return nil, fmt.Errorf("unable to load index bundle: %w", err)
	}

	index_data := make([]byte, index_bundle.Size())
	if _, err := index_bundle.ReadAt(index_data, 0); err != nil {
		return nil, fmt.Errorf("unable to read index bundle: %w", err)
	}

	p := 0

	bundle_count := binary.LittleEndian.Uint32(index_data[p:])
	p += 4

	bundles := make([]bundle_index_info, bundle_count)
	for i := range bundles {
		name_len := int(binary.LittleEndian.Uint32(index_data[p:]))
		p += 4

		name := string(index_data[p : p+name_len])
		p += name_len

		size := binary.LittleEndian.Uint32(index_data[p:])
		p += 4

		bundles[i] = bundle_index_info{
			Name:             name,
			UncompressedSize: size,
		}
	}

	file_count := binary.LittleEndian.Uint32(index_data[p:])
	p += 4

	filemap := make(map[uint64]bundle_file_info, file_count)
	for i := uint32(0); i < file_count; i++ {
		hash := binary.LittleEndian.Uint64(index_data[p+0:])
		bfi := bundle_file_info{
			bundleId: binary.LittleEndian.Uint32(index_data[p+8:]),
			offset:   binary.LittleEndian.Uint32(index_data[p+12:]),
			size:     binary.LittleEndian.Uint32(index_data[p+16:]),
		}
		p += 20
		if _, exists := filemap[hash]; exists {
			panic("duplicate filemap hash")
		}
		filemap[hash] = bfi
	}

	pathrep_count := binary.LittleEndian.Uint32(index_data[p:])
	p += 4

	pathmap := make(map[uint64]bundle_pathrep, pathrep_count)
	for i := uint32(0); i < pathrep_count; i++ {
		hash := binary.LittleEndian.Uint64(index_data[p+0:])
		pr := bundle_pathrep{
			offset:        binary.LittleEndian.Uint32(index_data[p+8:]),
			size:          binary.LittleEndian.Uint32(index_data[p+12:]),
			recursiveSize: binary.LittleEndian.Uint32(index_data[p+16:]),
		}
		p += 20
		if _, exists := pathmap[hash]; exists {
			panic("duplicate pathmap hash")
		}
		pathmap[hash] = pr
	}

	pathrep_bundle, err := LoadBundle(bytes.NewReader(index_data[p:]))
	if err != nil {
		return nil, fmt.Errorf("unable to read pathrep bundle: %w", err)
	}

	path_data := make([]byte, pathrep_bundle.Size())
	if _, err := pathrep_bundle.ReadAt(path_data, 0); err != nil {
		return nil, fmt.Errorf("unable to read pathrep bundle: %w", err)
	}

	for _, pr := range pathmap {
		data := path_data[pr.offset : pr.offset+pr.size]
		paths := dumpPathspec(data)
		for _, path := range paths {
			h := fnv.New64a()
			h.Write([]byte(strings.ToLower(path) + "++"))
			sum := h.Sum64()
			if fe, found := filemap[sum]; found {
				fe.path = path
				filemap[sum] = fe
			} else {
				panic("unable to map bundle path to file")
			}
		}
	}

	return &BundleIndex{
		bundles: bundles,
		files:   filemap,
	}, nil
}

func dumpPathspec(data []byte) []string {
	p := int(0)
	phase := 1
	names := make([]string, 0, 128)
	output := make([]string, 0, 128)

	for p < len(data) {
		n := int(binary.LittleEndian.Uint32(data[p:]))
		p += 4
		if n == 0 {
			phase = 1 - phase
			continue
		}

		str := readString(data, &p)
		if n-1 < len(names) {
			str = names[n-1] + str
		}
		if phase == 0 {
			names = append(names, str)
		} else {
			output = append(output, str)
		}
	}

	return output
}

func readString(data []byte, offset *int) string {
	p := *offset
	for p < len(data) && data[p] != 0 {
		p++
	}
	s := string(data[*offset:p])
	*offset = p + 1
	return s
}

func (g *File) OpenFileFromBundle(filename string) (ReadSeekerAt, error) {
	h := fnv.New64a()
	h.Write([]byte(strings.ToLower(filename) + "++"))
	fe, found := g.bundleIndex.files[h.Sum64()]
	if !found {
		return nil, fmt.Errorf("file not found")
	}

	bundlePath := "Bundles2/" + g.bundleIndex.bundles[fe.bundleId].Name + ".bundle.bin"
	bundleFile, err := g.Open(bundlePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open bundle %s containing %s: %w", bundlePath, filename, err)
	}
	bundle, err := LoadBundle(bundleFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open bundle %s containing %s: %w", bundlePath, filename, err)
	}

	return io.NewSectionReader(bundle, int64(fe.offset), int64(fe.size)), nil
}
