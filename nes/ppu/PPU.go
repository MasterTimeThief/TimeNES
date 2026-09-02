package ppu

import (
	"image/color"
	"mtt/timenes/common"
	"mtt/timenes/nes/cartridge"
	"mtt/timenes/nes/cartridge/mappers"
)

type PPU struct {
}

var WriteLatch bool        //PPU's w register
var TransferAddress uint16 //PPU's t register
var VRAMAddress uint16     //PPU's v register
var PPUReadBuffer byte

// $2000: PPUCTRL
var (
	PPUCTRL_NametableSelect    byte // PPUCTRL Bit 1 & 2
	PPUCTRL_VRAMInc32Mode      bool // PPUCTRL Bit 3
	PPUCTRL_SpritePatternTable bool // PPUCTRL Bit 4
	PPUCTRL_BGPatternTable     bool // PPUCTRL Bit 5
	PPUCTRL_Use8x16Sprites     bool // PPUCTRL Bit 6
	PPUCTRL_EnableNMI          bool // PPUCTRL Bit 8
)

// $2001: PPUMASK
var (
	PPUMASK_Greyscale      bool // PPUMASK Bit 0
	PPUMASK_8pxMaskBG      bool // PPUMASK Bit 1
	PPUMASK_8pxMaskSprites bool // PPUMASK Bit 2
	PPUMASK_RenderBG       bool // PPUMASK Bit 3
	PPUMASK_RenderSprites  bool // PPUMASK Bit 4
	PPUMASK_EmphasisRed    bool // PPUMASK Bit 5
	PPUMASK_EmphasisGreen  bool // PPUMASK Bit 6
	PPUMASK_EmphasisBlue   bool // PPUMASK Bit 7
)

// $2002: PPUSTATUS
var (
	PPUSTATUS_Overflow      bool // PPUSTATUS Bit 5
	PPUSTATUS_SpriteZeroHit bool // PPUSTATUS Bit 6
	PPUSTATUS_VBlank        bool // PPUSTATUS Bit 7
)

var PPUDot int      //The X position of the scanning beam
var PPUScanline int //The Y position of the scanning beam
var ppuShiftRegister_patternL, ppuShiftRegister_patternH, ppuShiftRegister_attributeL, ppuShiftRegister_attributeH uint16
var ppu8Step_patternLowBitPlane, ppu8Step_patternHighBitPlane, ppu8Step_attribute, ppu8Step_NextCharacter, ppu8Step_temp byte
var PPUScrollFineX byte
var DrawNewFrame bool = false
var FrameColorBuffer [61440]color.RGBA
var FrameColorBufferPos int = 0

var OAM [0x100]byte
var SecondaryOAM [0x20]byte
var ppuSpriteEvalTemp byte
var ppuOAMAddress byte
var OAMBusAddress byte
var ppuSecondaryOAMAddress, ppuSecondaryOAMSize uint16
var ppuSecondaryOAMFull, ppuScanlineContainsSpriteZero, ppuSpriteEvaluationOAMOverflowed bool
var ppuSpriteEvalTick int

var ppu_SpriteShiftRegisterL [8]byte
var ppu_SpriteShiftRegisterH [8]byte

var ppu_SpriteAttribute [8]byte
var ppu_SpritePattern [8]byte
var ppu_SpriteXposition [8]byte
var ppu_SpriteYposition [8]byte

var PPUAddressBus uint16
var PPUBus byte

var PPUOddFrame bool
var SuppressVBlank bool
var SuppressNMI bool

func NewPPU() *PPU {
	ppu := PPU{}
	return &ppu
}

func ResetPPU() {
	WriteLatch = false
	TransferAddress, VRAMAddress, PPUReadBuffer = 0, 0, 0

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

	PPUScrollFineX = 0
	DrawNewFrame = false
	FrameColorBufferPos = 0

	OAM = [0x100]byte{}
	SecondaryOAM = [0x20]byte{}
	ppuSpriteEvalTemp = 0
	ppuOAMAddress, OAMBusAddress, ppuSecondaryOAMAddress, ppuSecondaryOAMSize = 0, 0, 0, 0
	ppuSecondaryOAMFull, ppuScanlineContainsSpriteZero, ppuSpriteEvaluationOAMOverflowed = false, false, false
	ppuSpriteEvalTick = 0

	ppu_SpriteShiftRegisterL = [8]byte{}
	ppu_SpriteShiftRegisterH = [8]byte{}
	ppu_SpriteAttribute = [8]byte{}
	ppu_SpritePattern = [8]byte{}
	ppu_SpriteXposition = [8]byte{}
	ppu_SpriteYposition = [8]byte{}

	PPUAddressBus, PPUBus = 0, 0

	PPUOddFrame = true
}

