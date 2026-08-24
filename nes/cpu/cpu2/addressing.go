package cpu2

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

var FixHighByte bool

// Returns true if page boundary was crossed
func PageCrossingCheck(Address uint16, index byte) bool {
	return ((Address + uint16(index)) & 0xFF00) != (Address & 0xFF00)
}

func (cpu *CPU) SetAddressBusHigh(Value byte) {
	cpu.AddressBus &= 0x00FF
	cpu.AddressBus += uint16(Value) << 8
}

func (cpu *CPU) SetAddressBusLow(Value byte) {
	cpu.AddressBus &= 0xFF00
	cpu.AddressBus += uint16(Value)
}

func (cpu *CPU) SetPointerHigh(Value byte) {
	cpu.Pointer &= 0x00FF
	cpu.Pointer += uint16(Value) << 8
}

func (cpu *CPU) SetPointerLow(Value byte) {
	cpu.Pointer &= 0xFF00
	cpu.Pointer += uint16(Value)
}

func (cpu *CPU) SetTargetHigh(Value byte) {
	cpu.Target &= 0x00FF
	cpu.Target += uint16(Value) << 8
}

func (cpu *CPU) SetTargetLow(Value byte) {
	cpu.Target &= 0xFF00
	cpu.Target += uint16(Value)
}

// Addressing Modes

// Fetch the value at the program counter, store it in the DataLatch, and increment the Program Counter.
//
// 1 Step
func (cpu *CPU) GetAddress_Immediate() {
	cpu.DL = cpu.ReadFromPC()
	cpu.AddressBus = cpu.PC
}

// Fetch the value at the PC, and write to either the
// High byte or Low byte of the 16 bit address bus.
// Also increment the Program Counter.
//
// 2 Steps
func (cpu *CPU) GetAddress_Absolute() {
	switch cpu.subCycle {
	case 1:
		cpu.DL = cpu.ReadFromPC()
	case 2:
		cpu.AddressBus = (uint16(cpu.ReadFromPC())<<8 | uint16(cpu.DL))
	}
}

func (cpu *CPU) GetAddress_Indirect() {
	//AddressBus = uint16(ReadFromPC())
	//AddressBus = (uint16(ReadFromPC())<<8 | AddressBus)
	////Now read from HERE
	//low := bus.Read(AddressBus)
	//var high byte
	//if AddressBus&0x00FF == 0xFF {
	//	//Original NMOS Bug
	//	high = bus.Read(AddressBus & 0xFF00)
	//} else {
	//	high = bus.Read(AddressBus + 1)
	//}
	//AddressBus = BuildAddress(low, high)
}

// Fetch the High and Low byte values from the byte at the PC, then add X.
//
// 3-4 Steps
func (cpu *CPU) GetAddress_AbsoluteX(pbCheck bool) {
	// Some instructions will always take 4 cycles to determine the address,
	// and others will normally take 3, but take the extra cycle if a page boundary was crossed.
	if pbCheck {
		switch cpu.subCycle {
		case 1:
			cpu.DL = cpu.ReadFromPC()
		case 2:
			cpu.AddressBus = (uint16(cpu.ReadFromPC())<<8 | uint16(cpu.DL))
			cpu.TempAddress = cpu.AddressBus
			cpu.H = byte(cpu.AddressBus >> 8)

			if PageCrossingCheck(cpu.TempAddress, cpu.X) {
				FixHighByte = true
			} else {
				cpu.subCycle++
				FixHighByte = false
			}
			cpu.AddressBus = (cpu.AddressBus & 0xFF00) | ((cpu.AddressBus + uint16(cpu.X)) & 0xFF)
		case 3:
			cpu.DL = cpu.ReadFromAB()
			cpu.H = byte(cpu.AddressBus >> 8)
			cpu.H++
			if FixHighByte {
				cpu.AddressBus += 0x100
			}
		case 4:
			cpu.DL = cpu.ReadFromAB() // Dummy Read
		}
	} else {
		switch cpu.subCycle {
		case 1:
			cpu.DL = cpu.ReadFromPC()
		case 2:
			cpu.AddressBus = (uint16(cpu.ReadFromPC())<<8 | uint16(cpu.DL))
			cpu.TempAddress = cpu.AddressBus
			cpu.AddressBus = (cpu.AddressBus & 0xFF00) | ((cpu.AddressBus + uint16(cpu.X)) & 0xFF)
		case 3:
			cpu.DL = cpu.ReadFromAB()
			cpu.H = byte(cpu.AddressBus >> 8)
			cpu.H++
			if PageCrossingCheck(cpu.TempAddress, cpu.X) {
				cpu.AddressBus += 0x100
			}
		case 4:
			cpu.DL = cpu.ReadFromAB() // Dummy Read
		}
	}
}

