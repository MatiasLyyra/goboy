package goboy

// Memory is an interface for cpu to communicate with external devices (RAM, display etc.)
type Memory interface {
	Read(addr uint16) uint8
	Write(addr uint16, data uint8)
}

// Defines different memory boundaries for Gameboy
const (
	// ROM
	ROMStart     = 0x0000
	ROMEnd       = 0x3FFF
	ROMBankStart = 0x4000
	ROMBankEnd   = 0x7FFF

	ExtRAMStart = 0xA000
	ExtRAMEnd   = 0xBFFF

	// Video / GPU
	VideoRAMStart = 0x8000
	VideoRAMEnd   = 0x9FFF
	OAMStart      = 0xFE00
	OAMEnd        = 0xFE9F
	// Work RAM
	WRAMStart = 0xC000
	WRAMEnd   = 0xDFFF
	EchoStart = 0xE000
	EchoEnd   = 0xFDFF

	HRAMStart = 0xFF80
	HRAMEnd   = 0xFFFE

	// I/O Ports
	IOPortsStart = 0xFF00
	IOPortsEnd   = 0xFF7F

	ROMBankSize = ROMBankEnd - ROMBankStart + 1
	RAMBankSize = ExtRAMEnd - ExtRAMStart + 1
)

type MMU struct {
	ROM     Memory
	GPU     Memory
	WRAM    Memory
	HRAM    Memory
	IOPorts Memory
}

func (mmu *MMU) Read(addr uint16) uint8 {
	switch {
	case ROMStart <= addr && addr <= ROMBankEnd:
		return mmu.ROM.Read(addr)
	case VideoRAMStart <= addr && addr <= VideoRAMEnd:
		return mmu.GPU.Read(addr)
	case ExtRAMStart <= addr && addr <= ExtRAMEnd:
		return mmu.ROM.Read(addr)
	case WRAMStart <= addr && addr <= EchoEnd:
		return mmu.WRAM.Read(addr)
	case OAMStart <= addr && addr <= OAMEnd:
		return mmu.GPU.Read(addr)
	case IOPortsStart <= addr && addr <= IOPortsEnd:
		return mmu.IOPorts.Read(addr)
	case HRAMStart <= addr && addr <= HRAMEnd:
		return mmu.HRAM.Read(addr)
	case addr == 0xFFFF:
	}
	panic("Invalid addr")
}

func (mmu *MMU) Write(addr uint16, data uint8) {
	switch {
	case ROMStart <= addr && addr <= ROMBankEnd:
		mmu.ROM.Write(addr, data)
		return
	case VideoRAMStart <= addr && addr <= VideoRAMEnd:
		mmu.GPU.Write(addr, data)
		return
	case ExtRAMStart <= addr && addr <= ExtRAMEnd:
		mmu.ROM.Write(addr, data)
		return
	case WRAMStart <= addr && addr <= EchoEnd:
		mmu.WRAM.Write(addr, data)
		return
	case OAMStart <= addr && addr <= OAMEnd:
		mmu.GPU.Write(addr, data)
		return
	case IOPortsStart <= addr && addr <= IOPortsEnd:
		mmu.IOPorts.Write(addr, data)
		return
	case HRAMStart <= addr && addr <= HRAMEnd:
		mmu.HRAM.Write(addr, data)
		return
	case addr == 0xFFFF:
		return
	}
	panic("Invalid addr")
}
