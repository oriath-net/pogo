package ggpk

import (
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
	Hash     uint64
	BundleId uint32
	Offset   uint32
	Size     uint32
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

	filemap := make(map[uint64]bundle_file_info)
	for i := uint32(0); i < file_count; i++ {
		bfi := bundle_file_info{
			Hash:     binary.LittleEndian.Uint64(index_data[p+0:]),
			BundleId: binary.LittleEndian.Uint32(index_data[p+8:]),
			Offset:   binary.LittleEndian.Uint32(index_data[p+12:]),
			Size:     binary.LittleEndian.Uint32(index_data[p+16:]),
		}
		p += 20
		if _, exists := filemap[bfi.Hash]; exists {
			panic("duplicate hash")
		}
		filemap[bfi.Hash] = bfi
	}

	return &BundleIndex{
		bundles: bundles,
		files:   filemap,
	}, nil
}

func (g *File) OpenFileFromBundle(filename string) (ReadSeekerAt, error) {
	h := fnv.New64a()
	h.Write([]byte(strings.ToLower(filename) + "++"))
	fe, found := g.bundleIndex.files[h.Sum64()]
	if !found {
		return nil, fmt.Errorf("file not found")
	}

	bundlePath := "Bundles2/" + g.bundleIndex.bundles[fe.BundleId].Name + ".bundle.bin"
	bundleFile, err := g.Open(bundlePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open bundle %s containing %s: %w", bundlePath, filename, err)
	}
	bundle, err := LoadBundle(bundleFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open bundle %s containing %s: %w", bundlePath, filename, err)
	}

	return io.NewSectionReader(bundle, int64(fe.Offset), int64(fe.Size)), nil
}
