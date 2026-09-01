package cpu

//----------------------------------------
//	Access
//----------------------------------------

//	STA: Store Accumulator in Memory
//	A -> M
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) X85_STA_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.WriteToAB(cpu.A)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X95_STA_ZeroPage_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.PollInterrupts()
		cpu.WriteToAB(cpu.A)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X8D_STA_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.WriteToAB(cpu.A)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X9D_STA_Absolute_X() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteX(false)
	case 4:
		cpu.PollInterrupts()
		cpu.WriteToAB(cpu.A)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X99_STA_Absolute_Y() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(false)
	case 4:
		cpu.PollInterrupts()
		cpu.WriteToAB(cpu.A)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X81_STA_Indirect_X() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectX()
	case 5:
		cpu.PollInterrupts()
		cpu.WriteToAB(cpu.A)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X91_STA_Indirect_Y() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectY(false)
	case 5:
		cpu.PollInterrupts()
		cpu.WriteToAB(cpu.A)
		cpu.CompleteInstruction()
	}
}

//	LDA: Load Accumulator with Memory
//	M -> A
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) XA9_LDA_Immediate() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.A = cpu.ReadFromPC()
	cpu.SetZNFlags(cpu.A)
	cpu.CompleteInstruction()
}
func (cpu *CPU) XA5_LDA_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.A = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.A)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XB5_LDA_ZeroPage_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.PollInterrupts()
		cpu.A = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.A)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XAD_LDA_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.A = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.A)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XBD_LDA_Absolute_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteX(true)
	case 4:
		cpu.PollInterrupts()
		cpu.A = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.A)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XB9_LDA_Absolute_Y() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(true)
	case 4:
		cpu.PollInterrupts()
		cpu.A = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.A)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XA1_LDA_Indirect_X() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectX()
	case 5:
		cpu.PollInterrupts()
		cpu.A = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.A)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XB1_LDA_Indirect_Y() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectY(true)
	case 5:
		cpu.PollInterrupts()
		cpu.A = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.A)
		cpu.CompleteInstruction()
	}
}

//	STX: Store Index X in Memory
//	X -> M
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) X86_STX_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.WriteToAB(cpu.X)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X96_STX_ZeroPage_Y() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageY()
	case 3:
		cpu.PollInterrupts()
		cpu.WriteToAB(cpu.X)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X8E_STX_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.WriteToAB(cpu.X)
		cpu.CompleteInstruction()
	}
}

//	LDX: Load Index X with Memory
//	M -> X
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) XA2_LDX_Immediate() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.X = cpu.ReadFromPC()
	cpu.SetZNFlags(cpu.X)
	cpu.CompleteInstruction()
}
func (cpu *CPU) XA6_LDX_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.X = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.X)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XAE_LDX_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.X = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.X)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XB6_LDX_ZeroPage_Y() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageY()
	case 3:
		cpu.PollInterrupts()
		cpu.X = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.X)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XBE_LDX_Absolute_Y() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(true)
	case 4:
		cpu.PollInterrupts()
		cpu.X = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.X)
		cpu.CompleteInstruction()
	}
}

//	STY: Store Index Y in Memory
//	Y -> M
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) X84_STY_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.WriteToAB(cpu.Y)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X94_STY_ZeroPage_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.PollInterrupts()
		cpu.WriteToAB(cpu.Y)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X8C_STY_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.WriteToAB(cpu.Y)
		cpu.CompleteInstruction()
	}
}

//	LDY: Load Index Y with Memory
//	M -> Y
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) XA0_LDY_Immediate() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.Y = cpu.ReadFromPC()
	cpu.SetZNFlags(cpu.Y)
	cpu.CompleteInstruction()
}
func (cpu *CPU) XA4_LDY_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.Y = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.Y)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XAC_LDY_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.Y = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.Y)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XB4_LDY_ZeroPage_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.PollInterrupts()
		cpu.Y = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.Y)
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XBC_LDY_Absolute_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteX(true)
	case 4:
		cpu.PollInterrupts()
		cpu.Y = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.Y)
		cpu.CompleteInstruction()
	}
}

