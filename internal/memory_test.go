package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
