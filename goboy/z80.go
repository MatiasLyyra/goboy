package goboy

import "fmt"

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
	Memory *MMU

	Timer    int
	DivTimer int
}

func (cpu *CPU) HandleInterrupts() bool {
	if !cpu.EI {
		return false
	}
	ie := cpu.Memory.registers[AddrIE]
	ifReg := cpu.Memory.registers[AddrIF]
	ieVal := ie.Get()
	ifRegVal := ifReg.Get()
	for i := 0; i < 5; i++ {
		mask := uint8(1 << i)
		if ieVal&mask != 0 && ifRegVal&mask != 0 {
			cpu.EI = false
			cpu.Halt = false
			ifReg.RawSet(ifRegVal & ^mask)
			var intVector uint16
			switch i {
			case VBlankInt:
				intVector = 0x40
				fmt.Println("VBlank")
			case LCDStatInt:
				intVector = 0x48
				fmt.Println("LCDStatInt")
			case TimerInt:
				intVector = 0x50
				fmt.Println("TimerInt")
			case SerialInt:
				intVector = 0x58
			case JoypadInt:
				fmt.Println("JoypadInt")
				intVector = 0x60
			}
			cpu.Memory.Write(cpu.SP-1, uint8(cpu.PC>>8))
			cpu.Memory.Write(cpu.SP-2, uint8(cpu.PC))
			cpu.SP -= 2
			cpu.PC = intVector
			return true
		}
	}
	return false
}

func (cpu *CPU) RunSingleOpcode() int {
	cycles := 4
	if !cpu.Halt {
		opcode := cpu.Memory.Read(cpu.PC)
		cpu.PC++
		cycles = InstructionsTable[opcode](cpu)
	}
	cpu.updateTimers(cycles)
	cpu.HandleInterrupts()

	return cycles
}

func (cpu *CPU) updateTimers(cycles int) {
	var (
		div   = cpu.Memory.registers[AddrDIV]
		tima  = cpu.Memory.registers[AddrTIMA]
		tma   = cpu.Memory.registers[AddrTMA]
		tac   = cpu.Memory.registers[AddrTAC]
		ifReg = cpu.Memory.registers[AddrIF]
	)
	tacVal := tac.Get()
	// Timer is not enabled
	cpu.DivTimer += cycles
	divCount := cpu.DivTimer / 256
	cpu.DivTimer -= divCount * 256
	div.RawSet(div.Get() + uint8(divCount))
	if tacVal&0x04 == 0 {
		return
	}
	tacClock := 1024
	if timer := tacVal & 0x3; timer > 0 {
		tacClock = 16 << ((timer - 1) * 2)
	}
	cpu.Timer += cycles
	for cpu.Timer >= tacClock {
		cpu.Timer -= tacClock
		tima.RawSet(tima.Get() + 1)
		if tima.Get() == 0 {
			ifReg.RawSet(setBit(ifReg.Get(), TimerInt))
			tima.RawSet(tma.Get())
			cpu.Halt = false
		}
	}
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
