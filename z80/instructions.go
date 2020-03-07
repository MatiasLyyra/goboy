package z80

import "os"

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
		// 0x01: LD BC, d16
		// Load value d16 to BC
		func(cpu *CPU) int {
			cpu.C = cpu.Memory.Read(cpu.PC + 1)
			cpu.B = cpu.Memory.Read(cpu.PC)
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
			cpu.FZero = cpu.B == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.B&0xf == 0
			return 4
		},
		// 0x05: DEC B
		// Subtracts one from B
		func(cpu *CPU) int {
			cpu.B--
			cpu.FZero = cpu.B == 0
			cpu.FSub = true
			cpu.FHalfCarry = cpu.B&0xf == 0xf
			return 4
		},
		// 0x06: LD B, n
		// Loads value n to B
		func(cpu *CPU) int {
			cpu.B = cpu.Memory.Read(cpu.PC)
			cpu.PC++
			return 8
		},
		// 0x07: RLCA
		// Rotate A left through carry
		func(cpu *CPU) int {
			carry := cpu.A >> 7
			cpu.A = (cpu.A << 1) | carry
			cpu.FZero = false
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = carry == 1
			return 4
		},
		// 0x08: LD (a16), SP
		// Stores SP into address a16
		func(cpu *CPU) int {
			addr := uint16(cpu.Memory.Read(cpu.PC+1))<<8 | uint16(cpu.Memory.Read(cpu.PC))
			cpu.Memory.Write(addr, uint8(cpu.SP))
			cpu.Memory.Write(addr+1, uint8(cpu.SP>>8))
			cpu.PC += 2
			return 20
		},
		// 0x09: ADD HL, BC
		// HL = HL + BC
		func(cpu *CPU) int {
			// Make sure that temp doesn't overflow by making it uint32
			temp := uint32(cpu.HL()) + uint32(cpu.BC())
			cpu.FSub = false
			cpu.FHalfCarry = cpu.HL()&0xfff > uint16(temp&0xfff)
			cpu.FCarry = temp > 0xffff
			cpu.H = uint8(temp >> 8)
			cpu.L = uint8(temp)
			return 8
		},
		// 0x0A: LD A, (BC)
		// Load value pointed by BC to A
		func(cpu *CPU) int {
			cpu.A = cpu.Memory.Read(cpu.BC())
			return 8
		},
		// 0x0B: DEC BC
		// Substract one from BC
		func(cpu *CPU) int {
			temp := cpu.BC() - 1
			cpu.B = uint8(temp >> 8)
			cpu.C = uint8(temp)
			return 8
		},
		// 0x0C: INC C
		// Add one to C
		func(cpu *CPU) int {
			cpu.C++
			cpu.FZero = cpu.C == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.C&0xf == 0
			return 4
		},
		// 0x0D: DEC C
		// Subtract one from C
		func(cpu *CPU) int {
			cpu.C--
			cpu.FZero = cpu.C == 0
			cpu.FSub = true
			cpu.FHalfCarry = cpu.C&0xf == 0xf
			return 4
		},
		// 0x0E: LD C, d8
		// Loads value d8 to B
		func(cpu *CPU) int {
			cpu.C = cpu.Memory.Read(cpu.PC)
			cpu.PC++
			return 8
		},
		// 0x0f: RRCA
		// Rotate A right through carry
		func(cpu *CPU) int {
			carry := cpu.A & 1
			cpu.A = (cpu.A >> 1) | (carry << 7)
			cpu.FZero = false
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = carry == 1
			return 4
		},
		// 0x10f: RRCA
		// Rotate A right through carry
		func(cpu *CPU) int {
			cpu.PC++
			// TODO: Something better than just os.Exit...
			os.Exit(0)
			return 4
		},
		// 0x11: LD DE, d16
		// Load value d16 to DE
		func(cpu *CPU) int {
			cpu.D = cpu.Memory.Read(cpu.PC + 1)
			cpu.E = cpu.Memory.Read(cpu.PC)
			cpu.PC += 2
			return 12
		},
		// 0x12: LD (DE), A
		// Store A into memory pointed by DE
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.DE(), cpu.A)
			return 8
		},
		// 0x13: INC DE
		// Adds one to DE
		func(cpu *CPU) int {
			de := cpu.DE() + 1
			cpu.D = uint8(de >> 8)
			cpu.E = uint8(de)
			return 8
		},
		// 0x14: INC D
		// Adds one to D
		func(cpu *CPU) int {
			cpu.D++
			cpu.FZero = cpu.D == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.D&0xf == 0
			return 4
		},
		// 0x15: DEC D
		// Subtract one from D
		func(cpu *CPU) int {
			cpu.D--
			cpu.FZero = cpu.D == 0
			cpu.FSub = true
			cpu.FHalfCarry = cpu.D&0xf == 0xf
			return 4
		},
		// 0x16: LD D, d8
		// Load value d8 to D
		func(cpu *CPU) int {
			cpu.D = cpu.Memory.Read(cpu.PC)
			cpu.PC++
			return 8
		},
		// 0x17: RLA
		// Rotate A left through carry and insert previous carry at bit position 0
		func(cpu *CPU) int {
			var (
				prevCarry uint8
				carry     uint8
			)
			if cpu.FCarry {
				prevCarry = 1
			}
			carry = cpu.A >> 7 & 1
			cpu.A = (cpu.A << 1) | prevCarry
			cpu.FZero = false
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = carry == 1
			return 4
		},
		// 0x18: JR r8
		// Relative jump to r8
		func(cpu *CPU) int {
			relJump := int32(int8(cpu.Memory.Read(cpu.PC)))
			// Add one because of r8
			tempPC := int32(cpu.PC) + relJump + 1
			cpu.PC = uint16(tempPC)
			return 12
		},
		// 0x19: ADD HL, DE
		// HL = HL + DE
		func(cpu *CPU) int {
			// Make sure that temp doesn't overflow by making it uint32
			temp := uint32(cpu.HL()) + uint32(cpu.DE())
			cpu.FSub = false
			cpu.FHalfCarry = cpu.HL()&0xfff > uint16(temp&0xfff)
			cpu.FCarry = temp > 0xffff
			cpu.H = uint8(temp >> 8)
			cpu.L = uint8(temp)
			return 8
		},
		// 0x1A: LD A, (DE)
		// Load value pointed by BC to A
		func(cpu *CPU) int {
			cpu.A = cpu.Memory.Read(cpu.DE())
			return 8
		},
		// 0x1B: DEC DE
		// Substract one from DE
		func(cpu *CPU) int {
			temp := cpu.DE() - 1
			cpu.B = uint8(temp >> 8)
			cpu.C = uint8(temp)
			return 8
		},
		// 0x1C: INC E
		// Add one to E
		func(cpu *CPU) int {
			cpu.E++
			cpu.FZero = cpu.E == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.E&0xf == 0
			return 4
		},
		// 0x1D: DEC E
		// Subtract one from E
		func(cpu *CPU) int {
			cpu.E--
			cpu.FZero = cpu.E == 0
			cpu.FSub = true
			cpu.FHalfCarry = cpu.E&0xf == 0xf
			return 4
		},
		// 0x1E: LD E, d8
		// Loads value d8 to E
		func(cpu *CPU) int {
			cpu.E = cpu.Memory.Read(cpu.PC)
			cpu.PC++
			return 8
		},
		// 0x1F: RRA
		// Rotate A right through carry and insert previous carry at bit position 7
		func(cpu *CPU) int {
			var (
				prevCarry uint8
				carry     uint8
			)
			if cpu.FCarry {
				prevCarry = 1
			}
			carry = cpu.A & 1
			cpu.A = (cpu.A >> 1) | (prevCarry << 7)
			cpu.FZero = false
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = carry == 1
			return 4
		},
		// 0x20: JR NZ, r8
		// Relative jump to d8 if Z flag is not set
		func(cpu *CPU) int {
			if cpu.FZero {
				cpu.PC++
				return 8
			}
			relJump := int32(int8(cpu.Memory.Read(cpu.PC)))
			cpu.PC = uint16(int32(cpu.PC) + relJump + 1)
			return 12
		},
		// 0x21: LD HL, d16
		// Load value d16 to HL
		func(cpu *CPU) int {
			cpu.H = cpu.Memory.Read(cpu.PC + 1)
			cpu.L = cpu.Memory.Read(cpu.PC)
			cpu.PC += 2
			return 12
		},
		// 0x22: LD (HL+), A
		// Store A into memory pointed by HL and increment HL after it
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.HL(), cpu.A)
			temp := cpu.HL() + 1
			cpu.H = uint8(temp >> 8)
			cpu.L = uint8(temp)
			return 8
		},
		// 0x23: INC HL
		// Adds one to HL
		func(cpu *CPU) int {
			temp := cpu.HL() + 1
			cpu.H = uint8(temp >> 8)
			cpu.L = uint8(temp)
			return 8
		},
		// 0x24: INC H
		// Adds one to H
		func(cpu *CPU) int {
			cpu.H++
			cpu.FZero = cpu.H == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.H&0xf == 0
			return 4
		},
		// 0x25: DEC H
		// Subtract one from D
		func(cpu *CPU) int {
			cpu.H--
			cpu.FZero = cpu.H == 0
			cpu.FSub = true
			cpu.FHalfCarry = cpu.H&0xf == 0xf
			return 4
		},
		// 0x26: LD H, d8
		// Load value d8 to H
		func(cpu *CPU) int {
			cpu.H = cpu.Memory.Read(cpu.PC)
			cpu.PC++
			return 8
		},
		// 0x27: DAA
		// Conditionally adjust A for BCD representation
		func(cpu *CPU) int {
			temp := int(cpu.A)
			if !cpu.FSub {
				if cpu.FHalfCarry || (temp&0xf) > 9 {
					temp += 6
				}
				if cpu.FCarry || (temp>>4) > 9 {
					temp += 0x60
					cpu.FCarry = true
				}
			} else {
				if cpu.FHalfCarry {
					temp = (temp - 6) & 0xff
				}
				if cpu.FCarry {
					temp = (temp - 0x60) & 0xff
				}
			}
			cpu.FHalfCarry = false
			cpu.FZero = temp == 0
			cpu.A = uint8(temp)
			return 4
		},
		// 0x28: JR Z, r8
		// Relative jump to d8 if Z flag is set
		func(cpu *CPU) int {
			if !cpu.FZero {
				cpu.PC++
				return 8
			}
			relJump := int32(int8(cpu.Memory.Read(cpu.PC)))
			cpu.PC = uint16(int32(cpu.PC) + relJump + 1)
			return 12
		},
		// 0x29: ADD HL, HL
		// HL = HL + HL
		func(cpu *CPU) int {
			// Make sure that temp doesn't overflow by making it uint32
			temp := uint32(cpu.HL()) + uint32(cpu.HL())
			cpu.FSub = false
			cpu.FHalfCarry = cpu.HL()&0xfff > uint16(temp&0xfff)
			cpu.FCarry = temp > 0xffff
			cpu.H = uint8(temp >> 8)
			cpu.L = uint8(temp)
			return 8
		},
		// 0x2A: LD A, (HL+)
		func(cpu *CPU) int {
			cpu.A = cpu.Memory.Read(cpu.HL())
			temp := cpu.HL() + 1
			cpu.H = uint8(temp >> 8)
			cpu.L = uint8(temp)
			return 8
		},
		// 0x2B: DEC HL
		// Substract one from HL
		func(cpu *CPU) int {
			temp := cpu.HL() - 1
			cpu.H = uint8(temp >> 8)
			cpu.L = uint8(temp)
			return 8
		},
		// 0x2C: INC L
		// Add one to L
		func(cpu *CPU) int {
			cpu.L++
			cpu.FZero = cpu.L == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.L&0xf == 0
			return 4
		},
		// 0x2D: DEC L
		// Subtract one from L
		func(cpu *CPU) int {
			cpu.L--
			cpu.FZero = cpu.L == 0
			cpu.FSub = true
			cpu.FHalfCarry = cpu.L&0xf == 0xf
			return 4
		},
		// 0x2E: LD L, d8
		// Loads value d8 to L
		func(cpu *CPU) int {
			cpu.L = cpu.Memory.Read(cpu.PC)
			cpu.PC++
			return 8
		},
		// 0x2F: CPL
		// A = ~A
		func(cpu *CPU) int {
			cpu.A = ^cpu.A
			cpu.FSub = true
			cpu.FHalfCarry = true
			return 4
		},
		// 0x30: JR NC, r8
		// Relative jump to d8 if C flag is not set
		func(cpu *CPU) int {
			if cpu.FCarry {
				cpu.PC++
				return 8
			}
			relJump := int32(int8(cpu.Memory.Read(cpu.PC)))
			cpu.PC = uint16(int32(cpu.PC) + relJump + 1)
			return 12
		},
		// 0x31: LD SP, d16
		// Load value d16 to HL
		func(cpu *CPU) int {
			cpu.SP = uint16(cpu.Memory.Read(cpu.PC+1))<<8 | uint16(cpu.Memory.Read(cpu.PC))
			cpu.PC += 2
			return 12
		},
		// 0x32: LD (HL-), A
		// Store A into memory pointed by HL and decrement HL after it
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.HL(), cpu.A)
			temp := cpu.HL() - 1
			cpu.H = uint8(temp >> 8)
			cpu.L = uint8(temp)
			return 8
		},
		// 0x33: INC SP
		// Add one to SP
		func(cpu *CPU) int {
			cpu.SP++
			return 8
		},
		// 0x34: INC (HL)
		// Add one value pointed by HL
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL()) + 1
			cpu.FZero = val == 0
			cpu.FSub = false
			cpu.FHalfCarry = val&0xf == 0
			cpu.Memory.Write(cpu.HL(), val)
			return 12
		},
		// 0x35: DEC (HL)
		// Subtract one from value pointed by HL
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL()) - 1
			cpu.FZero = val == 0
			cpu.FSub = true
			cpu.FHalfCarry = val&0xf == 0xf
			cpu.Memory.Write(cpu.HL(), val)
			return 12
		},
		// 0x36: LD (HL), d8
		// Store value d8 to byte pointed by HL
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.HL(), cpu.Memory.Read(cpu.PC))
			cpu.PC++
			return 12
		},
		// 0x37: SCF
		// Set C flag
		func(cpu *CPU) int {
			cpu.FCarry = true
			return 4
		},
		// 0x38: JR C, r8
		// Relative jump to d8 if C flag is set
		func(cpu *CPU) int {
			if !cpu.FCarry {
				cpu.PC++
				return 8
			}
			relJump := int32(int8(cpu.Memory.Read(cpu.PC)))
			cpu.PC = uint16(int32(cpu.PC) + relJump + 1)
			return 12
		},
		// 0x39: ADD HL, SP
		// HL = HL + SP
		func(cpu *CPU) int {
			// Make sure that temp doesn't overflow by making it uint32
			temp := uint32(cpu.HL()) + uint32(cpu.SP)
			cpu.FSub = false
			cpu.FHalfCarry = cpu.HL()&0xfff > uint16(temp&0xfff)
			cpu.FCarry = temp > 0xffff
			cpu.H = uint8(temp >> 8)
			cpu.L = uint8(temp)
			return 8
		},
		// 0x3A: LD A, (HL-)
		// Load value pointed by HL and decrement HL after it
		func(cpu *CPU) int {
			cpu.A = cpu.Memory.Read(cpu.HL())
			temp := cpu.HL() - 1
			cpu.H = uint8(temp >> 8)
			cpu.L = uint8(temp)
			return 8
		},
		// 0x3B: DEC SP
		// Substract one from SP
		func(cpu *CPU) int {
			cpu.SP--
			return 8
		},
		// 0x3C: INC A
		// Add one to A
		func(cpu *CPU) int {
			cpu.A++
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf == 0
			return 4
		},
		// 0x3D: DEC A
		// Subtract one from A
		func(cpu *CPU) int {
			cpu.A--
			cpu.FZero = cpu.A == 0
			cpu.FSub = true
			cpu.FHalfCarry = cpu.A&0xf == 0xf
			return 4
		},
		// 0x2E: LD A, d8
		// Loads value d8 to A
		func(cpu *CPU) int {
			cpu.A = cpu.Memory.Read(cpu.PC)
			cpu.PC++
			return 8
		},
		// 0x2F: CCF
		// Flag C = ! Flag C
		func(cpu *CPU) int {
			cpu.FSub = true
			cpu.FHalfCarry = true
			cpu.FCarry = !cpu.FCarry
			return 4
		},
	}
}
