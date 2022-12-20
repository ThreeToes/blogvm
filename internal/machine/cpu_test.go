package machine

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
	t.Run("test push", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[SR].Value = 0x00000000
		registers.registerMap[R1].Value = 0x00
		bus.Write(0x100, 0x0AF00002)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0xFFDF), registers.registerMap[SP].Value)
		v, err := bus.Read(0xFFDF)
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x02), v)
	})
	t.Run("test pop", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[SR].Value = 0x00000000
		registers.registerMap[R1].Value = 0x00
		registers.registerMap[SP].Value = 0xFFDF
		bus.Write(0x100, 0x0BF00002)
		bus.Write(0xFFDF, 0x10)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0xFFE0), registers.registerMap[SP].Value)
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x10), registers.registerMap[R0].Value)
	})
	t.Run("test jmp", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[SR].Value = 0x00000000
		registers.registerMap[R1].Value = 0x00
		registers.registerMap[SP].Value = 0xFFDF
		bus.Write(0x100, 0x0CF01000)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x1000), registers.registerMap[PC].Value)
	})
	t.Run("test less", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[SR].Value = 0x00000000
		registers.registerMap[R1].Value = 0x0FFF
		registers.registerMap[SP].Value = 0xFFDF
		bus.Write(0x100, 0x0DF11000)
		bus.Write(0x102, 0x0DF10EEE)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x102), registers.registerMap[PC].Value)

		err = cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x103), registers.registerMap[PC].Value)
	})
	t.Run("test lte", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[SR].Value = 0x00000000
		registers.registerMap[R1].Value = 0x0FFF
		registers.registerMap[SP].Value = 0xFFDF
		bus.Write(0x100, 0x0EF11000)
		bus.Write(0x102, 0x0EF10EEE)
		bus.Write(0x103, 0x0EF10FFF)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x102), registers.registerMap[PC].Value)

		err = cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x103), registers.registerMap[PC].Value)

		err = cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x104), registers.registerMap[PC].Value)
	})
	t.Run("test gt", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[SR].Value = 0x00000000
		registers.registerMap[R1].Value = 0x0FFF
		registers.registerMap[SP].Value = 0xFFDF
		bus.Write(0x100, 0x0FF10EEE)
		bus.Write(0x102, 0x0FF11000)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x102), registers.registerMap[PC].Value)

		err = cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x103), registers.registerMap[PC].Value)
	})
	t.Run("test gte", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[SR].Value = 0x00000000
		registers.registerMap[R1].Value = 0x0FFF
		registers.registerMap[SP].Value = 0xFFDF
		bus.Write(0x100, 0x10F10EEE)
		bus.Write(0x102, 0x10F11000)
		bus.Write(0x103, 0x10F11001)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x102), registers.registerMap[PC].Value)

		err = cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x103), registers.registerMap[PC].Value)

		err = cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x104), registers.registerMap[PC].Value)
	})
	t.Run("test eq", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		registers.registerMap[SR].Value = 0x00000000
		registers.registerMap[R1].Value = 0x1000
		registers.registerMap[SP].Value = 0xFFDF
		bus.Write(0x100, 0x11F10EEE)
		bus.Write(0x102, 0x11F11000)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x102), registers.registerMap[PC].Value)

		err = cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x103), registers.registerMap[PC].Value)
	})
	t.Run("test call/return", func(t *testing.T) {
		registers := NewRegisterBank()
		bus := NewBus(NewMemory())
		cpu := NewCPU(registers, bus)

		bus.Write(0x100, 0x12F10200)
		bus.Write(0x200, 0x13F11000)

		err := cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x200), registers.registerMap[PC].Value)

		err = cpu.Tick()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0x101), registers.registerMap[PC].Value)
	})
}
