package cartridge

import (
	"mtt/timenes/common"
	"mtt/timenes/nes/cartridge/mappers"
	"os"
)

// Header Variables

var PRGROM_Size uint32         // Byte 4: Size of PRG ROM in 16KB Units
var CHRROM_Size uint32         // Byte 5: Size of CHR ROM in 8 KB units (value 0 means the board uses CHR RAM)
var IsNametableHorizontal bool // Byte 6, Bit 0: Nametable Arrangement (0: Vertical / 1: Horizontal)
var HasBatteryRAM bool         // Byte 6, Bit 1: Cartridge contains battery-backed PRG RAM ($6000-7FFF) or other persistent memory
var AltNametableLayout bool    // Byte 6, Bit 3: Alternative nametable layout
var NES2_Header bool           // Byte 7, Bits 2-3: If equal to 2, flags 8-15 are in NES 2.0 format
var MapperChipID byte          // Gotten from upper half of Bytes 6 and 7

var Header [0x10]byte
var PRGROM [0x80000]byte
var CHRROM [0x80000]byte
var IsCHRRAM bool

var RAM [0x800]byte
var VRAM [0x800]byte
var PaletteRAM [0x20]byte
var CartRAM [0x2000]byte
var CartVRAM [0x1000]byte //For mapper chips

func ResetCartridge() {
	Header = [0x10]byte{}
	PRGROM = [0x80000]byte{}
	CHRROM = [0x80000]byte{}
	IsCHRRAM = false

	//Reset RAM (Or don't, if I wanna do weird stuff...?)
	RAM = [0x800]byte{}
	VRAM = [0x800]byte{}
	PaletteRAM = [0x20]byte{}
	CartRAM = [0x2000]byte{}
	CartVRAM = [0x1000]byte{}

	PRGROM_Size = 0
	CHRROM_Size = 0
	IsNametableHorizontal = false
	HasBatteryRAM = false
	AltNametableLayout = false
	NES2_Header = false
	MapperChipID = 0
}

func LoadCartridge() {
	HeaderedROM, err := os.ReadFile(common.Filepath)
	common.Check(err)

	//TODO: Add a check to make sure the file loaded is indeed an NES game

	//Header info
	copy(Header[:], HeaderedROM[0x0:])
	PRGROM_Size = uint32(Header[4]) * uint32(0x4000)
	CHRROM_Size = uint32(Header[5]) * uint32(0x2000)
	IsNametableHorizontal = (Header[6] & 1) == 0
	HasBatteryRAM = (Header[6] & 0x02) != 0
	AltNametableLayout = (Header[6] & 0x08) != 0
	NES2_Header = ((Header[7] & 0xC) >> 2) == 2

	MapperChipID = (Header[6] >> 4) | (Header[7] & 0xF0)

	//size := uint16(Header[4])
	ROM_Endpoint := uint32(0x10 + (PRGROM_Size))
	CHR_Endpoint := uint32(ROM_Endpoint + uint32(CHRROM_Size))

	copy(PRGROM[:], HeaderedROM[0x10:ROM_Endpoint])
	if CHRROM_Size != 0 {
		copy(CHRROM[:], HeaderedROM[ROM_Endpoint:CHR_Endpoint])
	} else {

	}

	//Initialize any PRG-RAM from mapper chips
	if HasBatteryRAM {
		switch MapperChipID {
		case 1: //MMC1
			copy(mappers.MMC1_PRGRAM[:], HeaderedROM[0x10:])
		case 2: //UxROM
		case 3: //CNROM
			//Add support for Hayauchi Super Igo?
		case 4: //MMC3
		}
	}
	common.ROMLoaded = true
}
