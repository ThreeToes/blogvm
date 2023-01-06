package main

import (
	"flag"
	"fmt"
	"github.com/ThreeToes/blogvm/internal/assembler"
	"github.com/ThreeToes/blogvm/internal/machine"
	"os"
	"path/filepath"
)

type includeArgs []string

func (i *includeArgs) String() string {
	return fmt.Sprintln(*i)
}

func (i *includeArgs) Set(s string) error {
	*i = append(*i, s)
	return nil
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("must provide a command:")
		fmt.Println("\t* run - run a file")
		return
	}
	switch os.Args[1] {
	case "run":
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("error getting working directory: %v\n", err)
			return
		}
		execPath, err := os.Executable()
		if err != nil {
			fmt.Printf("error getting executable path: %v\n", err)
			return
		}
		libPath := filepath.Join(filepath.Dir(execPath), "lib")
		wdLib := filepath.Join(wd, "lib")
		includes := includeArgs{libPath, wd, wdLib}

		fs := flag.NewFlagSet("run", flag.ExitOnError)
		filePath := fs.String("file", "", "path to the file to run")
		flag.Var(&includes, "include", "add this folder to standard include paths")
		err = fs.Parse(os.Args[2:])
		if err != nil {
			fmt.Printf("could not parse args: %v\n", err)
			return
		}
		if *filePath == "" {
			fmt.Printf("file cannot be empty\n")
			fs.Usage()
			return
		}
		assembled, err := assembler.AssembleFile(*filePath, includes)
		if err != nil {
			fmt.Printf("could not assemble program: %v\n", err)
			return
		}

		registers := machine.NewRegisterBank()
		mem := machine.NewMemory()
		term := machine.NewTerminal()
		err = mem.Load(assembled)
		if err != nil {
			fmt.Printf("could not load assembled program: %v\n", err)
			return
		}
		bus := machine.NewBus(mem, term)
		cpu := machine.NewCPU(registers, bus)

		sr, err := registers.GetRegister(machine.SR)
		if err != nil {
			fmt.Printf("Could not get the staus register: %v", err)
			return
		}
		fmt.Println("Begin execution")
		fmt.Println("-------")
		for sr.Value&machine.STATUS_HALT == 0 {
			cpu.Tick()
		}
		fmt.Println()
		fmt.Println("-------")
		fmt.Println("Machine has halted")
	}
}
