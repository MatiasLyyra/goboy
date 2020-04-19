package goboy

import (
	"fmt"
)

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
	AddrSB       = 0xFF01
	AddrDIV      = 0xFF04
	AddrTIMA     = 0xFF05
	AddrTMA      = 0xFF06
	AddrTAC      = 0xFF07
	AddrJoy      = 0xFF00
)

// Defines different memory boundaries for Gameboy
const (
	// ROM
	ROMStart     = 0x0000
	ROMBootEnd   = 0x00FF
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
	WRAMSize     = WRAMEnd - WRAMStart + 1
	EchoSize     = EchoEnd - EchoStart + 1
	EchoOffset   = EchoStart - WRAMStart
	HRAMSize     = HRAMEnd - HRAMStart + 1
)

func NewMMU(cart *Cartridge) *MMU {
	mmu := &MMU{}
	mmu.BootEnabled = true
	mmu.Cartridge = cart
	mmu.Pad = &Joypad{
		mmu: mmu,
	}
	mmu.GPU = NewDisplay(mmu)
	mmu.WRAM = &GenericRAM{
		data:   make([]uint8, WRAMSize),
		offset: WRAMStart,
	}
	mmu.HRAM = &GenericRAM{
		data:   make([]uint8, HRAMSize),
		offset: HRAMStart,
	}
	mmu.registers = map[uint16]MemoryRegister{
		AddrLCDC:     NewRWRegister(0x91, 0),
		AddrLCDCStat: NewRWRegister(0, 0b111),
		// AddrLCDCStat: NewRWRegister(0, 0),
		AddrLYC:  NewRWRegister(0, 0),
		AddrLY:   NewRWRegister(0, 255),
		AddrIF:   NewRWRegister(0, 0),
		AddrIE:   NewRWRegister(0, 0),
		AddrSCX:  NewRWRegister(0, 0),
		AddrSCY:  NewRWRegister(0, 0),
		AddrBGP:  NewRWRegister(0, 0),
		AddrOBP0: NewRWRegister(0, 0),
		AddrOBP1: NewRWRegister(0, 0),
		AddrSB:   NewRWRegister(0, 0),
		// TODO: Should reset on write
		AddrDIV:  NewRWRegister(0, 0),
		AddrTIMA: NewRWRegister(0, 0),
		AddrTMA:  NewRWRegister(0, 0),
		AddrTAC:  NewRWRegister(0, 0),
		AddrDMA: CallbackRegister{
			fn: func(data uint8) {
				addr := uint16(data) << 8
				for i := OAMStart; i <= OAMEnd; i++ {
					mmu.Write(uint16(i), mmu.Read(addr))
					addr++
				}
			},
		},
		AddrWX:  NewRWRegister(0, 0),
		AddrWY:  NewRWRegister(0, 0),
		AddrJoy: mmu.Pad,
	}
	return mmu
}

type MMU struct {
	Cartridge   *Cartridge
	GPU         *Display
	Pad         *Joypad
	WRAM        Memory
	HRAM        Memory
	registers   map[uint16]MemoryRegister
	BootEnabled bool
}

func (mmu *MMU) Read(addr uint16) uint8 {
	if reg, found := mmu.registers[addr]; found {
		return reg.Get()
	}
	switch {
	case ROMStart <= addr && addr <= ROMBankEnd:
		// if addr <= ROMBootEnd && mmu.BootEnabled {
		// 	return boot2[addr]
		// }
		return mmu.Cartridge.Read(addr)
	case VideoRAMStart <= addr && addr <= VideoRAMEnd:
		return mmu.GPU.Read(addr)
	case ExtRAMStart <= addr && addr <= ExtRAMEnd:
		return mmu.Cartridge.Read(addr)
	case WRAMStart <= addr && addr <= EchoEnd:
		if addr >= EchoStart {
			addr -= EchoOffset
		}
		return mmu.WRAM.Read(addr)
	case OAMStart <= addr && addr <= OAMEnd:
		return mmu.GPU.Read(addr)
	// case IOPortsStart <= addr && addr <= IOPortsEnd:
	// 	return mmu.IOPorts.Read(addr)
	case HRAMStart <= addr && addr <= HRAMEnd:
		return mmu.HRAM.Read(addr)
	case addr == 0xFFFF:
	}
	return 0xFF
}

func (mmu *MMU) Write(addr uint16, data uint8) {
	if addr == 0xFF02 && data == 0x81 {
		fmt.Print(string(mmu.Read(AddrSB)))
	}
	if reg, found := mmu.registers[addr]; found {
		reg.Set(data)
		return
	}
	switch {
	case ROMStart <= addr && addr <= ROMBankEnd:
		mmu.Cartridge.Write(addr, data)
		return
	case VideoRAMStart <= addr && addr <= VideoRAMEnd:
		mmu.GPU.Write(addr, data)
		return
	case ExtRAMStart <= addr && addr <= ExtRAMEnd:
		mmu.Cartridge.Write(addr, data)
		return
	case WRAMStart <= addr && addr <= EchoEnd:
		if addr >= EchoStart {
			addr -= EchoOffset
		}
		mmu.WRAM.Write(addr, data)
		return
	case OAMStart <= addr && addr <= OAMEnd:
		mmu.GPU.Write(addr, data)
		return
	case HRAMStart <= addr && addr <= HRAMEnd:
		mmu.HRAM.Write(addr, data)
		return
	}
}

type GenericRAM struct {
	data   []uint8
	offset uint16
}

func (r *GenericRAM) Write(addr uint16, data uint8) {
	r.data[addr-r.offset] = data
}

func (r *GenericRAM) Read(addr uint16) uint8 {
	return r.data[addr-r.offset]
}
