package goboy

// MicroCodeFunc executes single instructions
type MicroCodeFunc func(cpu *CPU) int

// InstructionsTable contains all of the unprefixed instructions
var InstructionsTable [256]MicroCodeFunc

func init() {
	// TODO: Clean duplicate instructions into a function and combine them
	// like with
	InstructionsTable = [256]MicroCodeFunc{
		// 0x00: NOP
		// Do nothing
		func(cpu *CPU) int {
			return 4
		},
		// 0x01: LD BC, d16
		// Load value d16 to BC
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
		// 0x10: ???
		// ????
		func(cpu *CPU) int {
			cpu.PC++
			// TODO: Something better than just os.Exit...
			// fmt.Printf("exited at %x\n", cpu.PC)
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
				if cpu.FHalfCarry || temp&0xf > 9 {
					temp += 6
				}
				if cpu.FCarry || temp > 0x9F {
					temp += 0x60
					cpu.FCarry = true
				}
			} else {
				if cpu.FHalfCarry {
					temp -= 6
				}
				if cpu.FCarry {
					temp -= 0x60
				}
			}
			cpu.FHalfCarry = false
			cpu.FZero = temp&0xff == 0
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
		// 0x3E: LD A, d8
		// Loads value d8 to A
		func(cpu *CPU) int {
			cpu.A = cpu.Memory.Read(cpu.PC)
			cpu.PC++
			return 8
		},
		// 0x3F: CCF
		// Flag C = ! Flag C
		func(cpu *CPU) int {
			cpu.FSub = true
			cpu.FHalfCarry = true
			cpu.FCarry = !cpu.FCarry
			return 4
		},
		// 0x:40: LD B, B
		func(cpu *CPU) int {
			// NOOP
			return 4
		},
		// 0x41: LD B, C
		func(cpu *CPU) int {
			cpu.B = cpu.C
			return 4
		},
		// 0x42: LD B, D
		func(cpu *CPU) int {
			cpu.B = cpu.D
			return 4
		},
		// 0x43: LD B, E
		func(cpu *CPU) int {
			cpu.B = cpu.E
			return 4
		},
		// 0x44: LD B, H
		func(cpu *CPU) int {
			cpu.B = cpu.H
			return 4
		},
		// 0x45: LD B, L
		func(cpu *CPU) int {
			cpu.B = cpu.L
			return 4
		},
		// 0x46: LD B, (HL)
		func(cpu *CPU) int {
			cpu.B = cpu.Memory.Read(cpu.HL())
			return 8
		},
		// 0x47: LD B, A
		func(cpu *CPU) int {
			cpu.B = cpu.A
			return 4
		},
		// 0x:48: LD C, B
		func(cpu *CPU) int {
			cpu.C = cpu.B
			return 4
		},
		// 0x49: LD C, C
		func(cpu *CPU) int {
			// NOOP
			return 4
		},
		// 0x4A: LD C, D
		func(cpu *CPU) int {
			cpu.C = cpu.D
			return 4
		},
		// 0x4B: LD C, E
		func(cpu *CPU) int {
			cpu.C = cpu.E
			return 4
		},
		// 0x4C: LD C, H
		func(cpu *CPU) int {
			cpu.C = cpu.H
			return 4
		},
		// 0x4D: LD C, L
		func(cpu *CPU) int {
			cpu.C = cpu.L
			return 4
		},
		// 0x4E: LD C, (HL)
		func(cpu *CPU) int {
			cpu.C = cpu.Memory.Read(cpu.HL())
			return 8
		},
		// 0x4F: LD C, A
		func(cpu *CPU) int {
			cpu.C = cpu.A
			return 4
		},
		// 0x:50: LD D, B
		func(cpu *CPU) int {
			cpu.D = cpu.B
			return 4
		},
		// 0x51: LD D, C
		func(cpu *CPU) int {
			cpu.D = cpu.C
			return 4
		},
		// 0x52: LD D, D
		func(cpu *CPU) int {
			// NOOP
			return 4
		},
		// 0x53: LD D, E
		func(cpu *CPU) int {
			cpu.D = cpu.E
			return 4
		},
		// 0x54: LD D, H
		func(cpu *CPU) int {
			cpu.D = cpu.H
			return 4
		},
		// 0x55: LD D, L
		func(cpu *CPU) int {
			cpu.D = cpu.L
			return 4
		},
		// 0x56: LD D, (HL)
		func(cpu *CPU) int {
			cpu.D = cpu.Memory.Read(cpu.HL())
			return 8
		},
		// 0x57: LD D, A
		func(cpu *CPU) int {
			cpu.D = cpu.A
			return 4
		},
		// 0x:58: LD E, B
		func(cpu *CPU) int {
			cpu.E = cpu.B
			return 4
		},
		// 0x59: LD E, C
		func(cpu *CPU) int {
			cpu.E = cpu.C
			return 4
		},
		// 0x5A: LD E, D
		func(cpu *CPU) int {
			cpu.E = cpu.D
			return 4
		},
		// 0x5B: LD E, E
		func(cpu *CPU) int {
			// NOOP
			return 4
		},
		// 0x5C: LD E, H
		func(cpu *CPU) int {
			cpu.E = cpu.H
			return 4
		},
		// 0x5D: LD E, L
		func(cpu *CPU) int {
			cpu.E = cpu.L
			return 4
		},
		// 0x5E: LD E, (HL)
		func(cpu *CPU) int {
			cpu.E = cpu.Memory.Read(cpu.HL())
			return 8
		},
		// 0x5F: LD E, A
		func(cpu *CPU) int {
			cpu.E = cpu.A
			return 4
		},
		// 0x:60: LD H, B
		func(cpu *CPU) int {
			cpu.H = cpu.B
			return 4
		},
		// 0x61: LD H, C
		func(cpu *CPU) int {
			cpu.H = cpu.C
			return 4
		},
		// 0x62: LD H, D
		func(cpu *CPU) int {
			cpu.H = cpu.D
			return 4
		},
		// 0x63: LD H, E
		func(cpu *CPU) int {
			cpu.H = cpu.E
			return 4
		},
		// 0x64: LD H, H
		func(cpu *CPU) int {
			cpu.D = cpu.H
			return 4
		},
		// 0x65: LD H, L
		func(cpu *CPU) int {
			cpu.H = cpu.L
			return 4
		},
		// 0x66: LD H, (HL)
		func(cpu *CPU) int {
			cpu.H = cpu.Memory.Read(cpu.HL())
			return 8
		},
		// 0x67: LD H, A
		func(cpu *CPU) int {
			cpu.H = cpu.A
			return 4
		},
		// 0x:68: LD L, B
		func(cpu *CPU) int {
			cpu.L = cpu.B
			return 4
		},
		// 0x69: LD L, C
		func(cpu *CPU) int {
			cpu.L = cpu.C
			return 4
		},
		// 0x6A: LD L, D
		func(cpu *CPU) int {
			cpu.L = cpu.D
			return 4
		},
		// 0x6B: LD L, E
		func(cpu *CPU) int {
			cpu.L = cpu.E
			return 4
		},
		// 0x6C: LD L, H
		func(cpu *CPU) int {
			cpu.E = cpu.H
			return 4
		},
		// 0x6D: LD L, L
		func(cpu *CPU) int {
			// NOOP
			return 4
		},
		// 0x6E: LD L, (HL)
		func(cpu *CPU) int {
			cpu.L = cpu.Memory.Read(cpu.HL())
			return 8
		},
		// 0x6F: LD L, A
		func(cpu *CPU) int {
			cpu.L = cpu.A
			return 4
		},
		// 0x70: LD (HL), B
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.HL(), cpu.B)
			return 8
		},
		// 0x71: LD (HL), C
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.HL(), cpu.C)
			return 8
		},
		// 0x72: LD (HL), D
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.HL(), cpu.D)
			return 8
		},
		// 0x73: LD (HL), E
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.HL(), cpu.E)
			return 8
		},
		// 0x74: LD (HL), H
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.HL(), cpu.H)
			return 8
		},
		// 0x75: LD (HL), L
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.HL(), cpu.L)
			return 8
		},
		// 0x76: HALT
		func(cpu *CPU) int {
			cpu.Halt = true
			return 4
		},
		// 0x77: LD (HL), A
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.HL(), cpu.A)
			return 8
		},
		// 0x78: LD A, B
		func(cpu *CPU) int {
			cpu.A = cpu.B
			return 4
		},
		// 0x79: LD A, C
		func(cpu *CPU) int {
			cpu.A = cpu.C
			return 4
		},
		// 0x7A: LD A, D
		func(cpu *CPU) int {
			cpu.A = cpu.D
			return 4
		},
		// 0x7B: LD A, E
		func(cpu *CPU) int {
			cpu.A = cpu.E
			return 4
		},
		// 0x7C: LD A, H
		func(cpu *CPU) int {
			cpu.A = cpu.H
			return 4
		},
		// 0x7D: LD A, L
		func(cpu *CPU) int {
			cpu.A = cpu.L
			return 4
		},
		// 0x7E: LD A, (HL)
		func(cpu *CPU) int {
			cpu.A = cpu.Memory.Read(cpu.HL())
			return 8
		},
		// 0x7F: LD A, A
		func(cpu *CPU) int {
			// NOOP
			return 4
		},
		// 0x80: ADD A, B
		func(cpu *CPU) int {
			temp := uint16(cpu.A) + uint16(cpu.B)
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			return 4
		},
		// 0x81: ADD A, C
		func(cpu *CPU) int {
			temp := uint16(cpu.A) + uint16(cpu.C)
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			return 4
		},
		// 0x82: ADD A, D
		func(cpu *CPU) int {
			temp := uint16(cpu.A) + uint16(cpu.D)
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			return 4
		},
		// 0x83: ADD A, E
		func(cpu *CPU) int {
			temp := uint16(cpu.A) + uint16(cpu.E)
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			return 4
		},
		// 0x84: ADD A, H
		func(cpu *CPU) int {
			temp := uint16(cpu.A) + uint16(cpu.H)
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			return 4
		},
		// 0x85: ADD A, L
		func(cpu *CPU) int {
			temp := uint16(cpu.A) + uint16(cpu.L)
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			return 4
		},
		// 0x86: ADD A, (HL)
		func(cpu *CPU) int {
			temp := uint16(cpu.A) + uint16(cpu.Memory.Read(cpu.HL()))
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			return 8
		},
		// 0x87: ADD A, A
		func(cpu *CPU) int {
			temp := uint16(cpu.A) + uint16(cpu.A)
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			return 4
		},
		// 0x88: ADC A, B
		func(cpu *CPU) int {
			var carry uint16
			if cpu.FCarry {
				carry = 1
			}
			temp := uint16(cpu.A) + uint16(cpu.B) + carry
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			return 4
		},
		// 0x89: ADC A, C
		func(cpu *CPU) int {
			var carry uint16
			if cpu.FCarry {
				carry = 1
			}
			temp := uint16(cpu.A) + uint16(cpu.C) + carry
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			return 4
		},
		// 0x8A: ADC A, D
		func(cpu *CPU) int {
			var carry uint16
			if cpu.FCarry {
				carry = 1
			}
			temp := uint16(cpu.A) + uint16(cpu.D) + carry
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			return 4
		},
		// 0x8B: ADC A, E
		func(cpu *CPU) int {
			var carry uint16
			if cpu.FCarry {
				carry = 1
			}
			temp := uint16(cpu.A) + uint16(cpu.E) + carry
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			return 4
		},
		// 0x8C: ADC A, H
		func(cpu *CPU) int {
			var carry uint16
			if cpu.FCarry {
				carry = 1
			}
			temp := uint16(cpu.A) + uint16(cpu.H) + carry
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			return 4
		},
		// 0x8D: ADC A, L
		func(cpu *CPU) int {
			var carry uint16
			if cpu.FCarry {
				carry = 1
			}
			temp := uint16(cpu.A) + uint16(cpu.L) + carry
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			return 4
		},
		// 0x8E: ADC A, (HL)
		func(cpu *CPU) int {
			var carry uint16
			if cpu.FCarry {
				carry = 1
			}
			temp := uint16(cpu.A) + uint16(cpu.Memory.Read(cpu.HL())) + carry
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			return 8
		},
		// 0x8F: ADC A, A
		func(cpu *CPU) int {
			var carry uint16
			if cpu.FCarry {
				carry = 1
			}
			temp := uint16(cpu.A) + uint16(cpu.A) + carry
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			return 4
		},
		// 0x90: SUB A, B
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.B)
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = cpu.A&0xf < uint8(temp)&0xf
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			return 4
		},
		// 0x91: SUB A, C
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.C)
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = cpu.A&0xf < uint8(temp)&0xf
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			return 4
		},
		// 0x92: SUB A, D
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.D)
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = cpu.A&0xf < uint8(temp)&0xf
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			return 4
		},
		// 0x93: SUB A, E
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.E)
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = cpu.A&0xf < uint8(temp)&0xf
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			return 4
		},
		// 0x94: SUB A, H
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.E)
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = cpu.A&0xf < uint8(temp)&0xf
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			return 4
		},
		// 0x95: SUB A, L
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.E)
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = cpu.A&0xf < uint8(temp)&0xf
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			return 4
		},
		// 0x96: SUB A, (HL)
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.Memory.Read(cpu.HL()))
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = cpu.A&0xf < uint8(temp)&0xf
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			return 8
		},
		// 0x97: SUB A, A
		func(cpu *CPU) int {
			cpu.FZero = true
			cpu.FSub = true
			cpu.FHalfCarry = false
			cpu.FCarry = false
			cpu.A = 0
			return 4
		},
		// 0x98: SBC A, B
		func(cpu *CPU) int {
			var carry int
			if cpu.FCarry {
				carry = 1
			}
			temp := int(cpu.A) - int(cpu.B) - carry
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = int(cpu.A&0xf)-int(cpu.B&0xf)-carry < 0
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			return 4
		},
		// 0x99: SBC A, C
		func(cpu *CPU) int {
			var carry int
			if cpu.FCarry {
				carry = 1
			}
			temp := int(cpu.A) - int(cpu.C) - carry
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = int(cpu.A&0xf)-int(cpu.C&0xf)-carry < 0
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			return 4
		},
		// 0x9A: SBC A, D
		func(cpu *CPU) int {
			var carry int
			if cpu.FCarry {
				carry = 1
			}
			temp := int(cpu.A) - int(cpu.D) - carry
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = int(cpu.A&0xf)-int(cpu.D&0xf)-carry < 0
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			return 4
		},
		// 0x9B: SBC A, E
		func(cpu *CPU) int {
			var carry int
			if cpu.FCarry {
				carry = 1
			}
			temp := int(cpu.A) - int(cpu.E) - carry
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = int(cpu.A&0xf)-int(cpu.E&0xf)-carry < 0
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			return 4
		},
		// 0x9C: SBC A, H
		func(cpu *CPU) int {
			var carry int
			if cpu.FCarry {
				carry = 1
			}
			temp := int(cpu.A) - int(cpu.H) - carry
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = int(cpu.A&0xf)-int(cpu.H&0xf)-carry < 0
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			return 4
		},
		// 0x9D: SBC A, L
		func(cpu *CPU) int {
			var carry int
			if cpu.FCarry {
				carry = 1
			}
			temp := int(cpu.A) - int(cpu.L) - carry
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = int(cpu.A&0xf)-int(cpu.L&0xf)-carry < 0
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			return 4
		},
		// 0x9E: SBC A, (HL)
		func(cpu *CPU) int {
			var carry int
			if cpu.FCarry {
				carry = 1
			}
			hl := cpu.Memory.Read(cpu.HL())
			temp := int(cpu.A) - int(hl) - carry
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = int(cpu.A&0xf)-int(hl&0xf)-carry < 0
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			return 8
		},
		// 0x9D: SBC A, A
		func(cpu *CPU) int {
			var carry int
			if cpu.FCarry {
				carry = 1
			}
			temp := int(cpu.A) - int(cpu.A) - carry
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = int(cpu.A&0xf)-int(cpu.A&0xf)-carry < 0
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			return 4
		},
		// 0xA0: AND B
		func(cpu *CPU) int {
			cpu.A = cpu.A & cpu.B
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = true
			cpu.FCarry = false
			return 4
		},
		// 0xA1: AND C
		func(cpu *CPU) int {
			cpu.A = cpu.A & cpu.C
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = true
			cpu.FCarry = false
			return 4
		},
		// 0xA2: AND D
		func(cpu *CPU) int {
			cpu.A = cpu.A & cpu.D
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = true
			cpu.FCarry = false
			return 4
		},
		// 0xA3: AND E
		func(cpu *CPU) int {
			cpu.A = cpu.A & cpu.E
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = true
			cpu.FCarry = false
			return 4
		},
		// 0xA4: AND H
		func(cpu *CPU) int {
			cpu.A = cpu.A & cpu.H
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = true
			cpu.FCarry = false
			return 4
		},
		// 0xA5: AND L
		func(cpu *CPU) int {
			cpu.A = cpu.A & cpu.L
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = true
			cpu.FCarry = false
			return 4
		},
		// 0xA6: AND (HL)
		func(cpu *CPU) int {
			cpu.A = cpu.A & cpu.Memory.Read(cpu.HL())
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = true
			cpu.FCarry = false
			return 8
		},
		// 0xA7: AND A
		func(cpu *CPU) int {
			cpu.A = cpu.A & cpu.A
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = true
			cpu.FCarry = false
			return 4
		},
		// 0xA8: XOR B
		func(cpu *CPU) int {
			cpu.A = cpu.A ^ cpu.B
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			return 4
		},
		// 0xA9: XOR C
		func(cpu *CPU) int {
			cpu.A = cpu.A ^ cpu.C
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			return 4
		},
		// 0xAA: XOR D
		func(cpu *CPU) int {
			cpu.A = cpu.A ^ cpu.D
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			return 4
		},
		// 0xAB: XOR E
		func(cpu *CPU) int {
			cpu.A = cpu.A ^ cpu.E
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			return 4
		},
		// 0xAC: XOR H
		func(cpu *CPU) int {
			cpu.A = cpu.A ^ cpu.H
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			return 4
		},
		// 0xAD: XOR L
		func(cpu *CPU) int {
			cpu.A = cpu.A ^ cpu.L
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			return 4
		},
		// 0xAE: XOR (HL)
		func(cpu *CPU) int {
			cpu.A = cpu.A ^ cpu.Memory.Read(cpu.HL())
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			return 8
		},
		// 0xAF: XOR A
		func(cpu *CPU) int {
			cpu.A = cpu.A ^ cpu.A
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			return 4
		},
		// 0xB0: OR B
		func(cpu *CPU) int {
			cpu.A = cpu.A | cpu.B
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			return 4
		},
		// 0xB1: OR C
		func(cpu *CPU) int {
			cpu.A = cpu.A | cpu.C
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			return 4
		},
		// 0xB2: OR D
		func(cpu *CPU) int {
			cpu.A = cpu.A | cpu.D
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			return 4
		},
		// 0xB3: OR E
		func(cpu *CPU) int {
			cpu.A = cpu.A | cpu.E
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			return 4
		},
		// 0xB4: OR H
		func(cpu *CPU) int {
			cpu.A = cpu.A | cpu.H
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			return 4
		},
		// 0xB5: OR L
		func(cpu *CPU) int {
			cpu.A = cpu.A | cpu.L
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			return 4
		},
		// 0xB6: OR (HL)
		func(cpu *CPU) int {
			cpu.A = cpu.A | cpu.Memory.Read(cpu.HL())
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			return 8
		},
		// 0xB7: OR A
		func(cpu *CPU) int {
			cpu.A = cpu.A | cpu.A
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			return 4
		},
		// 0xB8: CP B
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.B)
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = temp&0xf > int(cpu.A&0xf)
			cpu.FCarry = temp < 0
			return 4
		},
		// 0xB9: CP C
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.C)
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = temp&0xf > int(cpu.A&0xf)
			cpu.FCarry = temp < 0
			return 4
		},
		// 0xBA: CP D
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.D)
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = temp&0xf > int(cpu.A&0xf)
			cpu.FCarry = temp < 0
			return 4
		},
		// 0xBB: CP E
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.E)
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = temp&0xf > int(cpu.A&0xf)
			cpu.FCarry = temp < 0
			return 4
		},
		// 0xBC: CP H
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.H)
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = temp&0xf > int(cpu.A&0xf)
			cpu.FCarry = temp < 0
			return 4
		},
		// 0xBD: CP L
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.L)
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = temp&0xf > int(cpu.A&0xf)
			cpu.FCarry = temp < 0
			return 4
		},
		// 0xBE: CP (HL)
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.Memory.Read(cpu.HL()))
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = temp&0xf > int(cpu.A&0xf)
			cpu.FCarry = temp < 0
			return 8
		},
		// 0xBF: CP A
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.A)
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = temp&0xf > int(cpu.A&0xf)
			cpu.FCarry = temp < 0
			return 4
		},
		// 0xC0: RET NZ
		func(cpu *CPU) int {
			if cpu.FZero {
				return 8
			}
			cpu.PC = uint16(cpu.Memory.Read(cpu.SP+1))<<8 | uint16(cpu.Memory.Read(cpu.SP))
			cpu.SP += 2
			return 20
		},
		// 0xC1: POP BC
		func(cpu *CPU) int {
			cpu.C = cpu.Memory.Read(cpu.SP)
			cpu.B = cpu.Memory.Read(cpu.SP + 1)
			cpu.SP += 2
			return 12
		},
		// 0xC2: JP NZ, a16
		func(cpu *CPU) int {
			if cpu.FZero {
				cpu.PC += 2
				return 12
			}
			cpu.PC = uint16(cpu.Memory.Read(cpu.PC+1))<<8 | uint16(cpu.Memory.Read(cpu.PC))
			return 16
		},
		// 0xC3: JP a16
		func(cpu *CPU) int {
			cpu.PC = uint16(cpu.Memory.Read(cpu.PC+1))<<8 | uint16(cpu.Memory.Read(cpu.PC))
			return 16
		},
		// 0xC4: CALL NZ, a16
		func(cpu *CPU) int {
			if cpu.FZero {
				cpu.PC += 2
				return 12
			}
			tempPC := cpu.PC + 2
			cpu.Memory.Write(cpu.SP-1, uint8(tempPC>>8))
			cpu.Memory.Write(cpu.SP-2, uint8(tempPC))
			cpu.SP -= 2
			cpu.PC = uint16(cpu.Memory.Read(cpu.PC+1))<<8 | uint16(cpu.Memory.Read(cpu.PC))
			return 24
		},
		// 0xC5: PUSH BC
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.SP-1, cpu.B)
			cpu.Memory.Write(cpu.SP-2, cpu.C)
			cpu.SP -= 2
			return 16
		},
		// 0xC6: ADD A, d8
		func(cpu *CPU) int {
			temp := uint16(cpu.A) + uint16(cpu.Memory.Read(cpu.PC))
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			cpu.PC++
			return 8
		},
		// 0xC7: RST 00H
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.SP-1, uint8(cpu.PC>>8))
			cpu.Memory.Write(cpu.SP-2, uint8(cpu.PC))
			cpu.SP -= 2
			cpu.PC = 0
			return 16
		},
		// 0xC8: RET Z
		func(cpu *CPU) int {
			if !cpu.FZero {
				return 8
			}
			cpu.PC = uint16(cpu.Memory.Read(cpu.SP+1))<<8 | uint16(cpu.Memory.Read(cpu.SP))
			cpu.SP += 2
			return 20
		},
		// 0xC9: RET
		func(cpu *CPU) int {
			cpu.PC = uint16(cpu.Memory.Read(cpu.SP+1))<<8 | uint16(cpu.Memory.Read(cpu.SP))
			cpu.SP += 2
			return 16
		},
		// 0xCA: JP Z, a16
		func(cpu *CPU) int {
			if !cpu.FZero {
				cpu.PC += 2
				return 12
			}
			cpu.PC = uint16(cpu.Memory.Read(cpu.PC+1))<<8 | uint16(cpu.Memory.Read(cpu.PC))
			return 16
		},
		// 0xCB: CB Prefix
		func(cpu *CPU) int {
			cycles := CBInstructionsTable[cpu.Memory.Read(cpu.PC)](cpu)
			cpu.PC++
			return cycles
		},
		// 0xCC: CALL Z, a16
		func(cpu *CPU) int {
			if !cpu.FZero {
				cpu.PC += 2
				return 12
			}
			tempPC := cpu.PC + 2
			cpu.Memory.Write(cpu.SP-1, uint8(tempPC>>8))
			cpu.Memory.Write(cpu.SP-2, uint8(tempPC))
			cpu.SP -= 2
			cpu.PC = uint16(cpu.Memory.Read(cpu.PC+1))<<8 | uint16(cpu.Memory.Read(cpu.PC))
			return 24
		},
		// 0xCD: CALL a16
		func(cpu *CPU) int {
			tempPC := cpu.PC + 2
			cpu.Memory.Write(cpu.SP-1, uint8(tempPC>>8))
			cpu.Memory.Write(cpu.SP-2, uint8(tempPC))
			cpu.SP -= 2
			cpu.PC = uint16(cpu.Memory.Read(cpu.PC+1))<<8 | uint16(cpu.Memory.Read(cpu.PC))
			return 24
		},
		// 0xCE: ADC A, d8
		func(cpu *CPU) int {
			var carry uint16
			if cpu.FCarry {
				carry = 1
			}
			temp := uint16(cpu.A) + uint16(cpu.Memory.Read(cpu.PC)) + carry
			uint8Val := uint8(temp)
			cpu.FZero = uint8Val == 0
			cpu.FSub = false
			cpu.FHalfCarry = cpu.A&0xf > uint8(temp)&0xf
			cpu.FCarry = temp > 0xff
			cpu.A = uint8Val
			cpu.PC++
			return 8
		},
		// 0xCF: RST 08H
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.SP-1, uint8(cpu.PC>>8))
			cpu.Memory.Write(cpu.SP-2, uint8(cpu.PC))
			cpu.SP -= 2
			cpu.PC = 0x8
			return 16
		},
		// 0xD0: RET NC
		func(cpu *CPU) int {
			if cpu.FCarry {
				return 8
			}
			cpu.PC = uint16(cpu.Memory.Read(cpu.SP+1))<<8 | uint16(cpu.Memory.Read(cpu.SP))
			cpu.SP += 2
			return 20
		},
		// 0xD1: POP DE
		func(cpu *CPU) int {
			cpu.E = cpu.Memory.Read(cpu.SP)
			cpu.D = cpu.Memory.Read(cpu.SP + 1)
			cpu.SP += 2
			return 12
		},
		// 0xD2: JP NC, a16
		func(cpu *CPU) int {
			if cpu.FCarry {
				cpu.PC += 2
				return 12
			}
			cpu.PC = uint16(cpu.Memory.Read(cpu.PC+1))<<8 | uint16(cpu.Memory.Read(cpu.PC))
			return 16
		},
		// 0xD3: NOP
		func(cpu *CPU) int {
			return 4
		},
		// 0xD4: CALL NC, a16
		func(cpu *CPU) int {
			if cpu.FCarry {
				cpu.PC += 2
				return 12
			}
			tempPC := cpu.PC + 2
			cpu.Memory.Write(cpu.SP-1, uint8(tempPC>>8))
			cpu.Memory.Write(cpu.SP-2, uint8(tempPC))
			cpu.SP -= 2
			cpu.PC = uint16(cpu.Memory.Read(cpu.PC+1))<<8 | uint16(cpu.Memory.Read(cpu.PC))
			return 24
		},
		// 0xD5: PUSH DE
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.SP-1, cpu.D)
			cpu.Memory.Write(cpu.SP-2, cpu.E)
			cpu.SP -= 2
			return 16
		},
		// 0xD6: SUB A, d8
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.Memory.Read(cpu.PC))
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = cpu.A&0xf < uint8(temp)&0xf
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			cpu.PC++
			return 8
		},
		// 0xD7: RST 10H
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.SP-1, uint8(cpu.PC>>8))
			cpu.Memory.Write(cpu.SP-2, uint8(cpu.PC))
			cpu.SP -= 2
			cpu.PC = 0x10
			return 16
		},
		// 0xD8: RET C
		func(cpu *CPU) int {
			if !cpu.FCarry {
				return 8
			}
			cpu.PC = uint16(cpu.Memory.Read(cpu.SP+1))<<8 | uint16(cpu.Memory.Read(cpu.SP))
			cpu.SP += 2
			return 20
		},
		// 0xD9: RETI
		func(cpu *CPU) int {
			cpu.EI = true
			cpu.PC = uint16(cpu.Memory.Read(cpu.SP+1))<<8 | uint16(cpu.Memory.Read(cpu.SP))
			cpu.SP += 2
			return 20
		},
		// 0xDA: JP C, a16
		func(cpu *CPU) int {
			if !cpu.FCarry {
				cpu.PC += 2
				return 12
			}
			cpu.PC = uint16(cpu.Memory.Read(cpu.PC+1))<<8 | uint16(cpu.Memory.Read(cpu.PC))
			return 16
		},
		// 0xDB: NOP
		func(cpu *CPU) int {
			return 4
		},
		// 0xDC: CALL Z, a16
		func(cpu *CPU) int {
			if !cpu.FCarry {
				cpu.PC += 2
				return 12
			}
			tempPC := cpu.PC + 2
			cpu.Memory.Write(cpu.SP-1, uint8(tempPC>>8))
			cpu.Memory.Write(cpu.SP-2, uint8(tempPC))
			cpu.SP -= 2
			cpu.PC = uint16(cpu.Memory.Read(cpu.PC+1))<<8 | uint16(cpu.Memory.Read(cpu.PC))
			return 24
		},
		// 0xDD: NOP
		func(cpu *CPU) int {
			return 4
		},
		// 0xDE: SBC A, d8
		func(cpu *CPU) int {
			var carry int
			if cpu.FCarry {
				carry = 1
			}
			tempVal := cpu.Memory.Read(cpu.PC)
			temp := int(cpu.A) - int(tempVal) - carry
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = int(cpu.A&0xf)-int(tempVal&0xf)-carry < 0
			cpu.FCarry = temp < 0
			cpu.A = uint8(temp)
			cpu.PC++
			return 8
		},
		// 0xDF: RST 18H
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.SP-1, uint8(cpu.PC>>8))
			cpu.Memory.Write(cpu.SP-2, uint8(cpu.PC))
			cpu.SP -= 2
			cpu.PC = 0x18
			return 16
		},
		// 0xE0: LDH (a8), A
		func(cpu *CPU) int {
			cpu.Memory.Write(0xff00+uint16(cpu.Memory.Read(cpu.PC)), cpu.A)
			cpu.PC++
			return 12
		},
		// 0xE1: POP HL
		func(cpu *CPU) int {
			cpu.L = cpu.Memory.Read(cpu.SP)
			cpu.H = cpu.Memory.Read(cpu.SP + 1)
			cpu.SP += 2
			return 12
		},
		// 0xE2: LD (C), A
		func(cpu *CPU) int {
			cpu.Memory.Write(0xff00+uint16(cpu.C), cpu.A)
			return 8
		},
		// 0xE3: NOP
		func(cpu *CPU) int {
			return 4
		},
		// 0xE4: NOP
		func(cpu *CPU) int {
			return 4
		},
		// 0xE5: PUSH HL
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.SP-1, cpu.H)
			cpu.Memory.Write(cpu.SP-2, cpu.L)
			cpu.SP -= 2
			return 16
		},
		// 0xE6: AND d8
		func(cpu *CPU) int {
			cpu.A = cpu.A & cpu.Memory.Read(cpu.PC)
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = true
			cpu.FCarry = false
			cpu.PC++
			return 8
		},
		// 0xE7: RST 20H
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.SP-1, uint8(cpu.PC>>8))
			cpu.Memory.Write(cpu.SP-2, uint8(cpu.PC))
			cpu.SP -= 2
			cpu.PC = 0x20
			return 16
		},
		// 0xE8: ADD SP, r8
		func(cpu *CPU) int {
			val := int(int8(cpu.Memory.Read(cpu.PC)))
			temp := val + int(cpu.SP)
			cpu.FZero = false
			cpu.FSub = false
			cpu.FCarry = (int(cpu.SP)^val^(temp&0xffff))&0x100 == 0x100
			cpu.FHalfCarry = (int(cpu.SP)^val^(temp&0xffff))&0x10 == 0x10
			cpu.PC++
			return 16
		},
		// 0xE9: JP (HL)
		func(cpu *CPU) int {
			cpu.PC = cpu.HL()
			return 4
		},
		// 0xEA: LD (a16), A
		func(cpu *CPU) int {
			cpu.Memory.Write(
				uint16(cpu.Memory.Read(cpu.PC+1))<<8|uint16(cpu.Memory.Read(cpu.PC)),
				cpu.A,
			)
			cpu.PC += 2
			return 16
		},
		// 0xEB: NOP
		func(cpu *CPU) int {
			return 4
		},
		// 0xEC: NOP
		func(cpu *CPU) int {
			return 4
		},
		// 0xED: NOP
		func(cpu *CPU) int {
			return 4
		},
		// 0xA8: XOR d8
		func(cpu *CPU) int {
			cpu.A = cpu.A ^ cpu.Memory.Read(cpu.PC)
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			cpu.PC++
			return 8
		},
		// 0xDF: RST 28H
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.SP-1, uint8(cpu.PC>>8))
			cpu.Memory.Write(cpu.SP-2, uint8(cpu.PC))
			cpu.SP -= 2
			cpu.PC = 0x28
			return 16
		},
		// 0xF0: LDH A, (a8)
		func(cpu *CPU) int {
			cpu.A = cpu.Memory.Read(0xff00 + uint16(cpu.Memory.Read(cpu.PC)))
			cpu.PC++
			return 12
		},
		// 0xF1: POP AF
		func(cpu *CPU) int {
			cpu.SetF(cpu.Memory.Read(cpu.SP))
			cpu.A = cpu.Memory.Read(cpu.SP + 1)
			cpu.SP += 2
			return 12
		},
		// 0xF2: LD A, (C)
		func(cpu *CPU) int {
			cpu.A = cpu.Memory.Read(0xff00 + uint16(cpu.C))
			return 8
		},
		// 0xF3: DI
		func(cpu *CPU) int {
			cpu.EI = false
			return 4
		},
		// 0xE4: NOP
		func(cpu *CPU) int {
			return 4
		},
		// 0xE5: PUSH AF
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.SP-1, cpu.A)
			cpu.Memory.Write(cpu.SP-2, cpu.F())
			cpu.SP -= 2
			return 16
		},
		// 0xF6: OR d8
		func(cpu *CPU) int {
			cpu.A = cpu.A | cpu.Memory.Read(cpu.PC)
			cpu.FZero = cpu.A == 0
			cpu.FSub = false
			cpu.FHalfCarry = false
			cpu.FCarry = false
			cpu.PC++
			return 4
		},
		// 0xF7: RST 20H
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.SP-1, uint8(cpu.PC>>8))
			cpu.Memory.Write(cpu.SP-2, uint8(cpu.PC))
			cpu.SP -= 2
			cpu.PC = 0x20
			return 16
		},
		// 0xF8: LD HL, SP+r8
		func(cpu *CPU) int {
			r8 := int(int8(cpu.Memory.Read(cpu.PC)))
			result := int(cpu.SP) + r8
			cpu.H = uint8(result >> 8)
			cpu.L = uint8(result)
			cpu.FZero = false
			cpu.FSub = false
			cpu.FHalfCarry = ((int(cpu.SP) ^ r8 ^ result) & 0x10) == 0x10
			cpu.FCarry = ((int(cpu.SP) ^ r8 ^ result) & 0x100) == 0x100
			cpu.PC++
			return 12
		},
		// 0xF9: LD SP, HL
		func(cpu *CPU) int {
			cpu.SP = cpu.AF()
			return 8
		},
		// 0xFA: LD A, (a16)
		func(cpu *CPU) int {
			cpu.A = cpu.Memory.Read(uint16(cpu.Memory.Read(cpu.PC+1))<<8 | uint16(cpu.Memory.Read(cpu.PC)))
			cpu.PC += 2
			return 16
		},
		// 0xFB: EI
		func(cpu *CPU) int {
			cpu.EI = true
			return 4
		},
		// 0xFC: NOP
		func(cpu *CPU) int {
			return 4
		},
		// 0xFD: NOP
		func(cpu *CPU) int {
			return 4
		},
		// 0xFE: CP d8
		func(cpu *CPU) int {
			temp := int(cpu.A) - int(cpu.Memory.Read(cpu.PC))
			cpu.FZero = temp == 0
			cpu.FSub = true
			cpu.FHalfCarry = temp&0xf > int(cpu.A&0xf)
			cpu.FCarry = temp < 0
			cpu.PC++
			return 8
		},
		// 0xFF: RST 38H
		func(cpu *CPU) int {
			cpu.Memory.Write(cpu.SP-1, uint8(cpu.PC>>8))
			cpu.Memory.Write(cpu.SP-2, uint8(cpu.PC))
			cpu.SP -= 2
			cpu.PC = 0x38
			return 4
		},
	}
}
