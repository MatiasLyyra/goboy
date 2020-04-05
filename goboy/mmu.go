package goboy

// Memory is an interface for cpu to communicate with external devices (RAM, display etc.)
type Memory interface {
	Read(addr uint16) uint8
	Write(addr uint16, data uint8)
}

const (
	VBlankInt = iota
	LCDStatInt
	TimerInt
	SerialInt
	JoypadInt
)

const (
	AddrLCDC     = 0xFF40
	AddrLCDCStat = 0xFF41
	AddrSCY      = 0xFF42
	AddrSCX      = 0xFF43
	AddrLY       = 0xFF44
	AddrLYC      = 0xFF45
	AddrWY       = 0xFF4A
	AddrWX       = 0xFF4B
	AddrBGP      = 0xFF47
	AddrOBP0     = 0xFF48
	AddrOBP1     = 0xFF49
	AddrDMA      = 0xFF46
	AddrIE       = 0xFFFF
	AddrIF       = 0xFF0F
)

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

	ROMBankSize  = ROMBankEnd - ROMBankStart + 1
	RAMBankSize  = ExtRAMEnd - ExtRAMStart + 1
	VideoRAMSize = VideoRAMEnd - VideoRAMStart + 1
	OAMSize      = OAMEnd - OAMStart + 1
)

type MMU struct {
	Cartridge Memory
	GPU       Memory
	WRAM      Memory
	HRAM      Memory
	registers map[uint16]MemoryRegister
}

func (mmu *MMU) Read(addr uint16) uint8 {
	if reg, found := mmu.registers[addr]; found {
		return reg.Get()
	}
	switch {
	case ROMStart <= addr && addr <= ROMBankEnd:
		return mmu.Cartridge.Read(addr)
	case VideoRAMStart <= addr && addr <= VideoRAMEnd:
		return mmu.GPU.Read(addr)
	case ExtRAMStart <= addr && addr <= ExtRAMEnd:
		return mmu.Cartridge.Read(addr)
	case WRAMStart <= addr && addr <= EchoEnd:
		return mmu.WRAM.Read(addr)
	case OAMStart <= addr && addr <= OAMEnd:
		return mmu.GPU.Read(addr)
	// case IOPortsStart <= addr && addr <= IOPortsEnd:
	// 	return mmu.IOPorts.Read(addr)
	case HRAMStart <= addr && addr <= HRAMEnd:
		return mmu.HRAM.Read(addr)
	case addr == 0xFFFF:
	}
	return 0x0
}

func (mmu *MMU) Write(addr uint16, data uint8) {
	if reg, found := mmu.registers[addr]; found {
		reg.Set(data)
		return
	}
	switch {
	case ROMStart <= addr && addr <= ROMBankEnd:
		mmu.Cartridge.Write(addr, data)
	case VideoRAMStart <= addr && addr <= VideoRAMEnd:
		mmu.GPU.Write(addr, data)
	case ExtRAMStart <= addr && addr <= ExtRAMEnd:
		mmu.Cartridge.Write(addr, data)
	case WRAMStart <= addr && addr <= EchoEnd:
		mmu.WRAM.Write(addr, data)
	case OAMStart <= addr && addr <= OAMEnd:
		mmu.GPU.Write(addr, data)
	case HRAMStart <= addr && addr <= HRAMEnd:
		mmu.HRAM.Write(addr, data)
	}
}

func (mmu *MMU) Apply(addr uint16, reg MemoryRegister) {
	mmu.registers[addr] = reg
}