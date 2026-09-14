package cpu

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

func (c *CPU) SetAddressBusHigh(Value byte) {
	c.AddressBus &= 0x00FF
	c.AddressBus += (uint16(Value) << 8)
}

func (c *CPU) SetAddressBusLow(Value byte) {
	c.AddressBus &= 0xFF00
	c.AddressBus += uint16(Value)
}

func (c *CPU) SetPointerHigh(Value byte) {
	c.Pointer &= 0x00FF
	c.Pointer += uint16(Value) << 8
}

func (c *CPU) SetPointerLow(Value byte) {
	c.Pointer &= 0xFF00
	c.Pointer += uint16(Value)
}

func (c *CPU) SetTargetHigh(Value byte) {
	c.Target &= 0x00FF
	c.Target += uint16(Value) << 8
}

func (c *CPU) SetTargetLow(Value byte) {
	c.Target &= 0xFF00
	c.Target += uint16(Value)
}

// Addressing Modes

// Fetch the value at the program counter, store it in the DataLatch, and increment the Program Counter.
//
// 1 Step
func (c *CPU) GetAddress_Immediate() {
	c.DL = c.ReadFromPC()
	c.AddressBus = c.PC
}

// Fetch the value at the PC, and write to either the
// High byte or Low byte of the 16 bit address bus.
// Also increment the Program Counter.
//
// 2 Steps
func (c *CPU) GetAddress_Absolute() {
	switch c.InstructionCycle {
	case 1:
		c.DL = c.ReadFromPC()
		c.SetAddressBusLow(c.DL)
	case 2:
		c.DL = c.ReadFromPC()
		c.SetAddressBusHigh(c.DL)
	}
}