//----------------------------------------
//	Transfer
//----------------------------------------

//	TAX: Transfer Accumulator to Index X
//	A -> X
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) XAA_TAX() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.X = cpu.A
	cpu.SetZNFlags(cpu.X)
	cpu.CompleteInstruction()
}

//	TAY: Transfer Accumulator to Index Y
//	A -> Y
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) XA8_TAY() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.Y = cpu.A
	cpu.SetZNFlags(cpu.Y)
	cpu.CompleteInstruction()
}

//	TXA: Transfer Index X to Accumulator
//	X -> A
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) X8A_TXA() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.A = cpu.X
	cpu.SetZNFlags(cpu.A)
	cpu.CompleteInstruction()
}

//	TYA: Transfer Index Y to Accumulator
//	Y -> A
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) X98_TYA() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.A = cpu.Y
	cpu.SetZNFlags(cpu.A)
	cpu.CompleteInstruction()
}

//----------------------------------------
//	Arithmetic
//----------------------------------------

//	ADC: Add Memory to Accumulator with Carry
//	A + M + C -> A, C
//	N	Z	C	I	D	V
//	+	+	+	-	-	+

func (cpu *CPU) X69_ADC_Immediate() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.Op_ADC(cpu.ReadFromPC())
	cpu.CompleteInstruction()
}
func (cpu *CPU) X65_ADC_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.Op_ADC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X75_ADC_ZeroPage_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.PollInterrupts()
		cpu.Op_ADC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X6D_ADC_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.Op_ADC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X7D_ADC_Absolute_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteX(true)
	case 4:
		cpu.PollInterrupts()
		cpu.Op_ADC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X79_ADC_Absolute_Y() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(true)
	case 4:
		cpu.PollInterrupts()
		cpu.Op_ADC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X61_ADC_Indirect_X() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectX()
	case 5:
		cpu.PollInterrupts()
		cpu.Op_ADC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X71_ADC_Indirect_Y() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectY(true)
	case 5:
		cpu.PollInterrupts()
		cpu.Op_ADC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}

//	SBC: Subtract Memory from Accumulator with Borrow
//	A - M - C̅ -> A
//	N	Z	C	I	D	V
//	+	+	+	-	-	+

func (cpu *CPU) XE9_SBC_Immediate() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.Op_SBC(cpu.ReadFromPC())
	cpu.CompleteInstruction()
}
func (cpu *CPU) XE5_SBC_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.Op_SBC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XED_SBC_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.Op_SBC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XF5_SBC_ZeroPage_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.PollInterrupts()
		cpu.Op_SBC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XFD_SBC_Absolute_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteX(true)
	case 4:
		cpu.PollInterrupts()
		cpu.Op_SBC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XF9_SBC_Absolute_Y() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(true)
	case 4:
		cpu.PollInterrupts()
		cpu.Op_SBC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XE1_SBC_Indirect_X() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectX()
	case 5:
		cpu.PollInterrupts()
		cpu.Op_SBC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XF1_SBC_Indirect_Y() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectY(true)
	case 5:
		cpu.PollInterrupts()
		cpu.Op_SBC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}

