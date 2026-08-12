package nes

import (
	"image/color"
	"mtt/timenes/common"
	"mtt/timenes/nes/cartridge"
	"mtt/timenes/nes/cartridge/mappers"
)

var WriteLatch bool        //PPU's w register
var TransferAddress uint16 //PPU's t register
var VRAMAddress uint16     //PPU's v register
var PPUReadBuffer byte
var NMILevelDetector, DoNMI bool

// $2000: PPUCTRL
var PPUCTRL_NametableSelect byte    // PPUCTRL Bit 1 & 2
var PPUCTRL_VRAMInc32Mode bool      // PPUCTRL Bit 3
var PPUCTRL_SpritePatternTable bool // PPUCTRL Bit 4
var PPUCTRL_BGPatternTable bool     // PPUCTRL Bit 5
var PPUCTRL_Use8x16Sprites bool     // PPUCTRL Bit 6
var PPUCTRL_EnableNMI bool          // PPUCTRL Bit 8

// $2001: PPUMASK
var PPUMASK_Greyscale bool      // PPUMASK Bit 0
var PPUMASK_8pxMaskBG bool      // PPUMASK Bit 1
var PPUMASK_8pxMaskSprites bool // PPUMASK Bit 2
var PPUMASK_RenderBG bool       // PPUMASK Bit 3
var PPUMASK_RenderSprites bool  // PPUMASK Bit 4
var PPUMASK_EmphasisRed bool    // PPUMASK Bit 5
var PPUMASK_EmphasisGreen bool  // PPUMASK Bit 6
var PPUMASK_EmphasisBlue bool   // PPUMASK Bit 7

// $2002: PPUSTATUS
var PPUSTATUS_Overflow bool      // PPUSTATUS Bit 5
var PPUSTATUS_SpriteZeroHit bool // PPUSTATUS Bit 6
var PPUSTATUS_VBlank bool        // PPUSTATUS Bit 7

var PPUDot int      //The X position of the scanning beam
var PPUScanline int //The Y position of the scanning beam
var ppuShiftRegister_patternL, ppuShiftRegister_patternH, ppuShiftRegister_attributeL, ppuShiftRegister_attributeH uint16
var ppu8Step_patternLowBitPlane, ppu8Step_patternHighBitPlane, ppu8Step_attribute, ppu8Step_NextCharacter, ppu8Step_temp byte
var PPUScrollFineX, ppuScrollFineY byte
var DrawNewFrame bool = false
var FrameColorBuffer [61440]color.RGBA
var FrameColorBufferPos int = 0

var OAM [0x100]byte
var SecondaryOAM [0x20]byte
var ppuSpriteEvalTemp byte
var ppuOAMAddress, ppuSecondaryOAMAddress, ppuSecondaryOAMSize uint16
var ppuSecondaryOAMFull, ppuScanlineContainsSpriteZero, ppuSpriteEvaluationOAMOverflowed bool
var ppuSpriteEvalTick int

var ppu_SpriteShiftRegisterL [8]byte
var ppu_SpriteShiftRegisterH [8]byte

var ppu_SpriteAttribute [8]byte
var ppu_SpritePattern [8]byte
var ppu_SpriteXposition [8]byte
var ppu_SpriteYposition [8]byte

