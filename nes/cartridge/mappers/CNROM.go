// Mapper 003: CNROM
package mappers

var CNROM_Register byte
var CNROM_PRGRAM [0x2000]byte

func CNROM_WriteToPRGRAM(Value byte, Addr uint16) {
	//PRG-RAM stuff, only used by one game afaik
	CNROM_PRGRAM[Addr&0x1FFF] = Value
}

func CNROM_Write(Value byte, Addr uint16) {
	// Bits 0-1 are for the CHR Bank to use
	// 4-5 are for diodes used in copy protection
	CNROM_Register = Value
}

func CNROM_FetchPPUAddress(ppuAddr uint16) uint16 {
	//Update PPU Address Bus
	cnromAddress := (ppuAddr & 0x1FFF) | uint16(CNROM_Register&3)<<13
	return cnromAddress
}
