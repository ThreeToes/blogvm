package main

import (
	"fmt"
	"vm_blog/internal"
)

func main() {
	// Load this in at mem[0x100]
	// This program will copy the value 0x05 into R0, 0x10 into
	// R1, and 0x1000 into R2, before executing an add on R0 & R1
	// and writing the contents into the address in R2. If all is
	// working, at the end of the program we should see the value
	// 0x15 at memory address 0x1000
	program := []uint32{
		0x03F00005,
		0x03F10010,
		0x03F21000,
		0x04010000,
		0x02120000,
		0x00000000,
	}

	registers := internal.NewRegisterBank()
	mem := internal.NewMemory()
	for i, p := range program {
		err := mem.Write(uint32(i+0x100), p)
		if err != nil {
			fmt.Printf("error writing to memory location 0x%x: %v\n",
				i, err)
			return
		}
	}
	bus := internal.NewBus(mem)
	cpu := internal.NewCPU(registers, bus)

	sr, err := registers.GetRegister(internal.SR)
	if err != nil {
		fmt.Printf("Could not get the staus register: %v", err)
		return
	}
	for sr.Value&internal.STATUS_HALT == 0 {
		cpu.Tick()
	}
	val, err := mem.Read(0x1000)
	if err != nil {
		fmt.Printf("could not read memory at 0x1000: %v\n", err)
		return
	}
	fmt.Printf("The value at address 0x1000 is 0x%x\n", val)
}
