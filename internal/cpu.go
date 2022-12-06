package internal

import "fmt"

const (
	HALT = iota
	READ
	WRITE
	COPY
	ADD
	SUB
	MUL
	DIV
	STAT
	SET
)

type CPU struct {
	registers *RegisterBank
	bus       *Bus
	Halted    bool
}

func (c *CPU) halt(_, _ *Register) {
	c.Halted = true
}

func (c *CPU) read(i1, i2 *Register) {
	val, err := c.bus.Read(i1.value)
	if err != nil {
		c.Halted = true
		return
	}
	i2.value = val
}

func (c *CPU) write(i1, i2 *Register) {
	err := c.bus.Write(i2.value, i1.value)
	if err != nil {
		c.Halted = true
		return
	}
}

func (c *CPU) copy(i1, i2 *Register) {
	i2.value = i1.value
}

func (c *CPU) add(i1, i2 *Register) {
	i1Val := uint64(i1.value)
	i2Val := uint64(i2.value)
	sum := i1Val + i2Val
	if sum > 0xFFFFFFFF {
		sr, err := c.registers.GetRegister(SR)
		// We'll panic here because if the status register doesn't work then our machine may as well crash
		if err != nil {
			panic(err)
		}
		sr.value = sr.value | STATUS_OVERFLOW
	}
	i2.value = uint32(sum & 0xFFFFFFFF)
}

func (c *CPU) sub(i1, i2 *Register) {
	i1Val := int64(i1.value)
	i2Val := int64(i2.value)
	diff := i1Val - i2Val
	if diff < 0 {
		sr, err := c.registers.GetRegister(SR)
		if err != nil {
			panic(err)
		}
		sr.value = sr.value | STATUS_UNDERFLOW
		diff = diff + 0xFFFFFFFF
	}
	i2.value = uint32(diff & 0xFFFFFFFF)
}

func (c *CPU) executeInstruction(instruction uint32) error {
	opcode := uint8(instruction >> 24)
	regIndex1 := uint8((instruction & 0x00F00000) >> 20)
	regIndex2 := uint8((instruction & 0x000F0000) >> 16)
	imm := instruction & 0x0000FFFF

	var i1, i2 *Register
	var err error
	if regIndex1 == 0xF {
		i1 = &Register{value: imm}
	} else {
		i1, err = c.registers.GetRegister(regIndex1)
		if err != nil {
			return err
		}
	}
	if regIndex2 == 0xF {
		i2 = &Register{value: imm}
	} else {
		i2, err = c.registers.GetRegister(regIndex2)
		if err != nil {
			return err
		}
	}

	switch opcode {
	case HALT:
		c.halt(i1, i2)
	case READ:
		c.read(i1, i2)
	case WRITE:
		c.write(i1, i2)
	case COPY:
		c.copy(i1, i2)
	case ADD:
		c.add(i1, i2)
	case SUB:
		c.sub(i1, i2)
	}

	return nil
}

func (c *CPU) Tick() error {
	if c.Halted {
		return fmt.Errorf("cannot tick on a Halted machine")
	}
	ir, err := c.registers.GetRegister(IR)
	if err != nil {
		return err
	}
	pc, err := c.registers.GetRegister(PC)
	if err != nil {
		return err
	}
	ir.value, err = c.bus.Read(pc.value)
	if err != nil {
		return err
	}
	pc.value++
	err = c.executeInstruction(ir.value)
	if err != nil {
		c.Halted = true
	}
	return err
}

func NewCPU(registers *RegisterBank, bus *Bus) *CPU {
	return &CPU{
		registers: registers,
		bus:       bus,
		Halted:    false,
	}
}
