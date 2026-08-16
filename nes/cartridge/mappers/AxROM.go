// Mapper 007: AxROM
package mappers

var AxROM_Register byte

// Clears the MMC1 Shift register to the default state
func AxROM_Write(Value byte, Addr uint16) {
	AxROM_Register = Value
}

func AxROM_FetchCPUAddress(Addr uint16, PRGLength uint32) uint32 {
	BankSelect := (AxROM_Register & 0x07)
	tempAddr := uint32(Addr & 0x7FFF)
	newAddr := ((0x8000 * uint32(BankSelect)) + tempAddr) & (PRGLength - 1)
	return newAddr
}

func AxROM_FetchNametable(Addr uint16) int {
	Arrangement := (AxROM_Register & 0x10) >> 4
	switch Arrangement {
	case 0: //Screen A only
		return int(Addr & 0x3FF)
	case 1: //Screen B Only
		return int(Addr&0x3FF) + 0x400

	default:
		return int(Addr)
	}
}
