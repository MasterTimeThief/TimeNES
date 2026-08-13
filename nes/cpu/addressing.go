package cpu

import "mtt/timenes/nes/bus"

/* Addressing Modes (Taken from NESDEV.org)
 	Abbr	Name				Cycles		Formula
	d,x		Zero page indexed	4			val = PEEK((arg + X) % 256)
	d,y		Zero page indexed	4			val = PEEK((arg + Y) % 256)
	a,x		Absolute indexed	4+			val = PEEK(arg + X)
	a,y		Absolute indexed	4+			val = PEEK(arg + Y)
	(d,x)	Indexed indirect	6			val = PEEK(PEEK((arg + X) % 256) + PEEK((arg + X + 1) % 256) * 256)
	(d),y	Indirect indexed	5+			val = PEEK(PEEK(arg) + PEEK((arg + 1) % 256) * 256 + Y)

	Other Addressing Modes

	Abbr	Name			Notes
			Implicit		Instructions like RTS or CLC have no address operand, the destination of results are implied.
	A		Accumulator		Many instructions can operate on the accumulator, e.g. LSR A. Some assemblers will treat no operand as an implicit A where applicable.
	#v		Immediate		Uses the 8-bit operand itself as the value for the operation, rather than fetching a value from a memory address.
	d		Zero page		Fetches the value from an 8-bit address on the zero page.
	a		Absolute		Fetches the value from a 16-bit address anywhere in memory.
	label	Relative		Branch instructions (e.g. BEQ, BCS) have a relative addressing mode that specifies an 8-bit signed offset relative to the current PC.
	(a)		Indirect		The JMP instruction has a special indirect addressing mode that can jump to the address stored in a 16-bit pointer anywhere in memory.

*/

func ReadOperands_AbsoluteAddressed(isJMP bool) {
	AddressBus = uint16(ReadFromPC())
	AddressBus = (uint16(ReadFromPC())<<8 | AddressBus)
	if !isJMP {
		//MasterClockTick("abs add, not jmp")
	}
	//return AddressBus
}

func ReadOperands_IndirectAddressed() {
	AddressBus = uint16(ReadFromPC())
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
	AddressBus = BuildAddress(indL, indH)
}

func ReadOperands_AbsoluteAddressed_XIndexed(pbCheck bool) {
	low := ReadFromPC()
	high := ReadFromPC()
	AddressBus = BuildAddress(low, high)
	AddressBus += uint16(X)
	if pbCheck && byte((AddressBus&0xFF00)>>4) != high {
		////MasterClockTick("abs add x")
		CPU_Cycles++ //Extra cycle for crossing page boundary
	}
	//return AddressBus
}

func ReadOperands_AbsoluteAddressed_YIndexed(pbCheck bool) {
	low := ReadFromPC()
	high := ReadFromPC()
	AddressBus = BuildAddress(low, high)
	AddressBus += uint16(Y)
	if pbCheck && byte((AddressBus&0xFF00)>>4) != high {
		////MasterClockTick("abs add y")
		CPU_Cycles++ //Extra cycle for crossing page boundary
	}
	//return AddressBus
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
