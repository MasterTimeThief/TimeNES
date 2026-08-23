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

func PageCrossingCheck(pbCheck bool, Address uint16, lo, hi byte) {
	if pbCheck && byte((Address&0xFF00)>>8) != hi {
		CPU_Cycles++ //Extra cycle for crossing page boundary
	}
}

func SetAddressBusHigh(Value byte) {
	AddressBus &= 0x00FF
	AddressBus += uint16(Value) << 8
}

func SetAddressBusLow(Value byte) {
	AddressBus &= 0xFF00
	AddressBus += uint16(Value)
}

func SetPointerHigh(Value byte) {
	Pointer &= 0x00FF
	Pointer += uint16(Value) << 8
}

func SetPointerLow(Value byte) {
	Pointer &= 0xFF00
	Pointer += uint16(Value)
}

func SetTargetHigh(Value byte) {
	Target &= 0x00FF
	Target += uint16(Value) << 8
}

func SetTargetLow(Value byte) {
	Target &= 0xFF00
	Target += uint16(Value)
}

func (cpu *CPU) ReadOperands_AbsoluteAddressed(isJMP bool) {
	AddressBus = uint16(cpu.ReadFromPC())
	AddressBus = (uint16(cpu.ReadFromPC())<<8 | AddressBus)
	if !isJMP {
		//MasterClockTick("abs add, not jmp")
	}
}

func (cpu *CPU) ReadOperands_IndirectAddressed() {
	AddressBus = uint16(cpu.ReadFromPC())
	AddressBus = (uint16(cpu.ReadFromPC())<<8 | AddressBus)
	//Now read from HERE
	low := bus.Read(AddressBus)
	var high byte
	if AddressBus&0x00FF == 0xFF {
		//Original NMOS Bug
		high = bus.Read(AddressBus & 0xFF00)
	} else {
		high = bus.Read(AddressBus + 1)
	}
	AddressBus = cpu.BuildAddress(low, high)
}

func (cpu *CPU) ReadOperands_AbsoluteAddressed_XIndexed(pbCheck bool) {
	low := cpu.ReadFromPC()
	high := cpu.ReadFromPC()
	AddressBus = cpu.BuildAddress(low, high) + uint16(cpu.X)
	PageCrossingCheck(pbCheck, AddressBus, low, high)
}

func (cpu *CPU) ReadOperands_AbsoluteAddressed_YIndexed(pbCheck bool) {
	low := cpu.ReadFromPC()
	high := cpu.ReadFromPC()
	AddressBus = cpu.BuildAddress(low, high) + uint16(cpu.Y)
	PageCrossingCheck(pbCheck, AddressBus, low, high)
}

func (cpu *CPU) ReadOperands_IndirectAddressed_XIndexed() {
	TempAddress := cpu.ReadFromPC() + cpu.X
	low := bus.Read(uint16(TempAddress))
	high := bus.Read(uint16(TempAddress + 1))
	AddressBus = cpu.BuildAddress(low, high)
}

func (cpu *CPU) ReadOperands_IndirectAddressed_YIndexed(pbCheck bool) {
	TempAddress := cpu.ReadFromPC()
	low := bus.Read(uint16(TempAddress))
	high := bus.Read(uint16(TempAddress + 1))
	AddressBus = cpu.BuildAddress(low, high) + uint16(cpu.Y)
	PageCrossingCheck(pbCheck, AddressBus, low, high)
}

func (cpu *CPU) ReadOperands_ZeroPageAddressed() {
	AddressBus = cpu.BuildAddress(cpu.ReadFromPC(), 0x00)
}

func (cpu *CPU) ReadOperands_ZeroPageAddressed_XIndexed() {
	AddressBus = cpu.BuildAddress(cpu.ReadFromPC()+cpu.X, 0x00)
}

func (cpu *CPU) ReadOperands_ZeroPageAddressed_YIndexed() {
	AddressBus = cpu.BuildAddress(cpu.ReadFromPC()+cpu.Y, 0x00)
}
