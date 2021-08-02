package bundle

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"io"
	"sort"
	"strings"
)

type bundleIndex struct {
	bundles []string
	files   []bundleFileInfo
}

type bundleFileInfo struct {
	path     string
	bundleId uint32
	offset   uint32
	size     uint32
}

type bundlePathrep struct {
	offset        uint32
	size          uint32
	recursiveSize uint32
}

func loadBundleIndex(indexFile io.ReaderAt) (bundleIndex, error) {
	indexBundle, err := openBundle(indexFile)
	if err != nil {
		return bundleIndex{}, fmt.Errorf("unable to load index bundle: %w", err)
	}

	indexData := make([]byte, indexBundle.Size())
	if _, err := indexBundle.ReadAt(indexData, 0); err != nil {
		return bundleIndex{}, fmt.Errorf("unable to read index bundle: %w", err)
	}

	p := 0

	bundleCount := binary.LittleEndian.Uint32(indexData[p:])
	p += 4

	bundles := make([]string, bundleCount)
	for i := range bundles {
		nameLen := int(binary.LittleEndian.Uint32(indexData[p:]))
		p += 4

		name := string(indexData[p : p+nameLen])
		p += nameLen

		// skip uncompressed size -- available elsewhere
		p += 4

		bundles[i] = name
	}

	fileCount := binary.LittleEndian.Uint32(indexData[p:])
	p += 4

	files := make([]bundleFileInfo, fileCount)
	filemap := make(map[uint64]int, fileCount)
	for i := 0; i < int(fileCount); i++ {
		hash := binary.LittleEndian.Uint64(indexData[p+0:])
		files[i] = bundleFileInfo{
			bundleId: binary.LittleEndian.Uint32(indexData[p+8:]),
			offset:   binary.LittleEndian.Uint32(indexData[p+12:]),
			size:     binary.LittleEndian.Uint32(indexData[p+16:]),
		}
		p += 20
		if _, exists := filemap[hash]; exists {
			panic("duplicate filemap hash")
		}
		filemap[hash] = i
	}

	pathrepCount := binary.LittleEndian.Uint32(indexData[p:])
	p += 4

	pathmap := make(map[uint64]bundlePathrep, pathrepCount)
	for i := uint32(0); i < pathrepCount; i++ {
		hash := binary.LittleEndian.Uint64(indexData[p+0:])
		pr := bundlePathrep{
			offset:        binary.LittleEndian.Uint32(indexData[p+8:]),
			size:          binary.LittleEndian.Uint32(indexData[p+12:]),
			recursiveSize: binary.LittleEndian.Uint32(indexData[p+16:]),
		}
		p += 20
		if _, exists := pathmap[hash]; exists {
			panic("duplicate pathmap hash")
		}
		pathmap[hash] = pr
	}

	pathrepBundle, err := openBundle(bytes.NewReader(indexData[p:]))
	if err != nil {
		return bundleIndex{}, fmt.Errorf("unable to read pathrep bundle: %w", err)
	}

	pathData := make([]byte, pathrepBundle.Size())
	if _, err := pathrepBundle.ReadAt(pathData, 0); err != nil {
		return bundleIndex{}, fmt.Errorf("unable to read pathrep bundle: %w", err)
	}

	for _, pr := range pathmap {
		data := pathData[pr.offset : pr.offset+pr.size]
		paths := readPathspec(data)
		for _, path := range paths {
			h := fnv.New64a()
			h.Write([]byte(strings.ToLower(path) + "++"))
			sum := h.Sum64()
			if fe, found := filemap[sum]; found {
				files[fe].path = path
			} else {
				panic("unable to map bundle path to file")
			}
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].path < files[j].path
	})

	return bundleIndex{
		bundles: bundles,
		files:   files,
	}, nil
}

func readPathspec(data []byte) []string {
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

		str := readPathspecString(data, &p)
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

func readPathspecString(data []byte, offset *int) string {
	p := *offset
	for p < len(data) && data[p] != 0 {
		p++
	}
	s := string(data[*offset:p])
	*offset = p + 1
	return s
}
