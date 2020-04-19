package goboy

import (
	"sort"
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
	X      int
	Y      int
	TileID uint8
	Flags  uint8
}

func (s Sprite) Upper8x16() uint8 {
	return s.TileID & 0xFE
}

func (s Sprite) Lower8x16() uint8 {
	return s.TileID | 0x01
}

func NewDisplay(mmu *MMU) *Display {
	d := &Display{
		mmu: mmu,
	}
	// // TODO: Read the actual values from memory
	// // These should be initially zero
	// for i := 0; i < 4; i++ {
	// 	d.spritePalettes[0][i] = uint8(i)
	// 	d.spritePalettes[1][i] = uint8(i)
	// 	d.bgPalette[i] = uint8(i)
	// }
	return d
}

type Display struct {
	mmu *MMU

	VRAM [VideoRAMSize]uint8
	oam  [OAMSize]uint8

	cycles         int
	row            int
	spritePalettes [2][4]uint8
	bgPalette      [4]uint8
	priorityBuffer [160]uint8
	spriteBuffer   [160 * 144]uint8
}

func (d *Display) Run(cycles int) bool {
	lcdc := d.mmu.registers[AddrLCDC]
	if lcdc.Get()&(1<<7) == 0 {
		return false
	}
	var hasDrawn bool
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
		// if lcdcStat.Get()&(1<<6) != 0 {
		ifReg.RawSet(setBit(ifReg.Get(), LCDStatInt))
		// }
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
			// if lcdcStat.Get()&(1<<4) != 0 {
			ifReg.RawSet(setBit(ifReg.Get(), VBlankInt))
			// }
			currentMode = ModeVBlank
			hasDrawn = true
		}
	} else if rowCycles < OAMDuration {
		if currentMode != ModeOAM {
			lcdcStat.RawSet((lcdcStat.Get() & ^uint8(0x3)) | ModeOAM)
			// Request LCD STAT interrupt if OAM Interrupts are enabled
			// if lcdcStat.Get()&(1<<5) != 0 {
			ifReg.RawSet(setBit(ifReg.Get(), LCDStatInt))
			// }
			currentMode = ModeOAM
		}
	} else if rowCycles < OAMDuration+TransferDuration {
		if currentMode != ModeTransfer {
			lcdcStat.RawSet((lcdcStat.Get() & ^uint8(0x3)) | ModeTransfer)
			currentMode = ModeTransfer
		}
	} else {
		if currentMode != ModeHBlank {
			currentMode = ModeHBlank
			lcdcStat.RawSet((lcdcStat.Get() & ^uint8(0x3)) | ModeHBlank)
			var (
				bgp  = d.mmu.registers[AddrBGP].Get()
				obp0 = d.mmu.registers[AddrOBP0].Get()
				obp1 = d.mmu.registers[AddrOBP1].Get()
			)
			// Start of HBlank, draw row into screen buffer
			// Request LCD STAT interrupt if HBlank Interrupts are enabled
			for i := range d.priorityBuffer {
				d.priorityBuffer[i] = prioUndrawn
			}
			d.bgPalette[0] = bgp & 3
			d.bgPalette[1] = (bgp & (3 << 2)) >> 2
			d.bgPalette[2] = (bgp & (3 << 4)) >> 4
			d.bgPalette[3] = (bgp & (3 << 6)) >> 6

			d.spritePalettes[0][1] = (obp0 & (3 << 2)) >> 2
			d.spritePalettes[0][2] = (obp0 & (3 << 4)) >> 4
			d.spritePalettes[0][3] = (obp0 & (3 << 6)) >> 6

			d.spritePalettes[1][1] = (obp1 & (3 << 2)) >> 2
			d.spritePalettes[1][2] = (obp1 & (3 << 4)) >> 4
			d.spritePalettes[1][3] = (obp1 & (3 << 6)) >> 6

			d.drawBackground(row)
			d.drawSpriteRow(row)
			// if lcdcStat.Get()&(1<<3) != 0 {
			ifReg.RawSet(setBit(ifReg.Get(), LCDStatInt))
			// }
			// sink <- d.spriteBuffer
		}
	}
	lcdcStat.RawSet((lcdcStat.Get() & ^uint8(0x3)) | currentMode)
	ly.RawSet(uint8(row))
	return hasDrawn
}

func (d *Display) Read(addr uint16) uint8 {
	if VideoRAMStart <= addr && addr <= VideoRAMEnd {
		return d.VRAM[addr-VideoRAMStart]
	}
	if OAMStart <= addr && addr <= OAMEnd {
		return d.oam[addr-OAMStart]
	}
	panic("Invalid addr")
}

func (d *Display) Write(addr uint16, data uint8) {
	if VideoRAMStart <= addr && addr <= VideoRAMEnd {
		d.VRAM[addr-VideoRAMStart] = data
		return
	}
	if OAMStart <= addr && addr <= OAMEnd {
		d.oam[addr-OAMStart] = data
		return
	}
	panic("Invalid addr")
}

