package internal

import "fmt"

type Memory struct {
	mem [0xFFE0]uint32
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

func NewMemory() *Memory {
	return &Memory{
		mem: [0xFFE0]uint32{},
	}
}
