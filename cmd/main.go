package main

import (
	"flag"
	"fmt"
	"github.com/ThreeToes/blogvm/internal/assembler"
	"github.com/ThreeToes/blogvm/internal/machine"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("must provide a command:")
		fmt.Println("\t* run - run a file")
		return
	}
	switch os.Args[1] {
	case "run":
		fs := flag.NewFlagSet("run", flag.ExitOnError)
		filePath := fs.String("file", "", "path to the file to run")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			fmt.Printf("could not parse args: %v\n", err)
			return
		}
		if *filePath == "" {
			fmt.Printf("file cannot be empty\n")
			fs.Usage()
			return
		}
		assembled, err := assembler.AssembleFile(*filePath)
		if err != nil {
			fmt.Printf("could not assemble program: %v\n", err)
			return
		}

		registers := machine.NewRegisterBank()
		mem := machine.NewMemory()
		err = mem.Load(assembled)
		if err != nil {
			fmt.Printf("could not load assembled program: %v\n", err)
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
		fmt.Printf("machine has halted\n")
	}
}
