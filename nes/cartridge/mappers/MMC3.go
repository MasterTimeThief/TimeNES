// Mapper 004: MMC3
package mappers

// ┌───────────────┐
// │ $8000 - $9FFF ├───────────────────────┬── One is swappable,
// ├───────────────┤                       │   One is fixed to 2nd-To-Last PRG Bank
// │ $A000 - $BFFF ├─ Always Swappable     │   Determined by MMC3_PRGBankMode
// ├───────────────┤                       │
// │ $C000 - $DFFF ├───────────────────────┘
// ├───────────────┤
// │ $E000 - $FFFF ├─ Always fixed to last PRG Bank
// └───────────────┘

var (
	MMC3_PRGOffset [4]uint32
	MMC3_CHROffset [8]uint32
	MMC3_Register  [8]byte
)

var ( // Bank Select ($8000-$9FFE, even)
	MMC3_BankSelect  byte
	MMC3_PRGBankMode bool
	MMC3_CHRBankMode bool
)

// Bank data ($8001-$9FFF, odd)
var MMC3_BankData byte
var MMC3_IsHorizontalNametable bool

var ( // PRG RAM protect ($A001-$BFFF, odd)
	MMC3_WriteProtect bool
	MMC3_PRGRAMEnable bool
)

// IRQ latch
var (
	MMC3_IRQCounter     byte
	MMC3_IRQReloadValue byte
	MMC3_IRQReloadFlag  bool
	MMC3_IRQEnabled     bool
	MMC3_PPUA12         bool
	MMC3_PPUA12Prev     bool
	MMC3_DoIRQ          bool
)

//

var ( //PRG Banks
	MMC3_LastBank      uint32
	MMC3_2ndToLastBank uint32

//	MMC3_BankAddress_8000 uint32
//	MMC3_BankAddress_A000 uint32
//	MMC3_BankAddress_C000 uint32
)

//
//var ( //CHR Banks
//	MMC3_CHRBankAddress0 uint32
//	MMC3_CHRBankAddress1 uint32
//	MMC3_CHRBankAddress2 uint32
//	MMC3_CHRBankAddress3 uint32
//	MMC3_CHRBankAddress4 uint32
//	MMC3_CHRBankAddress5 uint32
//)

var MMC3_PRGRAM [0x2000]byte

func MMC3_InitRegisters(PRGLength uint32) {
	MMC3_LastBank = PRGLength - 0x2000
	MMC3_2ndToLastBank = PRGLength - 0x4000
	MMC3_UpdateRegister()
}

func MMC3_WriteToPRGRAM(Value byte, Addr uint16) {
	if MMC3_PRGRAMEnable && !MMC3_WriteProtect {
		MMC3_PRGRAM[Addr&0x1FFF] = Value
	}
}

func MMC3_Write(Value byte, Addr uint16) {
	if Addr < 0xA000 {
		if (Addr & 1) == 0 { // Even
			// Bank Select ($8000-$9FFE, even)
			MMC3_BankSelect = Value & 0x7
			MMC3_PRGBankMode = (Value & 0x40) != 0
			MMC3_CHRBankMode = (Value & 0x80) != 0
			MMC3_UpdateRegister()

			// Also check for MMC6 behavior here
		} else { //Odd
			// Bank data ($8001-$9FFF, odd)
			//MMC3_BankData = Value
			//MMC3_BankSwap(Value)
			MMC3_Register[MMC3_BankSelect] = Value
			MMC3_UpdateRegister()
		}
	} else if Addr < 0xC000 {
		if (Addr & 1) == 0 { // Even
			// Nametable arrangement ($A000-$BFFE, even)
			MMC3_IsHorizontalNametable = (Value & 1) == 0
		} else { //Odd
			// PRG RAM protect ($A001-$BFFF, odd)
			MMC3_WriteProtect = (Value & 0x40) != 0
			MMC3_PRGRAMEnable = (Value & 0x80) != 0
			// Also check for MMC6 behavior here
		}
	} else if Addr < 0xE000 {
		if (Addr & 1) == 0 { // Even
			// IRQ latch ($C000-$DFFE, even)
			MMC3_IRQReloadValue = Value
		} else { //Odd
			// IRQ reload ($C001-$DFFF, odd)
			MMC3_IRQReloadFlag = true
		}
	} else {
		if (Addr & 1) == 0 { // Even
			// IRQ disable ($E000-$FFFE, even)
			MMC3_IRQEnabled = false
		} else { //Odd
			// IRQ enable ($E001-$FFFF, odd)
			MMC3_IRQEnabled = true
		}
	}
}