func (c *CPU) GetAddress_Indirect() {
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
func (c *CPU) GetAddress_AbsoluteX(pbCheck bool) {
	// Some instructions will always take 4 cycles to determine the address,
	// and others will normally take 3, but take the extra cycle if a page boundary was crossed.
	if pbCheck {
		switch c.InstructionCycle {
		case 1:
			c.DL = c.ReadFromPC()
		case 2:
			c.AddressBus = (uint16(c.ReadFromPC())<<8 | uint16(c.DL))
			c.TempAddress = c.AddressBus
			c.H = byte(c.AddressBus >> 8)

			if PageCrossingCheck(c.TempAddress, c.X) {
				FixHighByte = true
			} else {
				c.InstructionCycle++
				FixHighByte = false
			}
			c.AddressBus = (c.AddressBus & 0xFF00) | ((c.AddressBus + uint16(c.X)) & 0xFF)
		case 3:
			c.DL = c.ReadFromAB()
			c.H = byte(c.AddressBus >> 8)
			c.H++
			if FixHighByte {
				c.AddressBus += 0x100
			}
		case 4:
			c.DL = c.ReadFromAB() // Dummy Read
		}
	} else {
		switch c.InstructionCycle {
		case 1:
			c.DL = c.ReadFromPC()
		case 2:
			c.AddressBus = (uint16(c.ReadFromPC())<<8 | uint16(c.DL))
			c.TempAddress = c.AddressBus
			c.AddressBus = (c.AddressBus & 0xFF00) | ((c.AddressBus + uint16(c.X)) & 0xFF)
		case 3:
			c.DL = c.ReadFromAB()
			c.H = byte(c.AddressBus >> 8)
			c.H++
			if PageCrossingCheck(c.TempAddress, c.X) {
				c.AddressBus += 0x100
			}
		case 4:
			c.DL = c.ReadFromAB() // Dummy Read
		}
	}
}

// Fetch the High and Low byte values from the byte at the PC, then add Y.
//
// 3-4 Steps
func (c *CPU) GetAddress_AbsoluteY(pbCheck bool) {
	// Some instructions will always take 4 cycles to determine the address,
	// and others will normally take 3, but take the extra cycle if a page boundary was crossed.
	if pbCheck {
		switch c.InstructionCycle {
		case 1:
			c.DL = c.ReadFromPC()
		case 2:
			c.AddressBus = (uint16(c.ReadFromPC())<<8 | uint16(c.DL))
			c.TempAddress = c.AddressBus
			c.H = byte(c.AddressBus >> 8)

			if PageCrossingCheck(c.TempAddress, c.Y) {
				FixHighByte = true
			} else {
				c.InstructionCycle++
				FixHighByte = false
			}

			c.AddressBus = (c.AddressBus & 0xFF00) | ((c.AddressBus + uint16(c.Y)) & 0xFF)
		case 3:
			c.DL = c.ReadFromAB()
			c.H = byte(c.AddressBus >> 8)
			c.H++
			if FixHighByte {
				c.AddressBus += 0x100
			}
		case 4:
			c.DL = c.ReadFromAB() // Dummy Read
		}
	} else {
		switch c.InstructionCycle {
		case 1:
			c.DL = c.ReadFromPC()
		case 2:
			c.AddressBus = (uint16(c.ReadFromPC())<<8 | uint16(c.DL))
			c.TempAddress = c.AddressBus
			c.AddressBus = (c.AddressBus & 0xFF00) | ((c.AddressBus + uint16(c.Y)) & 0xFF)
		case 3:
			c.DL = c.ReadFromAB() // Dummy read
			c.H = byte(c.AddressBus >> 8)
			c.H++
			if PageCrossingCheck(c.TempAddress, c.Y) {
				c.AddressBus += 0x100
			}
		case 4:
			c.DL = c.ReadFromAB() // Dummy Read
		}
	}
}

// Fetch the value from the PC, then using
// that value as an 8-bit address on the zero page,
// add the X register, then set the High byte and
// Low byte of the Address Bus from there.
//
// 4 Steps
func (c *CPU) GetAddress_IndirectX() {
	switch c.InstructionCycle {
	case 1: // Fetch pointer address
		c.AddressBus = uint16(c.ReadFromPC())
	case 2: // Add X
		c.ReadFromAB() // Dummy Read
		c.AddressBus = (c.AddressBus + uint16(c.X)) & 0xFF
	case 3: // Fetch address low
		c.DL = c.ReadFromAB()
	case 4: // fetch address high
		c.AddressBus = (c.AddressBus + 1) & 0xFF
		c.AddressBus = (uint16(c.ReadFromAB())<<8 | uint16(c.DL))
	}
}

// Fetch the value from the PC.
// use that 8 bit location on the
// zero page to fetch the High and
// Low byte of the new Address Bus location,
// then add Y to that.
//
// 3-4 Steps
func (c *CPU) GetAddress_IndirectY(pbCheck bool) {

	// Some instructions will always take 4 cycles to determine the address,
	// and others will normally take 3, but take the extra cycle if a page boundary was crossed.
	if pbCheck {
		switch c.InstructionCycle {
		case 1: // Fetch pointer address
			c.AddressBus = uint16(c.ReadFromPC())
		case 2: // fetch address low
			c.DL = c.ReadFromAB()
		case 3: // fetch address high, add Y to low byte
			c.AddressBus = (c.AddressBus + 1) & 0xFF
			c.AddressBus = (uint16(c.ReadFromAB())<<8 | uint16(c.DL))
			c.TempAddress = c.AddressBus
			c.H = byte(c.AddressBus >> 8)
			if !PageCrossingCheck(c.TempAddress, c.Y) {
				c.InstructionCycle++
			}
			c.AddressBus = (c.AddressBus & 0xFF00) | ((c.AddressBus + uint16(c.Y)) & 0xFF)
		case 4: // increment high byte
			c.DL = c.ReadFromAB() // Dummy read
			c.H = byte(c.AddressBus >> 8)
			c.H++
			c.AddressBus += 0x100
		}
	} else {
		switch c.InstructionCycle {
		case 1: // Fetch pointer address
			c.AddressBus = uint16(c.ReadFromPC())
		case 2: // fetch address low
			c.DL = c.ReadFromAB()
		case 3: // fetch address high, add Y to low byte
			c.AddressBus = (c.AddressBus + 1) & 0xFF
			c.TempAddress = (uint16(c.ReadFromAB())<<8 | uint16(c.DL))
			c.AddressBus = (c.TempAddress & 0xFF00) | ((c.TempAddress + uint16(c.Y)) & 0xFF)
		case 4: // increment high byte
			c.DL = c.ReadFromAB() // Dummy read
			c.H = byte(c.AddressBus >> 8)
			c.H++
			if PageCrossingCheck(c.TempAddress, c.Y) {
				c.AddressBus += 0x100
			}
		}
	}

}

// Fetch the value at the PC, and this 8 bit value
// replaces the contents of the 16 bit address bus.
//
// 1 Step
func (c *CPU) GetAddress_ZeroPage() {
	c.AddressBus = uint16(c.ReadFromPC())
}

// Fetch the value from the PC, then add X to that.
//
// 2 Steps
func (c *CPU) GetAddress_ZeroPageX() {
	switch c.InstructionCycle {
	case 1: // Fetch address
		c.AddressBus = uint16(c.ReadFromPC())
	case 2: // Dummy read, and add X
		c.DL = c.ReadFromAB()
		c.AddressBus = (c.AddressBus + uint16(c.X)) & 0xFF
	}
}

// Fetch the value from the PC, then add Y to that.
//
// 2 Steps
func (c *CPU) GetAddress_ZeroPageY() {
	switch c.InstructionCycle {
	case 1: // Fetch address
		c.AddressBus = uint16(c.ReadFromPC())
	case 2: // Dummy read, and add Y
		c.DL = c.ReadFromAB()
		c.AddressBus = (c.AddressBus + uint16(c.Y)) & 0xFF
	}
}
