// Mapper 003: CNROM
package cnrom

var CNROM_Register byte

func WriteToCNROM(Value byte) {
	// Bits 0-1 are for the CHR Bank to use
	// 4-5 are for diodes used in copy protection
	CNROM_Register = Value
}

func UpdatePPUAddressBus(ppuAddr uint16) uint16 {
	cnromAddress := (ppuAddr & 0x1FFF) | uint16(CNROM_Register&3)<<13
	return cnromAddress
}