var CPUCyclesLastFrame int

func PPU_Cycle() {

	if PPUDot == 1 && PPUScanline == 241 {
		if !SuppressVBlank {
			PPUSTATUS_VBlank = true
		}
		DrawNewFrame = true
		//fmt.Println("Frame Cycles: " + fmt.Sprintf("%d", (common.CPU_TotalCycles-CPUCyclesLastFrame)))
		CPUCyclesLastFrame = common.CPU_TotalCycles
		PPUOddFrame = !PPUOddFrame
	} else if PPUDot == 1 && PPUScanline == 261 {
		PPUSTATUS_VBlank = false
		PPUSTATUS_Overflow = false
		PPUSTATUS_SpriteZeroHit = false
		SuppressVBlank = false
		SuppressNMI = false
	}

	SpriteEvaluation()

	if !PPUMASK_RenderBG && !PPUMASK_RenderSprites {
		PPUAddressBus = VRAMAddress // the address bus is always v when rendering is disabled.
	}

	if cartridge.MapperChipID == 4 {
		mappers.MMC3_ClockIRQ(PPUAddressBus)
	}

	//Get tile data for buffering
	RenderTileData()

	//Increment / Reset scroll
	UpdateScroll()

	//Drawing the next pixel
	RenderNextPixel()

	PPUDot++
	if PPUDot > 340 {
		PPUDot = 0
		PPUScanline++
		if PPUScanline > 261 {
			PPUScanline = 0
		}
	}

	// If rendering is enabled, skip the first dot of every odd frame
	if PPUDot == 340 && PPUScanline == 261 && PPUOddFrame && (PPUMASK_RenderBG || PPUMASK_RenderSprites) {
		PPUDot = 0
		PPUScanline = 0
	}

	if PPUBus != 0 {
		DecayPPUDataBus()
	}
}

func PPUNextTile8Steps() {
	//What part of the 8-step process to run this cycle
	cycleTick := byte((PPUDot - 1) & 7)
	switch cycleTick {
	case 0: // Load the shift registers, and get the address for the nametable byte
		ppuShiftRegister_patternL = ((ppuShiftRegister_patternL & 0xFF00) | uint16(ppu8Step_patternLowBitPlane))
		ppuShiftRegister_patternH = ((ppuShiftRegister_patternH & 0xFF00) | uint16(ppu8Step_patternHighBitPlane))
		ppuShiftRegister_attributeL = ((ppuShiftRegister_attributeL & 0xFF00) | common.Ternary((ppu8Step_attribute&1) == 1, 0xFF, 0x00))
		ppuShiftRegister_attributeH = ((ppuShiftRegister_attributeH & 0xFF00) | common.Ternary((ppu8Step_attribute&2) == 2, 0xFF, 0x00))
		PPUAddressBus = (0x2000 + (VRAMAddress & 0x0FFF))
		ppu8Step_temp = ReadPPU()
	case 1: // Save the nametable byte
		ppu8Step_NextCharacter = ppu8Step_temp
	case 2: // Get the address for the attribute byte
		PPUAddressBus = (0x23C0 | (VRAMAddress & 0x0C00) | ((VRAMAddress >> 4) & 0x38) | ((VRAMAddress >> 2) & 0x07))
		ppu8Step_temp = ReadPPU()
	case 3: // Save the attribute byte
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

func UpdateScroll() {
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
		ppuSpriteEvaluationOAMOverflowed = false
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
			ppuSpriteEvalTemp = OAM[ppuOAMAddress]
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
							if (PPUMASK_RenderBG || PPUMASK_RenderSprites) && ppuSpriteEvalTick%4 == 0 && (ppuSpriteEvalTemp < 240) {
								PPUSTATUS_Overflow = true
							}
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
			//Set this object's pattern in the array
			ppu_SpritePattern[ppuSecondaryOAMAddress/4] = SecondaryOAM[ppuSecondaryOAMAddress]
			ppuSecondaryOAMAddress++
		case 2:
			//Set this object's attributes in the array
			ppu_SpriteAttribute[ppuSecondaryOAMAddress/4] = SecondaryOAM[ppuSecondaryOAMAddress]
			ppuSecondaryOAMAddress++
		case 3:
			//Set this object's X position in the array
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
			return cartridge.CHRROM[mappers.CNROM_FetchPPUAddress(PPUAddressBus)]
		case 4: //MMC3
			return cartridge.CHRROM[mappers.MMC3_FetchPPUAddress(PPUAddressBus, cartridge.CHRROM_Size)]
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
		case 4: //MMC3
			return ReadFromNametable(PPUAddressBus, mappers.MMC3_IsHorizontalNametable)
			//return cartridge.CartVRAM[int(PPUAddressBus&0xFFF)]
		case 7: // AxROM
			return cartridge.VRAM[mappers.AxROM_FetchNametable(PPUAddressBus)]
		default:
			return ReadFromNametable(PPUAddressBus, cartridge.IsNametableHorizontal)
		}
	} else {
		//Read from Palette RAM
		if (PPUAddressBus & 3) == 0 {
			return cartridge.PaletteRAM[PPUAddressBus&0x0F]
		} else {
			return cartridge.PaletteRAM[PPUAddressBus&0x1F]
		}
	}
}

