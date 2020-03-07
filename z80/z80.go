package z80

type CPUMemory interface {
	Read(addr uint16) uint8
	Write(addr uint16, data uint8)
}

type CPU struct {
	// General purpose registers
	A     uint8
	FSign bool
	FZero bool
	// F5 Not in use
	FHalfCarry bool
	// F3 Not in use
	FParityOverflow bool
	FSub            bool
	FCarry          bool

	B uint8
	C uint8

	D uint8
	E uint8

	H uint8
	L uint8

	// Special purpose registers
	I  uint8
	R  uint8
	IX uint16
	IY uint16
	SP uint16
	PC uint16

	Memory CPUMemory
}

func (cpu *CPU) BC() uint16 {
	return uint16(cpu.B)<<8 | uint16(cpu.C)
}

func (cpu *CPU) DE() uint16 {
	return uint16(cpu.D)<<8 | uint16(cpu.E)
}

func (cpu *CPU) HL() uint16 {
	return uint16(cpu.H)<<8 | uint16(cpu.L)
}