//	INC: Increment Memory by One
//	M + 1 -> M
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) XE6_INC_ZeroPage() {
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
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XEE_INC_Absolute() {
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
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XF6_INC_ZeroPage_X() {
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
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XFE_INC_Absolute_X() {
	// CPU_Cycles = 7
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_AbsoluteX(false)
	case 5:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 6:
		cpu.PollInterrupts()
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}

//	DEC: Decrement Memory by One
//	M - 1 -> M
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) XC6_DEC_ZeroPage() {
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
		cpu.Op_DEC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XCE_DEC_Absolute() {
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
		cpu.Op_DEC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XD6_DEC_ZeroPage_X() {
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
		cpu.Op_DEC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XDE_DEC_Absolute_X() {
	// CPU_Cycles = 7
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_AbsoluteX(false)
	case 5:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 6:
		cpu.PollInterrupts()
		cpu.Op_DEC(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}

//	INX: Increment Index X by One
//	X + 1 -> X
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) XE8_INX() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.ReadFromAB() // Dummy Read
	cpu.X++
	cpu.SetZNFlags(cpu.X)
	cpu.CompleteInstruction()
}

//	DEX: Decrement Index X by One
//	X - 1 -> X
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) XCA_DEX() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.ReadFromAB() // Dummy Read
	cpu.X--
	cpu.SetZNFlags(cpu.X)
	cpu.CompleteInstruction()
}

//	INY: Increment Index Y by One
//	Y + 1 -> Y
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) XC8_INY() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.ReadFromAB() // Dummy Read
	cpu.Y++
	cpu.SetZNFlags(cpu.Y)
	cpu.CompleteInstruction()
}

//	DEY: Decrement Index Y by One
//	Y - 1 -> Y
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) X88_DEY() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.ReadFromAB() // Dummy Read
	cpu.Y--
	cpu.SetZNFlags(cpu.Y)
	cpu.CompleteInstruction()
}

//----------------------------------------
//	Shift
//----------------------------------------

//	ASL: Shift Left One Bit (Memory or Accumulator)
//	C <- [76543210] <- 0
//	N	Z	C	I	D	V
//	+	+	+	-	-	-

