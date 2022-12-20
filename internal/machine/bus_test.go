package machine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBus_Write(t *testing.T) {
	t.Run("device written correctly", func(t *testing.T) {
		m := NewMemory()
		bus := NewBus(m)
		err := bus.Write(0x100, 0xCCCC)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		assert.Equal(t, uint32(0xCCCC), m.mem[0x100])
	})
	t.Run("no mapped device", func(t *testing.T) {
		bus := NewBus()
		err := bus.Write(0x100, 0xCCCC)
		if !assert.Error(t, err) {
			t.FailNow()
		}
	})
}

func TestBus_Read(t *testing.T) {
	t.Run("device read correctly", func(t *testing.T) {
		m := NewMemory()
		m.mem[0x100] = 0xFFFF
		bus := NewBus(m)
		got, err := bus.Read(0x100)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		assert.Equal(t, uint32(0xFFFF), got)
	})
	t.Run("no mapped device", func(t *testing.T) {
		bus := NewBus()
		got, err := bus.Read(0x100)
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.Equal(t, uint32(0), got)
	})
}
