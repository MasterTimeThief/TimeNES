package cpu

/*
	switch cpu.subCycle{
	case 1:
	case 2:
		cpu.PollInterrupts()
		cpu.CompleteInstruction()
	}}
*/

//----------------------------------------
//	NOP Codes (unofficial)
//----------------------------------------

func (cpu *CPU) X1A_NOP_Implied() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.Read(cpu.PC)
	cpu.CompleteInstruction()
}
func (cpu *CPU) X3A_NOP_Implied() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.Read(cpu.PC)
	cpu.CompleteInstruction()
}
func (cpu *CPU) X5A_NOP_Implied() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.Read(cpu.PC)
	cpu.CompleteInstruction()
}
func (cpu *CPU) X7A_NOP_Implied() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.Read(cpu.PC)
	cpu.CompleteInstruction()
}
func (cpu *CPU) XDA_NOP_Implied() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.Read(cpu.PC)
	cpu.CompleteInstruction()
}
func (cpu *CPU) XFA_NOP_Implied() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.Read(cpu.PC)
	cpu.CompleteInstruction()
}

func (cpu *CPU) X80_NOP_Immediate() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.GetAddress_Immediate()
	cpu.CompleteInstruction()
}
func (cpu *CPU) X82_NOP_Immediate() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.GetAddress_Immediate()
	cpu.CompleteInstruction()
}
func (cpu *CPU) X89_NOP_Immediate() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.GetAddress_Immediate()
	cpu.CompleteInstruction()
}
func (cpu *CPU) XC2_NOP_Immediate() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.GetAddress_Immediate()
	cpu.CompleteInstruction()
}
func (cpu *CPU) XE2_NOP_Immediate() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.GetAddress_Immediate()
	cpu.CompleteInstruction()
}

func (cpu *CPU) X04_NOP_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.ReadFromAB() // Dummy Read
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X44_NOP_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.ReadFromAB() // Dummy Read
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X64_NOP_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.ReadFromAB() // Dummy Read
		cpu.CompleteInstruction()
	}
}

func (cpu *CPU) X14_NOP_ZeroPage_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.PollInterrupts()
		cpu.ReadFromAB() // Dummy Read
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X34_NOP_ZeroPage_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.PollInterrupts()
		cpu.ReadFromAB() // Dummy Read
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X54_NOP_ZeroPage_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.PollInterrupts()
		cpu.ReadFromAB() // Dummy Read
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X74_NOP_ZeroPage_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.PollInterrupts()
		cpu.ReadFromAB() // Dummy Read
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XD4_NOP_ZeroPage_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.PollInterrupts()
		cpu.ReadFromAB() // Dummy Read
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XF4_NOP_ZeroPage_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.PollInterrupts()
		cpu.ReadFromAB() // Dummy Read
		cpu.CompleteInstruction()
	}
}

func (cpu *CPU) X0C_NOP_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.ReadFromAB() // Dummy Read
		cpu.CompleteInstruction()
	}
}

func (cpu *CPU) X1C_NOP_Absolute_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_AbsoluteX(true)
	case 3:
		cpu.PollInterrupts()
		cpu.ReadFromAB() // Dummy Read
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X3C_NOP_Absolute_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_AbsoluteX(true)
	case 3:
		cpu.PollInterrupts()
		cpu.ReadFromAB() // Dummy Read
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X5C_NOP_Absolute_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_AbsoluteX(true)
	case 3:
		cpu.PollInterrupts()
		cpu.ReadFromAB() // Dummy Read
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X7C_NOP_Absolute_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_AbsoluteX(true)
	case 3:
		cpu.PollInterrupts()
		cpu.ReadFromAB() // Dummy Read
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XDC_NOP_Absolute_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_AbsoluteX(true)
	case 3:
		cpu.PollInterrupts()
		cpu.ReadFromAB() // Dummy Read
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XFC_NOP_Absolute_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_AbsoluteX(true)
	case 3:
		cpu.PollInterrupts()
		cpu.ReadFromAB() // Dummy Read
		cpu.CompleteInstruction()
	}
}

// SAX: A AND X -> M

