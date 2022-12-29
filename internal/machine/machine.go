package machine

// NewMachine creates a new, default machine
func NewMachine() *CPU {
	return NewCPU(NewRegisterBank(), NewBus(NewMemory()))
}