/*
func MMC3_BankSwap(NewBank byte) {
	switch MMC3_BankSelect {
	case 0: // Select 2 KB CHR bank at PPU $0000-$07FF (or $1000-$17FF)
		MMC3_CHRBankAddress0 = (0x800 * (uint32(NewBank & 0xFE)))
	case 1: // Select 2 KB CHR bank at PPU $0800-$0FFF (or $1800-$1FFF)
		MMC3_CHRBankAddress1 = (0x800 * (uint32(NewBank & 0xFE)))
	case 2: // Select 1 KB CHR bank at PPU $1000-$13FF (or $0000-$03FF)
		MMC3_CHRBankAddress2 = (0x400 * uint32(NewBank))
	case 3: // Select 1 KB CHR bank at PPU $1400-$17FF (or $0400-$07FF)
		MMC3_CHRBankAddress3 = (0x400 * uint32(NewBank))
	case 4: // Select 1 KB CHR bank at PPU $1800-$1BFF (or $0800-$0BFF)
		MMC3_CHRBankAddress4 = (0x400 * uint32(NewBank))
	case 5: // Select 1 KB CHR bank at PPU $1C00-$1FFF (or $0C00-$0FFF)
		MMC3_CHRBankAddress5 = (0x400 * uint32(NewBank))
	case 6: // Select 8 KB PRG ROM bank at $8000-$9FFF (or $C000-$DFFF)
		if MMC3_PRGBankMode {
			MMC3_BankAddress_C000 = (0x2000 * uint32(NewBank&0x3F))
		} else {
			MMC3_BankAddress_8000 = (0x2000 * uint32(NewBank&0x3F))
		}
	case 7: // Select 8 KB PRG ROM bank at $A000-$BFFF
		MMC3_BankAddress_A000 = (0x2000 * uint32(NewBank&0x3F))
	}
}
*/

func MMC3_FetchCPUAddress(Addr uint16, PRGLength uint32) uint32 {
	//tempAddr := Addr & 0x1FFF
	//if Addr < 0xA000 { //$8000-$9FFF
	//	return (MMC3_BankAddress_8000 + uint32(tempAddr)) //& (PRGLength - 1)
	//} else if Addr < 0xC000 { //$A000-$BFFF
	//	return (MMC3_BankAddress_A000 + uint32(tempAddr)) //& (PRGLength - 1)
	//} else if Addr < 0xE000 { //$C000-$DFFF
	//	return (MMC3_BankAddress_C000 + uint32(tempAddr)) //& (PRGLength - 1)
	//} else { //$E000-$FFFF
	//	return (MMC3_LastBank + uint32(tempAddr)) //& (PRGLength - 1)
	//}
	bank := (Addr - 0x8000) / 0x2000
	offset := uint32(Addr & 0x1FFF)
	return MMC3_PRGOffset[bank] + offset
}

func MMC3_FetchPPUAddress(Addr uint16, CHRLength uint32) uint32 {
	//if MMC3_CHRBankMode {
	//	if Addr < 0x0400 { // 1KB Bank
	//		return MMC3_CHRBankAddress2 + uint32(Addr&0x3FF)
	//	} else if Addr < 0x0800 { // 1KB Bank
	//		return MMC3_CHRBankAddress3 + uint32(Addr&0x3FF)
	//	} else if Addr < 0x0C00 { // 1KB Bank
	//		return MMC3_CHRBankAddress4 + uint32(Addr&0x3FF)
	//	} else if Addr < 0x1000 { // 1KB Bank
	//		return MMC3_CHRBankAddress5 + uint32(Addr&0x3FF)
	//	} else if Addr < 0x1800 { // 2KB Bank
	//		return MMC3_CHRBankAddress0 + uint32(Addr&0x7FF)
	//	} else { // 2KB Bank
	//		return MMC3_CHRBankAddress1 + uint32(Addr&0x7FF)
	//	}
	//} else {
	//	if Addr < 0x0800 { // 2KB Bank
	//		return MMC3_CHRBankAddress0 + uint32(Addr&0x7FF)
	//	} else if Addr < 0x1000 { // 2KB Bank
	//		return MMC3_CHRBankAddress1 + uint32(Addr&0x7FF)
	//	} else if Addr < 0x1400 { // 1KB Bank
	//		return MMC3_CHRBankAddress2 + uint32(Addr&0x3FF)
	//	} else if Addr < 0x1800 { // 1KB Bank
	//		return MMC3_CHRBankAddress3 + uint32(Addr&0x3FF)
	//	} else if Addr < 0x1C00 { // 1KB Bank
	//		return MMC3_CHRBankAddress4 + uint32(Addr&0x3FF)
	//	} else { // 1KB Bank
	//		return MMC3_CHRBankAddress5 + uint32(Addr&0x3FF)
	//	}
	//}
	bank := Addr / 0x400
	offset := uint32(Addr & 0x3FF)
	return MMC3_CHROffset[bank] + offset
}

