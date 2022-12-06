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
	SR
	PC
	IR
	IMMEDIATE
)

const (
	STATUS_HALT uint32 = 1 << iota
	STATUS_OVERFLOW
	STATUS_UNDERFLOW
	STATUS_DIVIDE_BY_ZERO
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
			R0: {0x00},
			R1: {0x00},
			R2: {0x00},
			R3: {0x00},
			SR: {0x00},
			PC: {0x100},
			IR: {0x00},
		},
	}
}