func (d *Display) drawBackground(row int) {
	var (
		scx         = 0 // d.mmu.registers[AddrSCX].Get()
		scy         = 0 // d.mmu.registers[AddrSCY].Get()
		lcdc        = d.mmu.registers[AddrLCDC]
		useLowerMap = lcdc.Get()&(1<<3) == 0
		tileData    []uint8
	)
	if useLowerMap {
		tileData = d.backgroundTileMap1()
	} else {
		tileData = d.backgroundTileMap2()
	}
	var i int
outer:
	for {
		pixelPosX := (i + int(scx)) % 256
		pixelPosY := (row + int(scy)) % 256
		tileX := pixelPosX / 8
		tileY := pixelPosY / 8
		tileID := tileData[tileY*32+tileX]
		tile := d.GetTile(tileID, false)
		tileRowStart := ((row) % 8) * 2
		pixels := getPixelRow([2]uint8{tile[tileRowStart], tile[tileRowStart+1]}, false)
		for _, val := range pixels {
			d.spriteBuffer[row*160+i] = d.bgPalette[val]

			if val != 0 {
				d.priorityBuffer[i] = prioBackground
			}
			i++
			if i >= 160 {
				break outer
			}
		}
	}
}

func (d *Display) drawSpriteRow(row int) {
	var spriteHeight int = 8
	var longSprites = d.mmu.registers[AddrLCDC].Get()&(1<<2) != 0
	// Check LCDC register bit 2 for current sprite size
	// 0 = 8x8, 1 = 8x16
	if longSprites {
		spriteHeight = 16
	}
	// Filter sprites that are visible in the current scanline
	sprites := d.FilterSprites(func(sprite Sprite) bool {
		top := int(sprite.Y)
		lower := int(sprite.Y) + spriteHeight
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
		xMirror := sprite.Flags&(1<<5) != 0
		yMirror := sprite.Flags&(1<<6) != 0
		if longSprites {
			// Check if we are drawing upper or lower portion of the sprite
			if yMirror {
				if row >= int(sprite.Y)+8 {
					tileID &= 0xFE
				} else {
					tileID |= 0x01
				}
			} else {
				if row < int(sprite.Y)+8 {
					tileID &= 0xFE
				} else {
					tileID |= 0x01
				}
			}

		}
		tile := d.GetTile(tileID, true)
		tileRowStart := ((row - (int(sprite.Y) - 16)) % 8) * 2
		// Check flag bit 7 if we are drawing sprite always above bg
		aboveBG := sprite.Flags&(1<<7) == 0
		spritePaletteID := (sprite.Flags & (1 << 4)) >> 4

		// fmt.Println(tileRowStart)
		if yMirror {
			tileRowStart = 14 - tileRowStart
		}
		pixels := getPixelRow([2]uint8{tile[tileRowStart], tile[tileRowStart+1]}, xMirror)
		for i, pixelVal := range pixels {
			pixelX := int(sprite.X) + i
			idx := row*160 + pixelX
			if pixelX < 0 || pixelX >= 160 {
				continue
			}
			// Check if that we aren't drawing above another sprite that was already drawn
			if d.priorityBuffer[pixelX] != prioSprite && (aboveBG || d.priorityBuffer[pixelX] == prioUndrawn) {
				if pixelVal != 0 {
					d.priorityBuffer[pixelX] = prioSprite
					d.spriteBuffer[idx] = d.spritePalettes[spritePaletteID][pixelVal]
				}
			}
		}
	}
}

func getPixelRow(rawTileRow [2]uint8, mirror bool) (pixels [8]uint8) {

	tileRow1, tileRow2 := rawTileRow[0], rawTileRow[1]
	for i := 0; i < 8; i++ {
		idx := i
		if mirror {
			idx = 7 - i
		}
		pixels[idx] = (tileRow1 & (1 << (7 - i))) >> (7 - i)
		pixels[idx] |= (tileRow2 & (1 << (7 - i))) >> (7 - i) << 1
		// pixels[i] &= 3
	}
	return
}

func (d *Display) GetTile(id uint8, spriteData bool) []uint8 {
	if spriteData || d.mmu.registers[AddrLCDC].Get()&(1<<4) != 0 {
		memoryLoc := int(id) * 16
		return d.tilePatternTable1()[memoryLoc : memoryLoc+16]
	}
	memoryLoc := 0x800 + int(int8(id))*16
	return d.tilePatternTable2()[memoryLoc : memoryLoc+16]
}

func (d *Display) tilePatternTable1() []uint8 {
	// Tile Pattern Table 1 starts at 0x8000 and ends at 0x8FFF
	// Return slice from 0x0 to 0x1000 (exclusive)
	return d.VRAM[0x0000:0x1000]
}
func (d *Display) tilePatternTable2() []uint8 {
	// Tile Pattern Table 2 starts at 0x8800 and ends at 0x97FF
	// Return slice from 0x0800 to 0x1800 (exclusive)
	return d.VRAM[0x0800:0x1800]
}

func (d *Display) backgroundTileMap1() []uint8 {
	return d.VRAM[0x1800:0x1C00]
}

func (d *Display) backgroundTileMap2() []uint8 {
	return d.VRAM[0x1C00:0x2000]
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
		Y:      int(spriteData[0]) - 16,
		X:      int(spriteData[1]) - 8,
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

func (d *Display) ScreenBuffer() []uint8 {
	return d.spriteBuffer[:]
}