// Fetch the High and Low byte values from the byte at the PC, then add Y.
//
// 3-4 Steps
func (cpu *CPU) GetAddress_AbsoluteY(pbCheck bool) {
	// Some instructions will always take 4 cycles to determine the address,
	// and others will normally take 3, but take the extra cycle if a page boundary was crossed.
	if pbCheck {
		switch cpu.subCycle {
		case 1:
			cpu.DL = cpu.ReadFromPC()
		case 2:
			cpu.AddressBus = (uint16(cpu.ReadFromPC())<<8 | uint16(cpu.DL))
			cpu.TempAddress = cpu.AddressBus
			cpu.H = byte(cpu.AddressBus >> 8)

			if PageCrossingCheck(cpu.TempAddress, cpu.Y) {
				FixHighByte = true
			} else {
				cpu.subCycle++
				FixHighByte = false
			}

			cpu.AddressBus = (cpu.AddressBus & 0xFF00) | ((cpu.AddressBus + uint16(cpu.Y)) & 0xFF)
		case 3:
			cpu.DL = cpu.ReadFromAB()
			cpu.H = byte(cpu.AddressBus >> 8)
			cpu.H++
			if FixHighByte {
				cpu.AddressBus += 0x100
			}
		case 4:
			cpu.DL = cpu.ReadFromAB() // Dummy Read
		}
	} else {
		switch cpu.subCycle {
		case 1:
			cpu.DL = cpu.ReadFromPC()
		case 2:
			cpu.AddressBus = (uint16(cpu.ReadFromPC())<<8 | uint16(cpu.DL))
			cpu.TempAddress = cpu.AddressBus
			cpu.AddressBus = (cpu.AddressBus & 0xFF00) | ((cpu.AddressBus + uint16(cpu.Y)) & 0xFF)
		case 3:
			cpu.DL = cpu.ReadFromAB() // Dummy read
			cpu.H = byte(cpu.AddressBus >> 8)
			cpu.H++
			if PageCrossingCheck(cpu.TempAddress, cpu.Y) {
				cpu.AddressBus += 0x100
			}
		case 4:
			cpu.DL = cpu.ReadFromAB() // Dummy Read
		}
	}
}

// Fetch the value from the PC, then using
// that value as an 8-bit address on the zero page,
// add the X register, then set the High byte and
// Low byte of the Address Bus from there.
//
// 4 Steps
func (cpu *CPU) GetAddress_IndirectX() {
	switch cpu.subCycle {
	case 1: // Fetch pointer address
		cpu.AddressBus = uint16(cpu.ReadFromPC())
	case 2: // Add X
		cpu.ReadFromAB() // Dummy Read
		cpu.AddressBus = (cpu.AddressBus + uint16(cpu.X)) & 0xFF
	case 3: // Fetch address low
		cpu.DL = cpu.ReadFromAB()
	case 4: // fetch address high
		cpu.AddressBus = (cpu.AddressBus + 1) & 0xFF
		cpu.AddressBus = (uint16(cpu.ReadFromAB())<<8 | uint16(cpu.DL))
	}
}

// Fetch the value from the PC.
// use that 8 bit location on the
// zero page to fetch the High and
// Low byte of the new Address Bus location,
// then add Y to that.
//
// 3-4 Steps
func (cpu *CPU) GetAddress_IndirectY(pbCheck bool) {

	// Some instructions will always take 4 cycles to determine the address,
	// and others will normally take 3, but take the extra cycle if a page boundary was crossed.
	if pbCheck {
		switch cpu.subCycle {
		case 1: // Fetch pointer address
			cpu.AddressBus = uint16(cpu.ReadFromPC())
		case 2: // fetch address low
			cpu.DL = cpu.ReadFromAB()
		case 3: // fetch address high, add Y to low byte
			cpu.AddressBus = (cpu.AddressBus + 1) & 0xFF
			cpu.AddressBus = (uint16(cpu.ReadFromAB())<<8 | uint16(cpu.DL))
			cpu.TempAddress = cpu.AddressBus
			cpu.H = byte(cpu.AddressBus >> 8)
			if !PageCrossingCheck(cpu.TempAddress, cpu.Y) {
				cpu.subCycle++
			}
			cpu.AddressBus = (cpu.AddressBus & 0xFF00) | ((cpu.AddressBus + uint16(cpu.Y)) & 0xFF)
		case 4: // increment high byte
			cpu.DL = cpu.ReadFromAB() // Dummy read
			cpu.H = byte(cpu.AddressBus >> 8)
			cpu.H++
			cpu.AddressBus += 0x100
		}
	} else {
		switch cpu.subCycle {
		case 1: // Fetch pointer address
			cpu.AddressBus = uint16(cpu.ReadFromPC())
		case 2: // fetch address low
			cpu.DL = cpu.ReadFromAB()
		case 3: // fetch address high, add Y to low byte
			cpu.AddressBus = (cpu.AddressBus + 1) & 0xFF
			cpu.TempAddress = (uint16(cpu.ReadFromAB())<<8 | uint16(cpu.DL))
			cpu.AddressBus = (cpu.TempAddress & 0xFF00) | ((cpu.TempAddress + uint16(cpu.Y)) & 0xFF)
		case 4: // increment high byte
			cpu.DL = cpu.ReadFromAB() // Dummy read
			cpu.H = byte(cpu.AddressBus >> 8)
			cpu.H++
			if PageCrossingCheck(cpu.TempAddress, cpu.Y) {
				cpu.AddressBus += 0x100
			}
		}
	}

}

// Fetch the value at the PC, and this 8 bit value
// replaces the contents of the 16 bit address bus.
//
// 1 Step
func (cpu *CPU) GetAddress_ZeroPage() {
	cpu.AddressBus = uint16(cpu.ReadFromPC())
}

// Fetch the value from the PC, then add X to that.
//
// 2 Steps
func (cpu *CPU) GetAddress_ZeroPageX() {
	switch cpu.subCycle {
	case 1: // Fetch address
		cpu.AddressBus = uint16(cpu.ReadFromPC())
	case 2: // Dummy read, and add X
		cpu.DL = cpu.ReadFromAB()
		cpu.AddressBus = (cpu.AddressBus + uint16(cpu.X)) & 0xFF
	}
}

// Fetch the value from the PC, then add Y to that.
//
// 2 Steps
func (cpu *CPU) GetAddress_ZeroPageY() {
	switch cpu.subCycle {
	case 1: // Fetch address
		cpu.AddressBus = uint16(cpu.ReadFromPC())
	case 2: // Dummy read, and add Y
		cpu.DL = cpu.ReadFromAB()
		cpu.AddressBus = (cpu.AddressBus + uint16(cpu.Y)) & 0xFF
	}
}
