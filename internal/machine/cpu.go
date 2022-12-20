package machine

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
	PUSH
	POP
	JMP
	LESS
	LTE
	GT
	GTE
	EQ
	CALL
	RETURN
)

type CPU struct {
	registers *RegisterBank
	bus       *Bus
}

func (c *CPU) halt(_, _ *Register) {
	sr, err := c.registers.GetRegister(SR)
	if err != nil {
		panic(err)
	}
	sr.Value = sr.Value | STATUS_HALT
}

func (c *CPU) read(i1, i2 *Register) {
	val, err := c.bus.Read(i1.Value)
	if err != nil {
		sr, err := c.registers.GetRegister(SR)
		if err != nil {
			panic(err)
		}
		sr.Value = sr.Value | STATUS_MEMORY_ERROR
		return
	}
	i2.Value = val
}

func (c *CPU) write(i1, i2 *Register) {
	err := c.bus.Write(i2.Value, i1.Value)
	if err != nil {
		sr, err := c.registers.GetRegister(SR)
		if err != nil {
			panic(err)
		}
		sr.Value = sr.Value | STATUS_MEMORY_ERROR
	}
}

func (c *CPU) copy(i1, i2 *Register) {
	i2.Value = i1.Value
}

func (c *CPU) add(i1, i2 *Register) {
	i1Val := uint64(i1.Value)
	i2Val := uint64(i2.Value)
	sum := i1Val + i2Val
	if sum > 0xFFFFFFFF {
		sr, err := c.registers.GetRegister(SR)
		// We'll panic here because if the status register doesn't work then our machine may as well crash
		if err != nil {
			panic(err)
		}
		sr.Value = sr.Value | STATUS_OVERFLOW
	}
	i2.Value = uint32(sum & 0xFFFFFFFF)
}

func (c *CPU) sub(i1, i2 *Register) {
	i1Val := int64(i1.Value)
	i2Val := int64(i2.Value)
	diff := i1Val - i2Val
	if diff < 0 {
		sr, err := c.registers.GetRegister(SR)
		if err != nil {
			panic(err)
		}
		sr.Value = sr.Value | STATUS_UNDERFLOW
		diff = diff + 0xFFFFFFFF
	}
	i2.Value = uint32(diff & 0xFFFFFFFF)
}

func (c *CPU) mul(i1, i2 *Register) {
	i1Val := int64(i1.Value)
	i2Val := int64(i2.Value)

	product := i1Val * i2Val

	if product > 0xFFFFFFFF {
		sr, err := c.registers.GetRegister(SR)
		// We'll panic here because if the status register doesn't work then our machine may as well crash
		if err != nil {
			panic(err)
		}
		sr.Value = sr.Value | STATUS_OVERFLOW
	}
	i2.Value = uint32(product & 0xFFFFFFFF)
}

func (c *CPU) div(i1, i2 *Register) {
	i1Val := i1.Value
	i2Val := i2.Value

	if i2Val == 0 {
		sr, err := c.registers.GetRegister(SR)
		// We'll panic here because if the status register doesn't work then our machine may as well crash
		if err != nil {
			panic(err)
		}
		sr.Value = sr.Value | STATUS_DIVIDE_BY_ZERO
		return
	}

	i2.Value = i1Val / i2Val
}

func (c *CPU) stat(i1, i2 *Register) {
	var bit uint32 = 1 << (i1.Value - 1)
	sr, err := c.registers.GetRegister(SR)
	if err != nil {
		panic(err)
	}
	// Should be either 0 or 1
	i2.Value = (sr.Value & bit) / bit
}

func (c *CPU) set(i1, i2 *Register) {
	sr, err := c.registers.GetRegister(SR)
	if err != nil {
		panic(err)
	}

	var bit uint32 = 1 << (i1.Value - 1)
	if i2.Value > 0 {
		sr.Value = sr.Value | bit
	} else {
		sr.Value = sr.Value ^ bit
	}
}

func (c *CPU) push(i1, _ *Register) {
	sp, err := c.registers.GetRegister(SP)
	if err != nil {
		panic(err)
	}

	sp.Value--
	err = c.bus.Write(sp.Value, i1.Value)
	if err != nil {
		sr, err := c.registers.GetRegister(SR)
		if err != nil {
			panic(err)
		}
		sr.Value = sr.Value | STATUS_MEMORY_ERROR
		return
	}
}

func (c *CPU) pop(_, i2 *Register) {
	sp, err := c.registers.GetRegister(SP)
	if err != nil {
		panic(err)
	}

	v, err := c.bus.Read(sp.Value)
	if err != nil {
		sr, err := c.registers.GetRegister(SR)
		if err != nil {
			panic(err)
		}
		sr.Value = sr.Value | STATUS_MEMORY_ERROR
		return
	}
	i2.Value = v
	sp.Value++
}

