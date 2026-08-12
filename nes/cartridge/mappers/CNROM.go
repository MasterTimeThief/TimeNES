// Mapper 003: CNROM
package mappers

var CNROM_Register byte

func CNROM_Write(Value byte, Addr uint16) {
	// Bits 0-1 are for the CHR Bank to use
	// 4-5 are for diodes used in copy protection
	if Addr < 0x8000 {
		//PRG-RAM stuff, only used by one game afaik
	} else {
		CNROM_Register = Value
	}
}

func CNROM_ReadAddress(ppuAddr uint16) uint16 {
	//Update PPU Address Bus
	cnromAddress := (ppuAddr & 0x1FFF) | uint16(CNROM_Register&3)<<13
	return cnromAddress
}
