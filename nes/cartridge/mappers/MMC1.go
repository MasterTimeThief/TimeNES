// Mapper 001: MMC1
package mappers

var MMC1_PRGRAM [0x2000]byte

var MMC1_ShiftRegister byte = 0x10
var MMC1_Control byte = 0x0C
var MMC1_CHRBank0 byte
var MMC1_CHRBank1 byte
var MMC1_PRGBank byte

var lastWriteCycle int

func MMC1_WriteToPRGRAM(Value byte, Addr uint16) {
	if (MMC1_PRGBank & 0x10) == 0 {
		MMC1_PRGRAM[Addr&0x1FFF] = Value
	}
}

func MMC1_Write(Value byte, Addr uint16, Cycle int) {

	//if Addr < 0x8000 { //PRG-RAM
	//	if (MMC1_PRGBank & 0x10) == 0 {
	//		MMC1_PRGRAM[Addr&0x1FFF] = Value
	//	}
	//} else {
	if (Value & 0x80) != 0 { //If bit 7 is set, clear the shift register
		MMC1_Reset()
	} else { //Add to the shift register
		if Cycle-lastWriteCycle <= 1 {
			lastWriteCycle = Cycle
			return
		} else {
			shiftCheck := MMC1_ShiftRegister & 1
			MMC1_ShiftRegister = (MMC1_ShiftRegister >> 1) | ((Value & 1) << 4)
			MMC1_ShiftRegister &= 0x1F //Probably redundant, but just in case
			if shiftCheck != 0 {       //This was write 5
				bankSelect := (Addr >> 13) & 3
				switch bankSelect {
				case 0: //$8000-$9FFF
					MMC1_Control = MMC1_ShiftRegister
				case 1: //$A000-$BFFF
					MMC1_CHRBank0 = MMC1_ShiftRegister
				case 2: //$C000-$DFFF
					MMC1_CHRBank1 = MMC1_ShiftRegister
				case 3: //$E000-$FFFF
					MMC1_PRGBank = MMC1_ShiftRegister
				}
				MMC1_ShiftRegister = 0x10
			}
		}
	}
	//}
}

func MMC1_FetchCPUAddress(Addr uint16, PRGLength uint32) uint32 {
	BankMode := (MMC1_Control & 0xC) >> 2
	var tempAddr uint32
	switch BankMode {
	case 0, 1:
		// switch 32 KB at $8000, ignoring low bit of bank number
		tempAddr = uint32(Addr & 0x7FFF)
		return ((0x8000 * uint32(MMC1_PRGBank&0x0E)) + tempAddr) & (PRGLength - 1)
	case 2:
		// fix first bank at $8000 and switch 16 KB bank at $C000
		tempAddr = uint32(Addr & 0x3FFF)
		if Addr >= 0xC000 {
			return (0x4000 * uint32(MMC1_PRGBank)) + tempAddr
		} else {
			return tempAddr /*| (uint32(MMC1_PRGBank & 0x08)<<13)*/
		}
	case 3:
		// fix last bank at $C000 and switch 16 KB bank at $8000
		tempAddr = uint32(Addr & 0x3FFF)
		if Addr >= 0xC000 {
			return PRGLength - 0x4000 + tempAddr
		} else {
			return (0x4000*uint32(MMC1_PRGBank&0x0F) + tempAddr) & (PRGLength - 1)
		}
	}
	return 0
}

func MMC1_FetchPPUAddress(Addr uint16, CHRLength uint32) uint32 {
	// bit 4 of Mapper_1_Control controls how the pattern tables are swapped. if set, 2 banks of 4Kib. Otherwise, 1 8Kib bank
	if (MMC1_Control & 0x10) != 0 {
		// with the MMC1 chip, you can swap out the pattern tables.
		// address < 0x1000 is the first pattern table, else, the second pattern table.
		// if the final write for the MMC1 shift register was in the $A000 - $BFFF, this updates Mapper_1_CHR0
		// if the final write for the MMC1 shift register was in the $B000 - $CFFF, this updates Mapper_1_CHR1

		if Addr < 0x1000 {
			return (uint32(uint32(MMC1_CHRBank0&0x1F)*0x1000) + uint32(Addr)) & uint32(CHRLength-1)
		} else {
			Addr &= 0xFFF
			return (uint32(uint32(MMC1_CHRBank1&0x1F)*0x1000) + uint32(Addr)) & uint32(CHRLength-1)
		}
	} else { // one swappable bank that changes both pattern tables.
		// this uses the value written to Mapper_1_CHR0
		return (uint32(uint32(MMC1_CHRBank0&0xFE)*0x2000) + uint32(Addr)) & uint32(CHRLength-1)
	}
}

func MMC1_FetchNametable(Addr uint16) int {
	Arrangement := MMC1_Control & 0x3
	switch Arrangement {
	case 0: //Screen A only
		return int(Addr & 0x3FF)
	case 1: //Screen B Only
		return int(Addr&0x3FF) + 0x400
	case 2: //Horizontal
		return int(Addr & 0x7FF)
	case 3: //Vertical
		return int(Addr&0x3FF) | int((Addr&0x800)>>1)
	default:
		return int(Addr)
	}
}

func MMC1_Reset() {
	MMC1_ShiftRegister = 0x10
	MMC1_Control |= 0x0C
}

func MMC1_GetNametableArrangement() byte {
	return MMC1_Control & 0x3
}
