package main

import (
	"fmt"
	"github.com/ThreeToes/blogvm/internal/assembler"
	"github.com/ThreeToes/blogvm/internal/machine"
)

var program = `COPY 0x01 R0
COPY 0x0A R1
LOOP COPY 0x01 R2
MUL R1 R0
EQ 0x01 R1
JMP END
SUB R1 R2
COPY R2 R1
COPY 0x01 R2
JMP LOOP
END WRITE R0 0x1000
HALT
`

func main() {
	assembled, err := assembler.AssembleString(program)
	if err != nil {
		fmt.Printf("could not assemble program: %v", err)
		return
	}

	registers := machine.NewRegisterBank()
	mem := machine.NewMemory()
	err = mem.Load(assembled)
	if err != nil {
		fmt.Printf("could not load assembled program: %v", err)
		return
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
