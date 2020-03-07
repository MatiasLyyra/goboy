package z80

type MicroCodeFunc func(cpu *CPU) int

var InstructionsTable [256]MicroCodeFunc

func NOP(cpu *CPU) int {
	cpu.PC++
	return 4
}
func init() {
	InstructionsTable = [256]MicroCodeFunc{
		// 0x00: NOP
		// Do nothing
		func(cpu *CPU) int {
			return 4
		},
		// 0x01: LD BC, nn
		// Load value nn to BC
		func(cpu *CPU) int {
			cpu.C = cpu.Memory.Read(cpu.PC)
			cpu.B = cpu.Memory.Read(cpu.PC + 1)
			cpu.PC += 2
			return 12
		},
		// 0x02: LD (BC), A
		// Store A into memory pointed by BC
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.BC(), cpu.A)
			return 8
		},
		// 0x03: INC BC
		// Adds one to BC
		func(cpu *CPU) int {
			bc := cpu.BC() + 1
			cpu.B = uint8(bc >> 8)
			cpu.C = uint8(bc)
			return 8
		},
		// 0x04: INC B
		// Adds one to B
		func(cpu *CPU) int {
			cpu.B++
			cpu.FSub = false
			cpu.FHalfCarry = cpu.B&0xf == 0
			cpu.FZero = cpu.B == 0
			return 4
		},
		// 0x05: DEC B
		// Subtracts one from B
		func(cpu *CPU) int {
			cpu.B--
			cpu.FSub = true
			cpu.FHalfCarry = cpu.B&0xf == 0xf
			cpu.FZero = cpu.B == 0
			return 4
		},
		// 0x06: LD B, n
		// Subtracts one from B
		func(cpu *CPU) int {
			cpu.B--
			cpu.FSub = true
			cpu.FHalfCarry = cpu.B&0xf == 0xf
			cpu.FZero = cpu.B == 0
			return 8
		},
		// 0x07: RCLA
		// Rotate A left through carry
		func(cpu *CPU) int {
			carry := cpu.A >> 7
			cpu.A = (cpu.A << 1) | carry
			cpu.FCarry = carry == 1
			cpu.FZero = false
			cpu.FSub = false
			cpu.FHalfCarry = false
			return 4
		},
	}
}
