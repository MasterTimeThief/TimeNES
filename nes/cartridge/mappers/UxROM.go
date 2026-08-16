// Mapper 002: UxROM
package mappers

var UxROM_Register byte

func UxROM_Write(Value byte, Addr uint16) {
	UxROM_Register = Value
}

func UxROM_FetchCPUAddress(Addr uint16, PRGLength uint32) uint32 {
	tempAddr := uint32(Addr & 0x3FFF)
	if Addr >= 0xC000 {
		return PRGLength - 0x4000 + tempAddr
	} else {
		return (0x4000*uint32(UxROM_Register&0x0F) + tempAddr) & (PRGLength - 1)
	}
}
