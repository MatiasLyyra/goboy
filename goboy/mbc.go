package goboy

type MBC0 struct {
	rom [ROMBankSize]byte
}

func (mbc *MBC0) Read(addr uint16) uint8 {
	return mbc.rom[uint(addr)-ROMBankStart]
}

func (mbc *MBC0) Write(addr uint16, data uint8) {
	// No writing on my lawn!
}

type MBC1 struct {
	ramEnabled    uint8
	romBankNumber uint8
	ramBankNumber uint8
	romModeSelect uint8

	rom [128 * ROMBankSize]byte
	ram [4 * RAMBankSize]byte
}

func (mbc *MBC1) Read(addr uint16) uint8 {
	if ROMBankStart <= addr && addr <= ROMBankEnd {
		return mbc.rom[uint(mbc.SelectedROM())*ROMBankSize+(uint(addr)-ROMBankStart)]
	} else if mbc.RAMEnabled() && ExtRAMStart <= addr && addr <= ExtRAMEnd {
		return mbc.ram[uint(mbc.SelectedRAM())*RAMBankSize+(uint(addr)-ExtRAMStart)]
	}
	return 0
}

func (mbc *MBC1) RAMEnabled() bool {
	return mbc.ramEnabled&0xF == 0xA
}

func (mbc *MBC1) SelectedROM() uint8 {
	romBank := mbc.romBankNumber & 0x1F
	if mbc.RomModeSelected() {
		romBank &= (mbc.ramBankNumber & 0x3) << 4
	}
	return romBank
}

func (mbc *MBC1) SelectedRAM() uint8 {
	if !mbc.RomModeSelected() {
		return mbc.ramBankNumber & 0x3
	}
	return 0
}

func (mbc *MBC1) RomModeSelected() bool {
	return mbc.romModeSelect != 0x1
}

func (mbc *MBC1) Write(addr uint16, data uint8) {
	if addr <= 0x1FFF {
		mbc.ramEnabled = data
		return
	}
	if addr <= 0x3FFF {
		// Selecting zero bank will select bank 1
		if data&0x1F == 0 {
			data++
		}
		mbc.romBankNumber = data
		return
	}
	if addr <= 0x5FFF {
		mbc.ramBankNumber = data
		return
	}
	if addr <= 0x7FFF {
		mbc.romModeSelect = data
		return
	}

}
