package goboy

import (
	"fmt"
	"io"
)

type CGB uint8
type CartrigeType uint8
type ROMBanks uint8
type ExtRAMBanks uint8
type DestinationCode uint8

const (
	GB      CGB = 0x00
	NonCGB  CGB = 0x80
	OnlyCGB CGB = 0xC0
)

func (f CGB) String() string {
	switch f {
	case GB:
		return "Old Gameboy cartridge"
	case NonCGB:
		return "CGB cartridge with backwards compatibility"
	case OnlyCGB:
		return "Only CGB cartridge"
	default:
		panic("Invalid CGBFlag")
	}
	return ""
}

const (
	Banks_0   ROMBanks = 0x00
	Banks_4   ROMBanks = 0x01
	Banks_8   ROMBanks = 0x02
	Banks_16  ROMBanks = 0x03
	Banks_32  ROMBanks = 0x04
	Banks_64  ROMBanks = 0x05
	Banks_128 ROMBanks = 0x06
	Banks_256 ROMBanks = 0x07
	Banks_72  ROMBanks = 0x52
	Banks_80  ROMBanks = 0x53
	Banks_96  ROMBanks = 0x54
)

func (banks ROMBanks) Count() int {
	if banks <= Banks_256 {
		return 2 << int(banks)
	} else if banks == Banks_72 {
		return 72
	} else if banks == Banks_80 {
		return 80
	} else if banks == Banks_96 {
		return 96
	}
	panic("Invalid ROMBanks")
}

func (banks ROMBanks) String() string {
	switch banks {
	case Banks_0:
		return "Banks 0"
	case Banks_4:
		return "Banks 4"
	case Banks_8:
		return "Banks 8"
	case Banks_16:
		return "Banks 16"
	case Banks_32:
		return "Banks 32"
	case Banks_64:
		return "Banks 64"
	case Banks_128:
		return "Banks 128"
	case Banks_256:
		return "Banks 256"
	case Banks_72:
		return "Banks 72"
	case Banks_80:
		return "Banks 80"
	case Banks_96:
		return "Banks 96"
	default:
		panic("Invalid RomBanks")
	}
	return ""
}

const (
	CART_ROM_ONLY                CartrigeType = 0x00
	CART_MBC1                    CartrigeType = 0x01
	CART_MBC1_RAM                CartrigeType = 0x02
	CART_MBC1_RAM_BATTERY        CartrigeType = 0x03
	CART_MBC2                    CartrigeType = 0x05
	CART_MBC2_BATTERY            CartrigeType = 0x06
	CART_ROM_RAM                 CartrigeType = 0x08
	CART_ROM_RAM_BATTERY         CartrigeType = 0x09
	CART_MMM01                   CartrigeType = 0x0B
	CART_MMM01_RAM               CartrigeType = 0x0C
	CART_MMM01_RAM_BATTERY       CartrigeType = 0x0D
	CART_MBC3_TIMER_BATTERY      CartrigeType = 0x0F
	CART_MBC3_TIMER_RAM_BATTERY  CartrigeType = 0x10
	CART_MBC3                    CartrigeType = 0x11
	CART_MBC3_RAM                CartrigeType = 0x12
	CART_MBC3_RAM_BATTERY        CartrigeType = 0x13
	CART_MBC4                    CartrigeType = 0x15
	CART_MBC4_RAM                CartrigeType = 0x16
	CART_MBC4_RAM_BATTERY        CartrigeType = 0x17
	CART_MBC5                    CartrigeType = 0x19
	CART_MBC5_RAM                CartrigeType = 0x1A
	CART_MBC5_RAM_BATTERY        CartrigeType = 0x1B
	CART_MBC5_RUMBLE             CartrigeType = 0x1C
	CART_MBC5_RUMBLE_RAM         CartrigeType = 0x1D
	CART_MBC5_RUMBLE_RAM_BATTERY CartrigeType = 0x1E
	CART_POCKET_CAMERA           CartrigeType = 0xFC
	CART_BANDAI_TAMA5            CartrigeType = 0xFD
	CART_HuC3                    CartrigeType = 0xFE
	CART_HuC1_RAM_BATTERY        CartrigeType = 0xFF
)