func (cpu *CPU) X0A_ASL() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.flag_Carry = cpu.A > 127
	cpu.A <<= 1
	cpu.SetZNFlags(cpu.A)
	cpu.CompleteInstruction()
}
func (cpu *CPU) X06_ASL_ZeroPage() {
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
		cpu.Op_ASL()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X0E_ASL_Absolute() {
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
		cpu.Op_ASL()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X16_ASL_ZeroPage_X() {
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
		cpu.Op_ASL()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X1E_ASL_Absolute_X() {
	// CPU_Cycles = 7
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_AbsoluteX(false)
	case 5:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 6:
		cpu.PollInterrupts()
		cpu.Op_ASL()
		cpu.CompleteInstruction()
	}
}

//	LSR: Shift One Bit Right (Memory or Accumulator)
//	0 -> [76543210] -> C
//	N	Z	C	I	D	V
//	0	+	+	-	-	-

func (cpu *CPU) X4A_LSR() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.flag_Carry = (cpu.A & 1) != 0
	cpu.A >>= 1
	cpu.SetZNFlags(cpu.A)
	cpu.CompleteInstruction()
}
func (cpu *CPU) X46_LSR_ZeroPage() {
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
		cpu.Op_LSR()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X4E_LSR_Absolute() {
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
		cpu.Op_LSR()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X56_LSR_ZeroPage_X() {
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
		cpu.Op_LSR()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X5E_LSR_Absolute_X() {
	// CPU_Cycles = 7
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_AbsoluteX(false)
	case 5:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 6:
		cpu.PollInterrupts()
		cpu.Op_LSR()
		cpu.CompleteInstruction()
	}
}

//	ROL: Rotate One Bit Left (Memory or Accumulator)
//	C <- [76543210] <- C
//	N	Z	C	I	D	V
//	+	+	+	-	-	-

func (cpu *CPU) X2A_ROL() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	futureCarry := (cpu.A >= 0x80)
	cpu.A <<= 1
	if cpu.flag_Carry {
		cpu.A |= 1
	}
	cpu.flag_Carry = futureCarry
	cpu.SetZNFlags(cpu.A)
	cpu.CompleteInstruction()
}
func (cpu *CPU) X26_ROL_ZeroPage() {
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
		cpu.Op_ROL()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X2E_ROL_Absolute() {
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
		cpu.Op_ROL()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X36_ROL_ZeroPage_X() {
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
		cpu.Op_ROL()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X3E_ROL_Absolute_X() {
	// CPU_Cycles = 7
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_AbsoluteX(false)
	case 5:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 6:
		cpu.PollInterrupts()
		cpu.Op_ROL()
		cpu.CompleteInstruction()
	}
}

//	ROR: Rotate One Bit Right (Memory or Accumulator)
//	C -> [76543210] -> C
//	N	Z	C	I	D	V
//	+	+	+	-	-	-

func (cpu *CPU) X6A_ROR() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	futureCarry := (cpu.A & 1) != 0
	cpu.A >>= 1
	if cpu.flag_Carry {
		cpu.A |= 0x80
	}
	cpu.flag_Carry = futureCarry
	cpu.SetZNFlags(cpu.A)
	cpu.CompleteInstruction()
}
func (cpu *CPU) X66_ROR_ZeroPage() {
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
		cpu.Op_ROR()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X6E_ROR_Absolute() {
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
		cpu.Op_ROR()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X76_ROR_ZeroPage_X() {
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
		cpu.Op_ROR()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X7E_ROR_Absolute_X() {
	// CPU_Cycles = 7
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_AbsoluteX(false)
	case 5:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 6:
		cpu.PollInterrupts()
		cpu.Op_ROR()
		cpu.CompleteInstruction()
	}
}

//----------------------------------------
//	Bitwise
//----------------------------------------

//	AND: AND Memory with Accumulator
//	A AND M -> A
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) X29_AND_Immediate() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.Op_AND(cpu.ReadFromPC())
	cpu.CompleteInstruction()
}
func (cpu *CPU) X25_AND_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.Op_AND(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X2D_AND_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.Op_AND(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X35_AND_ZeroPage_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.PollInterrupts()
		cpu.Op_AND(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X3D_AND_Absolute_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteX(true)
	case 4:
		cpu.PollInterrupts()
		cpu.Op_AND(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X39_AND_Absolute_Y() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(true)
	case 4:
		cpu.PollInterrupts()
		cpu.Op_AND(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X21_AND_Indirect_X() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectX()
	case 5:
		cpu.PollInterrupts()
		cpu.Op_AND(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X31_AND_Indirect_Y() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectY(true)
	case 5:
		cpu.PollInterrupts()
		cpu.Op_AND(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}

//	ORA: OR Memory with Accumulator
//	A OR M -> A
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) X09_ORA_Immediate() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.Op_ORA(cpu.ReadFromPC())
	cpu.CompleteInstruction()
}
func (cpu *CPU) X05_ORA_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.Op_ORA(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X0D_ORA_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.Op_ORA(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X15_ORA_ZeroPage_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.PollInterrupts()
		cpu.Op_ORA(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X1D_ORA_Absolute_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteX(true)
	case 4:
		cpu.PollInterrupts()
		cpu.Op_ORA(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X19_ORA_Absolute_Y() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(true)
	case 4:
		cpu.PollInterrupts()
		cpu.Op_ORA(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X01_ORA_Indirect_X() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectX()
	case 5:
		cpu.PollInterrupts()
		cpu.Op_ORA(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X11_ORA_Indirect_Y() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectY(true)
	case 5:
		cpu.PollInterrupts()
		cpu.Op_ORA(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}

//	EOR: Exclusive-OR Memory with Accumulator
//	A EOR M -> A
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) X49_EOR_Immediate() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.Op_EOR(cpu.ReadFromPC())
	cpu.CompleteInstruction()
}
func (cpu *CPU) X45_EOR_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.Op_EOR(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X4D_EOR_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.Op_EOR(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X55_EOR_ZeroPage_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.PollInterrupts()
		cpu.Op_EOR(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X5D_EOR_Absolute_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteX(true)
	case 4:
		cpu.PollInterrupts()
		cpu.Op_EOR(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X59_EOR_Absolute_Y() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(true)
	case 4:
		cpu.PollInterrupts()
		cpu.Op_EOR(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X41_EOR_Indirect_X() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectX()
	case 5:
		cpu.PollInterrupts()
		cpu.Op_EOR(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X51_EOR_Indirect_Y() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectY(true)
	case 5:
		cpu.PollInterrupts()
		cpu.Op_EOR(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}

//	BIT: Test Bits in Memory with Accumulator
//	bits 7 and 6 of operand are transfered to bit 7 and 6 of SR (N,V);
//	the zero-flag is set according to the result of the operand AND
//	the accumulator (set, if the result is zero, unset otherwise).
//	This allows a quick check of a few bits at once without affecting
//	any of the registers, other than the status register (SR).
//
//	A AND M -> Z, M[7] -> N, M[6] -> V
//
//	N	Z	C	I	D	V
//	M7	+	-	-	-	M6

func (cpu *CPU) X24_BIT_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.Op_BIT(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X2C_BIT_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.Op_BIT(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}

//----------------------------------------
//	Compare
//----------------------------------------

//	CMP: Compare Memory with Accumulator
//	A - M
//	N	Z	C	I	D	V
//	+	+	+	-	-	-

func (cpu *CPU) XC9_CMP_Immediate() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.Op_CMP(cpu.ReadFromPC())
	cpu.CompleteInstruction()
}
func (cpu *CPU) XC5_CMP_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.Op_CMP(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XCD_CMP_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.Op_CMP(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XD5_CMP_ZeroPage_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_ZeroPageX()
	case 3:
		cpu.PollInterrupts()
		cpu.Op_CMP(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XDD_CMP_Absolute_X() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteX(true)
	case 4:
		cpu.PollInterrupts()
		cpu.Op_CMP(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XD9_CMP_Absolute_Y() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.GetAddress_AbsoluteY(true)
	case 4:
		cpu.PollInterrupts()
		cpu.Op_CMP(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XC1_CMP_Indirect_X() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectX()
	case 5:
		cpu.PollInterrupts()
		cpu.Op_CMP(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XD1_CMP_Indirect_Y() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_IndirectY(true)
	case 5:
		cpu.PollInterrupts()
		cpu.Op_CMP(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}

//	CPX: Compare Memory and Index X
//	X - M
//	N	Z	C	I	D	V
//	+	+	+	-	-	-

func (cpu *CPU) XE0_CPX_Immediate() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.Op_CPX(cpu.ReadFromPC())
	cpu.CompleteInstruction()
}
func (cpu *CPU) XE4_CPX_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.Op_CPX(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XEC_CPX_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.Op_CPX(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}

//	CPY: Compare Memory and Index Y
//	Y - M
//	N	Z	C	I	D	V
//	+	+	+	-	-	-

func (cpu *CPU) XC0_CPY_Immediate() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.Op_CPY(cpu.ReadFromPC())
	cpu.CompleteInstruction()
}
func (cpu *CPU) XC4_CPY_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
		cpu.PollInterrupts()
		cpu.Op_CPY(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) XCC_CPY_Absolute() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.PollInterrupts()
		cpu.Op_CPY(cpu.ReadFromAB())
		cpu.CompleteInstruction()
	}
}

//----------------------------------------
//	Branch
//----------------------------------------

//	BCC: Branch on Carry Clear
//	branch on C = 0
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) X90_BCC() {
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.Branch(!cpu.flag_Carry)
	}
}

//	BCS: Branch on Carry Set
//	branch on C = 1
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) XB0_BCS() {
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.Branch(cpu.flag_Carry)
	}
}

//	BEQ: Branch on Result Zero
//	branch on Z = 1
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) XF0_BEQ() {
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.Branch(cpu.flag_Zero)
	}
}

//	BNE: Branch on Result not Zero
//	branch on Z = 0
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) XD0_BNE() {
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.Branch(!cpu.flag_Zero)
	}
}

//	BPL: Branch on Result Plus
//	branch on N = 0
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) X10_BPL() {
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.Branch(!cpu.flag_Negative)
	}
}

//	BMI: Branch on Result Minus
//	branch on N = 1
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) X30_BMI() {
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.Branch(cpu.flag_Negative)
	}
}

//	BVC: Branch on Overflow Clear
//	branch on V = 0
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) X50_BVC() {
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.Branch(!cpu.flag_Overflow)
	}
}

//	BVS: Branch on Overflow Set
//	branch on V = 1
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) X70_BVS() {
	switch cpu.subCycle {
	case 1, 2, 3:
		cpu.Branch(cpu.flag_Overflow)
	}
}

//----------------------------------------
//	Jump
//----------------------------------------

//	JMP: Jump to New Location
//	operand 1st byte -> PCL
//	operand 2nd byte -> PCH
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) X4C_JMP() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_Absolute()
	case 2:
		cpu.PollInterrupts()
		cpu.GetAddress_Absolute()
		cpu.PC = cpu.AddressBus
		cpu.CompleteInstruction()
	}
}

func (cpu *CPU) X6C_JMP_Indirect() {
	// CPU_Cycles = 5
	switch cpu.subCycle {
	case 1, 2:
		cpu.GetAddress_Absolute()
	case 3:
		cpu.SB = cpu.ReadFromAB()
	case 4:
		cpu.PollInterrupts()
		cpu.DL = cpu.bus.Read((cpu.AddressBus & 0xFF00) | ((cpu.AddressBus + 1) & 0xFF))
		cpu.PC = cpu.BuildAddress(cpu.SB, cpu.DL)
		cpu.CompleteInstruction()
	}
}

//	JSR: Jump to New Location Saving Return Address
//	push (PC+2),
//	operand 1st byte -> PCL
//	operand 2nd byte -> PCH
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) X20_JSR() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1: // fetch the byte that will be PC low
		cpu.GetAddress_Immediate()
	case 2: // transfer stack pointer to address bus, and alu to stack pointer. I'm just reusing `dl` here, but this instruction actually uses the Arithmetic Logic Unit for this.
		cpu.AddressBus = uint16(cpu.SP) | 0x100
		//cpu.SP = cpu.DL
		cpu.ReadFromAB() // Dummy Read
	case 3: // push PC high to stack via address bus
		cpu.Push(byte(cpu.PC / 0x100))
	case 4: // push PC low to stack via address bus
		cpu.Push(byte(cpu.PC))
	case 5: // fetch PC High, transfer stack pointer to PC low, address bus to stack pointer.
		cpu.PollInterrupts()
		cpu.PC = cpu.BuildAddress(cpu.DL, cpu.ReadFromPC())
		cpu.CompleteInstruction()
	}
}

