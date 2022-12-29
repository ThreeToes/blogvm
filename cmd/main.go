package main

import (
	"fmt"
	"github.com/ThreeToes/blogvm/internal/machine"
)

func main() {
	// Load this in at mem[0x100]
	// This program will calculate the value of 10! (10*9*8...) and store it at address 0x1000
	program := []uint32{
		0x03F00001, //0x100: COPY 0x01 R0
		0x03F1000A, //0x101: COPY 0x0A R1
		0x03F20001, //0x102: COPY 0x01 R2
		0x06100000, //0x103: MUL R1 R0
		0x11F10001, //0x104: EQ 0x01 R1
		0x0CF0010A, //0x105: JMP 0x10A
		0x05120000, //0x106: SUB R1 R2
		0x03210000, //0x107: COPY R2 R1
		0x03F20001, //0x108: COPY 0x01 R2
		0x0CF00102, //0x109: JMP 0x102
		0x020F1000, //0x10A: WRITE R0 0x1000
		0x00000000, //0x10B: HALT
	}

	registers := machine.NewRegisterBank()
	mem := machine.NewMemory()
	for i, p := range program {
		err := mem.Write(uint32(i+0x100), p)
		if err != nil {
			fmt.Printf("error writing to memory location 0x%x: %v\n",
				i, err)
			return
		}
	}
	bus := machine.NewBus(mem)
	cpu := machine.NewCPU(registers, bus)

	sr, err := registers.GetRegister(machine.SR)
	if err != nil {
		fmt.Printf("Could not get the staus register: %v", err)
		return
	}
	for sr.Value&machine.STATUS_HALT == 0 {
		cpu.Tick()
	}
	val, err := mem.Read(0x1000)
	if err != nil {
		fmt.Printf("could not read memory at 0x1000: %v\n", err)
		return
	}
	// Value in memory should be 3628800
	fmt.Printf("The value at address 0x1000 is %d\n", val)
}
