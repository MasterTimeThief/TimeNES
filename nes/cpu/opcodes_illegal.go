package cpu

//----------------------------------------
//	NOP Codes (unofficial)
//----------------------------------------

func (c *CPU) X1A_NOP_Implied() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Read(c.PC)
	c.CompleteInstruction()
}
func (c *CPU) X3A_NOP_Implied() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Read(c.PC)
	c.CompleteInstruction()
}
func (c *CPU) X5A_NOP_Implied() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Read(c.PC)
	c.CompleteInstruction()
}
func (c *CPU) X7A_NOP_Implied() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Read(c.PC)
	c.CompleteInstruction()
}
func (c *CPU) XDA_NOP_Implied() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Read(c.PC)
	c.CompleteInstruction()
}
func (c *CPU) XFA_NOP_Implied() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Read(c.PC)
	c.CompleteInstruction()
}

func (c *CPU) X80_NOP_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.GetAddress_Immediate()
	c.CompleteInstruction()
}
func (c *CPU) X82_NOP_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.GetAddress_Immediate()
	c.CompleteInstruction()
}
func (c *CPU) X89_NOP_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.GetAddress_Immediate()
	c.CompleteInstruction()
}
func (c *CPU) XC2_NOP_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.GetAddress_Immediate()
	c.CompleteInstruction()
}
func (c *CPU) XE2_NOP_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.GetAddress_Immediate()
	c.CompleteInstruction()
}

func (c *CPU) X04_NOP_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.ReadFromAB() // Dummy Read
		c.CompleteInstruction()
	}
}
func (c *CPU) X44_NOP_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.ReadFromAB() // Dummy Read
		c.CompleteInstruction()
	}
}
func (c *CPU) X64_NOP_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.ReadFromAB() // Dummy Read
		c.CompleteInstruction()
	}
}

func (c *CPU) X14_NOP_ZeroPage_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.PollInterrupts()
		c.ReadFromAB() // Dummy Read
		c.CompleteInstruction()
	}
}
func (c *CPU) X34_NOP_ZeroPage_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.PollInterrupts()
		c.ReadFromAB() // Dummy Read
		c.CompleteInstruction()
	}
}
func (c *CPU) X54_NOP_ZeroPage_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.PollInterrupts()
		c.ReadFromAB() // Dummy Read
		c.CompleteInstruction()
	}
}
func (c *CPU) X74_NOP_ZeroPage_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.PollInterrupts()
		c.ReadFromAB() // Dummy Read
		c.CompleteInstruction()
	}
}
func (c *CPU) XD4_NOP_ZeroPage_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.PollInterrupts()
		c.ReadFromAB() // Dummy Read
		c.CompleteInstruction()
	}
}
func (c *CPU) XF4_NOP_ZeroPage_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.PollInterrupts()
		c.ReadFromAB() // Dummy Read
		c.CompleteInstruction()
	}
}

func (c *CPU) X0C_NOP_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(true)
	case 4:
		c.PollInterrupts()
		c.ReadFromAB() // Dummy Read
		c.CompleteInstruction()
	}
}

func (c *CPU) X1C_NOP_Absolute_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(true)
	case 4:
		c.PollInterrupts()
		c.ReadFromAB() // Dummy Read
		c.CompleteInstruction()
	}
}
func (c *CPU) X3C_NOP_Absolute_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(true)
	case 4:
		c.PollInterrupts()
		c.ReadFromAB() // Dummy Read
		c.CompleteInstruction()
	}
}
func (c *CPU) X5C_NOP_Absolute_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(true)
	case 4:
		c.PollInterrupts()
		c.ReadFromAB() // Dummy Read
		c.CompleteInstruction()
	}
}
func (c *CPU) X7C_NOP_Absolute_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(true)
	case 4:
		c.PollInterrupts()
		c.ReadFromAB() // Dummy Read
		c.CompleteInstruction()
	}
}
func (c *CPU) XDC_NOP_Absolute_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(true)
	case 4:
		c.PollInterrupts()
		c.ReadFromAB() // Dummy Read
		c.CompleteInstruction()
	}
}
func (c *CPU) XFC_NOP_Absolute_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(true)
	case 4:
		c.PollInterrupts()
		c.ReadFromAB() // Dummy Read
		c.CompleteInstruction()
	}
}