func WritePPU(Value byte) {
	if VRAMAddress < 0x2000 {
		//Write to pattern table. (If the cartridge supports it)
		if cartridge.CHRROM_Size == 0 {
			switch cartridge.MapperChipID {
			case 1: //MMC1
				//if AltNametableLayout {
				//	mappers.MMC1_VRAM[mappers.MMC1_FetchPPUAddress(VRAMAddress, CHRROM_Size)] = Value
				//} else {
				cartridge.CHRROM[mappers.MMC1_FetchPPUAddress(VRAMAddress, cartridge.CHRROM_Size)] = Value
				//}
			case 3: //CNROM
				cartridge.CHRROM[mappers.CNROM_FetchPPUAddress(VRAMAddress)] = Value
			case 4: //MMC3
				cartridge.CHRROM[mappers.MMC3_FetchPPUAddress(VRAMAddress, cartridge.CHRROM_Size)] = Value
			default:
				cartridge.CHRROM[VRAMAddress] = Value
			}
			//For Pattern Table viewing
			common.PendingPatternUpdate = true
		}
		//else, nothing happens because it's CHR-ROM
	} else if VRAMAddress < 0x3F00 {
		//Write to the Nametables

		switch cartridge.MapperChipID {
		case 1: //MMC1
			switch mappers.MMC1_GetNametableArrangement() {
			case 0: //Screen A only
				cartridge.CartVRAM[int(VRAMAddress&0x3FF)] = Value
			case 1: //Screen B Only
				cartridge.CartVRAM[int(VRAMAddress&0x3FF)+0x400] = Value
			case 2: //Horizontal
				cartridge.CartVRAM[int(VRAMAddress&0x7FF)] = Value
			case 3: //Vertical
				cartridge.CartVRAM[int(VRAMAddress&0x3FF)|int(VRAMAddress&0x800)>>1] = Value
			}
			//cartridge.CartVRAM[VRAMAddress&0xFFF] = Value
		//case 3: //CNROM
		case 4: //MMC3
			WriteToNametable(Value, mappers.MMC3_IsHorizontalNametable)
			//cartridge.CartVRAM[int(VRAMAddress&0xFFF)] = Value
		case 7:
			cartridge.VRAM[mappers.AxROM_FetchNametable(PPUAddressBus)] = Value
		default:
			WriteToNametable(Value, cartridge.IsNametableHorizontal)
		}

	} else {
		//Write to Palette RAM
		if (VRAMAddress & 3) == 0 {
			cartridge.PaletteRAM[VRAMAddress&0x0F] = Value
		} else {
			cartridge.PaletteRAM[VRAMAddress&0x1F] = Value
		}
	}
}

func ReadFromNametable(Addr uint16, isHoriz bool) byte {
	if cartridge.AltNametableLayout {
	} else {
		if cartridge.IsNametableHorizontal {
			// Horizontal Mirroring
			return cartridge.VRAM[int(Addr&0x3FF)|(int(Addr&0x800)>>1)]
		} else {
			//Vertical Mirroring
			return cartridge.VRAM[int(Addr&0x7FF)]
		}
	}
	return 0
}

func WriteToNametable(Value byte, isHoriz bool) {
	if cartridge.AltNametableLayout {
	} else {
		if cartridge.IsNametableHorizontal {
			// Horizontal Mirroring
			cartridge.VRAM[int(VRAMAddress&0x3FF)|int(VRAMAddress&0x800)>>1] = Value
		} else {
			//Vertical Mirroring
			cartridge.VRAM[int(VRAMAddress&0x7FF)] = Value
		}
	}
}

var PPUBusDecay [8]int
var PPUBusDecayConstant int = (29781 * 3 * 60)

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
	PPUBus = (Value & 0xE0) | (PPUBus & 0x1F)
	//Only update the decay constant on the top 3 bits
	for i := 5; i < 8; i++ {
		PPUBusDecay[i] = PPUBusDecayConstant
	}
}

func UpdatePPUBus2007Palette(Value byte) {
	PPUBus = (Value & 0x3F) | (PPUBus & 0xC0)
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

func GetOAMBusAddress() byte {
	return OAMBusAddress
}

func SetOAMBusAddress(Value byte) {
	OAMBusAddress = Value
}
