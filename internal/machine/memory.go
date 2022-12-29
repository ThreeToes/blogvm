package machine

import (
	"fmt"
	"github.com/ThreeToes/blogvm/internal/executable"
)

const maxMemorySize = 0xFFE1
const maxMemoryAddress = 0xFFE0

type Memory struct {
	mem [maxMemorySize]uint32
}

func (m *Memory) MemoryRange() *MemoryRange {
	return &MemoryRange{
		Start: 0x0000,
		End:   0xFFE0,
	}
}

func (m *Memory) Read(address uint32) (uint32, error) {
	if address > 0xFFE0 {
		return 0, fmt.Errorf("address %x out of range", address)
	}
	return m.mem[address], nil
}

func (m *Memory) Write(address, value uint32) error {
	if address > 0xFFE0 {
		return fmt.Errorf("address %x out of range", address)
	}
	m.mem[address] = value
	return nil
}

func (m *Memory) Load(l *executable.LoadableFile) error {
	for i := uint32(0); i < l.BlockCount; i++ {
		b := l.Blocks[i]
		if b.Address+b.BlockSize > maxMemoryAddress {
			return fmt.Errorf("address %d is not valid for block size %d", b.Address, b.BlockSize)
		}
		for j := uint32(0); j < b.BlockSize; j++ {
			m.mem[b.Address+j] = b.Words[j]
		}
	}
	return nil
}

func NewMemory() *Memory {
	return &Memory{
		mem: [maxMemorySize]uint32{},
	}
}
