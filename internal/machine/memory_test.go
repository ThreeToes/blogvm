package machine

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"vm_blog/internal/executable"
)

func TestNewMemory(t *testing.T) {
	got := NewMemory()
	assert.Len(t, got.mem, 0xFFE1)
	for _, v := range got.mem {
		assert.Equal(t, uint32(0), v)
	}
}

func TestWriteRead(t *testing.T) {
	mem := NewMemory()
	assert.Equal(t, uint32(0), mem.mem[0x100])
	err := mem.Write(0x100, 0xCCCC)
	if !assert.NoError(t, err) {
		t.Errorf("could not write to memory: %v", err)
		t.FailNow()
	}
	read, err := mem.Read(0x100)
	if !assert.NoError(t, err) {
		t.Errorf("could not read from memory: %v", err)
		t.FailNow()
	}
	assert.Equal(t, uint32(0xCCCC), read)
}

func TestMemory_Load(t *testing.T) {
	mem := NewMemory()
	file := &executable.LoadableFile{
		BlockCount: 3,
		Flags:      0,
		Blocks: []*executable.MemoryBlock{
			{
				Address:   0x100,
				BlockSize: 0x01,
				Words: []uint32{
					0x1234,
				},
			},
			{
				Address:   0x1500,
				BlockSize: 0x02,
				Words: []uint32{
					0x1234,
					0x4567,
				},
			},
			{
				Address:   0x300,
				BlockSize: 0x03,
				Words: []uint32{
					0x1234,
					0x4567,
					0x89,
				},
			},
		},
	}
	mem.Load(file)
	for _, b := range file.Blocks {
		for i := uint32(0); i < b.BlockSize; i++ {
			read, err := mem.Read(b.Address + i)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, b.Words[i], read)
		}
	}
}
