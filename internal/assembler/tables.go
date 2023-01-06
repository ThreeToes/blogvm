package assembler

import (
	"fmt"
	"strconv"
	"strings"
)

type assemblable interface {
	calculateSize(sourceLine string) uint32
	assemble(sourceLine string, symbolTable symbolTableType) ([]uint32, error)
}

type opCode struct {
	mnemonic     string
	opcode       uint8
	hasI1        bool
	hasI2        bool
	allowSymbols bool
}

func (o *opCode) instructionMask() uint32 {
	return uint32(o.opcode) << 24
}

func (o *opCode) calculateSize(sourceLine string) uint32 {
	// NOTE: could change in the future
	return 1
}

func (o *opCode) assemble(sourceLine string, symbolTable symbolTableType) ([]uint32, error) {
	curIdx := 1
	instruction := uint32(o.opcode) << 24
	cols := strings.Split(sourceLine, " ")
	if cols[0] != o.mnemonic {
		curIdx = 2
	}
	if o.hasI1 {
		if len(cols) <= curIdx {
			return nil, fmt.Errorf("not enough args supplied")
		}
		arg := cols[curIdx]
		if reg, ok := registerTable[arg]; ok {
			nibble := uint32(reg.nibble)
			instruction = instruction | (nibble << 20)
		} else if symbol, ok := symbolTable[arg]; ok && o.allowSymbols {
			// Try resolving a symbol
			instruction = instruction | (0xF << 20) | (symbol.line & 0xFFFF)
		} else {
			// Try to parse immediate data
			p, err := parseLiteral(arg)
			if err != nil {
				return nil, fmt.Errorf("unrecognised symbol %q", arg)
			}
			instruction = instruction | (uint32(0xF) << 20) | (p & 0xFFFF)
		}
		curIdx++
	}
	if o.hasI2 {

		if len(cols) <= curIdx {
			return nil, fmt.Errorf("not enough args supplied")
		}
		arg := cols[curIdx]
		if reg, ok := registerTable[arg]; ok {
			nibble := uint32(reg.nibble)
			instruction = instruction | (nibble << 16)
		} else if symbol, ok := symbolTable[arg]; ok && o.allowSymbols {
			instruction = instruction | (0xF << 16) | (symbol.line & 0xFFFF)
		} else {
			p, err := parseLiteral(arg)
			if err != nil {
				return nil, err
			}
			instruction = instruction | (uint32(0xF) << 16) | (p & 0xFFFF)
		}
	}

	return []uint32{instruction}, nil
}

type opcodeTableType map[string]*opCode

var opcodeTable = opcodeTableType{
	"HALT": {
		mnemonic: "HALT",
		opcode:   0x00,
		hasI1:    false,
		hasI2:    false,
	},
	"READ": {
		mnemonic:     "READ",
		opcode:       0x01,
		hasI1:        true,
		hasI2:        true,
		allowSymbols: true,
	},
	"WRITE": {
		mnemonic:     "WRITE",
		opcode:       0x02,
		hasI1:        true,
		hasI2:        true,
		allowSymbols: true,
	},
	"COPY": {
		mnemonic: "COPY",
		opcode:   0x03,
		hasI1:    true,
		hasI2:    true,
	},
	"ADD": {
		mnemonic: "ADD",
		opcode:   0x04,
		hasI1:    true,
		hasI2:    true,
	},
	"SUB": {
		mnemonic: "SUB",
		opcode:   0x05,
		hasI1:    true,
		hasI2:    true,
	},
	"MUL": {
		mnemonic: "MUL",
		opcode:   0x06,
		hasI1:    true,
		hasI2:    true,
	},
	"DIV": {
		mnemonic: "DIV",
		opcode:   0x07,
		hasI1:    true,
		hasI2:    true,
	},
	"STAT": {
		mnemonic: "STAT",
		opcode:   0x08,
		hasI1:    true,
		hasI2:    true,
	},
	"SET": {
		mnemonic: "SET",
		opcode:   0x09,
		hasI1:    true,
		hasI2:    true,
	},
	"PUSH": {
		mnemonic: "PUSH",
		opcode:   0x0A,
		hasI1:    true,
		hasI2:    false,
	},
	"POP": {
		mnemonic: "POP",
		opcode:   0x0B,
		hasI1:    false,
		hasI2:    true,
	},
	"JMP": {
		mnemonic:     "JMP",
		opcode:       0x0C,
		hasI1:        true,
		hasI2:        false,
		allowSymbols: true,
	},
	"LESS": {
		mnemonic: "LESS",
		opcode:   0x0D,
		hasI1:    true,
		hasI2:    true,
	},
	"LTE": {
		mnemonic: "LTE",
		opcode:   0x0E,
		hasI1:    true,
		hasI2:    true,
	},
	"GT": {
		mnemonic: "GT",
		opcode:   0x0F,
		hasI1:    true,
		hasI2:    true,
	},
	"GTE": {
		mnemonic: "GTE",
		opcode:   0x10,
		hasI1:    true,
		hasI2:    true,
	},
	"EQ": {
		mnemonic: "EQ",
		opcode:   0x11,
		hasI1:    true,
		hasI2:    true,
	},
	"CALL": {
		mnemonic:     "CALL",
		opcode:       0x12,
		hasI1:        true,
		hasI2:        false,
		allowSymbols: true,
	},
	"RETURN": {
		mnemonic: "RETURN",
		opcode:   0x13,
		hasI1:    false,
		hasI2:    false,
	},
}

