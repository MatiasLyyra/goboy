package debug

import (
	"github.com/MatiasLyyra/goboy/goboy"
)

type Debugger struct {
	CPU         *goboy.CPU
	Breakpoints map[uint16]struct{}
}

func (d Debugger) StepToNextBreakpoint(sink chan<- [160 * 144]uint8) {
	for {
		if _, found := d.Breakpoints[d.CPU.PC]; !found {
			cycles := d.CPU.RunSingleOpcode()
			d.CPU.Memory.GPU.Run(cycles, sink)
		} else {
			break
		}
	}
}

func (d Debugger) TranslateOpcode(addr uint16) DecodedInsturction {
	opcode := d.CPU.Memory.Read(addr)
	if opcode == 0xCB {
		return DecodedInsturction{
			Instruction: bitInstructions[d.CPU.Memory.Read(addr+1)],
			Addr:        addr,
		}
	}
	instr := instructions[opcode]
	switch instr.Len {
	case 1:
		return DecodedInsturction{
			Instruction: instr,
			Addr:        addr,
		}
	case 2:
		data := d.CPU.Memory.Read(addr + 1)
		return DecodedInsturction{
			Instruction: instr,
			Data:        uint16(data),
			Addr:        addr,
			HasData:     true,
		}
	case 3:
		low := d.CPU.Memory.Read(addr + 1)
		high := d.CPU.Memory.Read(addr + 2)
		val := uint16(low) | uint16(high)<<8
		return DecodedInsturction{
			Instruction: instr,
			Data:        val,
			Addr:        addr,
			HasData:     true,
		}
	default:
		panic("Invalid opcode length")
	}
}

func (d Debugger) DecodeROM() ([]DecodedInsturction, map[uint16]int) {
	// addrToSearch := []uint16{0x0, 0x8, 0x10, 0x18, 0x20, 0x28, 0x30, 0x38, 0x100}
	addrToSearch := []uint16{d.CPU.PC}
	decoded := make(map[uint16]DecodedInsturction)
	var addr uint16
	for len(addrToSearch) > 0 {
		lastAddrIndex := len(addrToSearch) - 1
		addr = addrToSearch[lastAddrIndex]
		addrToSearch = addrToSearch[:lastAddrIndex]
		if _, found := decoded[addr]; found {
			continue
		}
		op := d.TranslateOpcode(addr)
		decoded[addr] = op
		// if !op.Stop {
		if op.Jump {
			if op.Relative {
				addrToSearch = append(addrToSearch, uint16(int(addr)+2+int(int8(op.Data))))
			} else {
				addrToSearch = append(addrToSearch, op.Data)
			}
		}
		addrToSearch = append(addrToSearch, op.Addr+uint16(op.Len))
		// }
	}
	decodedArray := make([]DecodedInsturction, 0, len(decoded))
	lookup := make(map[uint16]int)
	for i := 0; i < 0x10000; i++ {
		if op, found := decoded[uint16(i)]; found {
			decodedArray = append(decodedArray, op)
			lookup[uint16(i)] = len(decodedArray) - 1
		}
	}
	return decodedArray, lookup
}

func (d Debugger) ToggleBreakpoint(addr uint16) {
	if _, found := d.Breakpoints[addr]; found {
		delete(d.Breakpoints, addr)
	} else {
		d.Breakpoints[addr] = struct{}{}
	}
}