//	RTS: Return from Subroutine
//	pull PC, PC+1 -> PC
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) X60_RTS() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_Immediate() // Dummy read to advance PC
	case 2:
		cpu.SP--
		cpu.Pull() // Dummy Read (Pull)
	case 3: // Target byte low
		cpu.DL = cpu.Pull()
	case 4: // Target byte high
		cpu.PC = cpu.BuildAddress(cpu.DL, cpu.Pull())
	case 5:
		cpu.PollInterrupts()
		cpu.ReadFromPC() // Dummy read to advance PC
		cpu.CompleteInstruction()
	}
}

//	BRK: Force Break
//	BRK initiates a software interrupt similar to a hardware
//	interrupt (IRQ). The return address pushed to the stack is
//	PC+2, providing an extra byte of spacing for a break mark
//	(identifying a reason for the break.)
//	The status register will be pushed to the stack with the break
//	flag set to 1. However, when retrieved during RTI or by a PLP
//	instruction, the break flag will be ignored.
//	The interrupt disable flag is not set automatically.
//
//	interrupt,
//	push PC+2, push SR
//	N	Z	C	I	D	V
//	-	-	-	1	-	-

func (cpu *CPU) X00_BRK() {
	// CPU_Cycles = 7
	switch cpu.subCycle {
	case 1:
		if cpu.BreakSource == Break_Software {
			cpu.GetAddress_Immediate() // Dummy read that increments PC
		} else {
			cpu.ReadFromAB() // Dummy read that does not increment PC
		}
	case 2:
		if cpu.BreakSource != Break_Reset {
			cpu.Push(byte(cpu.PC >> 8))
		} else {
			cpu.ResetReadPush()
		}
	case 3:
		if cpu.BreakSource != Break_Reset {
			cpu.Push(byte(cpu.PC))
		} else {
			cpu.ResetReadPush()
		}
	case 4:
		if cpu.BreakSource != Break_Reset {
			cpu.flag_B = false
			if cpu.BreakSource == Break_Software {
				cpu.flag_B = true
			}
			cpu.PushFlags()
		} else {
			cpu.ResetReadPush()
		}
		cpu.PollInterrupts()
	case 5:
		if cpu.BreakSource == Break_NMI {
			cpu.PC = (cpu.PC & 0xFF00) | uint16(cpu.bus.Read(0xFFFA))
		} else if cpu.BreakSource == Break_Reset {
			cpu.PC = (cpu.PC & 0xFF00) | uint16(cpu.bus.Read(0xFFFC))
		} else {
			cpu.PC = (cpu.PC & 0xFF00) | uint16(cpu.bus.Read(0xFFFE))
		}
	case 6:
		if cpu.BreakSource == Break_NMI {
			cpu.PC = (cpu.PC & 0xFF) | (uint16(cpu.bus.Read(0xFFFB)) << 8)
		} else if cpu.BreakSource == Break_Reset {
			cpu.PC = (cpu.PC & 0xFF) | (uint16(cpu.bus.Read(0xFFFD)) << 8)
		} else {
			cpu.PC = (cpu.PC & 0xFF) | (uint16(cpu.bus.Read(0xFFFF)) << 8)
		}
		cpu.BreakSource = Break_None
		cpu.flag_InterruptDisable = true
		cpu.CompleteInstruction()
	}
}

