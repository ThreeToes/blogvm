package machine

import "fmt"

// TerminalDevice is a bus device that backs directly onto a real terminal
type TerminalDevice struct{}

const (
	TERMINAL = uint32(0xFFE1) + iota
	TERMINAL_INT
	TERMINAL_X
	TERMINAL_Y
	__terminal_reserved1
)

func (t *TerminalDevice) MemoryRange() *MemoryRange {
	// Addresses:
	// * 0xFFE1 - Write a character to terminal or read a character
	// * 0xFFE2 - Write a number to the terminal
	// * 0xFFE3 - Cursor X position
	// * 0xFFE4 - Cursor Y position
	// * 0xFFE5 - reserved
	return &MemoryRange{
		Start: 0xFFE1,
		End:   0xFFE5,
	}
}

func (t *TerminalDevice) Read(address uint32) (uint32, error) {
	switch address {
	// TODO: read values
	}
	return 0, nil
}

func (t *TerminalDevice) Write(address, value uint32) error {
	switch address {
	case TERMINAL:
		fmt.Print(rune(value))
	case TERMINAL_INT:
		fmt.Printf("%d", value)
	}
	return nil
}

func NewTerminal() *TerminalDevice {
	return &TerminalDevice{}
}