//----------------------------------------
// SAX: A AND X -> M
//----------------------------------------

func (c *CPU) X87_SAX_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.WriteToAB((c.A & c.X))
		c.CompleteInstruction()
	}
}
func (c *CPU) X97_SAX_ZeroPage_Y() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageY()
	case 3:
		c.PollInterrupts()
		c.WriteToAB((c.A & c.X))
		c.CompleteInstruction()
	}
}
func (c *CPU) X8F_SAX_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.WriteToAB((c.A & c.X))
		c.CompleteInstruction()
	}
}
func (c *CPU) X83_SAX_Indirect_X() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectX()
	case 5:
		c.PollInterrupts()
		c.WriteToAB((c.A & c.X))
		c.CompleteInstruction()
	}
}

//----------------------------------------
// LAX: LDA + LDX
//----------------------------------------

func (c *CPU) XA7_LAX_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.A = c.ReadFromAB()
		c.X = c.A
		c.SetZNFlags(c.X)
		c.CompleteInstruction()
	}
}
func (c *CPU) XB7_LAX_ZeroPage_Y() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageY()
	case 3:
		c.PollInterrupts()
		c.A = c.ReadFromAB()
		c.X = c.A
		c.SetZNFlags(c.X)
		c.CompleteInstruction()
	}
}
func (c *CPU) XAF_LAX_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.A = c.ReadFromAB()
		c.X = c.A
		c.SetZNFlags(c.X)
		c.CompleteInstruction()
	}
}
func (c *CPU) XBF_LAX_Absolute_Y() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(true)
	case 4:
		c.PollInterrupts()
		c.A = c.ReadFromAB()
		c.X = c.A
		c.SetZNFlags(c.X)
		c.CompleteInstruction()
	}
}
func (c *CPU) XA3_LAX_Indirect_X() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectX()
	case 5:
		c.PollInterrupts()
		c.A = c.ReadFromAB()
		c.X = c.A
		c.SetZNFlags(c.X)
		c.CompleteInstruction()
	}
}
func (c *CPU) XB3_LAX_Indirect_Y() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectY(true)
	case 5:
		c.PollInterrupts()
		c.A = c.ReadFromAB()
		c.X = c.A
		c.SetZNFlags(c.X)
		c.CompleteInstruction()
	}
}

//----------------------------------------
// SLO: ASL + ORA
//----------------------------------------

func (c *CPU) X07_SLO_ZeroPage() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.DL = c.ReadFromAB()
	case 3:
		c.WriteToAB(c.DL) // Dummy write
	case 4:
		c.PollInterrupts()
		c.Op_SLO()
		c.CompleteInstruction()
	}
}
func (c *CPU) X17_SLO_ZeroPage_X() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.DL = c.ReadFromAB()
	case 4:
		c.WriteToAB(c.DL) // Dummy write
	case 5:
		c.PollInterrupts()
		c.Op_SLO()
		c.CompleteInstruction()
	}
}
func (c *CPU) X0F_SLO_Absolute() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.DL = c.ReadFromAB()
	case 4:
		c.WriteToAB(c.DL) // Dummy write
	case 5:
		c.PollInterrupts()
		c.Op_SLO()
		c.CompleteInstruction()
	}
}
func (c *CPU) X1F_SLO_Absolute_X() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(false)
	case 4:
		c.DL = c.ReadFromAB()
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_SLO()
		c.CompleteInstruction()
	}
}
func (c *CPU) X1B_SLO_Absolute_Y() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(false)
	case 4:
		c.DL = c.ReadFromAB()
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_SLO()
		c.CompleteInstruction()
	}
}
func (c *CPU) X03_SLO_Indirect_X() {
	// CPU_Cycles = 8
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectX()
	case 5:
		c.DL = c.ReadFromAB()
	case 6:
		c.WriteToAB(c.DL) // Dummy write
	case 7:
		c.PollInterrupts()
		c.Op_SLO()
		c.CompleteInstruction()
	}
}
func (c *CPU) X13_SLO_Indirect_Y() {
	// CPU_Cycles = 8
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectY(false)
	case 5:
		c.DL = c.ReadFromAB()
	case 6:
		c.WriteToAB(c.DL) // Dummy write
	case 7:
		c.PollInterrupts()
		c.Op_SLO()
		c.CompleteInstruction()
	}
}