//	RTI: Return from Interrupt
//	The status register is pulled with the break flag
//	and bit 5 ignored. Then PC is pulled from the stack.
//
//	pull SR, pull PC
//	N	Z	C	I	D	V
//	from stack

func (cpu *CPU) X40_RTI() {
	// CPU_Cycles = 6
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_Immediate()
	case 2:
		cpu.AddressBus = uint16(cpu.SP) | 0x100
		cpu.ReadFromAB() //Dummy read
	case 3:
		status := cpu.Pull()
		cpu.flag_Carry = (status & 1) != 0
		cpu.flag_Zero = (status & 2) != 0
		cpu.flag_InterruptDisable = (status & 4) != 0
		cpu.flag_Decimal = (status & 8) != 0
		cpu.flag_Overflow = (status & 64) != 0
		cpu.flag_Negative = (status & 128) != 0
	case 4:
		cpu.DL = cpu.Pull()
	case 5:
		cpu.PollInterrupts()
		cpu.PC = cpu.BuildAddress(cpu.DL, cpu.Pull())
		cpu.CompleteInstruction()
	}
}

//----------------------------------------
//	Stack
//----------------------------------------

//	PHA: Push Accumulator on Stack
//	push A
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) X48_PHA() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.DL = cpu.ReadFromAB() // Dummy read
	case 2:
		cpu.PollInterrupts()
		cpu.Push(cpu.A)
		cpu.CompleteInstruction()
	}
}

