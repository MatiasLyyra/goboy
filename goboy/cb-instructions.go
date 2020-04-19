package goboy

// CBInstructionsTable contains all of the cb prefixed instructions
var CBInstructionsTable [256]MicroCodeFunc

func rlc(val uint8, cpu *CPU) uint8 {
	carry := val >> 7
	val = (val << 1) | carry
	cpu.FZero = val == 0
	cpu.FSub = false
	cpu.FHalfCarry = false
	cpu.FCarry = carry == 1
	return val
}

func rrc(val uint8, cpu *CPU) uint8 {
	carry := val & 1
	val = (val >> 1) | (carry << 7)
	cpu.FZero = val == 0
	cpu.FSub = false
	cpu.FHalfCarry = false
	cpu.FCarry = carry == 1
	return val
}

func rl(val uint8, cpu *CPU) uint8 {
	var (
		prevCarry uint8
		carry     uint8
	)
	if cpu.FCarry {
		prevCarry = 1
	}
	carry = val >> 7 & 1
	val = (val << 1) | prevCarry
	cpu.FZero = val == 0
	cpu.FSub = false
	cpu.FHalfCarry = false
	cpu.FCarry = carry == 1
	return val
}

func rr(val uint8, cpu *CPU) uint8 {
	var (
		prevCarry uint8
		carry     uint8
	)
	if cpu.FCarry {
		prevCarry = 1
	}
	carry = val & 1
	val = (val >> 1) | (prevCarry << 7)
	cpu.FZero = val == 0
	cpu.FSub = false
	cpu.FHalfCarry = false
	cpu.FCarry = carry == 1
	return val
}

func sla(val uint8, cpu *CPU) uint8 {
	carry := val >> 7
	val = val << 1
	cpu.FZero = val == 0
	cpu.FSub = false
	cpu.FHalfCarry = false
	cpu.FCarry = carry == 1
	return val
}

func sra(val uint8, cpu *CPU) uint8 {
	carry := val & 1
	if val&0x80 != 0 {
		val = (val >> 1) | 0x80
	} else {
		val = val >> 1
	}
	cpu.FZero = val == 0
	cpu.FSub = false
	cpu.FHalfCarry = false
	cpu.FCarry = carry == 1
	return val
}

func swap(val uint8, cpu *CPU) uint8 {
	l := val & 0xf
	h := val >> 4
	val = (l << 4) | h
	cpu.FZero = val == 0
	cpu.FSub = false
	cpu.FHalfCarry = false
	cpu.FCarry = false
	return val
}

func srl(val uint8, cpu *CPU) uint8 {
	carry := val & 1
	val = val >> 1
	cpu.FZero = val == 0
	cpu.FSub = false
	cpu.FHalfCarry = false
	cpu.FCarry = carry == 1
	return val
}

func testBit(val, bit uint8, cpu *CPU) {
	cpu.FZero = val&(1<<bit) == 0
	cpu.FSub = false
	cpu.FHalfCarry = true
}

func resetBit(val, bit uint8) uint8 {
	mask := uint8(^(1 << bit))
	return mask & val
}

func setBit(val, bit uint8) uint8 {
	return val | (1 << bit)
}