//----------------------------------------
// DCP: DEC + CMP
//----------------------------------------

func (c *CPU) XC7_DCP_ZeroPage() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.DL = c.ReadFromAB()
	case 3:
		c.WriteToAB(c.DL) // Dummy write
	case 4:
		c.PollInterrupts()
		c.Op_DCP()
		c.CompleteInstruction()
	}
}
func (c *CPU) XD7_DCP_ZeroPage_X() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.DL = c.ReadFromAB()
	case 4:
		c.WriteToAB(c.DL) // Dummy write
	case 5:
		c.PollInterrupts()
		c.Op_DCP()
		c.CompleteInstruction()
	}
}
func (c *CPU) XCF_DCP_Absolute() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.DL = c.ReadFromAB()
	case 4:
		c.WriteToAB(c.DL) // Dummy write
	case 5:
		c.PollInterrupts()
		c.Op_DCP()
		c.CompleteInstruction()
	}
}
func (c *CPU) XDF_DCP_Absolute_X() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(false)
	case 4:
		c.DL = c.ReadFromAB()
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_DCP()
		c.CompleteInstruction()
	}
}
func (c *CPU) XDB_DCP_Absolute_Y() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(false)
	case 4:
		c.DL = c.ReadFromAB()
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_DCP()
		c.CompleteInstruction()
	}
}
func (c *CPU) XC3_DCP_Indirect_X() {
	// CPU_Cycles = 8
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectX()
	case 5:
		c.DL = c.ReadFromAB()
	case 6:
		c.WriteToAB(c.DL) // Dummy write
	case 7:
		c.PollInterrupts()
		c.Op_DCP()
		c.CompleteInstruction()
	}
}
func (c *CPU) XD3_DCP_Indirect_Y() {
	// CPU_Cycles = 8
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectY(false)
	case 5:
		c.DL = c.ReadFromAB()
	case 6:
		c.WriteToAB(c.DL) // Dummy write
	case 7:
		c.PollInterrupts()
		c.Op_DCP()
		c.CompleteInstruction()
	}
}

//----------------------------------------
// SHA: Stores A AND X AND (high-byte of addr. + 1) at addr.
//----------------------------------------

func (c *CPU) X9F_SHA_Absolute_Y() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(false)
	case 4:
		c.PollInterrupts()
		if (c.TempAddress & 0xFF00) != (c.AddressBus & 0xFF00) {
			// If the page boundary was crossed, this code has gone "unstable"
			c.UnstableAddressBus(c.X)
		}
		c.WriteToAB(c.A & (c.X | c.Magic) & c.H)
		c.CompleteInstruction()
	}
}
func (c *CPU) X93_SHA_Indirect_Y() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectY(false)
	case 5:
		c.PollInterrupts()
		if (c.TempAddress & 0xFF00) != (c.AddressBus & 0xFF00) {
			// If the page boundary was crossed, this code has gone "unstable"
			c.UnstableAddressBus(c.X)
		}
		c.WriteToAB(c.A & (c.X | c.Magic) & c.H)
		c.CompleteInstruction()
	}
}

//----------------------------------------
// SHX: Stores X AND (high-byte of addr. + 1) at addr.
//----------------------------------------

func (c *CPU) X9E_SHX_Absolute_Y() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(false)
	case 4:
		c.PollInterrupts()
		if (c.TempAddress & 0xFF00) != (c.AddressBus & 0xFF00) {
			// If the page boundary was crossed, this code has gone "unstable"
			c.UnstableAddressBus(c.X)
		}
		c.WriteToAB(c.X & c.H)
		c.CompleteInstruction()
	}
}

