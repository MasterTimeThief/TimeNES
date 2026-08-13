package cpu

import "mtt/timenes/nes/bus"

func ReadOperands_AbsoluteAddressed(isJMP bool) uint16 {
	AddressBus := uint16(ReadFromPC())
	//MasterClockTick("Abs add")
	AddressBus = (uint16(ReadFromPC())<<8 | AddressBus)
	if !isJMP {
		//MasterClockTick("abs add, not jmp")
	}
	return AddressBus
}

func ReadOperands_IndirectAddressed() uint16 {
	AddressBus := uint16(ReadFromPC())
	AddressBus = (uint16(ReadFromPC())<<8 | AddressBus)
	//Now read from HERE
	indL := bus.Read(AddressBus)
	var indH byte
	if AddressBus&0x00FF == 0xFF {
		//Original NMOS Bug
		indH = bus.Read(AddressBus & 0xFF00)
	} else {
		indH = bus.Read(AddressBus + 1)
	}
	//MasterClockTick("ind add")
	return BuildAddress(indL, indH)
}

func ReadOperands_AbsoluteAddressed_XIndexed(pbCheck bool) uint16 {
	low := ReadFromPC()
	high := ReadFromPC()
	AddressBus = BuildAddress(low, high)
	AddressBus += uint16(X)
	if pbCheck && byte((AddressBus&0xFF00)>>4) != high {
		////MasterClockTick("abs add x")
		CPU_Cycles++ //Extra cycle for crossing page boundary
	}

	return AddressBus
}

func ReadOperands_AbsoluteAddressed_YIndexed(pbCheck bool) uint16 {
	low := ReadFromPC()
	high := ReadFromPC()
	AddressBus = BuildAddress(low, high)
	AddressBus += uint16(Y)
	if pbCheck && byte((AddressBus&0xFF00)>>4) != high {
		////MasterClockTick("abs add y")
		CPU_Cycles++ //Extra cycle for crossing page boundary
	}
	return AddressBus
}

func ReadOperands_IndirectAddressed_XIndexed() uint16 {
	Addr := ReadFromPC() + X
	////MasterClockTick("ind x")
	TempAddress := Addr
	Addr = bus.Read(uint16(TempAddress)) //Low byte of new address
	TempAddress++
	AddressBus = (uint16(bus.Read(uint16(TempAddress)))<<8 | uint16(Addr)) //High byte
	return AddressBus
}

func ReadOperands_IndirectAddressed_YIndexed(pbCheck bool) uint16 {
	Addr := ReadFromPC()
	TempAddress := Addr
	Addr = bus.Read(uint16(TempAddress)) //Low byte of new address
	TempAddress++
	AddressBus = (uint16(bus.Read(uint16(TempAddress)))<<8 | uint16(Addr)) //High byte
	AddressBus += uint16(Y)
	if pbCheck && bus.Read(uint16(TempAddress)) != byte((AddressBus&0xFF00)>>4) {
		////MasterClockTick("ind add y")
		CPU_Cycles++
	}
	return AddressBus
}

func ReadOperands_ZeroPageAddressed() uint16 {
	AddressBus := ReadFromPC()
	return BuildAddress(AddressBus, 0x00)
}

func ReadOperands_ZeroPageAddressed_XIndexed() uint16 {
	AddressBus := ReadFromPC()
	////MasterClockTick("Zp x")
	return BuildAddress(AddressBus+X, 0x00)
}

func ReadOperands_ZeroPageAddressed_YIndexed() uint16 {
	AddressBus := ReadFromPC()
	////MasterClockTick("ZP Y")
	return BuildAddress(AddressBus+Y, 0x00)
}
