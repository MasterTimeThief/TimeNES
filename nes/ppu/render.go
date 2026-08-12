package ppu

import (
	"mtt/timenes/common"
	"mtt/timenes/nes/cartridge"
)

func RenderTileData() {
	if PPUScanline < 240 || PPUScanline == 261 {
		if (PPUDot > 0 && PPUDot <= 256) || (PPUDot > 320 && PPUDot <= 336) {
			//If this is a visible pixel, or preparing the start of the next scanline
			if PPUMASK_RenderBG || PPUMASK_RenderSprites {
				//If rendering is enabled
				if PPUMASK_RenderBG { //If rendering the background, update the shift registers for the background
					ppuShiftRegister_patternL = ppuShiftRegister_patternL << 1     //Shift 1 bit to the left
					ppuShiftRegister_patternH = ppuShiftRegister_patternH << 1     //Shift 1 bit to the left
					ppuShiftRegister_attributeL = ppuShiftRegister_attributeL << 1 //Shift 1 bit to the left
					ppuShiftRegister_attributeH = ppuShiftRegister_attributeH << 1 //Shift 1 bit to the left
				}
				if PPUMASK_RenderBG || PPUMASK_RenderSprites { //If rendering at all, let's decrement the X position of the objects
					if PPUDot > 1 && PPUDot <= 256 { //Don't decrement until dot 1
						for i := 0; i < 8; i++ {
							if ppu_SpriteXposition[i] > 0 {
								ppu_SpriteXposition[i]-- //Decrement the position of all objects in secondary OAM. When this is zero, the PPU can draw it
							} else {
								ppu_SpriteShiftRegisterL[i] = byte(ppu_SpriteShiftRegisterL[i] << 1) //Shift 1 bit to the left
								ppu_SpriteShiftRegisterH[i] = byte(ppu_SpriteShiftRegisterH[i] << 1) //Shift 1 bit to the left
							}
						}
					}
				}
				PPU8Steps()
			}
		}
	}
}

func RenderNextPixel( /*g *Game*/ ) {

	if PPUScanline < 240 && PPUDot > 0 && PPUDot <= 256 {
		var PalHi byte = 0  //Which color palette to use?
		var PalLow byte = 0 //Index into a color palette
		if PPUMASK_RenderBG && (PPUDot > 8 || PPUMASK_8pxMaskBG) {
			col0 := byte((ppuShiftRegister_patternL >> (15 - PPUScrollFineX)) & 1)
			col1 := byte((ppuShiftRegister_patternH >> (15 - PPUScrollFineX)) & 1)
			PalLow = byte(uint16(col1)<<1 | uint16(col0))

			pal0 := byte((ppuShiftRegister_attributeL >> (15 - PPUScrollFineX)) & 1)
			pal1 := byte((ppuShiftRegister_attributeH >> (15 - PPUScrollFineX)) & 1)
			PalHi = byte(uint16(pal1)<<1 | uint16(pal0))

			if PalLow == 0 && PalHi != 0 { //Color 0 of all palettes are mirrors of color 0 of palette 0
				PalHi = 0
			}
		}

		if PPUScanline >= 238 && PPUDot == 255 {
			common.Print("")
		}
		var SpritePalHi byte = 0        //Which color palette to use
		var SpritePalLow byte = 0       //Index into a color palette
		var SpritePriority bool = false //Is the sprite in front or behind the BG?
		if PPUMASK_RenderSprites && (PPUDot > 8 || PPUMASK_8pxMaskSprites) {
			for i := 0; i < 8; i++ {
				if ppu_SpriteXposition[i] == 0 && i <= int(ppuSecondaryOAMSize/4) { //If the sprite X position == 0 (The x position is decremented every ppu cycle)
					SpixelL := ((ppu_SpriteShiftRegisterL[i]) & 0x80) != 0
					SpixelH := ((ppu_SpriteShiftRegisterH[i]) & 0x80) != 0
					SpritePalLow = 0
					if SpixelL {
						SpritePalLow = 1
					}
					if SpixelH {
						SpritePalLow |= 2
					}

					SpritePalHi = byte((ppu_SpriteAttribute[i] & 0x03) | 0x04)
					SpritePriority = ((ppu_SpriteAttribute[i] >> 5) & 1) == 0
				} else {
					continue
				}
				if SpritePalLow != 0 {

					if i == 0 && ppuScanlineContainsSpriteZero && PalLow != 0 && PPUMASK_RenderBG && PPUDot < 256 {
						ppuScanlineContainsSpriteZero = false
						PPUSTATUS_SpriteZeroHit = true
					}
					break
				}
			}
		}

		if (SpritePriority && SpritePalLow != 0) || PalLow == 0 {
			PalLow = SpritePalLow
			PalHi = SpritePalHi
			if PalLow == 0 {
				PalHi = 0
			}
		}

		color := Palette[cartridge.PaletteRAM[(PalHi*4)+PalLow]&0x3F]

		//RenderPixel(color)
		FrameColorBuffer[FrameColorBufferPos] = color
		FrameColorBufferPos++
		//RenderNTSCPixel(PPUDot, pixel uint16, ppuCycleCounter int)

		//g.gameScreen.Set(PPUDot-1, PPUScanline, color)
	}
}
