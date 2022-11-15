package internal

import "fmt"

const (
	R0 uint8 = iota
	R1
	R2
	R3
	__reserved1
	__reserved2
	__reserved3
	__reserved4
	__reserved5
	__reserved6
	__reserved7
	__reserved8
	__reserved9
	PC
	IR
	IMMEDIATE
)

type Register struct {
	value uint32
}

type RegisterBank struct {
	registerMap map[uint8]*Register
}

func (r *RegisterBank) GetRegister(name uint8) (*Register, error) {
	if reg, ok := r.registerMap[name]; ok {
		return reg, nil
	}
	return nil, fmt.Errorf("no such register")
}

func NewRegisterBank() *RegisterBank {
	return &RegisterBank{
		registerMap: map[uint8]*Register{
			R0: &Register{0x00},
			R1: &Register{0x00},
			R2: &Register{0x00},
			R3: &Register{0x00},
			PC: &Register{0x100},
			IR: &Register{0x00},
		},
	}
}