//----------------------------------------
// SHY: Stores Y AND (high-byte of addr. + 1) at addr.
//----------------------------------------

func (c *CPU) X9C_SHY_Absolute_X() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(false)
	case 4:
		c.PollInterrupts()
		if (c.TempAddress & 0xFF00) != (c.AddressBus & 0xFF00) {
			// If the page boundary was crossed, this code has gone "unstable"
			c.UnstableAddressBus(c.Y)
		}
		c.WriteToAB(c.Y & c.H)
		c.CompleteInstruction()
	}
}

//----------------------------------------
// TAS (XAS, SHS): Puts A AND X in SP
// and stores A AND X AND (high-byte of addr. + 1) at addr.
// A AND X -> SP, A AND X AND (H+1) -> M
//----------------------------------------

func (c *CPU) X9B_TAS_Absolute_Y() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(false)
	case 4:
		c.PollInterrupts()
		if (c.TempAddress & 0xFF00) != (c.AddressBus & 0xFF00) {
			// If the page boundary was crossed, this code has gone "unstable"
			c.UnstableAddressBus(c.Y)
		}
		c.SP = c.A & c.X
		c.WriteToAB(c.A & (c.X | c.Magic) & c.H)
		c.CompleteInstruction()
	}
}

//----------------------------------------
// LAS (LAR, LAE): LDA/TSX oper
// M AND SP -> A, X, SP
//----------------------------------------

func (c *CPU) XBB_LAS_Absolute_Y() {
	// CPU_Cycles = 4+
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(true)
	case 4:
		c.PollInterrupts()
		c.DL = (c.ReadFromAB() & c.SP)
		c.A = c.DL
		c.X = c.DL
		c.SP = c.DL
		c.SetZNFlags(c.A)
		c.CompleteInstruction()
	}
}

//----------------------------------------
// RLA: ROL + AND
//----------------------------------------

func (c *CPU) X27_RLA_ZeroPage() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.DL = c.ReadFromAB()
	case 3:
		c.WriteToAB(c.DL) // Dummy write
	case 4:
		c.PollInterrupts()
		c.Op_RLA()
		c.CompleteInstruction()
	}
}
func (c *CPU) X37_RLA_ZeroPage_X() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.DL = c.ReadFromAB()
	case 4:
		c.WriteToAB(c.DL) // Dummy write
	case 5:
		c.PollInterrupts()
		c.Op_RLA()
		c.CompleteInstruction()
	}
}
func (c *CPU) X2F_RLA_Absolute() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.DL = c.ReadFromAB()
	case 4:
		c.WriteToAB(c.DL) // Dummy write
	case 5:
		c.PollInterrupts()
		c.Op_RLA()
		c.CompleteInstruction()
	}
}
func (c *CPU) X3F_RLA_Absolute_X() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(false)
	case 4:
		c.DL = c.ReadFromAB()
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_RLA()
		c.CompleteInstruction()
	}
}
func (c *CPU) X3B_RLA_Absolute_Y() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(false)
	case 4:
		c.DL = c.ReadFromAB()
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_RLA()
		c.CompleteInstruction()
	}
}
func (c *CPU) X23_RLA_Indirect_X() {
	// CPU_Cycles = 8
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectX()
	case 5:
		c.DL = c.ReadFromAB()
	case 6:
		c.WriteToAB(c.DL) // Dummy write
	case 7:
		c.PollInterrupts()
		c.Op_RLA()
		c.CompleteInstruction()
	}
}
func (c *CPU) X33_RLA_Indirect_Y() {
	// CPU_Cycles = 8
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectY(false)
	case 5:
		c.DL = c.ReadFromAB()
	case 6:
		c.WriteToAB(c.DL) // Dummy write
	case 7:
		c.PollInterrupts()
		c.Op_RLA()
		c.CompleteInstruction()
	}
}

//----------------------------------------
// SRE: LSR + EOR
//----------------------------------------

