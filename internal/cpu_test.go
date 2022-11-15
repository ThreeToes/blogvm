package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCPU(t *testing.T) {
	t.Run("test halt", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		bus.Write(0x100, 0x00)
		err := cpu.Tick()
		assert.NoError(t, err)
		assert.True(t, cpu.Halted)
	})
	t.Run("test read", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[R0].value = 0x1000
		bus.Write(0x100, 0x01010000)
		bus.Write(0x1000, 0xFFFF)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.False(t, cpu.Halted)
		assert.Equal(t, uint32(0xFFFF), registers.registerMap[R1].value)
	})
	t.Run("test write", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[R0].value = 0xFFFF
		registers.registerMap[R1].value = 0x1000
		bus.Write(0x100, 0x02010000)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.False(t, cpu.Halted)
		val, err := bus.Read(0x1000)
		assert.NoError(t, err)
		assert.Equal(t, uint32(0xFFFF), val)
	})
	t.Run("test copy", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[R0].value = 0xFFFF
		bus.Write(0x100, 0x03010000)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.False(t, cpu.Halted)
		assert.Equal(t, uint32(0xFFFF), registers.registerMap[R1].value)
	})
	t.Run("test add", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[R1].value = 0x10
		bus.Write(0x100, 0x04F10003)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.False(t, cpu.Halted)
		assert.Equal(t, uint32(0x13), registers.registerMap[R1].value)
	})
}
