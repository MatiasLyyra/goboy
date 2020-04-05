package goboy

import "sort"

// Available colors
const (
	// 0xFFFFFF
	ColorWhite = iota
	// 0xAAAAAA
	ColorLightGray
	// 0x555555
	ColorDarkGray
	// 0x000000
	ColorBlack
)

const (
	prioUndrawn = iota
	prioBackground
	prioSprite
)

// Cycle:
// OAM (Mode 2): [0, 80) => Transfer (Mode 3): [80, 252) => HBlank (Mode 0) [252, 456)
// Repeated for 144 times and then VBlank for 4560 cycles
const (
	ModeHBlank = iota
	ModeVBlank
	ModeOAM
	ModeTransfer

	HBlankDuration   = 204
	VBlankDuration   = 4560
	OAMDuration      = 80
	TransferDuration = 172
)

type Sprite struct {
	ID     uint8
	X      uint8
	Y      uint8
	TileID uint8
	Flags  uint8
}

func (s Sprite) Upper8x16() uint8 {
	return s.TileID & 0xFE
}

func (s Sprite) Lower8x16() uint8 {
	return s.TileID | 0x01
}

type Display struct {
	mmu *MMU

	vram [VideoRAMSize]uint8
	oam  [OAMSize]uint8

	cycles         int
	row            int
	spritePalettes [2][4]uint8
	bgPalette      [4]uint8
	priorityBuffer [160 * 144]uint8
	spriteBuffer   [160 * 144]uint8
}

func (d *Display) Run(cycles int) {
	var (
		lcdcStat = d.mmu.registers[AddrLCDCStat]
		ifReg    = d.mmu.registers[AddrIF]
		ly       = d.mmu.registers[AddrLY]
		lyc      = d.mmu.registers[AddrLYC]
	)
	d.cycles += cycles

	currentMode := lcdcStat.Get() & 0x3
	rowCycles := d.cycles % (OAMDuration + TransferDuration + HBlankDuration)
	row := d.cycles / (OAMDuration + TransferDuration + HBlankDuration)
	if uint8(row) == lyc.Get() {
		// Set Coincidence flag
		lcdcStat.RawSet(setBit(lcdcStat.Get(), 2))
		// Request LCD STAT interrupt if Coincidence Interrupts are enabled
		if lcdcStat.Get()&(1<<6) != 0 {
			ifReg.RawSet(setBit(ifReg.Get(), LCDStatInt))
		}
	} else {
		// Unset coincidence flag
		lcdcStat.RawSet(resetBit(lcdcStat.Get(), 2))
	}
	// VBlank
	if row >= 144 {
		// End of VBlank
		if row > 153 {
			row = 0
			d.cycles = 0
		}
		if currentMode != ModeVBlank {
			lcdcStat.RawSet((lcdcStat.Get() & ^uint8(0x3)) | ModeVBlank)
			if lcdcStat.Get()&(1<<4) != 0 {
				ifReg.RawSet(setBit(ifReg.Get(), VBlankInt))
			}
		}
	} else if rowCycles < OAMDuration {
		if currentMode != ModeOAM {
			lcdcStat.RawSet((lcdcStat.Get() & ^uint8(0x3)) | ModeOAM)
			// Request LCD STAT interrupt if OAM Interrupts are enabled
			if lcdcStat.Get()&(1<<5) != 0 {
				ifReg.RawSet(setBit(ifReg.Get(), LCDStatInt))
			}
		}
	} else if rowCycles < OAMDuration+TransferDuration {
		if currentMode != ModeTransfer {
			lcdcStat.RawSet((lcdcStat.Get() & ^uint8(0x3)) | ModeTransfer)
		}
	} else {
		if currentMode != ModeHBlank {
			lcdcStat.RawSet((lcdcStat.Get() & ^uint8(0x3)) | ModeHBlank)
			// Start of HBlank, draw row into screen buffer
			// Request LCD STAT interrupt if HBlank Interrupts are enabled
			d.drawRow(row)
			if lcdcStat.Get()&(1<<3) != 0 {
				ifReg.RawSet(setBit(ifReg.Get(), LCDStatInt))
			}
		}
	}
	ly.RawSet(uint8(row))
}

func (d *Display) Read(addr uint16) uint8 {
	if VideoRAMStart <= addr && addr <= VideoRAMEnd {
		return d.vram[addr-VideoRAMStart]
	}
	if OAMStart <= addr && addr <= OAMEnd {
		return d.oam[addr-OAMStart]
	}
	panic("Invalid addr")
}

func (d *Display) Write(addr uint16, data uint8) {
	if VideoRAMStart <= addr && addr <= VideoRAMEnd {
		d.vram[addr-VideoRAMStart] = data
		return
	}
	if OAMStart <= addr && addr <= OAMEnd {
		d.oam[addr-OAMStart] = data
		return
	}
	panic("Invalid addr")
}