var Palette = [64]color.RGBA{
	{R: 0x65, G: 0x65, B: 0x65, A: 0xFF},
	{R: 0x00, G: 0x2A, B: 0x84, A: 0xFF},
	{R: 0x15, G: 0x13, B: 0xA2, A: 0xFF},
	{R: 0x3A, G: 0x01, B: 0x9E, A: 0xFF},
	{R: 0x59, G: 0x00, B: 0x7A, A: 0xFF},
	{R: 0x6A, G: 0x00, B: 0x3E, A: 0xFF},
	{R: 0x68, G: 0x08, B: 0x00, A: 0xFF},
	{R: 0x53, G: 0x1D, B: 0x00, A: 0xFF},
	{R: 0x32, G: 0x34, B: 0x00, A: 0xFF},
	{R: 0x0D, G: 0x46, B: 0x00, A: 0xFF},
	{R: 0x00, G: 0x4F, B: 0x00, A: 0xFF},
	{R: 0x00, G: 0x4C, B: 0x09, A: 0xFF},
	{R: 0x00, G: 0x3F, B: 0x4B, A: 0xFF},
	{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
	{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
	{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},

	{R: 0xAE, G: 0xAE, B: 0xAE, A: 0xFF},
	{R: 0x17, G: 0x5F, B: 0xD6, A: 0xFF},
	{R: 0x43, G: 0x41, B: 0xFF, A: 0xFF},
	{R: 0x75, G: 0x29, B: 0xFA, A: 0xFF},
	{R: 0x9E, G: 0x1D, B: 0xCA, A: 0xFF},
	{R: 0xB4, G: 0x20, B: 0x7B, A: 0xFF},
	{R: 0xB1, G: 0x33, B: 0x22, A: 0xFF},
	{R: 0x96, G: 0x4E, B: 0x00, A: 0xFF},
	{R: 0x6A, G: 0x6C, B: 0x00, A: 0xFF},
	{R: 0x39, G: 0x84, B: 0x00, A: 0xFF},
	{R: 0x0F, G: 0x90, B: 0x00, A: 0xFF},
	{R: 0x00, G: 0x8D, B: 0x33, A: 0xFF},
	{R: 0x00, G: 0x7B, B: 0x8C, A: 0xFF},
	{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
	{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
	{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},

	{R: 0xFE, G: 0xFE, B: 0xFE, A: 0xFF},
	{R: 0x66, G: 0xAF, B: 0xFF, A: 0xFF},
	{R: 0x93, G: 0x90, B: 0xFF, A: 0xFF},
	{R: 0xC5, G: 0x78, B: 0xFF, A: 0xFF},
	{R: 0xEE, G: 0x6C, B: 0xFF, A: 0xFF},
	{R: 0xFF, G: 0x6F, B: 0xCA, A: 0xFF},
	{R: 0xFF, G: 0x82, B: 0x71, A: 0xFF},
	{R: 0xE6, G: 0x9E, B: 0x25, A: 0xFF},
	{R: 0xBA, G: 0xBC, B: 0x00, A: 0xFF},
	{R: 0x88, G: 0xD5, B: 0x01, A: 0xFF},
	{R: 0x5E, G: 0xE1, B: 0x32, A: 0xFF},
	{R: 0x47, G: 0xDD, B: 0x82, A: 0xFF},
	{R: 0x4A, G: 0xCB, B: 0xDC, A: 0xFF},
	{R: 0x4E, G: 0x4E, B: 0x4E, A: 0xFF},
	{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
	{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},

	{R: 0xFE, G: 0xFE, B: 0xFE, A: 0xFF},
	{R: 0xC0, G: 0xDE, B: 0xFF, A: 0xFF},
	{R: 0xD2, G: 0xD1, B: 0xFF, A: 0xFF},
	{R: 0xE7, G: 0xC7, B: 0xFF, A: 0xFF},
	{R: 0xF8, G: 0xC2, B: 0xFF, A: 0xFF},
	{R: 0xFF, G: 0xC3, B: 0xE9, A: 0xFF},
	{R: 0xFF, G: 0xCB, B: 0xC4, A: 0xFF},
	{R: 0xF5, G: 0xD7, B: 0xA5, A: 0xFF},
	{R: 0xE2, G: 0xE3, B: 0x94, A: 0xFF},
	{R: 0xCE, G: 0xED, B: 0x96, A: 0xFF},
	{R: 0xBC, G: 0xF2, B: 0xAA, A: 0xFF},
	{R: 0xB3, G: 0xF1, B: 0xCB, A: 0xFF},
	{R: 0xB4, G: 0xE9, B: 0xF0, A: 0xFF},
	{R: 0xB6, G: 0xB6, B: 0xB6, A: 0xFF},
	{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
	{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
}

func ResetPPU() {
	WriteLatch = false
	TransferAddress, VRAMAddress, PPUReadBuffer = 0, 0, 0
	NMILevelDetector, DoNMI = false, false

	// $2000: PPUCTRL
	PPUCTRL_NametableSelect = 0
	PPUCTRL_VRAMInc32Mode = false
	PPUCTRL_SpritePatternTable = false
	PPUCTRL_BGPatternTable = false
	PPUCTRL_Use8x16Sprites = false
	PPUCTRL_EnableNMI = false

	// $2001: PPUMASK
	PPUMASK_Greyscale = false
	PPUMASK_8pxMaskBG = false
	PPUMASK_8pxMaskSprites = false
	PPUMASK_RenderBG = false
	PPUMASK_RenderSprites = false
	PPUMASK_EmphasisRed = false
	PPUMASK_EmphasisGreen = false
	PPUMASK_EmphasisBlue = false

	// $2002: PPUSTATUS
	PPUSTATUS_Overflow = false
	PPUSTATUS_SpriteZeroHit = false
	PPUSTATUS_VBlank = false

	PPUDot, PPUScanline = 0, 0

	ppuShiftRegister_patternL, ppuShiftRegister_patternH, ppuShiftRegister_attributeL, ppuShiftRegister_attributeH = 0, 0, 0, 0

	ppu8Step_patternLowBitPlane = 0
	ppu8Step_patternHighBitPlane = 0
	ppu8Step_attribute = 0
	ppu8Step_NextCharacter = 0
	ppu8Step_temp = 0

	PPUScrollFineX, ppuScrollFineY = 0, 0
	DrawNewFrame = false

	OAM = [0x100]byte{}
	SecondaryOAM = [0x20]byte{}
	ppuSpriteEvalTemp = 0
	ppuOAMAddress, ppuSecondaryOAMAddress, ppuSecondaryOAMSize = 0, 0, 0
	ppuSecondaryOAMFull, ppuScanlineContainsSpriteZero, ppuSpriteEvaluationOAMOverflowed = false, false, false
	ppuSpriteEvalTick = 0

	ppu_SpriteShiftRegisterL = [8]byte{}
	ppu_SpriteShiftRegisterH = [8]byte{}
	ppu_SpriteAttribute = [8]byte{}
	ppu_SpritePattern = [8]byte{}
	ppu_SpriteXposition = [8]byte{}
	ppu_SpriteYposition = [8]byte{}
}

func Emulate_PPU( /*g *Game*/ ) {

	if PPUDot == 1 && PPUScanline == 241 {
		PPUSTATUS_VBlank = true
		DrawNewFrame = true
	} else if PPUDot == 1 && PPUScanline == 261 {
		PPUSTATUS_VBlank = false
		PPUSTATUS_Overflow = false
		PPUSTATUS_SpriteZeroHit = false
	}

	SpriteEvaluation()

	if !PPUMASK_RenderBG && !PPUMASK_RenderSprites {
		PPUAddressBus = VRAMAddress // the address bus is always v when rendering is disabled.
	}

	PPURender()

	//Increment / Reset scroll
	if (PPUScanline < 240 || PPUScanline == 261) && (PPUMASK_RenderBG || PPUMASK_RenderSprites) {
		//If this is a visible scanline and rendering sprites / background is enabled
		if PPUDot == 256 { //The Y Scroll is incremented on dot 256
			PPU_IncrementScrollY()
		} else if PPUDot == 257 { //The X scroll is reset on dot 257
			PPU_ResetXScroll()
		}
		if PPUDot >= 280 && PPUDot <= 304 && PPUScanline == 261 { //numbers from the nesdev wiki
			PPU_ResetYScroll() //The Y scroll is reset on every dot from 280 through 304 on the pre-render scanline
		}
	}

	//Drawing
	DrawScreen()

	PPUDot++
	if PPUDot > 341 {
		PPUDot = 0
		PPUScanline++
		if PPUScanline > 261 {
			PPUScanline = 0
		}
	}

	if PPUBus != 0 {
		DecayPPUDataBus()
	}
}

func PPURender() {
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

func DrawScreen( /*g *Game*/ ) {

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

func PPU8Steps() {
	//What part of the 8-step process to run this cycle
	cycleTick := byte((PPUDot - 1) & 7)
	switch cycleTick {
	case 0:
		ppuShiftRegister_patternL = ((ppuShiftRegister_patternL & 0xFF00) | uint16(ppu8Step_patternLowBitPlane))
		ppuShiftRegister_patternH = ((ppuShiftRegister_patternH & 0xFF00) | uint16(ppu8Step_patternHighBitPlane))
		ppuShiftRegister_attributeL = ((ppuShiftRegister_attributeL & 0xFF00) | common.Ternary((ppu8Step_attribute&1) == 1, 0xFF, 0x00))
		ppuShiftRegister_attributeH = ((ppuShiftRegister_attributeH & 0xFF00) | common.Ternary((ppu8Step_attribute&2) == 2, 0xFF, 0x00))
		PPUAddressBus = (0x2000 + (VRAMAddress & 0x0FFF))
		ppu8Step_temp = ReadPPU()
	case 1:
		ppu8Step_NextCharacter = ppu8Step_temp
	case 2:
		PPUAddressBus = (0x23C0 | (VRAMAddress & 0x0C00) | ((VRAMAddress >> 4) & 0x38) | ((VRAMAddress >> 2) & 0x07))
		ppu8Step_temp = ReadPPU()
	case 3:
		ppu8Step_attribute = ppu8Step_temp
		//1 byte of attribute data covers 4 tiles. determine which tile this is for
		if (VRAMAddress & 3) >= 2 { //If this is in the right tile
			ppu8Step_attribute = byte(ppu8Step_attribute >> 2)
		}
		if (((VRAMAddress & 0b0000001111100000) >> 5) & 3) >= 2 { //If this is in the bottom tile
			ppu8Step_attribute = byte(ppu8Step_attribute >> 4)
		}
		ppu8Step_attribute = byte(ppu8Step_attribute & 3)
	case 4:
		PPUAddressBus = (((VRAMAddress & 0b0111000000000000) >> 12) | (uint16(ppu8Step_NextCharacter) * 16) | common.Ternary(PPUCTRL_BGPatternTable, 0x1000, 0))
		ppu8Step_temp = ReadPPU()
	case 5:
		ppu8Step_patternLowBitPlane = ppu8Step_temp
		PPUAddressBus += 8
	case 6:
		ppu8Step_temp = ReadPPU()
	case 7:
		ppu8Step_patternHighBitPlane = ppu8Step_temp
		//Increment VRAM with scrolling
		if (VRAMAddress & 0x001F) == 31 {
			VRAMAddress &= 0xFFE0 //Reset the scroll
			VRAMAddress ^= 0x0400 //Crossing into next nametable
		} else {
			VRAMAddress++
		}
	}
}

func PPU_IncrementScrollY() {
	if (VRAMAddress & 0x7000) != 0x7000 {
		VRAMAddress += 0x1000
	} else {
		VRAMAddress &= 0x0FFF
		y := uint16((VRAMAddress & 0x03E0) >> 5)
		if y == 29 {
			y = 0 // Reset the Y value and also flip some other bit in the 'v' register
			VRAMAddress ^= 0x0800
		} else {
			y++ //Increment the Y value
			y &= 0x1F
		}
		VRAMAddress = (uint16(VRAMAddress&0xFC1F) | uint16(y)<<5)
	}
}

func PPU_ResetXScroll() {
	VRAMAddress = ((VRAMAddress & 0b0111101111100000) | (TransferAddress & 0b0000010000011111))
}

func PPU_ResetYScroll() {
	VRAMAddress = ((VRAMAddress & 0b0000010000011111) | (TransferAddress & 0b0111101111100000))
}

func SpriteEvaluation() {
	if PPUDot == 0 { //Step 0: Reset Secondary OAM count
		ppuSecondaryOAMAddress = 0
		ppuSecondaryOAMFull = false
	} else if PPUDot > 0 && PPUDot <= 64 { //Step 1: Clear Secondary OAM
		if (PPUDot & 1) == 1 {
			//Odd PPU cycles load the value $FF
			ppuSpriteEvalTemp = 0xFF
		} else {
			//Even PPU cycles store the value in secondaryOAM
			SecondaryOAM[ppuSecondaryOAMAddress] = ppuSpriteEvalTemp
			ppuSecondaryOAMAddress++
			ppuSecondaryOAMAddress &= 0x1F //Keep this limited from $00 to $1F
		}
	} else if PPUDot > 64 && PPUDot <= 256 { //Step 2: Load OAM into Secondary OAM (If not full)
		if (PPUDot & 1) == 1 {
			//Odd PPU cycles load the value from OAM
			ppuSpriteEvalTemp = OAM[ppuOAMAddress&0xFF] //CHECK IF THIS IS LEGIT
		} else {
			if !ppuSpriteEvaluationOAMOverflowed {
				//Even PPU cycles store the value in secondaryOAM
				if !ppuSecondaryOAMFull { //If SecondaryOAM is not full yet
					// As long as secondaryOAM isn't full, this write *always* occurs, regardless of evaluation
					SecondaryOAM[ppuSecondaryOAMAddress] = ppuSpriteEvalTemp
				}
				if ppuSpriteEvalTick == 0 {
					//Reading index 0 of an object's set of 4 bytes
					if (PPUScanline-int(ppuSpriteEvalTemp) >= 0) && (PPUScanline-int(ppuSpriteEvalTemp) < int(common.Ternary(PPUCTRL_Use8x16Sprites, 16, 8))) {
						//This object *is* on this scanline!
						if !ppuSecondaryOAMFull {
							ppuSecondaryOAMAddress++ //Increment this for the next write to Secondary OAM
							ppuOAMAddress++          //Increment this for the next ream of Object Attribute Memory
							if PPUDot == 66 {
								// Rather than verifying that this is OAM index 0,
								// the PPU sets this flag if we found an object on this scanline
								// during PPUDot 66, which would be the PPU cycle evaluating index 0
								ppuScanlineContainsSpriteZero = true
							}
						} else {
							PPUSTATUS_Overflow = true
						}
						ppuSpriteEvalTick++
					} else {
						ppuOAMAddress += 4
					}
				} else { //If ppuSpriteEvalTick != 0
					// Reading index 1, 2, or 3 of an object's OAM data.
					// We're not going to be making any checks for if things are on this scanline,
					// so we can just simply increment the OAM address
					ppuSecondaryOAMAddress++ //Increment this for the next write to Secondary OAM
					ppuOAMAddress++          //Increment this for the next ream of Object Attribute Memory
					if ppuSecondaryOAMAddress == 0x20 {
						ppuSecondaryOAMFull = true
					}
					ppuSpriteEvalTick++
					ppuSpriteEvalTick &= 3 // Wrap around to tick 0 after tick 3
				}
				if ppuOAMAddress == 0 {
					// If we overflow the OAM address, we want to stop running the sprite evaluation checks until dot 257
					ppuSpriteEvaluationOAMOverflowed = true
				}
			}
		}
	} else if PPUDot > 256 && PPUDot <= 320 { //Step 3:
		ppuOAMAddress = 0 //This is set to $00 during every one of these cycles
		if PPUDot == 257 {
			ppuSecondaryOAMSize = ppuSecondaryOAMAddress
			ppuSecondaryOAMAddress = 0
			ppuSpriteEvalTick = 0
		}

		switch ppuSpriteEvalTick {
		case 0:
			//Set this object's Y position in the array
			ppu_SpriteYposition[ppuSecondaryOAMAddress/4] = SecondaryOAM[ppuSecondaryOAMAddress]
			ppuSecondaryOAMAddress++
		case 1:
			//Set this object's Y position in the array
			ppu_SpritePattern[ppuSecondaryOAMAddress/4] = SecondaryOAM[ppuSecondaryOAMAddress]
			ppuSecondaryOAMAddress++
		case 2:
			//Set this object's Y position in the array
			ppu_SpriteAttribute[ppuSecondaryOAMAddress/4] = SecondaryOAM[ppuSecondaryOAMAddress]
			ppuSecondaryOAMAddress++
		case 3:
			//Set this object's Y position in the array
			ppu_SpriteXposition[ppuSecondaryOAMAddress/4] = SecondaryOAM[ppuSecondaryOAMAddress]
		case 4:
			PPUAddressBus = ppuFindSpritePatternData(ppuSecondaryOAMAddress / 4)
		case 5:
			ppuSpriteEvalTemp = ReadPPU()
			if PPUScanline == 261 {
				ppuSpriteEvalTemp = 0 //Clear this if this is the pre-render line
			}
			if ((ppu_SpriteAttribute[ppuSecondaryOAMAddress/4] >> 6) & 1) == 1 { //Attributes are set up to flip X
				// Real nice way to change the order of bits from 76543210 to 01234567
				ppuSpriteEvalTemp = byte(((ppuSpriteEvalTemp & 0xF0) >> 4) | ((ppuSpriteEvalTemp & 0xF) << 4))
				ppuSpriteEvalTemp = byte(((ppuSpriteEvalTemp & 0xCC) >> 2) | ((ppuSpriteEvalTemp & 0x33) << 2))
				ppuSpriteEvalTemp = byte(((ppuSpriteEvalTemp & 0xAA) >> 1) | ((ppuSpriteEvalTemp & 0x55) << 1))
			}
			ppu_SpriteShiftRegisterL[ppuSecondaryOAMAddress/4] = ppuSpriteEvalTemp
		case 6:
			PPUAddressBus += 8
		case 7:
			ppuSpriteEvalTemp = ReadPPU()
			if PPUScanline == 261 {
				ppuSpriteEvalTemp = 0 //Clear this if this is the pre-render line
			}
			if ((ppu_SpriteAttribute[ppuSecondaryOAMAddress/4] >> 6) & 1) == 1 { //Attributes are set up to flip X
				// Real nice way to change the order of bits from 76543210 to 01234567
				ppuSpriteEvalTemp = byte(((ppuSpriteEvalTemp & 0xF0) >> 4) | ((ppuSpriteEvalTemp & 0xF) << 4))
				ppuSpriteEvalTemp = byte(((ppuSpriteEvalTemp & 0xCC) >> 2) | ((ppuSpriteEvalTemp & 0x33) << 2))
				ppuSpriteEvalTemp = byte(((ppuSpriteEvalTemp & 0xAA) >> 1) | ((ppuSpriteEvalTemp & 0x55) << 1))
			}
			ppu_SpriteShiftRegisterH[ppuSecondaryOAMAddress/4] = ppuSpriteEvalTemp
			ppuSecondaryOAMAddress++
		}
		ppuSpriteEvalTick++
		ppuSpriteEvalTick &= 7 // And reset at 8
	}
}

func ppuFindSpritePatternData(SecondaryOAMSlot uint16) uint16 {
	if !PPUCTRL_Use8x16Sprites { //8x8 sprites
		// The address is $0000 or $1000 depending on the pattern table
		// plus the pattern value from OAM * 16
		// plus the number of scanlines from the top of the object
		// if the attributes are set to flip Y, it's 7 - the number of scanlines from the top of the object
		if ((ppu_SpriteAttribute[SecondaryOAMSlot] >> 7) & 1) == 0 { //Attributes are not set up to flip Y
			return uint16(common.Ternary(PPUCTRL_SpritePatternTable, 0x1000, 0) + (uint16(ppu_SpritePattern[SecondaryOAMSlot]) << 4) + uint16(PPUScanline-int(ppu_SpriteYposition[SecondaryOAMSlot])))
		} else { //Attributes are set up to flip Y
			return uint16(common.Ternary(PPUCTRL_SpritePatternTable, 0x1000, 0) + (uint16(ppu_SpritePattern[SecondaryOAMSlot]) << 4) + uint16((7-(PPUScanline-int(ppu_SpriteYposition[SecondaryOAMSlot])))&7))
		}
	} else { //8x16 sprites
		// in 8x16 mode, instead of using ppu_SpritePattern to deternime which pattern table to fetch data from...
		// these sprites instead use bit 0 of the object's pattern information from OAM

		// The address is $0000 or $1000 depending on the pattern table
		// plus (the pattern value from OAM, clearing bit 0) * 16
		// plus the number of scanlines from the top of the object
		// if the attributes are set to flip Y, it's 7 - the number of scanlines from the top of the object

		//If we're drawing the bottom half of the sprite, add 16
		if ((ppu_SpriteAttribute[SecondaryOAMSlot] >> 7) & 1) == 0 { //Attributes are not set up to flip Y
			if PPUScanline-int(ppu_SpriteYposition[SecondaryOAMSlot]) < 8 {
				return uint16(common.Ternary((ppu_SpritePattern[SecondaryOAMSlot]&1) == 1, 0x1000, 0) | (uint16(ppu_SpritePattern[SecondaryOAMSlot]&0xFE) << 4) + uint16(PPUScanline-int(ppu_SpriteYposition[SecondaryOAMSlot])))
			} else {
				return uint16(common.Ternary((ppu_SpritePattern[SecondaryOAMSlot]&1) == 1, 0x1000, 0) + ((uint16(ppu_SpritePattern[SecondaryOAMSlot]&0xFE) << 4) + 16) + uint16((PPUScanline-int(ppu_SpriteYposition[SecondaryOAMSlot]))&7))
			}
		} else { //Attributes are set up to flip Y
			if PPUScanline-int(ppu_SpriteYposition[SecondaryOAMSlot]) < 8 {
				return uint16(common.Ternary((ppu_SpritePattern[SecondaryOAMSlot]&1) == 1, 0x1000, 0) + ((uint16(ppu_SpritePattern[SecondaryOAMSlot]&0xFE) << 4) + 16) + uint16(((PPUScanline-int(ppu_SpriteYposition[SecondaryOAMSlot]))&7)+7))
			} else {
				return uint16(common.Ternary((ppu_SpritePattern[SecondaryOAMSlot]&1) == 1, 0x1000, 0) + ((uint16(ppu_SpritePattern[SecondaryOAMSlot]&0xFE) << 4) + 7) + uint16((PPUScanline-int(ppu_SpriteYposition[SecondaryOAMSlot]))&7))
			}
		}
	}
}

func ReadPPU( /*Address uint16*/ ) byte {
	if PPUAddressBus < 0x2000 {
		//Read from pattern table.
		switch cartridge.MapperChipID {
		case 1: //MMC1
			return cartridge.CHRROM[mappers.MMC1_FetchPPUAddress(PPUAddressBus, cartridge.CHRROM_Size)]
		case 3: //CNROM
			return cartridge.CHRROM[mappers.CNROM_ReadAddress(PPUAddressBus)]
		//case 4: //MMC3
		//	return 0
		default:
			return cartridge.CHRROM[PPUAddressBus]
		}

		//else, nothing happens
	} else if PPUAddressBus < 0x3F00 {
		//Read from the Nametables
		switch cartridge.MapperChipID {
		case 1: //MMC1
			return cartridge.CartVRAM[mappers.MMC1_FetchNametable(PPUAddressBus)]
		//case 3: //CNROM
		//case 4: //MMC3
		default:
			if !cartridge.AltNametableLayout {
				if cartridge.IsNametableHorizontal {
					// Horizontal Mirroring
					return cartridge.VRAM[int(PPUAddressBus&0x3FF)|(int(PPUAddressBus&0x800)>>1)]
				} else {
					//Vertical Mirroring
					return cartridge.VRAM[int(PPUAddressBus&0x7FF)]
				}
			}
		}
		return 0

	} else {
		//Read from Palette RAM
		if (PPUAddressBus & 3) == 0 {
			return cartridge.PaletteRAM[PPUAddressBus&0x0F]
		} else {
			return cartridge.PaletteRAM[PPUAddressBus&0x1F]
		}
	}
}

var PPUBusDecay [8]int
var PPUBusDecayConstant int = 1786830

func UpdatePPUBus(Value byte) byte {
	PPUBus = Value
	//PPU decay buffer code here
	for i := 0; i < 8; i++ {
		PPUBusDecay[i] = PPUBusDecayConstant
	}
	//???
	return PPUBus
}

func UpdatePPUBus2002(Value byte) {
	PPUBus = Value
	//Only update the decay constant on the top 3 bits
	for i := 5; i < 8; i++ {
		PPUBusDecay[i] = PPUBusDecayConstant
	}
}

func UpdatePPUBus2007Palette(Value byte) {
	PPUBus = Value
	//Only update the decay constant on the bottom 6 bits
	for i := 0; i < 6; i++ {
		PPUBusDecay[i] = PPUBusDecayConstant
	}
}

func DecayPPUDataBus() {
	DecayBitmask := [8]byte{0xFE, 0xFD, 0xFB, 0xF7, 0xEF, 0xDF, 0xBF, 0x7F}

	for i := range PPUBusDecay {
		if PPUBusDecay[i] > 0 {
			PPUBusDecay[i]--
			if PPUBusDecay[i] == 0 {
				PPUBus &= DecayBitmask[i]
			}
		}
	}
}