func (cpu *CPU) X87_SAX_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.WriteToAB((cpu.A & cpu.X))
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X97_SAX_ZeroPage_Y() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageY()
	case 3:
		cpu.PollInterrupts()
		cpu.WriteToAB((cpu.A & cpu.X))
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X8F_SAX_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.WriteToAB((cpu.A & cpu.X))
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X83_SAX_Indirect_X() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectX()
	case 5:
		cpu.PollInterrupts()
		cpu.WriteToAB((cpu.A & cpu.X))
		cpu.CompleteInstruction()
	}
}

// LAX: LDA + LDX

func (cpu *CPU) XA7_LAX_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.A = cpu.ReadFromAB()
		cpu.X = cpu.A
		cpu.SetZNFlags(cpu.X)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XB7_LAX_ZeroPage_Y() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageY()
	case 3:
		cpu.PollInterrupts()
		cpu.A = cpu.ReadFromAB()
		cpu.X = cpu.A
		cpu.SetZNFlags(cpu.X)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XAF_LAX_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.A = cpu.ReadFromAB()
		cpu.X = cpu.A
		cpu.SetZNFlags(cpu.X)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XBF_LAX_Absolute_Y() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(true)
	case 4:
		cpu.PollInterrupts()
		cpu.A = cpu.ReadFromAB()
		cpu.X = cpu.A
		cpu.SetZNFlags(cpu.X)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XA3_LAX_Indirect_X() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectX()
	case 5:
		cpu.PollInterrupts()
		cpu.A = cpu.ReadFromAB()
		cpu.X = cpu.A
		cpu.SetZNFlags(cpu.X)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XB3_LAX_Indirect_Y() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectY(true)
	case 5:
		cpu.PollInterrupts()
		cpu.A = cpu.ReadFromAB()
		cpu.X = cpu.A
		cpu.SetZNFlags(cpu.X)
		cpu.CompleteInstruction()
	}
}

// SLO: ASL + ORA

