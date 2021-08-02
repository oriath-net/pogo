package bundle

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/oriath-net/gooz"
)

type bundle struct {
	data        io.ReaderAt
	size        int64
	granularity int64 // size of each chunk of uncompressed data, usually 256KiB
	blocks      []bundleBlock
}

// descriptions of compressed blocks relative to bundle.data
type bundleBlock struct {
	offset int64
	length int64
}

type bundleHead struct {
	UncompressedSize             uint32
	TotalPayloadSize             uint32
	HeadPayloadSize              uint32
	FirstFileEncode              uint32
	_                            uint32
	UncompressedSize2            int64
	TotalPayloadSize2            int64
	BlockCount                   uint32
	UncompressedBlockGranularity uint32
	_                            [4]uint32
}

func openBundle(r io.ReaderAt) (*bundle, error) {
	rs := io.NewSectionReader(r, 0, 1<<24)

	var bh bundleHead
	if err := binary.Read(rs, binary.LittleEndian, &bh); err != nil {
		return nil, fmt.Errorf("failed to read bundle head: %w", err)
	}

	blockSizes := make([]uint32, bh.BlockCount)
	if err := binary.Read(rs, binary.LittleEndian, &blockSizes); err != nil {
		return nil, fmt.Errorf("failed to read bundle block sizes: %w", err)
	}

	blocks := make([]bundleBlock, bh.BlockCount)
	p := int64(binary.Size(bh) + binary.Size(blockSizes))
	for i := range blockSizes {
		sz := int64(blockSizes[i])
		blocks[i] = bundleBlock{offset: p, length: sz}
		p += sz
	}

	b := bundle{
		data:        r,
		size:        bh.UncompressedSize2,
		granularity: int64(bh.UncompressedBlockGranularity),
		blocks:      blocks,
	}

	// do a quick sanity check here
	if b.granularity == 0 {
		return nil, fmt.Errorf("granularity is 0?!")
	}

	expectedBlocks := b.size / b.granularity
	if b.size%b.granularity > 0 {
		expectedBlocks += 1
	}

	if int(expectedBlocks) != len(blocks) {
		return nil, fmt.Errorf(
			"got %d blocks of size %d for %d bytes data",
			len(blocks),
			b.granularity,
			b.size,
		)
	}

	return &b, nil
}

func (b *bundle) Size() int64 {
	return b.size
}

func (b *bundle) ReadAt(p []byte, off int64) (int, error) {
	if off+int64(len(p)) > b.size {
		// FIXME: This could be handled more gracefully
		return 0, fmt.Errorf("read outside bounds of file")
	}

	// Temporary buffers for compressed and decompressed data
	ibuf := make([]byte, b.granularity+64)
	obuf := make([]byte, b.granularity)

	n := 0
	for n < len(p) {
		blkId := int(off / b.granularity)
		blkOff := int(off % b.granularity)
		blk := &b.blocks[blkId]

		rawSize := int(b.granularity)
		if blkId == len(b.blocks)-1 {
			rawSize = int(b.size - int64(blkId)*b.granularity)
		}

		oodleBlk := ibuf[:blk.length]
		if n, err := b.data.ReadAt(oodleBlk, blk.offset); n != len(oodleBlk) {
			return 0, err
		}

		_, err := gooz.Decompress(oodleBlk, obuf[:rawSize])
		if err != nil {
			return 0, fmt.Errorf("decompression failed: %w", err)
		}

		copied := copy(p[n:], obuf[blkOff:])
		n += copied
		off += int64(copied)
	}

	return n, nil
}