func (d *Display) drawRow(row int) {
	var spriteHeight int = 8
	var longSprites = d.mmu.registers[AddrLCDC].Get()&(1<<2) != 0
	// Check LCDC register bit 2 for current sprite size
	// 0 = 8x8, 1 = 8x16
	if longSprites {
		spriteHeight = 16
	}
	// Filter sprites that are visible in the current scanline
	sprites := d.FilterSprites(func(sprite Sprite) bool {
		top := int(sprite.Y) - spriteHeight
		lower := int(sprite.Y)
		return top <= row && row < lower
	})
	// Sort sprites according to their X values
	// Use stable sort because we also want to keep memory location order
	// if x-values are equal
	sort.SliceStable(sprites, func(i, j int) bool {
		return sprites[i].X < sprites[j].X
	})
	// Note: sprite X and Y are location of lower right corner of the sprite
	for i := 0; i < len(sprites) && i < 10; i++ {
		sprite := sprites[i]
		// TODO: Sprite mirroring
		tileID := sprite.TileID
		if longSprites {
			// Check if we are drawing upper or lower portion of the sprite
			if row < int(sprite.Y)-8 {
				tileID &= 0xFE
			} else {
				tileID |= 0x01
			}
		}
		tile := d.GetTile(tileID, true)
		tileRowStart := ((row - (int(sprite.Y) - 16)) % 8) * 2
		// Check flag bit 7 if we are drawing sprite always above bg
		aboveBG := sprite.Flags&(1<<7) == 0
		spritePaletteID := (sprite.Flags & (1 << 4)) >> 4
		if len(tile) != 16 || tileRowStart < 0 {
			// TODO: Remove this debug sanity check
			panic("Invalid tile length")
		}
		pixels := getPixelRow([2]uint8{tile[tileRowStart], tile[tileRowStart+1]})
		for i, pixelVal := range pixels {
			pixelX := int(sprite.X) - 8 + i
			idx := row*144 + pixelX

			// Check if that we aren't drawing above another sprite that was already drawn
			//
			if d.priorityBuffer[idx] != prioSprite && (aboveBG || d.priorityBuffer[idx] == prioUndrawn) && pixelX < 160 {
				if pixelVal != 0 {
					d.priorityBuffer[idx] = prioSprite
				}
				d.spriteBuffer[idx] = d.spritePalettes[spritePaletteID][pixelVal]
			}
		}
	}
}

func getPixelRow(rawTileRow [2]uint8) (pixels [8]uint8) {
	tileRow1, tileRow2 := rawTileRow[0], rawTileRow[1]
	pixels[0] = (tileRow1 & (3 << 6)) >> 6
	pixels[1] = (tileRow1 & (3 << 4)) >> 4
	pixels[2] = (tileRow1 & (3 << 2)) >> 2
	pixels[3] = tileRow1 & 3
	pixels[4] = (tileRow2 & (3 << 6)) >> 6
	pixels[5] = (tileRow2 & (3 << 4)) >> 4
	pixels[6] = (tileRow2 & (3 << 2)) >> 2
	pixels[7] = tileRow2 & 3
	return
}

func (d *Display) GetTile(id uint8, spriteData bool) []uint8 {
	if spriteData || d.mmu.registers[AddrLCDC].Get()&(1<<4) != 0 {
		return d.tilePatternTable1()[id*16 : id*16+16]
	}
	memoryLoc := 0x800 + int(int8(id))
	return d.tilePatternTable2()[memoryLoc : memoryLoc+16]
}

func (d *Display) tilePatternTable1() []uint8 {
	// Tile Pattern Table 1 starts at 0x8000 and ends at 0x8FFF
	// Return slice from 0x0 to 0x1000 (exclusive)
	return d.vram[0x0000:0x1000]
}
func (d *Display) tilePatternTable2() []uint8 {
	// Tile Pattern Table 2 starts at 0x8800 and ends at 0x97FF
	// Return slice from 0x0800 to 0x1800 (exclusive)
	return d.vram[0x0800:0x1800]
}

func (d *Display) GetSprite(id int) Sprite {
	// One sprite takes 4 bytes
	if id >= 40 {
		panic("Invalid sprite id")
	}
	startAddr := id * 4
	spriteData := d.oam[startAddr : startAddr+4]
	return Sprite{
		ID:     uint8(id),
		X:      spriteData[0],
		Y:      spriteData[1],
		TileID: spriteData[2],
		Flags:  spriteData[3],
	}
}

func (d *Display) FilterSprites(filterFunc func(Sprite) bool) []Sprite {
	var sprites []Sprite
	for i := 0; i < 40; i++ {
		sprite := d.GetSprite(i)
		if filterFunc(sprite) {
			sprites = append(sprites, sprite)
		}
	}
	return sprites
}