//	PLA: Pull Accumulator from Stack
//	pull A
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) X68_PLA() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1:
		cpu.AddressBus = cpu.PC
		cpu.ReadFromAB() // Dummy Read
	case 2:
		cpu.AddressBus = uint16(cpu.SP) | 0x100
		cpu.ReadFromAB() //Dummy read
	case 3:
		cpu.PollInterrupts()
		cpu.A = cpu.Pull()
		cpu.SetZNFlags(cpu.A)
		cpu.CompleteInstruction()
	}
}

//	PHP: Push Processor Status on Stack
//	The status register will be pushed with the break
//	flag and bit 5 set to 1.
//	push SR
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) X08_PHP() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.ReadFromAB() // Dummy read
	case 2:
		cpu.PollInterrupts()
		cpu.flag_B = true
		cpu.PushFlags()
		cpu.CompleteInstruction()
	}
}

//	PLP: Pull Processor Status from Stack
//	The status register will be pulled with the break
//	flag and bit 5 ignored.
//	pull SR
//	N	Z	C	I	D	V
//	from stack

func (cpu *CPU) X28_PLP() {
	// CPU_Cycles = 4
	switch cpu.subCycle {
	case 1:
		cpu.AddressBus = cpu.PC
		cpu.ReadFromAB() // Dummy Read
	case 2:
		cpu.AddressBus = uint16(cpu.SP) | 0x100
		cpu.ReadFromAB() //Dummy read
	case 3:
		cpu.PollInterrupts()
		cpu.PullFlags()
		cpu.CompleteInstruction()
	}
}