func (c *CPU) jmp(i1, _ *Register) {
	pc, err := c.registers.GetRegister(PC)

	if err != nil {
		panic(err)
	}
	pc.Value = i1.Value
}

func (c *CPU) less(i1, i2 *Register) {
	if !(i1.Value < i2.Value) {
		pc, err := c.registers.GetRegister(PC)
		if err != nil {
			panic(err)
		}
		pc.Value++
	}
}

func (c *CPU) lte(i1, i2 *Register) {
	if !(i1.Value <= i2.Value) {
		pc, err := c.registers.GetRegister(PC)
		if err != nil {
			panic(err)
		}
		pc.Value++
	}
}

func (c *CPU) gt(i1, i2 *Register) {
	if !(i1.Value > i2.Value) {
		pc, err := c.registers.GetRegister(PC)
		if err != nil {
			panic(err)
		}
		pc.Value++
	}
}

func (c *CPU) gte(i1, i2 *Register) {
	if !(i1.Value >= i2.Value) {
		pc, err := c.registers.GetRegister(PC)
		if err != nil {
			panic(err)
		}
		pc.Value++
	}
}

func (c *CPU) eq(i1, i2 *Register) {
	if !(i1.Value == i2.Value) {
		pc, err := c.registers.GetRegister(PC)
		if err != nil {
			panic(err)
		}
		pc.Value++
	}
}

func (c *CPU) call(i1, _ *Register) {
	sp, err := c.registers.GetRegister(SP)
	if err != nil {
		panic(err)
	}
	pc, err := c.registers.GetRegister(PC)
	if err != nil {
		panic(err)
	}

	sp.Value--
	err = c.bus.Write(sp.Value, pc.Value)
	if err != nil {
		sr, err := c.registers.GetRegister(SR)
		if err != nil {
			panic(err)
		}
		sr.Value = sr.Value | STATUS_MEMORY_ERROR
		return
	}

	pc.Value = i1.Value
}

func (c *CPU) ret(_, _ *Register) {
	sp, err := c.registers.GetRegister(SP)
	if err != nil {
		panic(err)
	}

	val, err := c.bus.Read(sp.Value)
	sp.Value++
	if err != nil {
		sr, err := c.registers.GetRegister(SR)
		if err != nil {
			panic(err)
		}
		sr.Value = sr.Value | STATUS_MEMORY_ERROR
		return
	}

	pc, err := c.registers.GetRegister(PC)
	if err != nil {
		panic(err)
	}
	pc.Value = val
}

func (c *CPU) executeInstruction(instruction uint32) error {
	opcode := uint8(instruction >> 24)
	regIndex1 := uint8((instruction & 0x00F00000) >> 20)
	regIndex2 := uint8((instruction & 0x000F0000) >> 16)
	imm := instruction & 0x0000FFFF

	var i1, i2 *Register
	var err error
	if regIndex1 == 0xF {
		i1 = &Register{Value: imm}
	} else {
		i1, err = c.registers.GetRegister(regIndex1)
		if err != nil {
			return err
		}
	}
	if regIndex2 == 0xF {
		i2 = &Register{Value: imm}
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
	case MUL:
		c.mul(i1, i2)
	case DIV:
		c.div(i1, i2)
	case STAT:
		c.stat(i1, i2)
	case SET:
		c.set(i1, i2)
	case PUSH:
		c.push(i1, i2)
	case POP:
		c.pop(i1, i2)
	case JMP:
		c.jmp(i1, i2)
	case LESS:
		c.less(i1, i2)
	case LTE:
		c.lte(i1, i2)
	case GT:
		c.gt(i1, i2)
	case GTE:
		c.gte(i1, i2)
	case EQ:
		c.eq(i1, i2)
	case CALL:
		c.call(i1, i2)
	case RETURN:
		c.ret(i1, i2)
	default:
		// Halt the machine if we can't figure out the instruction
		c.set(&Register{1}, &Register{1})
		return fmt.Errorf("unrecognised opcode '%x'", opcode)
	}

	return nil
}

func (c *CPU) Tick() error {
	sr, err := c.registers.GetRegister(SR)
	if err != nil {
		return err
	}
	if sr.Value&STATUS_HALT > 0 {
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
	ir.Value, err = c.bus.Read(pc.Value)
	if err != nil {
		return err
	}
	pc.Value++
	err = c.executeInstruction(ir.Value)
	if err != nil {
		sr.Value = sr.Value | STATUS_HALT
	}
	return err
}

func NewCPU(registers *RegisterBank, bus *Bus) *CPU {
	return &CPU{
		registers: registers,
		bus:       bus,
	}
}
