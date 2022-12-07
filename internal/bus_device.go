package internal

type MemoryRange struct {
	Start uint32
	End   uint32
}

// BusDevice is an interface for devices we attach to the bus
type BusDevice interface {
	// MemoryRange gives the memory range of the device
	MemoryRange() *MemoryRange
	// Read takes an address and returns the Value at address
	Read(address uint32) (uint32, error)
	// Write writes Value to address
	Write(address, value uint32) error
}
