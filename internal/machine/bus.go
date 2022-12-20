package machine

import "fmt"

type Bus struct {
	devices []BusDevice
}

func (b *Bus) Read(address uint32) (uint32, error) {
	for _, d := range b.devices {
		memRange := d.MemoryRange()
		if memRange.Start <= address && address <= memRange.End {
			return d.Read(address)
		}
	}
	return 0, fmt.Errorf("bus read: unmapped address %x", address)
}

func (b *Bus) Write(address, value uint32) error {
	for _, d := range b.devices {
		memRange := d.MemoryRange()
		if memRange.Start <= address && address <= memRange.End {
			return d.Write(address, value)
		}
	}
	return fmt.Errorf("bus write: unmapped address %x", address)
}

func NewBus(devices ...BusDevice) *Bus {
	return &Bus{
		devices: devices,
	}
}