func MMC3_UpdateRegister() {
	if MMC3_PRGBankMode {
		MMC3_PRGOffset[0] = MMC3_2ndToLastBank
		MMC3_PRGOffset[1] = MMC3_GetPRGOffset(MMC3_Register[7])
		MMC3_PRGOffset[2] = MMC3_GetPRGOffset(MMC3_Register[6])
		MMC3_PRGOffset[3] = MMC3_LastBank
	} else {
		MMC3_PRGOffset[0] = MMC3_GetPRGOffset(MMC3_Register[6])
		MMC3_PRGOffset[1] = MMC3_GetPRGOffset(MMC3_Register[7])
		MMC3_PRGOffset[2] = MMC3_2ndToLastBank
		MMC3_PRGOffset[3] = MMC3_LastBank
	}

	if MMC3_CHRBankMode {
		MMC3_CHROffset[0] = MMC3_GetCHROffset(MMC3_Register[2])
		MMC3_CHROffset[1] = MMC3_GetCHROffset(MMC3_Register[3])
		MMC3_CHROffset[2] = MMC3_GetCHROffset(MMC3_Register[4])
		MMC3_CHROffset[3] = MMC3_GetCHROffset(MMC3_Register[5])
		MMC3_CHROffset[4] = MMC3_GetCHROffset(MMC3_Register[0] & 0xFE)
		MMC3_CHROffset[5] = MMC3_GetCHROffset(MMC3_Register[0] | 1)
		MMC3_CHROffset[6] = MMC3_GetCHROffset(MMC3_Register[1] & 0xFE)
		MMC3_CHROffset[7] = MMC3_GetCHROffset(MMC3_Register[1] | 1)
	} else {
		MMC3_CHROffset[0] = MMC3_GetCHROffset(MMC3_Register[0] & 0xFE)
		MMC3_CHROffset[1] = MMC3_GetCHROffset(MMC3_Register[0] | 1)
		MMC3_CHROffset[2] = MMC3_GetCHROffset(MMC3_Register[1] & 0xFE)
		MMC3_CHROffset[3] = MMC3_GetCHROffset(MMC3_Register[1] | 1)
		MMC3_CHROffset[4] = MMC3_GetCHROffset(MMC3_Register[2])
		MMC3_CHROffset[5] = MMC3_GetCHROffset(MMC3_Register[3])
		MMC3_CHROffset[6] = MMC3_GetCHROffset(MMC3_Register[4])
		MMC3_CHROffset[7] = MMC3_GetCHROffset(MMC3_Register[5])
	}

}

func MMC3_GetPRGOffset(offset byte) uint32 {
	return uint32(0x2000 * uint32(offset&0x3F))
}

func MMC3_GetCHROffset(offset byte) uint32 {
	return (0x400 * uint32(offset))
}

func MMC3_ClockIRQ(Addr uint16) {
	MMC3_CheckA12Pin(Addr)
	if !MMC3_PPUA12Prev && MMC3_PPUA12 && MMC3_PPUA12PrevCount == 0 {
		//Check for IRQ
		if MMC3_IRQCounter == 0 && MMC3_IRQEnabled {
			MMC3_DoIRQ = true
		}

		if MMC3_IRQCounter == 0 || MMC3_IRQReloadFlag {
			MMC3_IRQCounter = MMC3_IRQReloadValue
			MMC3_IRQReloadFlag = false
		} else {
			MMC3_IRQCounter--
		}
	}
}

var (
	MMC3_PPUA12PrevCount int
)

func MMC3_CheckA12Pin(Addr uint16) {
	MMC3_PPUA12Prev = MMC3_PPUA12
	MMC3_PPUA12 = (Addr & 0b0001000000000000) != 0
	if MMC3_PPUA12 {
		if MMC3_PPUA12PrevCount > 0 {
			MMC3_PPUA12PrevCount--
		}
	} else {
		MMC3_PPUA12PrevCount = 3
	}
}