func (cpu *CPU) X07_SLO_ZeroPage() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.DL = cpu.ReadFromAB()
	case 3:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 4:
		cpu.PollInterrupts()
		cpu.Op_SLO()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X17_SLO_ZeroPage_X() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.DL = cpu.ReadFromAB()
	case 4:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 5:
		cpu.PollInterrupts()
		cpu.Op_SLO()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X0F_SLO_Absolute() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.DL = cpu.ReadFromAB()
	case 4:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 5:
		cpu.PollInterrupts()
		cpu.Op_SLO()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X1F_SLO_Absolute_X() {
	// CPU_Cycles = 7
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteX(false)
	case 4:
		cpu.DL = cpu.ReadFromAB()
	case 5:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 6:
		cpu.PollInterrupts()
		cpu.Op_SLO()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X1B_SLO_Absolute_Y() {
	// CPU_Cycles = 7
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(false)
	case 4:
		cpu.DL = cpu.ReadFromAB()
	case 5:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 6:
		cpu.PollInterrupts()
		cpu.Op_SLO()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X03_SLO_Indirect_X() {
	// CPU_Cycles = 8
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectX()
	case 5:
		cpu.DL = cpu.ReadFromAB()
	case 6:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 7:
		cpu.PollInterrupts()
		cpu.Op_SLO()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X13_SLO_Indirect_Y() {
	// CPU_Cycles = 8
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectY(false)
	case 5:
		cpu.DL = cpu.ReadFromAB()
	case 6:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 7:
		cpu.PollInterrupts()
		cpu.Op_SLO()
		cpu.CompleteInstruction()
	}
}

// DCP: DEC + CMP

func (cpu *CPU) XC7_DCP_ZeroPage() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.DL = cpu.ReadFromAB()
	case 3:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 4:
		cpu.PollInterrupts()
		cpu.Op_DCP()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XD7_DCP_ZeroPage_X() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.DL = cpu.ReadFromAB()
	case 4:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 5:
		cpu.PollInterrupts()
		cpu.Op_DCP()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XCF_DCP_Absolute() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.DL = cpu.ReadFromAB()
	case 4:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 5:
		cpu.PollInterrupts()
		cpu.Op_DCP()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XDF_DCP_Absolute_X() {
	// CPU_Cycles = 7
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteX(false)
	case 4:
		cpu.DL = cpu.ReadFromAB()
	case 5:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 6:
		cpu.PollInterrupts()
		cpu.Op_DCP()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XDB_DCP_Absolute_Y() {
	// CPU_Cycles = 7
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(false)
	case 4:
		cpu.DL = cpu.ReadFromAB()
	case 5:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 6:
		cpu.PollInterrupts()
		cpu.Op_DCP()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XC3_DCP_Indirect_X() {
	// CPU_Cycles = 8
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectX()
	case 5:
		cpu.DL = cpu.ReadFromAB()
	case 6:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 7:
		cpu.PollInterrupts()
		cpu.Op_DCP()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XD3_DCP_Indirect_Y() {
	// CPU_Cycles = 8
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectY(false)
	case 5:
		cpu.DL = cpu.ReadFromAB()
	case 6:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 7:
		cpu.PollInterrupts()
		cpu.Op_DCP()
		cpu.CompleteInstruction()
	}
}

// SHA: Stores A AND X AND (high-byte of addr. + 1) at addr.

func (cpu *CPU) X9F_SHA_Absolute_Y() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(false)
	case 4:
		cpu.PollInterrupts()
		if (cpu.TempAddress & 0xFF00) != (cpu.AddressBus & 0xFF00) {
			// If the page boundary was crossed, this code has gone "unstable"
			cpu.AddressBus = cpu.AddressBus | ((cpu.AddressBus>>8)&uint16(cpu.X))<<8
		}
		cpu.WriteToAB(cpu.A & (cpu.X | cpu.Magic) & cpu.H)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X93_SHA_Indirect_Y() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectY(false)
	case 5:
		cpu.PollInterrupts()
		if (cpu.TempAddress & 0xFF00) != (cpu.AddressBus & 0xFF00) {
			// If the page boundary was crossed, this code has gone "unstable"
			cpu.AddressBus = cpu.AddressBus | ((cpu.AddressBus>>8)&uint16(cpu.X))<<8
		}
		cpu.WriteToAB(cpu.A & (cpu.X | cpu.Magic) & cpu.H)
		cpu.CompleteInstruction()
	}
}

// SHX: Stores X AND (high-byte of addr. + 1) at addr.

func (cpu *CPU) X9E_SHX_Absolute_Y() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(false)
	case 4:
		cpu.PollInterrupts()
		if (cpu.TempAddress & 0xFF00) != (cpu.AddressBus & 0xFF00) {
			// If the page boundary was crossed, this code has gone "unstable"
			cpu.AddressBus = cpu.AddressBus | ((cpu.AddressBus>>8)&uint16(cpu.X))<<8
		}
		cpu.WriteToAB(cpu.X & cpu.H)
		cpu.CompleteInstruction()
	}
}

// SHY: Stores Y AND (high-byte of addr. + 1) at addr.

func (cpu *CPU) X9C_SHY_Absolute_X() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteX(false)
	case 4:
		cpu.PollInterrupts()
		if (cpu.TempAddress & 0xFF00) != (cpu.AddressBus & 0xFF00) {
			// If the page boundary was crossed, this code has gone "unstable"
			cpu.AddressBus = cpu.AddressBus | ((cpu.AddressBus>>8)&uint16(cpu.Y))<<8
		}
		cpu.WriteToAB(cpu.Y & cpu.H)
		cpu.CompleteInstruction()
	}
}

// TAS (XAS, SHS): Puts A AND X in SP and stores A AND X AND (high-byte of addr. + 1) at addr.
// A AND X -> SP, A AND X AND (H+1) -> M

func (cpu *CPU) X9B_TAS_Absolute_Y() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(false)
	case 4:
		cpu.PollInterrupts()
		if (cpu.TempAddress & 0xFF00) != (cpu.AddressBus & 0xFF00) {
			// If the page boundary was crossed, this code has gone "unstable"
			cpu.AddressBus = cpu.AddressBus | ((cpu.AddressBus>>8)&uint16(cpu.Y))<<8
		}
		cpu.Push(cpu.A & cpu.X)
		cpu.WriteToAB(cpu.A & (cpu.X | cpu.Magic) & cpu.H)
		cpu.CompleteInstruction()
	}
}

// LAS (LAR, LAE): LDA/TSX oper
// M AND SP -> A, X, SP

