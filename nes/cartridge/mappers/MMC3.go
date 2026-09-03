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
	MMC3_M2Count        int
	MMC3_IRQPending     bool
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
		if (Addr & 1) == 0 { // Bank Select ($8000-$9FFE, even)
			MMC3_BankSelect = Value & 0x7
			MMC3_PRGBankMode = (Value & 0x40) != 0
			MMC3_CHRBankMode = (Value & 0x80) != 0
			MMC3_UpdateRegister()
			// Also check for MMC6 behavior here
		} else { // Bank data ($8001-$9FFF, odd)
			MMC3_Register[MMC3_BankSelect] = Value
			MMC3_UpdateRegister()
		}
	} else if Addr < 0xC000 {
		if (Addr & 1) == 0 { // Nametable arrangement ($A000-$BFFE, even)
			MMC3_IsHorizontalNametable = (Value & 1) == 0
		} else { // PRG RAM protect ($A001-$BFFF, odd)
			MMC3_WriteProtect = (Value & 0x40) != 0
			MMC3_PRGRAMEnable = (Value & 0x80) != 0
			// Also check for MMC6 behavior here
		}
	} else if Addr < 0xE000 {
		if (Addr & 1) == 0 { // IRQ latch ($C000-$DFFE, even)
			MMC3_IRQReloadValue = Value
		} else { // IRQ reload ($C001-$DFFF, odd)
			MMC3_IRQReloadFlag = true
		}
	} else {
		if (Addr & 1) == 0 { // IRQ disable ($E000-$FFFE, even)
			MMC3_IRQEnabled = false
			MMC3_IRQPending = false
		} else { // IRQ enable ($E001-$FFFF, odd)
			MMC3_IRQEnabled = true
		}
	}
}

func MMC3_FetchCPUAddress(Addr uint16, PRGLength uint32) uint32 {
	bank := (Addr - 0x8000) / 0x2000
	offset := uint32(Addr & 0x1FFF)
	return MMC3_PRGOffset[bank] + offset
}

func MMC3_FetchPPUAddress(Addr uint16, CHRLength uint32) uint32 {
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
	MMC3_PPUA12Prev = MMC3_PPUA12
	MMC3_PPUA12 = (Addr & 0b0001000000000000) != 0

	//Check for IRQ
	if !MMC3_PPUA12Prev && MMC3_PPUA12 && MMC3_M2Count == 3 {
		if MMC3_IRQReloadFlag {
			MMC3_IRQCounter = MMC3_IRQReloadValue
			MMC3_IRQReloadFlag = false
			if MMC3_IRQCounter == 0 && MMC3_IRQEnabled {
				MMC3_IRQPending = true
			}
		} else {
			MMC3_IRQCounter--
			if MMC3_IRQCounter == 0 && MMC3_IRQEnabled {
				MMC3_IRQPending = true
			} else if MMC3_IRQCounter == 0xFF {
				MMC3_IRQCounter = MMC3_IRQReloadValue
				MMC3_IRQReloadFlag = false
				if MMC3_IRQCounter == 0 && MMC3_IRQEnabled {
					MMC3_IRQPending = true
				}
			}
		}
	}
	if MMC3_PPUA12 {
		MMC3_M2Count = 0
	}
}

func MMC3_ClockM2(Addr uint16) {
	if (Addr & 0b0001000000000000) == 0 {
		if MMC3_M2Count < 3 {
			MMC3_M2Count++
		}
	}
}