func (ct CartrigeType) String() string {
	switch ct {
	case CART_ROM_ONLY:
		return "ROM_ONLY"
	case CART_MBC1:
		return "MBC1"
	case CART_MBC1_RAM:
		return "MBC1_RAM"
	case CART_MBC1_RAM_BATTERY:
		return "MBC1_RAM_BATTERY"
	case CART_MBC2:
		return "MBC2"
	case CART_MBC2_BATTERY:
		return "MBC2_BATTERY"
	case CART_ROM_RAM:
		return "ROM_RAM"
	case CART_ROM_RAM_BATTERY:
		return "ROM_RAM_BATTERY"
	case CART_MMM01:
		return "MMM01"
	case CART_MMM01_RAM:
		return "MMM01_RAM"
	case CART_MMM01_RAM_BATTERY:
		return "MMM01_RAM_BATTERY"
	case CART_MBC3_TIMER_BATTERY:
		return "MBC3_TIMER_BATTERY"
	case CART_MBC3_TIMER_RAM_BATTERY:
		return "MBC3_TIMER_RAM_BATTERY"
	case CART_MBC3:
		return "MBC3"
	case CART_MBC3_RAM:
		return "MBC3_RAM"
	case CART_MBC3_RAM_BATTERY:
		return "MBC3_RAM_BATTERY"
	case CART_MBC4:
		return "MBC4"
	case CART_MBC4_RAM:
		return "MBC4_RAM"
	case CART_MBC4_RAM_BATTERY:
		return "MBC4_RAM_BATTERY"
	case CART_MBC5:
		return "MBC5"
	case CART_MBC5_RAM:
		return "MBC5_RAM"
	case CART_MBC5_RAM_BATTERY:
		return "MBC5_RAM_BATTERY"
	case CART_MBC5_RUMBLE:
		return "MBC5_RUMBLE"
	case CART_MBC5_RUMBLE_RAM:
		return "MBC5_RUMBLE_RAM"
	case CART_MBC5_RUMBLE_RAM_BATTERY:
		return "MBC5_RUMBLE_RAM_BATTERY"
	case CART_POCKET_CAMERA:
		return "POCKET_CAMERA"
	case CART_BANDAI_TAMA5:
		return "BANDAI_TAMA5"
	case CART_HuC3:
		return "HuC3"
	case CART_HuC1_RAM_BATTERY:
		return "HuC1_RAM_BATTERY"
	default:
		panic("Invalid CartridgeType")
	}
}

const (
	ExtRAMNone ExtRAMBanks = 0x00
	ExtRAM2KB  ExtRAMBanks = 0x01
	ExtRAM8KB  ExtRAMBanks = 0x02
	ExtRAM32KB ExtRAMBanks = 0x03
)

func (ram ExtRAMBanks) String() string {
	switch ram {
	case ExtRAMNone:
		return "External RAM 0KB"
	case ExtRAM2KB:
		return "External RAM 2KB"
	case ExtRAM8KB:
		return "External RAM 8KB"
	case ExtRAM32KB:
		return "External RAM 32KB"
	default:
		panic("Invalid ExtRAMBanks")
	}
}

const (
	DestinationJP    DestinationCode = 0x00
	DestinationNonJP DestinationCode = 0x01
)

func (code DestinationCode) String() string {
	switch code {
	case DestinationJP:
		return "Japanese"
	case DestinationNonJP:
		return "Non-Japanese"
	default:
		panic("Invalid DestinationCode")
	}
}

type Cartridge struct {
	Bank0 [ROMBankSize]byte
	MBC   Memory
}

func (r *Cartridge) String() string {
	return fmt.Sprintf(
		`Title: %v
GB Mode: %v
Cartridge type: %v
ROM Banks: %v
RAM Banks: %v
Destination code: %v`,
		r.Title(), r.GCBFlag(), r.Cartridge(), r.ROMBanks(), r.ExtRAMBanks(), r.DestinationCode(),
	)
}

func LoadCartridge(r io.Reader) (*Cartridge, error) {
	rom := Cartridge{}
	// TODO: Error handling
	r.Read(rom.Bank0[:])
	switch rom.Cartridge() {
	case CART_ROM_ONLY:
		mbc := MBC0{}
		r.Read(mbc.rom[:])
		rom.MBC = &mbc
	case CART_MBC1, CART_MBC1_RAM, CART_MBC1_RAM_BATTERY:
		mbc := MBC1{}
		r.Read(mbc.rom[:])
		rom.MBC = &mbc
	}
	return &rom, nil
}

func (rom *Cartridge) Read(addr uint16) uint8 {
	if addr <= ROMEnd {
		return rom.Bank0[addr]
	}
	return rom.MBC.Read(addr)
}

func (rom *Cartridge) Title() string {
	titleBytes := rom.Bank0[0x134:0x144]
	title := ""
	for _, titleByte := range titleBytes {
		if titleByte == 32 || ('A' <= titleByte && titleByte <= 'z') {
			title += string(titleByte)
		} else {
			break
		}
	}
	return title
}

func (r *Cartridge) GCBFlag() CGB {
	flag := CGB(r.Bank0[0x143])
	if flag == NonCGB || flag == OnlyCGB {
		return flag
	}
	return GB
}

func (r *Cartridge) Cartridge() CartrigeType {
	return CartrigeType(r.Bank0[0x147])
}

func (r *Cartridge) ROMBanks() ROMBanks {
	return ROMBanks(r.Bank0[0x148])
}

func (r *Cartridge) ExtRAMBanks() ExtRAMBanks {
	return ExtRAMBanks(r.Bank0[0x148])
}

func (r *Cartridge) DestinationCode() DestinationCode {
	return DestinationCode(r.Bank0[0x14A])
}