func (cpu *CPU) XBB_LAS_Absolute_Y() {
	// CPU_Cycles = 4+
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(true)
	case 4:
		cpu.PollInterrupts()
		cpu.DL = (cpu.ReadFromAB() & cpu.SP)
		cpu.A = cpu.DL
		cpu.X = cpu.DL
		cpu.SP = cpu.DL
		cpu.SetZNFlags(cpu.A)
		cpu.CompleteInstruction()
	}
}

// RLA: ROL + AND

func (cpu *CPU) X27_RLA_ZeroPage() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.DL = cpu.ReadFromAB()
	case 3:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 4:
		cpu.PollInterrupts()
		cpu.Op_RLA()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X37_RLA_ZeroPage_X() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.DL = cpu.ReadFromAB()
	case 4:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 5:
		cpu.PollInterrupts()
		cpu.Op_RLA()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X2F_RLA_Absolute() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.DL = cpu.ReadFromAB()
	case 4:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 5:
		cpu.PollInterrupts()
		cpu.Op_RLA()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X3F_RLA_Absolute_X() {
	// CPU_Cycles = 7
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteX(false)
	case 4:
		cpu.DL = cpu.ReadFromAB()
	case 5:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 6:
		cpu.PollInterrupts()
		cpu.Op_RLA()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X3B_RLA_Absolute_Y() {
	// CPU_Cycles = 7
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(false)
	case 4:
		cpu.DL = cpu.ReadFromAB()
	case 5:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 6:
		cpu.PollInterrupts()
		cpu.Op_RLA()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X23_RLA_Indirect_X() {
	// CPU_Cycles = 8
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectX()
	case 5:
		cpu.DL = cpu.ReadFromAB()
	case 6:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 7:
		cpu.PollInterrupts()
		cpu.Op_RLA()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X33_RLA_Indirect_Y() {
	// CPU_Cycles = 8
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectY(false)
	case 5:
		cpu.DL = cpu.ReadFromAB()
	case 6:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 7:
		cpu.PollInterrupts()
		cpu.Op_RLA()
		cpu.CompleteInstruction()
	}
}

// SRE: LSR + EOR

func (cpu *CPU) X47_SRE_ZeroPage() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.DL = cpu.ReadFromAB()
	case 3:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 4:
		cpu.PollInterrupts()
		cpu.Op_SRE()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X57_SRE_ZeroPage_X() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.DL = cpu.ReadFromAB()
	case 4:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 5:
		cpu.PollInterrupts()
		cpu.Op_SRE()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X4F_SRE_Absolute() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.DL = cpu.ReadFromAB()
	case 4:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 5:
		cpu.PollInterrupts()
		cpu.Op_SRE()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X5F_SRE_Absolute_X() {
	// CPU_Cycles = 7
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteX(false)
	case 4:
		cpu.DL = cpu.ReadFromAB()
	case 5:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 6:
		cpu.PollInterrupts()
		cpu.Op_SRE()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X5B_SRE_Absolute_Y() {
	// CPU_Cycles = 7
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(false)
	case 4:
		cpu.DL = cpu.ReadFromAB()
	case 5:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 6:
		cpu.PollInterrupts()
		cpu.Op_SRE()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X43_SRE_Indirect_X() {
	// CPU_Cycles = 8
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectX()
	case 5:
		cpu.DL = cpu.ReadFromAB()
	case 6:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 7:
		cpu.PollInterrupts()
		cpu.Op_SRE()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X53_SRE_Indirect_Y() {
	// CPU_Cycles = 8
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectY(false)
	case 5:
		cpu.DL = cpu.ReadFromAB()
	case 6:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 7:
		cpu.PollInterrupts()
		cpu.Op_SRE()
		cpu.CompleteInstruction()
	}
}

/*
	// RRA: ROR + ADC

	func (cpu *CPU) X67_RRA_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		Op_ROR(AddressBus)
		Op_ADC(cpu.ReadFromAB())
		// CPU_Cycles = 5
	func (cpu *CPU) X77_RRA_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		Op_ROR(AddressBus)
		Op_ADC(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) X6F_RRA_Absolute(){
		cpu.GetAddress_Absolute()
		Op_ROR(AddressBus)
		Op_ADC(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) X7F_RRA_Absolute_X(){
		cpu.GetAddress_AbsoluteX(false)
		Op_ROR(AddressBus)
		Op_ADC(cpu.ReadFromAB())
		// CPU_Cycles = 7
	func (cpu *CPU) X7B_RRA_Absolute_Y(){
		cpu.GetAddress_AbsoluteY(false)
		Op_ROR(AddressBus)
		Op_ADC(cpu.ReadFromAB())
		// CPU_Cycles = 7
	func (cpu *CPU) X63_RRA_Indirect_X(){
		cpu.GetAddress_IndirectX()
		Op_ROR(AddressBus)
		Op_ADC(cpu.ReadFromAB())
		// CPU_Cycles = 8
	func (cpu *CPU) X73_RRA_Indirect_Y(){
		cpu.GetAddress_IndirectY(false)
		Op_ROR(AddressBus)
		Op_ADC(cpu.ReadFromAB())
		// CPU_Cycles = 8
/*
	// ISC: INC + SBC

	func (cpu *CPU) XE7_ISC_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		// CPU_Cycles = 5
	func (cpu *CPU) XF7_ISC_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) XEF_ISC_Absolute(){
		cpu.GetAddress_Absolute()
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) XFF_ISC_Absolute_X(){
		cpu.GetAddress_AbsoluteX(false)
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		// CPU_Cycles = 7
	func (cpu *CPU) XFB_ISC_Absolute_Y(){
		cpu.GetAddress_AbsoluteY(false)
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		// CPU_Cycles = 7
	func (cpu *CPU) XE3_ISC_Indirect_X(){
		cpu.GetAddress_IndirectX()
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		// CPU_Cycles = 8
	func (cpu *CPU) XF3_ISC_Indirect_Y(){
		cpu.GetAddress_IndirectY(false)
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		// CPU_Cycles = 8
/*
	// Immediates (unofficial)

	func (cpu *CPU) X0BANC_Immediate(){
		//AND + Set Carry as ASL
		cpu.Op_AND(cpu.ReadFromPC())
		flag_Carry = flag_Negative
		// CPU_Cycles = 2
	func (cpu *CPU) X2BANC2_Immediate(){
		//AND + Set Carry as ROL
		//Same as $0B
		cpu.Op_AND(cpu.ReadFromPC())
		flag_Carry = flag_Negative
		// CPU_Cycles = 2
	func (cpu *CPU) X4BALR_Immediate(){
		//AND + LSR
		cpu.Op_AND(cpu.ReadFromPC())
		flag_Carry = (A & 1) != 0
		A >>= 1
		cpu.SetZNFlags(cpu.A)
		// CPU_Cycles = 2
	func (cpu *CPU) X6BARR_Immediate(){
		//AND + ROR
		cpu.Op_AND(cpu.ReadFromPC())
		flag_Overflow = A == 0

		A >>= 1
		if flag_Carry {
			A |= 0x80
		}
		flag_Carry = ((A & 0x40) >> 6) == 1
		flag_Overflow = (((A & 0x20) >> 5) ^ ((A & 0x40) >> 6)) == 1

		cpu.SetZNFlags(cpu.A)
		// CPU_Cycles = 2
	func (cpu *CPU) X8BANE_Immediate(){
		//Highly unstable
		A = X & cpu.ReadFromPC()
		cpu.SetZNFlags(cpu.A)
		// CPU_Cycles = 2
	func (cpu *CPU) XABLXA_Immediate(){
		//Highly unstable
		cpu.A = cpu.ReadFromPC()
		cpu.X = cpu.A
		cpu.SetZNFlags(cpu.A)
		// CPU_Cycles = 2
	func (cpu *CPU) XCBSBX_Immediate(){
		//(A AND X) - oper -> X
		X = (cpu.A & cpu.X) - cpu.ReadFromPC()
		Op_CMP(cpu.X)
		cpu.SetZNFlags(cpu.X)
		// CPU_Cycles = 2
	func (cpu *CPU) XEBSBC_Immediate(){
		cpu.Op_SBC(cpu.ReadFromPC())
		// CPU_Cycles = 2

	default:
		fmt.Println("Unknown Opcode: " + fmt.Sprintf("%02X", opcode))

*/