func (c *CPU) X47_SRE_ZeroPage() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.DL = c.ReadFromAB()
	case 3:
		c.WriteToAB(c.DL) // Dummy write
	case 4:
		c.PollInterrupts()
		c.Op_SRE()
		c.CompleteInstruction()
	}
}
func (c *CPU) X57_SRE_ZeroPage_X() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.DL = c.ReadFromAB()
	case 4:
		c.WriteToAB(c.DL) // Dummy write
	case 5:
		c.PollInterrupts()
		c.Op_SRE()
		c.CompleteInstruction()
	}
}
func (c *CPU) X4F_SRE_Absolute() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.DL = c.ReadFromAB()
	case 4:
		c.WriteToAB(c.DL) // Dummy write
	case 5:
		c.PollInterrupts()
		c.Op_SRE()
		c.CompleteInstruction()
	}
}
func (c *CPU) X5F_SRE_Absolute_X() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(false)
	case 4:
		c.DL = c.ReadFromAB()
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_SRE()
		c.CompleteInstruction()
	}
}
func (c *CPU) X5B_SRE_Absolute_Y() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(false)
	case 4:
		c.DL = c.ReadFromAB()
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_SRE()
		c.CompleteInstruction()
	}
}
func (c *CPU) X43_SRE_Indirect_X() {
	// CPU_Cycles = 8
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectX()
	case 5:
		c.DL = c.ReadFromAB()
	case 6:
		c.WriteToAB(c.DL) // Dummy write
	case 7:
		c.PollInterrupts()
		c.Op_SRE()
		c.CompleteInstruction()
	}
}
func (c *CPU) X53_SRE_Indirect_Y() {
	// CPU_Cycles = 8
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectY(false)
	case 5:
		c.DL = c.ReadFromAB()
	case 6:
		c.WriteToAB(c.DL) // Dummy write
	case 7:
		c.PollInterrupts()
		c.Op_SRE()
		c.CompleteInstruction()
	}
}

//----------------------------------------
// RRA: ROR + ADC
//----------------------------------------

func (c *CPU) X67_RRA_ZeroPage() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.DL = c.ReadFromAB()
	case 3:
		c.WriteToAB(c.DL) // Dummy write
	case 4:
		c.PollInterrupts()
		c.Op_RRA()
		c.CompleteInstruction()
	}
}

func (c *CPU) X77_RRA_ZeroPage_X() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.DL = c.ReadFromAB()
	case 4:
		c.WriteToAB(c.DL) // Dummy write
	case 5:
		c.PollInterrupts()
		c.Op_RRA()
		c.CompleteInstruction()
	}
}
func (c *CPU) X6F_RRA_Absolute() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.DL = c.ReadFromAB()
	case 4:
		c.WriteToAB(c.DL) // Dummy write
	case 5:
		c.PollInterrupts()
		c.Op_RRA()
		c.CompleteInstruction()
	}
}
func (c *CPU) X7F_RRA_Absolute_X() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(false)
	case 4:
		c.DL = c.ReadFromAB()
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_RRA()
		c.CompleteInstruction()
	}
}
func (c *CPU) X7B_RRA_Absolute_Y() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(false)
	case 4:
		c.DL = c.ReadFromAB()
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_RRA()
		c.CompleteInstruction()
	}
}
func (c *CPU) X63_RRA_Indirect_X() {
	// CPU_Cycles = 8
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectX()
	case 5:
		c.DL = c.ReadFromAB()
	case 6:
		c.WriteToAB(c.DL) // Dummy write
	case 7:
		c.PollInterrupts()
		c.Op_RRA()
		c.CompleteInstruction()
	}
}
func (c *CPU) X73_RRA_Indirect_Y() {
	// CPU_Cycles = 8
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectY(false)
	case 5:
		c.DL = c.ReadFromAB()
	case 6:
		c.WriteToAB(c.DL) // Dummy write
	case 7:
		c.PollInterrupts()
		c.Op_RRA()
		c.CompleteInstruction()
	}
}

//----------------------------------------
// ISC: INC + SBC
//----------------------------------------

