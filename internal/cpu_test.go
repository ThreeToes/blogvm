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
		sr, err := registers.GetRegister(SR)
		if !assert.NoError(t, err) {
			return
		}
		assert.Greater(t, sr.Value&STATUS_HALT, uint32(0))
	})
	t.Run("test read", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[R0].Value = 0x1000
		bus.Write(0x100, 0x01010000)
		bus.Write(0x1000, 0xFFFF)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0xFFFF), registers.registerMap[R1].Value)
	})
	t.Run("test write", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[R0].Value = 0xFFFF
		registers.registerMap[R1].Value = 0x1000
		bus.Write(0x100, 0x02010000)

		err := cpu.Tick()
		assert.NoError(t, err)
		val, err := bus.Read(0x1000)
		assert.NoError(t, err)
		assert.Equal(t, uint32(0xFFFF), val)
	})
	t.Run("test copy", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[R0].Value = 0xFFFF
		bus.Write(0x100, 0x03010000)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0xFFFF), registers.registerMap[R1].Value)
	})
	t.Run("test add", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[R1].Value = 0x10
		bus.Write(0x100, 0x04F10003)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x13), registers.registerMap[R1].Value)
		assert.Equal(t, uint32(0x0), registers.registerMap[SR].Value)
	})
	t.Run("test add with overflow", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[R1].Value = 0xFFFFFFFE
		bus.Write(0x100, 0x04F10003)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x1), registers.registerMap[R1].Value)
		assert.Equal(t, STATUS_OVERFLOW, registers.registerMap[SR].Value)
	})
	t.Run("test sub", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[R1].Value = 0x02
		bus.Write(0x100, 0x05F10003)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x1), registers.registerMap[R1].Value)
		assert.Equal(t, uint32(0x0), registers.registerMap[SR].Value)
	})
	t.Run("test sub with underflow check", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[R1].Value = 0x05
		bus.Write(0x100, 0x05F10003)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0xFFFFFFFD), registers.registerMap[R1].Value)
		assert.Equal(t, STATUS_UNDERFLOW, registers.registerMap[SR].Value)
	})
	t.Run("test mul", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[R1].Value = 0x10
		bus.Write(0x100, 0x06F10003)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x30), registers.registerMap[R1].Value)
		assert.Equal(t, uint32(0x0), registers.registerMap[SR].Value)
	})
	t.Run("test mul with overflow", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[R1].Value = 0xFFFFFFFE
		bus.Write(0x100, 0x06F10003)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0xFFFFFFFA), registers.registerMap[R1].Value)
		assert.Equal(t, STATUS_OVERFLOW, registers.registerMap[SR].Value)
	})
	t.Run("test div", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[R1].Value = 0x03
		bus.Write(0x100, 0x07F10009)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x03), registers.registerMap[R1].Value)
		assert.Equal(t, uint32(0x0), registers.registerMap[SR].Value)
	})
	t.Run("test div with divide by zero", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[R1].Value = 0x00
		bus.Write(0x100, 0x07F10003)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x00), registers.registerMap[R1].Value)
		assert.Equal(t, STATUS_DIVIDE_BY_ZERO, registers.registerMap[SR].Value)
	})
	t.Run("test stat flag unset", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[SR].Value = 0x00000000
		registers.registerMap[R1].Value = 0x01
		bus.Write(0x100, 0x08F10001)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x00), registers.registerMap[R1].Value)
	})
	t.Run("test stat flag set", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[SR].Value = 0x0000000A
		registers.registerMap[R1].Value = 0x00
		bus.Write(0x100, 0x08F10002)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x01), registers.registerMap[R1].Value)
	})
	t.Run("test set flag", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[SR].Value = 0x00000000
		registers.registerMap[R1].Value = 0x01
		bus.Write(0x100, 0x09F10001)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x01), registers.registerMap[SR].Value)
	})
	t.Run("test reset flag", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[SR].Value = 0x0000000A
		registers.registerMap[R1].Value = 0x00
		bus.Write(0x100, 0x09F10002)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x08), registers.registerMap[SR].Value)
	})
	t.Run("run invalid instruction", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		bus.Write(0x100, 0xFFF10002)

		err := cpu.Tick()
		assert.Error(t, err)
		assert.Equal(t, STATUS_HALT, registers.registerMap[SR].Value)
	})
}