func init() {
	CBInstructionsTable = [256]MicroCodeFunc{
		// 0x00: RLC B
		func(cpu *CPU) int {
			cpu.B = rlc(cpu.B, cpu)
			return 8
		},
		// 0x01: RLC C
		func(cpu *CPU) int {
			cpu.C = rlc(cpu.C, cpu)
			return 8
		},
		// 0x02: RLC D
		func(cpu *CPU) int {
			cpu.D = rlc(cpu.D, cpu)
			return 8
		},
		// 0x03: RLC E
		func(cpu *CPU) int {
			cpu.E = rlc(cpu.E, cpu)
			return 8
		},
		// 0x04: RLC H
		func(cpu *CPU) int {
			cpu.H = rlc(cpu.H, cpu)
			return 8
		},
		// 0x05: RLC L
		func(cpu *CPU) int {
			cpu.L = rlc(cpu.L, cpu)
			return 8
		},
		// 0x06: RLC (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = rlc(val, cpu)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0x07: RLC A
		func(cpu *CPU) int {
			cpu.A = rlc(cpu.A, cpu)
			return 8
		},
		// 0x08: RRC B
		func(cpu *CPU) int {
			cpu.B = rrc(cpu.B, cpu)
			return 8
		},
		// 0x09: RRC C
		func(cpu *CPU) int {
			cpu.C = rrc(cpu.C, cpu)
			return 8
		},
		// 0x0A: RRC D
		func(cpu *CPU) int {
			cpu.D = rrc(cpu.D, cpu)
			return 8
		},
		// 0x0B: RRC E
		func(cpu *CPU) int {
			cpu.E = rrc(cpu.E, cpu)
			return 8
		},
		// 0x0C: RRC H
		func(cpu *CPU) int {
			cpu.H = rrc(cpu.H, cpu)
			return 8
		},
		// 0x0D: RRC L
		func(cpu *CPU) int {
			cpu.L = rrc(cpu.L, cpu)
			return 8
		},
		// 0x0E: RRC (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = rrc(val, cpu)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0x0F: RRC A
		func(cpu *CPU) int {
			cpu.A = rrc(cpu.A, cpu)
			return 8
		},
		// 0x10: RL B
		func(cpu *CPU) int {
			cpu.B = rl(cpu.B, cpu)
			return 8
		},
		// 0x11: RL C
		func(cpu *CPU) int {
			cpu.C = rl(cpu.C, cpu)
			return 8
		},
		// 0x12: RL D
		func(cpu *CPU) int {
			cpu.D = rl(cpu.D, cpu)
			return 8
		},
		// 0x13: RL E
		func(cpu *CPU) int {
			cpu.E = rl(cpu.E, cpu)
			return 8
		},
		// 0x14: RL H
		func(cpu *CPU) int {
			cpu.H = rl(cpu.H, cpu)
			return 8
		},
		// 0x15: RL L
		func(cpu *CPU) int {
			cpu.L = rl(cpu.L, cpu)
			return 8
		},
		// 0x16: RL (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = rl(val, cpu)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0x17: RL A
		func(cpu *CPU) int {
			cpu.A = rl(cpu.A, cpu)
			return 8
		},
		// 0x18: RR B
		func(cpu *CPU) int {
			cpu.B = rr(cpu.B, cpu)
			return 8
		},
		// 0x19: RR C
		func(cpu *CPU) int {
			cpu.C = rr(cpu.C, cpu)
			return 8
		},
		// 0x1A: RR D
		func(cpu *CPU) int {
			cpu.D = rr(cpu.D, cpu)
			return 8
		},
		// 0x1B: RR E
		func(cpu *CPU) int {
			cpu.E = rr(cpu.E, cpu)
			return 8
		},
		// 0x1C: RR H
		func(cpu *CPU) int {
			cpu.H = rr(cpu.H, cpu)
			return 8
		},
		// 0x1D: RR L
		func(cpu *CPU) int {
			cpu.L = rr(cpu.L, cpu)
			return 8
		},
		// 0x1E: RR (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = rr(val, cpu)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0x1F: RR A
		func(cpu *CPU) int {
			cpu.A = rr(cpu.A, cpu)
			return 8
		},
		// 0x20: SLA B
		func(cpu *CPU) int {
			cpu.B = sla(cpu.B, cpu)
			return 8
		},
		// 0x21: SLA C
		func(cpu *CPU) int {
			cpu.C = sla(cpu.C, cpu)
			return 8
		},
		// 0x22: SLA D
		func(cpu *CPU) int {
			cpu.D = sla(cpu.D, cpu)
			return 8
		},
		// 0x23: SLA E
		func(cpu *CPU) int {
			cpu.E = sla(cpu.E, cpu)
			return 8
		},
		// 0x24: SLA H
		func(cpu *CPU) int {
			cpu.H = sla(cpu.H, cpu)
			return 8
		},
		// 0x25: SLA L
		func(cpu *CPU) int {
			cpu.L = sla(cpu.L, cpu)
			return 8
		},
		// 0x26: SLA (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = sla(val, cpu)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0x27: SLA A
		func(cpu *CPU) int {
			cpu.A = sla(cpu.A, cpu)
			return 8
		},
		// 0x28: SRA B
		func(cpu *CPU) int {
			cpu.B = sra(cpu.B, cpu)
			return 8
		},
		// 0x29: SRA C
		func(cpu *CPU) int {
			cpu.C = sra(cpu.C, cpu)
			return 8
		},
		// 0x2A: SRA D
		func(cpu *CPU) int {
			cpu.D = sra(cpu.D, cpu)
			return 8
		},
		// 0x2B: SRA E
		func(cpu *CPU) int {
			cpu.E = sra(cpu.E, cpu)
			return 8
		},
		// 0x2C: SRA H
		func(cpu *CPU) int {
			cpu.H = sra(cpu.H, cpu)
			return 8
		},
		// 0x2D: SRA L
		func(cpu *CPU) int {
			cpu.L = sra(cpu.L, cpu)
			return 8
		},
		// 0x2E: SRA (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = sra(val, cpu)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0x2F: SRA A
		func(cpu *CPU) int {
			cpu.A = sra(cpu.A, cpu)
			return 8
		},
		// 0x30: SWAP B
		func(cpu *CPU) int {
			cpu.B = swap(cpu.B, cpu)
			return 8
		},
		// 0x31: SWAP C
		func(cpu *CPU) int {
			cpu.C = swap(cpu.C, cpu)
			return 8
		},
		// 0x32: SWAP D
		func(cpu *CPU) int {
			cpu.D = swap(cpu.D, cpu)
			return 8
		},
		// 0x33: SWAP E
		func(cpu *CPU) int {
			cpu.E = swap(cpu.E, cpu)
			return 8
		},
		// 0x34: SWAP H
		func(cpu *CPU) int {
			cpu.H = swap(cpu.H, cpu)
			return 8
		},
		// 0x35: SWAP L
		func(cpu *CPU) int {
			cpu.L = swap(cpu.L, cpu)
			return 8
		},
		// 0x36: SWAP (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = swap(val, cpu)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0x37: SWAP A
		func(cpu *CPU) int {
			cpu.A = swap(cpu.A, cpu)
			return 8
		},
		// 0x38: SRL B
		func(cpu *CPU) int {
			cpu.B = srl(cpu.B, cpu)
			return 8
		},
		// 0x39: SRL C
		func(cpu *CPU) int {
			cpu.C = srl(cpu.C, cpu)
			return 8
		},
		// 0x3A: SRL D
		func(cpu *CPU) int {
			cpu.D = srl(cpu.D, cpu)
			return 8
		},
		// 0x3B: SRL E
		func(cpu *CPU) int {
			cpu.E = srl(cpu.E, cpu)
			return 8
		},
		// 0x3C: SRL H
		func(cpu *CPU) int {
			cpu.H = srl(cpu.H, cpu)
			return 8
		},
		// 0x3D: SRL L
		func(cpu *CPU) int {
			cpu.L = srl(cpu.L, cpu)
			return 8
		},
		// 0x3E: SRL (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = srl(val, cpu)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0x3F: SRL A
		func(cpu *CPU) int {
			cpu.A = srl(cpu.A, cpu)
			return 8
		},
		// 0x40: BIT 0, B
		func(cpu *CPU) int {
			testBit(cpu.B, 0, cpu)
			return 8
		},
		// 0x41: BIT 0, C
		func(cpu *CPU) int {
			testBit(cpu.C, 0, cpu)
			return 8
		},
		// 0x42: BIT 0, D
		func(cpu *CPU) int {
			testBit(cpu.D, 0, cpu)
			return 8
		},
		// 0x43: BIT 0, E
		func(cpu *CPU) int {
			testBit(cpu.E, 0, cpu)
			return 8
		},
		// 0x44: BIT 0, H
		func(cpu *CPU) int {
			testBit(cpu.H, 0, cpu)
			return 8
		},
		// 0x45: BIT 0, L
		func(cpu *CPU) int {
			testBit(cpu.L, 0, cpu)
			return 8
		},
		// 0x46: BIT 0, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			testBit(val, 0, cpu)
			return 12
		},
		// 0x47: BIT 0, A
		func(cpu *CPU) int {
			testBit(cpu.A, 0, cpu)
			return 8
		},
		// 0x48: BIT 1, B
		func(cpu *CPU) int {
			testBit(cpu.B, 1, cpu)
			return 8
		},
		// 0x49: BIT 1, C
		func(cpu *CPU) int {
			testBit(cpu.C, 1, cpu)
			return 8
		},
		// 0x4A: BIT 1, D
		func(cpu *CPU) int {
			testBit(cpu.D, 1, cpu)
			return 8
		},
		// 0x4B: BIT 1, E
		func(cpu *CPU) int {
			testBit(cpu.E, 1, cpu)
			return 8
		},
		// 0x4C: BIT 1, H
		func(cpu *CPU) int {
			testBit(cpu.H, 1, cpu)
			return 8
		},
		// 0x4D: BIT 1, L
		func(cpu *CPU) int {
			testBit(cpu.L, 1, cpu)
			return 8
		},
		// 0x4E: BIT 1, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			testBit(val, 1, cpu)
			return 12
		},
		// 0x4F: BIT 1, A
		func(cpu *CPU) int {
			testBit(cpu.A, 1, cpu)
			return 8
		},

		// 0x50: BIT 2, B
		func(cpu *CPU) int {
			testBit(cpu.B, 2, cpu)
			return 8
		},
		// 0x51: BIT 2, C
		func(cpu *CPU) int {
			testBit(cpu.C, 2, cpu)
			return 8
		},
		// 0x52: BIT 2, D
		func(cpu *CPU) int {
			testBit(cpu.D, 2, cpu)
			return 8
		},
		// 0x53: BIT 2, E
		func(cpu *CPU) int {
			testBit(cpu.E, 2, cpu)
			return 8
		},
		// 0x54: BIT 2, H
		func(cpu *CPU) int {
			testBit(cpu.H, 2, cpu)
			return 8
		},
		// 0x55: BIT 2, L
		func(cpu *CPU) int {
			testBit(cpu.L, 2, cpu)
			return 8
		},
		// 0x56: BIT 2, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			testBit(val, 2, cpu)
			return 12
		},
		// 0x57: BIT 2, A
		func(cpu *CPU) int {
			testBit(cpu.A, 2, cpu)
			return 8
		},
		// 0x58: BIT 3, B
		func(cpu *CPU) int {
			testBit(cpu.B, 3, cpu)
			return 8
		},
		// 0x59: BIT 3, C
		func(cpu *CPU) int {
			testBit(cpu.C, 3, cpu)
			return 8
		},
		// 0x5A: BIT 3, D
		func(cpu *CPU) int {
			testBit(cpu.D, 3, cpu)
			return 8
		},
		// 0x5B: BIT 3, E
		func(cpu *CPU) int {
			testBit(cpu.E, 3, cpu)
			return 8
		},
		// 0x5C: BIT 3, H
		func(cpu *CPU) int {
			testBit(cpu.H, 3, cpu)
			return 8
		},
		// 0x5D: BIT 3, L
		func(cpu *CPU) int {
			testBit(cpu.L, 3, cpu)
			return 8
		},
		// 0x5E: BIT 3, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			testBit(val, 3, cpu)
			return 12
		},
		// 0x5F: BIT 3, A
		func(cpu *CPU) int {
			testBit(cpu.A, 3, cpu)
			return 8
		},

		// 0x60: BIT 4, B
		func(cpu *CPU) int {
			testBit(cpu.B, 4, cpu)
			return 8
		},
		// 0x61: BIT 4, C
		func(cpu *CPU) int {
			testBit(cpu.C, 4, cpu)
			return 8
		},
		// 0x62: BIT 4, D
		func(cpu *CPU) int {
			testBit(cpu.D, 4, cpu)
			return 8
		},
		// 0x63: BIT 4, E
		func(cpu *CPU) int {
			testBit(cpu.E, 4, cpu)
			return 8
		},
		// 0x64: BIT 4, H
		func(cpu *CPU) int {
			testBit(cpu.H, 4, cpu)
			return 8
		},
		// 0x65: BIT 4, L
		func(cpu *CPU) int {
			testBit(cpu.L, 4, cpu)
			return 8
		},
		// 0x66: BIT 4, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			testBit(val, 4, cpu)
			return 12
		},
		// 0x67: BIT 4, A
		func(cpu *CPU) int {
			testBit(cpu.A, 4, cpu)
			return 8
		},
		// 0x68: BIT 5, B
		func(cpu *CPU) int {
			testBit(cpu.B, 5, cpu)
			return 8
		},
		// 0x69: BIT 5, C
		func(cpu *CPU) int {
			testBit(cpu.C, 5, cpu)
			return 8
		},
		// 0x6A: BIT 5, D
		func(cpu *CPU) int {
			testBit(cpu.D, 5, cpu)
			return 8
		},
		// 0x6B: BIT 5, E
		func(cpu *CPU) int {
			testBit(cpu.E, 5, cpu)
			return 8
		},
		// 0x6C: BIT 5, H
		func(cpu *CPU) int {
			testBit(cpu.H, 5, cpu)
			return 8
		},
		// 0x6D: BIT 5, L
		func(cpu *CPU) int {
			testBit(cpu.L, 5, cpu)
			return 8
		},
		// 0x6E: BIT 5, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			testBit(val, 5, cpu)
			return 12
		},
		// 0x6F: BIT 5, A
		func(cpu *CPU) int {
			testBit(cpu.A, 5, cpu)
			return 8
		},

		// 0x70: BIT 6, B
		func(cpu *CPU) int {
			testBit(cpu.B, 6, cpu)
			return 8
		},
		// 0x71: BIT 6, C
		func(cpu *CPU) int {
			testBit(cpu.C, 6, cpu)
			return 8
		},
		// 0x72: BIT 6, D
		func(cpu *CPU) int {
			testBit(cpu.D, 6, cpu)
			return 8
		},
		// 0x73: BIT 6, E
		func(cpu *CPU) int {
			testBit(cpu.E, 6, cpu)
			return 8
		},
		// 0x74: BIT 6, H
		func(cpu *CPU) int {
			testBit(cpu.H, 6, cpu)
			return 8
		},
		// 0x75: BIT 6, L
		func(cpu *CPU) int {
			testBit(cpu.L, 6, cpu)
			return 8
		},
		// 0x76: BIT 6, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			testBit(val, 6, cpu)
			return 12
		},
		// 0x77: BIT 6, A
		func(cpu *CPU) int {
			testBit(cpu.A, 6, cpu)
			return 8
		},
		// 0x78: BIT 7, B
		func(cpu *CPU) int {
			testBit(cpu.B, 7, cpu)
			return 8
		},
		// 0x79: BIT 7, C
		func(cpu *CPU) int {
			testBit(cpu.C, 7, cpu)
			return 8
		},
		// 0x7A: BIT 7, D
		func(cpu *CPU) int {
			testBit(cpu.D, 7, cpu)
			return 8
		},
		// 0x7B: BIT 7, E
		func(cpu *CPU) int {
			testBit(cpu.E, 7, cpu)
			return 8
		},
		// 0x7C: BIT 7, H
		func(cpu *CPU) int {
			testBit(cpu.H, 7, cpu)
			return 8
		},
		// 0x7D: BIT 7, L
		func(cpu *CPU) int {
			testBit(cpu.L, 7, cpu)
			return 8
		},
		// 0x7E: BIT 7, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			testBit(val, 7, cpu)
			return 12
		},
		// 0x7F: BIT 7, A
		func(cpu *CPU) int {
			testBit(cpu.A, 7, cpu)
			return 8
		},
		// 0x80: RES 0, B
		func(cpu *CPU) int {
			cpu.B = resetBit(cpu.B, 0)
			return 8
		},
		// 0x81: RES 0, C
		func(cpu *CPU) int {
			cpu.C = resetBit(cpu.C, 0)
			return 8
		},
		// 0x82: RES 0, D
		func(cpu *CPU) int {
			cpu.D = resetBit(cpu.D, 0)
			return 8
		},
		// 0x83: RES 0, E
		func(cpu *CPU) int {
			cpu.E = resetBit(cpu.E, 0)
			return 8
		},
		// 0x84: RES 0, H
		func(cpu *CPU) int {
			cpu.H = resetBit(cpu.H, 0)
			return 8
		},
		// 0x85: RES 0, L
		func(cpu *CPU) int {
			cpu.L = resetBit(cpu.L, 0)
			return 8
		},
		// 0x86: RES 0, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = resetBit(val, 0)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0x87: RES 0, A
		func(cpu *CPU) int {
			cpu.A = resetBit(cpu.A, 0)
			return 8
		},
		// 0x88: RES 1, B
		func(cpu *CPU) int {
			cpu.B = resetBit(cpu.B, 1)
			return 8
		},
		// 0x89: RES 1, C
		func(cpu *CPU) int {
			cpu.C = resetBit(cpu.C, 1)
			return 8
		},
		// 0x8A: RES 1, D
		func(cpu *CPU) int {
			cpu.D = resetBit(cpu.D, 1)
			return 8
		},
		// 0x8B: RES 1, E
		func(cpu *CPU) int {
			cpu.E = resetBit(cpu.E, 1)
			return 8
		},
		// 0x8C: RES 1, H
		func(cpu *CPU) int {
			cpu.H = resetBit(cpu.H, 1)
			return 8
		},
		// 0x8D: RES 1, L
		func(cpu *CPU) int {
			cpu.L = resetBit(cpu.L, 1)
			return 8
		},
		// 0x8E: RES 1, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = resetBit(val, 1)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0x8F: RES 1, A
		func(cpu *CPU) int {
			cpu.A = resetBit(cpu.A, 1)
			return 8
		},

		// 0x90: RES 2, B
		func(cpu *CPU) int {
			cpu.B = resetBit(cpu.B, 2)
			return 8
		},
		// 0x91: RES 2, C
		func(cpu *CPU) int {
			cpu.C = resetBit(cpu.C, 2)
			return 8
		},
		// 0x92: RES 2, D
		func(cpu *CPU) int {
			cpu.D = resetBit(cpu.D, 2)
			return 8
		},
		// 0x93: RES 2, E
		func(cpu *CPU) int {
			cpu.E = resetBit(cpu.E, 2)
			return 8
		},
		// 0x94: RES 2, H
		func(cpu *CPU) int {
			cpu.H = resetBit(cpu.H, 2)
			return 8
		},
		// 0x95: RES 2, L
		func(cpu *CPU) int {
			cpu.L = resetBit(cpu.L, 2)
			return 8
		},
		// 0x96: RES 2, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = resetBit(val, 2)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0x97: RES 2, A
		func(cpu *CPU) int {
			cpu.A = resetBit(cpu.A, 2)
			return 8
		},
		// 0x98: RES 3, B
		func(cpu *CPU) int {
			cpu.B = resetBit(cpu.B, 3)
			return 8
		},
		// 0x99: RES 3, C
		func(cpu *CPU) int {
			cpu.C = resetBit(cpu.C, 3)
			return 8
		},
		// 0x9A: RES 3, D
		func(cpu *CPU) int {
			cpu.D = resetBit(cpu.D, 3)
			return 8
		},
		// 0x9B: RES 3, E
		func(cpu *CPU) int {
			cpu.E = resetBit(cpu.E, 3)
			return 8
		},
		// 0x9C: RES 3, H
		func(cpu *CPU) int {
			cpu.H = resetBit(cpu.H, 3)
			return 8
		},
		// 0x9D: RES 3, L
		func(cpu *CPU) int {
			cpu.L = resetBit(cpu.L, 3)
			return 8
		},
		// 0x9E: RES 3, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = resetBit(val, 3)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0x9F: RES 3, A
		func(cpu *CPU) int {
			cpu.A = resetBit(cpu.A, 3)
			return 8
		},

		// 0xA0: RES 4, B
		func(cpu *CPU) int {
			cpu.B = resetBit(cpu.B, 4)
			return 8
		},
		// 0xA1: RES 4, C
		func(cpu *CPU) int {
			cpu.C = resetBit(cpu.C, 4)
			return 8
		},
		// 0xA2: RES 4, D
		func(cpu *CPU) int {
			cpu.D = resetBit(cpu.D, 4)
			return 8
		},
		// 0xA3: RES 4, E
		func(cpu *CPU) int {
			cpu.E = resetBit(cpu.E, 4)
			return 8
		},
		// 0xA4: RES 4, H
		func(cpu *CPU) int {
			cpu.H = resetBit(cpu.H, 4)
			return 8
		},
		// 0xA5: RES 4, L
		func(cpu *CPU) int {
			cpu.L = resetBit(cpu.L, 4)
			return 8
		},
		// 0xA6: RES 4, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = resetBit(val, 4)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0xA7: RES 4, A
		func(cpu *CPU) int {
			cpu.A = resetBit(cpu.A, 4)
			return 8
		},
		// 0xA8: RES 5, B
		func(cpu *CPU) int {
			cpu.B = resetBit(cpu.B, 5)
			return 8
		},
		// 0xA9: RES 5, C
		func(cpu *CPU) int {
			cpu.C = resetBit(cpu.C, 5)
			return 8
		},
		// 0xAA: RES 5, D
		func(cpu *CPU) int {
			cpu.D = resetBit(cpu.D, 5)
			return 8
		},
		// 0xAB: RES 5, E
		func(cpu *CPU) int {
			cpu.E = resetBit(cpu.E, 5)
			return 8
		},
		// 0xAC: RES 5, H
		func(cpu *CPU) int {
			cpu.H = resetBit(cpu.H, 5)
			return 8
		},
		// 0xAD: RES 5, L
		func(cpu *CPU) int {
			cpu.L = resetBit(cpu.L, 5)
			return 8
		},
		// 0xAE: RES 5, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = resetBit(val, 5)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0xAF: RES 5, A
		func(cpu *CPU) int {
			cpu.A = resetBit(cpu.A, 5)
			return 8
		},

		// 0xB0: RES 6, B
		func(cpu *CPU) int {
			cpu.B = resetBit(cpu.B, 6)
			return 8
		},
		// 0xB1: RES 6, C
		func(cpu *CPU) int {
			cpu.C = resetBit(cpu.C, 6)
			return 8
		},
		// 0xB2: RES 6, D
		func(cpu *CPU) int {
			cpu.D = resetBit(cpu.D, 6)
			return 8
		},
		// 0xB3: RES 6, E
		func(cpu *CPU) int {
			cpu.E = resetBit(cpu.E, 6)
			return 8
		},
		// 0xB4: RES 6, H
		func(cpu *CPU) int {
			cpu.H = resetBit(cpu.H, 6)
			return 8
		},
		// 0xB5: RES 6, L
		func(cpu *CPU) int {
			cpu.L = resetBit(cpu.L, 6)
			return 8
		},
		// 0xB6: RES 6, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = resetBit(val, 6)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0xB7: RES 6, A
		func(cpu *CPU) int {
			cpu.A = resetBit(cpu.A, 6)
			return 8
		},
		// 0xB8: RES 7, B
		func(cpu *CPU) int {
			cpu.B = resetBit(cpu.B, 7)
			return 8
		},
		// 0xB9: RES 7, C
		func(cpu *CPU) int {
			cpu.C = resetBit(cpu.C, 7)
			return 8
		},
		// 0xBA: RES 7, D
		func(cpu *CPU) int {
			cpu.D = resetBit(cpu.D, 7)
			return 8
		},
		// 0xBB: RES 7, E
		func(cpu *CPU) int {
			cpu.E = resetBit(cpu.E, 7)
			return 8
		},
		// 0xBC: RES 7, H
		func(cpu *CPU) int {
			cpu.H = resetBit(cpu.H, 7)
			return 8
		},
		// 0xBD: RES 7, L
		func(cpu *CPU) int {
			cpu.L = resetBit(cpu.L, 7)
			return 8
		},
		// 0xBE: RES 7, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = resetBit(val, 7)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0xBF: RES 7, A
		func(cpu *CPU) int {
			cpu.A = resetBit(cpu.A, 7)
			return 8
		},

		// 0xC0: SET 0, B
		func(cpu *CPU) int {
			cpu.B = setBit(cpu.B, 0)
			return 8
		},
		// 0xC1: SET 0, C
		func(cpu *CPU) int {
			cpu.C = setBit(cpu.C, 0)
			return 8
		},
		// 0xC2: SET 0, D
		func(cpu *CPU) int {
			cpu.D = setBit(cpu.D, 0)
			return 8
		},
		// 0xC3: SET 0, E
		func(cpu *CPU) int {
			cpu.E = setBit(cpu.E, 0)
			return 8
		},
		// 0xC4: SET 0, H
		func(cpu *CPU) int {
			cpu.H = setBit(cpu.H, 0)
			return 8
		},
		// 0xC5: SET 0, L
		func(cpu *CPU) int {
			cpu.L = setBit(cpu.L, 0)
			return 8
		},
		// 0xC6: SET 0, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = setBit(val, 0)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0xC7: SET 0, A
		func(cpu *CPU) int {
			cpu.A = setBit(cpu.A, 0)
			return 8
		},
		// 0xC8: SET 1, B
		func(cpu *CPU) int {
			cpu.B = setBit(cpu.B, 1)
			return 8
		},
		// 0xC9: SET 1, C
		func(cpu *CPU) int {
			cpu.C = setBit(cpu.C, 1)
			return 8
		},
		// 0xCA: SET 1, D
		func(cpu *CPU) int {
			cpu.D = setBit(cpu.D, 1)
			return 8
		},
		// 0xCB: SET 1, E
		func(cpu *CPU) int {
			cpu.E = setBit(cpu.E, 1)
			return 8
		},
		// 0xCC: SET 1, H
		func(cpu *CPU) int {
			cpu.H = setBit(cpu.H, 1)
			return 8
		},
		// 0xCD: SET 1, L
		func(cpu *CPU) int {
			cpu.L = setBit(cpu.L, 1)
			return 8
		},
		// 0xCE: SET 1, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = setBit(val, 1)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0xCF: SET 1, A
		func(cpu *CPU) int {
			cpu.A = setBit(cpu.A, 1)
			return 8
		},

		// 0xD0: SET 2, B
		func(cpu *CPU) int {
			cpu.B = setBit(cpu.B, 2)
			return 8
		},
		// 0xD1: SET 2, C
		func(cpu *CPU) int {
			cpu.C = setBit(cpu.C, 2)
			return 8
		},
		// 0xD2: SET 2, D
		func(cpu *CPU) int {
			cpu.D = setBit(cpu.D, 2)
			return 8
		},
		// 0xD3: SET 2, E
		func(cpu *CPU) int {
			cpu.E = setBit(cpu.E, 2)
			return 8
		},
		// 0xD4: SET 2, H
		func(cpu *CPU) int {
			cpu.H = setBit(cpu.H, 2)
			return 8
		},
		// 0xD5: SET 2, L
		func(cpu *CPU) int {
			cpu.L = setBit(cpu.L, 2)
			return 8
		},
		// 0xD6: SET 2, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = setBit(val, 2)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0xD7: SET 2, A
		func(cpu *CPU) int {
			cpu.A = setBit(cpu.A, 2)
			return 8
		},
		// 0xD8: SET 3, B
		func(cpu *CPU) int {
			cpu.B = setBit(cpu.B, 3)
			return 8
		},
		// 0xD9: SET 3, C
		func(cpu *CPU) int {
			cpu.C = setBit(cpu.C, 3)
			return 8
		},
		// 0xDA: SET 3, D
		func(cpu *CPU) int {
			cpu.D = setBit(cpu.D, 3)
			return 8
		},
		// 0xDB: SET 3, E
		func(cpu *CPU) int {
			cpu.E = setBit(cpu.E, 3)
			return 8
		},
		// 0xDC: SET 3, H
		func(cpu *CPU) int {
			cpu.H = setBit(cpu.H, 3)
			return 8
		},
		// 0xDD: SET 3, L
		func(cpu *CPU) int {
			cpu.L = setBit(cpu.L, 3)
			return 8
		},
		// 0xDE: SET 3, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = setBit(val, 3)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0xDF: SET 3, A
		func(cpu *CPU) int {
			cpu.A = setBit(cpu.A, 3)
			return 8
		},

		// 0xE0: SET 4, B
		func(cpu *CPU) int {
			cpu.B = setBit(cpu.B, 4)
			return 8
		},
		// 0xE1: SET 4, C
		func(cpu *CPU) int {
			cpu.C = setBit(cpu.C, 4)
			return 8
		},
		// 0xE2: SET 4, D
		func(cpu *CPU) int {
			cpu.D = setBit(cpu.D, 4)
			return 8
		},
		// 0xE3: SET 4, E
		func(cpu *CPU) int {
			cpu.E = setBit(cpu.E, 4)
			return 8
		},
		// 0xE4: SET 4, H
		func(cpu *CPU) int {
			cpu.H = setBit(cpu.H, 4)
			return 8
		},
		// 0xE5: SET 4, L
		func(cpu *CPU) int {
			cpu.L = setBit(cpu.L, 4)
			return 8
		},
		// 0xE6: SET 4, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = setBit(val, 4)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0xE7: SET 4, A
		func(cpu *CPU) int {
			cpu.A = setBit(cpu.A, 4)
			return 8
		},
		// 0xE8: SET 5, B
		func(cpu *CPU) int {
			cpu.B = setBit(cpu.B, 5)
			return 8
		},
		// 0xE9: SET 5, C
		func(cpu *CPU) int {
			cpu.C = setBit(cpu.C, 5)
			return 8
		},
		// 0xEA: SET 5, D
		func(cpu *CPU) int {
			cpu.D = setBit(cpu.D, 5)
			return 8
		},
		// 0xEB: SET 5, E
		func(cpu *CPU) int {
			cpu.E = setBit(cpu.E, 5)
			return 8
		},
		// 0xEC: SET 5, H
		func(cpu *CPU) int {
			cpu.H = setBit(cpu.H, 5)
			return 8
		},
		// 0xED: SET 5, L
		func(cpu *CPU) int {
			cpu.L = setBit(cpu.L, 5)
			return 8
		},
		// 0xEE: SET 5, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = setBit(val, 5)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0xEF: SET 5, A
		func(cpu *CPU) int {
			cpu.A = setBit(cpu.A, 5)
			return 8
		},

		// 0xF0: SET 6, B
		func(cpu *CPU) int {
			cpu.B = setBit(cpu.B, 6)
			return 8
		},
		// 0xF1: SET 6, C
		func(cpu *CPU) int {
			cpu.C = setBit(cpu.C, 6)
			return 8
		},
		// 0xF2: SET 6, D
		func(cpu *CPU) int {
			cpu.D = setBit(cpu.D, 6)
			return 8
		},
		// 0xF3: SET 6, E
		func(cpu *CPU) int {
			cpu.E = setBit(cpu.E, 6)
			return 8
		},
		// 0xF4: SET 6, H
		func(cpu *CPU) int {
			cpu.H = setBit(cpu.H, 6)
			return 8
		},
		// 0xF5: SET 6, L
		func(cpu *CPU) int {
			cpu.L = setBit(cpu.L, 6)
			return 8
		},
		// 0xF6: SET 6, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = setBit(val, 6)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0xF7: SET 6, A
		func(cpu *CPU) int {
			cpu.A = setBit(cpu.A, 6)
			return 8
		},
		// 0xF8: SET 7, B
		func(cpu *CPU) int {
			cpu.B = setBit(cpu.B, 7)
			return 8
		},
		// 0xF9: SET 7, C
		func(cpu *CPU) int {
			cpu.C = setBit(cpu.C, 7)
			return 8
		},
		// 0xFA: SET 7, D
		func(cpu *CPU) int {
			cpu.D = setBit(cpu.D, 7)
			return 8
		},
		// 0xFB: SET 7, E
		func(cpu *CPU) int {
			cpu.E = setBit(cpu.E, 7)
			return 8
		},
		// 0xFC: SET 7, H
		func(cpu *CPU) int {
			cpu.H = setBit(cpu.H, 7)
			return 8
		},
		// 0xFD: SET 7, L
		func(cpu *CPU) int {
			cpu.L = setBit(cpu.L, 7)
			return 8
		},
		// 0xFE: SET 7, (HL)
		func(cpu *CPU) int {
			val := cpu.Memory.Read(cpu.HL())
			val = setBit(val, 7)
			cpu.Memory.Write(cpu.HL(), val)
			return 16
		},
		// 0xFF: SET 7, A
		func(cpu *CPU) int {
			cpu.A = setBit(cpu.A, 7)
			return 8
		},
	}
}