func (o opcodeTableType) isMnemonic(mnem string) bool {
	_, ok := o[mnem]
	return ok
}

func (o opcodeTableType) reverseLookup(opcode uint8) *opCode {
	for _, v := range o {
		if v.opcode == opcode {
			return v
		}
	}
	return nil
}

type directive struct {
	mnemonic     string
	sizeCalc     func(sourceLine string) uint32
	assembleFunc func(sourceLine string, symbolTable symbolTableType) ([]uint32, error)
}

func (d *directive) calculateSize(sourceLine string) uint32 {
	return d.sizeCalc(sourceLine)
}

func (d *directive) assemble(sourceLine string, symbolTable symbolTableType) ([]uint32, error) {
	return d.assembleFunc(sourceLine, symbolTable)
}

type directiveTableType map[string]*directive

var directiveTable = directiveTableType{
	"WORD": {
		mnemonic: "WORD",
		sizeCalc: func(_ string) uint32 {
			// Always one line long
			return 1
		},
		assembleFunc: func(sourceLine string, _ symbolTableType) ([]uint32, error) {
			spl := strings.Split(sourceLine, " ")
			col := 1
			if len(spl) > 2 {
				col = 2
			}
			if len(spl) <= col {
				return nil, fmt.Errorf("not enough arguments to WORD directive")
			}
			arg := spl[col]
			p, err := parseLiteral(arg)
			if err != nil {
				return nil, err
			}
			return []uint32{p}, nil
		},
	},
	"STRING": {
		mnemonic: "STRING",
		sizeCalc: func(sourceLine string) uint32 {
			split := strings.SplitN(sourceLine, " ", 2)
			if split[0] != "STRING" {
				split = strings.SplitN(split[1], " ", 2)
			}
			if len(split) == 1 {
				return 1
			}
			return uint32(len(split[1]) + 1)
		},
		assembleFunc: func(sourceLine string, _ symbolTableType) ([]uint32, error) {
			split := strings.SplitN(sourceLine, " ", 2)
			if split[0] != "STRING" {
				split = strings.SplitN(split[1], " ", 2)
			}
			if len(split) == 1 {
				return []uint32{0x00}, nil
			}
			var ret []uint32
			for _, ch := range split[1] {
				ret = append(ret, uint32(ch))
			}
			ret = append(ret, 0x00)
			return ret, nil
		},
	},
	// Loads the address of a symbol into a register
	"ADDRESS": {
		mnemonic: "ADDRESS",
		sizeCalc: func(_ string) uint32 {
			return 1
		},
		assembleFunc: func(sourceLine string, symbolTable symbolTableType) ([]uint32, error) {
			split := strings.Split(sourceLine, " ")
			symbolIdx := 1
			if split[0] != "ADDRESS" {
				symbolIdx = 2
			}
			if len(split) <= symbolIdx+1 {
				return nil, fmt.Errorf("ADDRESS directive did not have enough arguments")
			}
			symbolName := split[symbolIdx]
			dest := split[symbolIdx+1]
			if symbol, ok := symbolTable[symbolName]; ok {
				instr := fmt.Sprintf("COPY %d %s", symbol.line, dest)
				return opcodeTable["COPY"].assemble(instr, symbolTable)
			}
			return nil, fmt.Errorf("unrecognised symbol %q", symbolName)
		},
	},
}

type register struct {
	mnemonic string
	nibble   uint8
}

type registerTableType map[string]*register

var registerTable = registerTableType{
	"R0": {
		mnemonic: "R0",
		nibble:   0x0,
	},
	"R1": {
		mnemonic: "R1",
		nibble:   0x1,
	},
	"R2": {
		mnemonic: "R2",
		nibble:   0x2,
	},
	"R3": {
		mnemonic: "R3",
		nibble:   0x3,
	},
	"SP": {
		mnemonic: "SP",
		nibble:   0xB,
	},
	"SR": {
		mnemonic: "SR",
		nibble:   0xC,
	},
	"PC": {
		mnemonic: "PC",
		nibble:   0xD,
	},
	"IR": {
		mnemonic: "IR",
		nibble:   0xE,
	},
}

func parseLiteral(arg string) (uint32, error) {

	// Try to parse immediate data
	base := 10
	stripCount := 0
	if strings.HasPrefix(arg, "0x") {
		base = 16
		stripCount = 2
	} else if strings.HasPrefix(arg, "0b") {
		base = 2
		stripCount = 2
	} else if strings.HasPrefix(arg, "0") {
		base = 8
		stripCount = 1
	}
	p, err := strconv.ParseUint(arg[stripCount:], base, 32)
	if err != nil {
		return 0, fmt.Errorf("unrecognised symbol %q", arg)
	}
	return uint32(p), nil
}