func (c *CPU) XE7_ISC_ZeroPage() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.DL = c.ReadFromAB()
	case 3:
		c.WriteToAB(c.DL) // Dummy write
	case 4:
		c.PollInterrupts()
		c.Op_ISC()
		c.CompleteInstruction()
	}
}
func (c *CPU) XF7_ISC_ZeroPage_X() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.DL = c.ReadFromAB()
	case 4:
		c.WriteToAB(c.DL) // Dummy write
	case 5:
		c.PollInterrupts()
		c.Op_ISC()
		c.CompleteInstruction()
	}
}
func (c *CPU) XEF_ISC_Absolute() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.DL = c.ReadFromAB()
	case 4:
		c.WriteToAB(c.DL) // Dummy write
	case 5:
		c.PollInterrupts()
		c.Op_ISC()
		c.CompleteInstruction()
	}
}
func (c *CPU) XFF_ISC_Absolute_X() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(false)
	case 4:
		c.DL = c.ReadFromAB()
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_ISC()
		c.CompleteInstruction()
	}
}
func (c *CPU) XFB_ISC_Absolute_Y() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(false)
	case 4:
		c.DL = c.ReadFromAB()
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_ISC()
		c.CompleteInstruction()
	}
}
func (c *CPU) XE3_ISC_Indirect_X() {
	// CPU_Cycles = 8
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectX()
	case 5:
		c.DL = c.ReadFromAB()
	case 6:
		c.WriteToAB(c.DL) // Dummy write
	case 7:
		c.PollInterrupts()
		c.Op_ISC()
		c.CompleteInstruction()
	}
}
func (c *CPU) XF3_ISC_Indirect_Y() {
	// CPU_Cycles = 8
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectY(false)
	case 5:
		c.DL = c.ReadFromAB()
	case 6:
		c.WriteToAB(c.DL) // Dummy write
	case 7:
		c.PollInterrupts()
		c.Op_ISC()
		c.CompleteInstruction()
	}
}

//----------------------------------------
// Immediates (unofficial)
//----------------------------------------

func (c *CPU) X0B_ANC_Immediate() { // AND + Set Carry as ASL
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Op_AND(c.ReadFromPC())
	c.flag_Carry = c.flag_Negative
	c.CompleteInstruction()
}
func (c *CPU) X2B_ANC_Immediate() { // AND + Set Carry as ROL (Same as $0B)
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.X0B_ANC_Immediate()
	c.CompleteInstruction()
}
func (c *CPU) X4B_ALR_Immediate() { // AND + LSR
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Op_AND(c.ReadFromPC())
	c.flag_Carry = (c.A & 1) != 0
	c.A >>= 1
	c.SetZNFlags(c.A)
	c.CompleteInstruction()
}
func (c *CPU) X6B_ARR_Immediate() { // AND + ROR
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Op_AND(c.ReadFromPC())
	c.flag_Overflow = c.A == 0

	c.A >>= 1
	if c.flag_Carry {
		c.A |= 0x80
	}
	c.flag_Carry = ((c.A & 0x40) >> 6) == 1
	c.flag_Overflow = (((c.A & 0x20) >> 5) ^ ((c.A & 0x40) >> 6)) == 1

	c.SetZNFlags(c.A)
	c.CompleteInstruction()
}
func (c *CPU) X8B_ANE_Immediate() { // Highly unstable
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.A = c.X & c.ReadFromPC()
	c.SetZNFlags(c.A)
	c.CompleteInstruction()
}
func (c *CPU) XAB_LXA_Immediate() { // Highly unstable
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.A = c.ReadFromPC()
	c.X = c.A
	c.SetZNFlags(c.A)
	c.CompleteInstruction()
}
func (c *CPU) XCB_SBX_Immediate() { // (A AND X) - oper -> X
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.GetAddress_Immediate()
	c.X = (c.A & c.X)
	c.flag_Carry = c.X >= c.DL
	c.X -= c.DL
	c.flag_Zero = (c.X == 0)
	c.flag_Negative = (c.X >= 0x80)
	c.CompleteInstruction()
}
func (c *CPU) XEB_SBC_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Op_SBC(c.ReadFromPC())
	c.CompleteInstruction()
}