//	TSX: Transfer Stack Pointer to Index X
//	SP -> X
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (cpu *CPU) XBA_TSX() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.X = cpu.SP
	cpu.ReadFromAB() // Dummy read
	cpu.SetZNFlags(cpu.X)
	cpu.CompleteInstruction()
}

//	TXS: Transfer Index X to Stack Register
//	X -> SP
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) X9A_TXS() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.SP = cpu.X
	cpu.ReadFromAB() // Dummy read
	cpu.CompleteInstruction()
}

//----------------------------------------
//	Flags
//----------------------------------------

//	CLC: Clear Carry Flag
//	0 -> C
//	N	Z	C	I	D	V
//	-	-	0	-	-	-

func (cpu *CPU) X18_CLC() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.flag_Carry = false
	cpu.ReadFromAB() // Dummy read
	cpu.CompleteInstruction()
}

//	SEC: Set Carry Flag
//	1 -> C
//	N	Z	C	I	D	V
//	-	-	1	-	-	-

func (cpu *CPU) X38_SEC() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.flag_Carry = true
	cpu.ReadFromAB() // Dummy read
	cpu.CompleteInstruction()
}

//	CLI: Clear Interrupt Disable Bit
//	0 -> I
//	N	Z	C	I	D	V
//	-	-	-	0	-	-

func (cpu *CPU) X58_CLI() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.flag_InterruptDisable = false
	cpu.ReadFromAB() // Dummy read
	cpu.CompleteInstruction()
}

//	SEI: Set Interrupt Disable Status
//	1 -> I
//	N	Z	C	I	D	V
//	-	-	-	1	-	-

func (cpu *CPU) X78_SEI() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.flag_InterruptDisable = true
	cpu.ReadFromAB() // Dummy read
	cpu.CompleteInstruction()
}

//	CLD: Clear Decimal Mode
//	0 -> D
//	N	Z	C	I	D	V
//	-	-	-	-	0	-

func (cpu *CPU) XD8_CLD() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.flag_Decimal = false
	cpu.ReadFromAB() // Dummy read
	cpu.CompleteInstruction()
}

//	SED: Set Decimal Flag
//	1 -> D
//	N	Z	C	I	D	V
//	-	-	-	-	1	-

func (cpu *CPU) XF8_SED() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.flag_Decimal = true
	cpu.ReadFromAB() // Dummy read
	cpu.CompleteInstruction()
}

//	CLV: Clear Overflow Flag
//	0 -> V
//	N	Z	C	I	D	V
//	-	-	-	-	-	0

func (cpu *CPU) XB8_CLV() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.flag_Overflow = false
	cpu.ReadFromAB() // Dummy read
	cpu.CompleteInstruction()
}

//----------------------------------------
//	Other
//----------------------------------------

//	NOP: No Operation
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (cpu *CPU) XEA_NOP() {
	// CPU_Cycles = 2
	cpu.PollInterrupts()
	cpu.ReadFromAB() // Dummy read
	cpu.CompleteInstruction()
}
