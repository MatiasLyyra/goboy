package goboy

// CPU represents internal state of the z80 cpu
type CPU struct {
	// General purpose registers
	A uint8

	B uint8
	C uint8

	D uint8
	E uint8

	H uint8
	L uint8

	// Flag registers
	FZero      bool
	FHalfCarry bool
	FSub       bool
	FCarry     bool

	// Special purpose registers
	SP uint16
	PC uint16

	// Misc
	Halt   bool
	EI     bool
	Memory Memory
}

func (cpu *CPU) RunSingleOpcode() int {
	opcode := cpu.Memory.Read(cpu.PC)
	// fmt.Println(InstructionStrings[opcode])
	cpu.PC++
	return InstructionsTable[opcode](cpu)
}

// F returns flags as uint8
func (cpu *CPU) F() uint8 {
	var (
		zero  uint8
		sub   uint8
		half  uint8
		carry uint8
	)
	if cpu.FZero {
		zero = 1 << 7
	}
	if cpu.FSub {
		sub = 1 << 6
	}
	if cpu.FHalfCarry {
		half = 1 << 5
	}
	if cpu.FCarry {
		carry = 1 << 4
	}
	return zero | sub | half | carry
}

// SetF sets cpu flags
func (cpu *CPU) SetF(flags uint8) {
	if (flags>>7)&1 == 1 {
		cpu.FZero = true
	} else {
		cpu.FZero = false
	}

	if (flags>>6)&1 == 1 {
		cpu.FSub = true
	} else {
		cpu.FSub = false
	}

	if (flags>>5)&1 == 1 {
		cpu.FHalfCarry = true
	} else {
		cpu.FHalfCarry = false
	}

	if (flags>>4)&1 == 1 {
		cpu.FCarry = true
	} else {
		cpu.FCarry = false
	}
}

// AF returns A and F registers combined into a word
func (cpu *CPU) AF() uint16 {
	return uint16(cpu.A)<<8 | uint16(cpu.F())
}

// BC returns B and C registers combined into a word
func (cpu *CPU) BC() uint16 {
	return uint16(cpu.B)<<8 | uint16(cpu.C)
}

// DE returns D and E registers combined into a word
func (cpu *CPU) DE() uint16 {
	return uint16(cpu.D)<<8 | uint16(cpu.E)
}

// HL returns H and L registers combined into a word
func (cpu *CPU) HL() uint16 {
	return uint16(cpu.H)<<8 | uint16(cpu.L)
}
