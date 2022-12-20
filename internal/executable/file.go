package executable

import (
	"fmt"
	"io"
)

// LoadableFile represents a file we can load into memory
type LoadableFile struct {
	BlockCount uint32
	Flags      uint32
	Blocks     []*MemoryBlock
}

type MemoryBlock struct {
	Address   uint32
	BlockSize uint32
	Words     []uint32
}

// Load loads a loadable file from a binary stream
func Load(bs io.ByteReader) (*LoadableFile, error) {
	blockCount, err := nextWord(bs)
	if err != nil {
		return nil, fmt.Errorf("error reading block count: %v", err)
	}

	flags, err := nextWord(bs)
	if err != nil {
		return nil, fmt.Errorf("error reading flags: %v", err)
	}

	blocks, err := loadBlocks(blockCount, bs)
	if err != nil {
		return nil, fmt.Errorf("error loading blocks: %v", err)
	}

	return &LoadableFile{
		BlockCount: blockCount,
		Flags:      flags,
		Blocks:     blocks,
	}, nil
}

// loadBlocks from a stream
func loadBlocks(blockCount uint32, bs io.ByteReader) ([]*MemoryBlock, error) {
	if blockCount == 0 {
		return nil, nil
	}
	blocks := make([]*MemoryBlock, blockCount)
	var err error
	for i := 0; i < int(blockCount); i++ {
		blocks[i], err = loadBlock(bs)
		if err != nil {
			return nil, fmt.Errorf("error loading block %d: %v", i, err)
		}
	}
	return blocks, nil
}

// loadBlock from a stream
func loadBlock(bs io.ByteReader) (*MemoryBlock, error) {
	address, err := nextWord(bs)
	if err != nil {
		return nil, fmt.Errorf("error reading address: %v", err)
	}

	blockSize, err := nextWord(bs)
	if err != nil {
		return nil, fmt.Errorf("error reading block size: %v", err)
	}

	words := make([]uint32, blockSize)
	for i := 0; i < int(blockSize); i++ {
		words[i], err = nextWord(bs)
		if err != nil {
			return nil, fmt.Errorf("error loading word %d: %v", i, err)
		}
	}

	return &MemoryBlock{
		Address:   address,
		BlockSize: blockSize,
		Words:     words,
	}, nil
}

// nextWord gets the next word in a binary stream
func nextWord(bs io.ByteReader) (uint32, error) {
	b0, err := bs.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("error reading first byte: %v", err)
	}
	b1, err := bs.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("error reading second byte: %v", err)
	}
	b2, err := bs.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("error reading third byte: %v", err)
	}
	b3, err := bs.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("error reading fourth byte: %v", err)
	}
	ret := uint32(b0) << 24
	ret |= uint32(b1) << 16
	ret |= uint32(b2) << 8
	ret |= uint32(b3)

	return ret, nil
}
