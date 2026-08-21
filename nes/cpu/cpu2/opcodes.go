package cpu2

/*
switch cpu.subCycle {
case 1:
case 2:
case 3:
case 4:
	cpu.CompleteInstruction()
}}
*/

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
	cpu.A = cpu.Y
	cpu.SetZNFlags(cpu.A)
	cpu.CompleteInstruction()
}

//Arithmetic

//	ADC: Add Memory to Accumulator with Carry
//	A + M + C -> A, C
//	N	Z	C	I	D	V
//	+	+	+	-	-	+

func (cpu *CPU) X69_ADC_Immediate() {
	// CPU_Cycles = 2
	cpu.Op_ADC(cpu.ReadFromPC())
	cpu.CompleteInstruction()
}
func (cpu *CPU) X65_ADC_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
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
	cpu.Op_SBC(cpu.ReadFromPC())
	cpu.CompleteInstruction()
}
func (cpu *CPU) XE5_SBC_ZeroPage() {
	// CPU_Cycles = 3
	switch cpu.subCycle {
	case 1:
		cpu.GetAddress_ZeroPage()
	case 2:
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
	cpu.ReadFromAB() // Dummy Read
	cpu.Y--
	cpu.SetZNFlags(cpu.Y)
	cpu.CompleteInstruction()
}

//Shift

//	ASL: Shift Left One Bit (Memory or Accumulator)
//	C <- [76543210] <- 0
//	N	Z	C	I	D	V
//	+	+	+	-	-	-

func (cpu *CPU) X0A_ASL() {
	// CPU_Cycles = 2
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
		cpu.Op_LSR()
		cpu.CompleteInstruction()
	}
}
func (cpu *CPU) X5E_LSR_Absolute_X() {
	cpu.GetAddress_AbsoluteX(false)
	// CPU_Cycles = 7
	switch cpu.subCycle {
	case 1, 2, 3, 4:
		cpu.GetAddress_AbsoluteX(false)
	case 5:
		cpu.WriteToAB(cpu.DL) // Dummy write
	case 6:
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
		cpu.Op_ROL()
		cpu.CompleteInstruction()
	}
}

/*
	//	ROR: Rotate One Bit Right (Memory or Accumulator)
	//	C -> [76543210] -> C
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	func (cpu *CPU) X6A_ROR(){
		futureCarry := (A & 1) != 0
		A >>= 1
		if flag_Carry {
			A |= 0x80
		}
		flag_Carry = futureCarry
		cpu.SetZNFlags(cpu.A)
		// CPU_Cycles = 2
	func (cpu *CPU) X66_ROR_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		Op_ROR(AddressBus)
		// CPU_Cycles = 5
	func (cpu *CPU) X6E_ROR_Absolute(){
		cpu.GetAddress_Absolute()
		Op_ROR(AddressBus)
		// CPU_Cycles = 6
	func (cpu *CPU) X76_ROR_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		Op_ROR(AddressBus)
		// CPU_Cycles = 6
	func (cpu *CPU) X7E_ROR_Absolute_X(){
		cpu.GetAddress_AbsoluteIndX(false)
		Op_ROR(AddressBus)
		// CPU_Cycles = 7
/*
	//Bitwise

	//	AND: AND Memory with Accumulator
	//	A AND M -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	func (cpu *CPU) X29_AND_Immediate(){
		Op_AND(cpu.ReadFromPC())
		// CPU_Cycles = 2
	func (cpu *CPU) X25_AND_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		Op_AND(cpu.ReadFromAB())
		// CPU_Cycles = 3
	func (cpu *CPU) X2D_AND_Absolute(){
		cpu.GetAddress_Absolute()
		Op_AND(cpu.ReadFromAB())
		// CPU_Cycles = 4
	func (cpu *CPU) X35_AND_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		Op_AND(cpu.ReadFromAB())
		// CPU_Cycles = 4
	func (cpu *CPU) X3D_AND_Absolute_X(){
		// CPU_Cycles = 4
		cpu.GetAddress_AbsoluteIndX(true)
		Op_AND(cpu.ReadFromAB())
	func (cpu *CPU) X39_AND_Absolute_Y(){
		// CPU_Cycles = 4
		cpu.GetAddress_AbsoluteIndY(true)
		Op_AND(cpu.ReadFromAB())
	func (cpu *CPU) X21_AND_Indirect_X(){
		cpu.GetAddress_IndirectX()
		Op_AND(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) X31_AND_Indirect_Y(){
		// CPU_Cycles = 5
		cpu.GetAddress_IndirectY(true)
		Op_AND(cpu.ReadFromAB())
/*
	//	ORA: OR Memory with Accumulator
	//	A OR M -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	func (cpu *CPU) X09ORA_Immediate(){
		Op_ORA(cpu.ReadFromPC())
		// CPU_Cycles = 2
	func (cpu *CPU) X05ORA_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		Op_ORA(cpu.ReadFromAB())
		// CPU_Cycles = 3
	func (cpu *CPU) X0DORA_Absolute(){
		cpu.GetAddress_Absolute()
		Op_ORA(cpu.ReadFromAB())
		// CPU_Cycles = 4
	func (cpu *CPU) X15ORA_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		Op_ORA(cpu.ReadFromAB())
		// CPU_Cycles = 4
	func (cpu *CPU) X1DORA_Absolute_X(){
		// CPU_Cycles = 4
		cpu.GetAddress_AbsoluteIndX(true)
		Op_ORA(cpu.ReadFromAB())
	func (cpu *CPU) X19ORA_Absolute_Y(){
		// CPU_Cycles = 4
		cpu.GetAddress_AbsoluteIndY(true)
		Op_ORA(cpu.ReadFromAB())
	func (cpu *CPU) X01ORA_Indirect_X(){
		cpu.GetAddress_IndirectX()
		Op_ORA(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) X11ORA_Indirect_Y(){
		// CPU_Cycles = 5
		cpu.GetAddress_IndirectY(true)
		Op_ORA(cpu.ReadFromAB())

	//	EOR: Exclusive-OR Memory with Accumulator
	//	A EOR M -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	func (cpu *CPU) X49EOR_Immediate(){
		Op_EOR(cpu.ReadFromPC())
		// CPU_Cycles = 2
	func (cpu *CPU) X45EOR_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		Op_EOR(cpu.ReadFromAB())
		// CPU_Cycles = 3
	func (cpu *CPU) X4DEOR_Absolute(){
		cpu.GetAddress_Absolute()
		Op_EOR(cpu.ReadFromAB())
		// CPU_Cycles = 4
	func (cpu *CPU) X55EOR_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		Op_EOR(cpu.ReadFromAB())
		// CPU_Cycles = 4
	func (cpu *CPU) X5DEOR_Absolute_X(){
		// CPU_Cycles = 4
		cpu.GetAddress_AbsoluteIndX(true)
		Op_EOR(cpu.ReadFromAB())
	func (cpu *CPU) X59EOR_Absolute_Y(){
		// CPU_Cycles = 4
		cpu.GetAddress_AbsoluteIndY(true)
		Op_EOR(cpu.ReadFromAB())
	func (cpu *CPU) X41EOR_Indirect_X(){
		cpu.GetAddress_IndirectX()
		Op_EOR(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) X51EOR_Indirect_Y(){
		// CPU_Cycles = 5
		cpu.GetAddress_IndirectY(true)
		Op_EOR(cpu.ReadFromAB())

	//	BITTest Bits in Memory with Accumulator
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

	func (cpu *CPU) X24BIT_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		Op_BIT(cpu.ReadFromAB())
		// CPU_Cycles = 3
	func (cpu *CPU) X2CBIT_Absolute(){
		cpu.GetAddress_Absolute()
		Op_BIT(cpu.ReadFromAB())
		// CPU_Cycles = 4

	//Compare

	//	CMP: Compare Memory with Accumulator
	//	A - M
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	func (cpu *CPU) XC9CMP_Immediate(){
		Op_CMP(cpu.ReadFromPC())
		// CPU_Cycles = 2
	func (cpu *CPU) XC5CMP_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		Op_CMP(cpu.ReadFromAB())
		// CPU_Cycles = 3
	func (cpu *CPU) XCDCMP_Absolute(){
		cpu.GetAddress_Absolute()
		Op_CMP(cpu.ReadFromAB())
		// CPU_Cycles = 4
	func (cpu *CPU) XD5CMP_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		Op_CMP(cpu.ReadFromAB())
		// CPU_Cycles = 4
	func (cpu *CPU) XDDCMP_Absolute_X(){
		// CPU_Cycles = 4
		cpu.GetAddress_AbsoluteIndX(true)
		Op_CMP(cpu.ReadFromAB())
	func (cpu *CPU) XD9CMP_Absolute_Y(){
		// CPU_Cycles = 4
		cpu.GetAddress_AbsoluteIndY(true)
		Op_CMP(cpu.ReadFromAB())
	func (cpu *CPU) XC1CMP_Indirect_X(){
		cpu.GetAddress_IndirectX()
		Op_CMP(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) XD1CMP_Indirect_Y(){
		// CPU_Cycles = 5
		cpu.GetAddress_IndirectY(true)
		Op_CMP(cpu.ReadFromAB())

	//	CPX: Compare Memory and Index X
	//	X - M
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	func (cpu *CPU) XE0CPX_Immediate(){
		Op_CPX(cpu.ReadFromPC())
		// CPU_Cycles = 2
	func (cpu *CPU) XE4CPX_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		Op_CPX(cpu.ReadFromAB())
		// CPU_Cycles = 3
	func (cpu *CPU) XECCPX_Absolute(){
		cpu.GetAddress_Absolute()
		Op_CPX(cpu.ReadFromAB())
		// CPU_Cycles = 4

	//	CPY: Compare Memory and Index Y
	//	Y - M
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	func (cpu *CPU) XC0CPY_Immediate(){
		Op_CPY(cpu.ReadFromPC())
		// CPU_Cycles = 2
	func (cpu *CPU) XC4CPY_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		Op_CPY(cpu.ReadFromAB())
		// CPU_Cycles = 3
	func (cpu *CPU) XCCCPY_Absolute(){
		cpu.GetAddress_Absolute()
		Op_CPY(cpu.ReadFromAB())
		// CPU_Cycles = 4

	//Branch

	//	BCC: Branch on Carry Clear
	//	branch on C = 0
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	func (cpu *CPU) X90BCC (Branch on Carry Clear)
		temp := cpu.ReadFromPC()
		Branch(!flag_Carry, temp)

	//	BCS: Branch on Carry Set
	//	branch on C = 1
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	func (cpu *CPU) XB0BCS (Branch on Carry Set)
		temp := cpu.ReadFromPC()
		Branch(flag_Carry, temp)

	//	BEQ: Branch on Result Zero
	//	branch on Z = 1
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	func (cpu *CPU) XF0BEQ (Branch on Equal)
		temp := cpu.ReadFromPC()
		Branch(flag_Zero, temp)

	//	BNE: Branch on Result not Zero
	//	branch on Z = 0
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	func (cpu *CPU) XD0BNE (Branch on Not Equal)
		temp := cpu.ReadFromPC()
		Branch(!flag_Zero, temp)

	//	BPL: Branch on Result Plus
	//	branch on N = 0
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	func (cpu *CPU) X10BPL (Branch on Plus)
		temp := cpu.ReadFromPC()
		Branch(!flag_Negative, temp)

	//	BMI: Branch on Result Minus
	//	branch on N = 1
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	func (cpu *CPU) X30BMI (Branch on Minus)
		temp := cpu.ReadFromPC()
		Branch(flag_Negative, temp)

	//	BVC: Branch on Overflow Clear
	//	branch on V = 0
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	func (cpu *CPU) X50BVC (Branch on Overflow Clear)
		temp := cpu.ReadFromPC()
		Branch(!flag_Overflow, temp)

	//	BVS: Branch on Overflow Set
	//	branch on V = 1
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	func (cpu *CPU) X70BVS (Branch on Overflow Set)
		temp := cpu.ReadFromPC()
		Branch(flag_Overflow, temp)

	//Jump

	//	JMP: Jump to New Location
	//	operand 1st byte -> PCL
	//	operand 2nd byte -> PCH
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	func (cpu *CPU) X4CJMP
		cpu.GetAddress_Absolute(true)
		PC = AddressBus
		// CPU_Cycles = 3

	func (cpu *CPU) X6CJMP Indirect
		cpu.GetAddress_Indirect()
		PC = AddressBus
		// CPU_Cycles = 5 //TODO: What the fuck

	//	JSR: Jump to New Location Saving Return Address
	//	push (PC+2),
	//	operand 1st byte -> PCL
	//	operand 2nd byte -> PCH
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	func (cpu *CPU) X20JSR
		temp_low := cpu.ReadFromPC()
		Push(byte(PC / 0x100))
		Push(byte(PC))
		temp_high := cpu.ReadFromPC()
		PC = BuildAddress(temp_low, temp_high)
		//MasterClockTick("jsr")
		// CPU_Cycles = 6

	//	RTS: Return from Subroutine
	//	pull PC, PC+1 -> PC
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	func (cpu *CPU) X60RTS
		temp_low := Pull()
		temp_high := Pull()
		//MasterClockTick("rts Pull1")
		//MasterClockTick("rts Pull2")
		PC = BuildAddress(temp_low, temp_high)
		PC++
		//MasterClockTick("rts pc++")
		// CPU_Cycles = 6

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

	func (cpu *CPU) X00BRK
		flag_B = false
		if !DoNMI {
			PC++
			flag_B = true
		}
		Push(byte(PC >> 8))
		Push(byte(PC))
		PushFlags()
		//flag_InterruptDisable = true

		PCL := bus.Read(common.Ternary(DoNMI, 0xFFFA, 0xFFFE))
		PCH := bus.Read(common.Ternary(DoNMI, 0xFFFB, 0xFFFF))
		PC = uint16((uint16(PCH) * 0x100) + uint16(PCL)) //BuildAddress(PCL, PCH)
		DoNMI = false
		// CPU_Cycles = 7

	//	RTI: Return from Interrupt
	//	The status register is pulled with the break flag
	//	and bit 5 ignored. Then PC is pulled from the stack.
	//
	//	pull SR, pull PC
	//	N	Z	C	I	D	V
	//	from stack

	func (cpu *CPU) X40RTI
		temp := Pull()
		flag_Carry = (temp & 1) != 0
		flag_Zero = (temp & 2) != 0
		flag_InterruptDisable = (temp & 4) != 0
		flag_Decimal = (temp & 8) != 0
		flag_Overflow = (temp & 64) != 0
		flag_Negative = (temp & 128) != 0
		//MasterClockTick("rti Pull1")
		//MasterClockTick("rti Pull2")
		temp_low := Pull()
		temp_high := Pull()
		PC = BuildAddress(temp_low, temp_high)
		// CPU_Cycles = 6

	//Stack

	//	PHA: Push Accumulator on Stack
	//	push A
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	func (cpu *CPU) X48PHA
		Push(cpu.A)
		// CPU_Cycles = 3

	//	PLA: Pull Accumulator from Stack
	//	pull A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	func (cpu *CPU) X68PLA
		A = Pull()
		//MasterClockTick("pla")
		cpu.SetZNFlags(cpu.A)
		// CPU_Cycles = 4

	//	PHP: Push Processor Status on Stack
	//	The status register will be pushed with the break
	//	flag and bit 5 set to 1.
	//	push SR
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-
	func (cpu *CPU) X08PHP
		flag_B = true
		PushFlags()
		// CPU_Cycles = 3

	//	PLP: Pull Processor Status from Stack
	//	The status register will be pulled with the break
	//	flag and bit 5 ignored.
	//	pull SR
	//	N	Z	C	I	D	V
	//	from stack

	func (cpu *CPU) X28PLP
		PullFlags()
		// CPU_Cycles = 4

	//	TSX: Transfer Stack Pointer to Index X
	//	SP -> X
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	func (cpu *CPU) XBATSX
		X = SP
		cpu.SetZNFlags(cpu.X)
		// CPU_Cycles = 2

	//	TXS: Transfer Index X to Stack Register
	//	X -> SP
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	func (cpu *CPU) X9ATXS
		SP = X
		// CPU_Cycles = 2

	//Flags

	//	CLC: Clear Carry Flag
	//	0 -> C
	//	N	Z	C	I	D	V
	//	-	-	0	-	-	-

	func (cpu *CPU) X18CLC
		flag_Carry = false
		//MasterClockTick("clc")
		// CPU_Cycles = 2

	//	SEC: Set Carry Flag
	//	1 -> C
	//	N	Z	C	I	D	V
	//	-	-	1	-	-	-

	func (cpu *CPU) X38SEC
		flag_Carry = true
		//MasterClockTick("sec")
		// CPU_Cycles = 2

	//	CLI: Clear Interrupt Disable Bit
	//	0 -> I
	//	N	Z	C	I	D	V
	//	-	-	-	0	-	-

	func (cpu *CPU) X58CLI
		flag_InterruptDisable = false
		//MasterClockTick("cli")
		// CPU_Cycles = 2

	//	SEI: Set Interrupt Disable Status
	//	1 -> I
	//	N	Z	C	I	D	V
	//	-	-	-	1	-	-

	func (cpu *CPU) X78SEI
		flag_InterruptDisable = true
		//MasterClockTick("sei")
		// CPU_Cycles = 2

	//	CLD: Clear Decimal Mode
	//	0 -> D
	//	N	Z	C	I	D	V
	//	-	-	-	-	0	-

	func (cpu *CPU) XD8CLD
		flag_Decimal = false
		//MasterClockTick("cld")
		// CPU_Cycles = 2

	//	SED: Set Decimal Flag
	//	1 -> D
	//	N	Z	C	I	D	V
	//	-	-	-	-	1	-

	func (cpu *CPU) XF8SED
		flag_Decimal = true
		//MasterClockTick("sed")
		// CPU_Cycles = 2

	//	CLV: Clear Overflow Flag
	//	0 -> V
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	0

	func (cpu *CPU) XB8CLV
		flag_Overflow = false
		//MasterClockTick("clv")
		// CPU_Cycles = 2

	//Other

	//	NOP: No Operation
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	func (cpu *CPU) XEANOP
		// CPU_Cycles = 2

	//----------------------------------------------------------------------
	//Unofficial Opcodes
	//----------------------------------------------------------------------

	//HLT Codes

	func (cpu *CPU) X02HLT
		CPU_Halted = true
	func (cpu *CPU) X12HLT
		CPU_Halted = true
	func (cpu *CPU) X22HLT
		CPU_Halted = true
	func (cpu *CPU) X32HLT
		CPU_Halted = true
	func (cpu *CPU) X42HLT
		CPU_Halted = true
	func (cpu *CPU) X52HLT
		CPU_Halted = true
	func (cpu *CPU) X62HLT
		CPU_Halted = true
	func (cpu *CPU) X72HLT
		CPU_Halted = true
	func (cpu *CPU) X92HLT
		CPU_Halted = true
	func (cpu *CPU) XB2HLT
		CPU_Halted = true
	func (cpu *CPU) XD2HLT
		CPU_Halted = true
	func (cpu *CPU) XF2HLT
		CPU_Halted = true

	//NOP Codes (unofficial)

	func (cpu *CPU) X1ANOP Implied
		// CPU_Cycles = 2
	func (cpu *CPU) X3ANOP Implied
		// CPU_Cycles = 2
	func (cpu *CPU) X5ANOP Implied
		// CPU_Cycles = 2
	func (cpu *CPU) X7ANOP Implied
		// CPU_Cycles = 2
	func (cpu *CPU) XDANOP Implied
		// CPU_Cycles = 2
	func (cpu *CPU) XFANOP Implied
		// CPU_Cycles = 2

	func (cpu *CPU) X80NOP_Immediate(){
		cpu.ReadFromPC()
		// CPU_Cycles = 2
	func (cpu *CPU) X82NOP_Immediate(){
		cpu.ReadFromPC()
		// CPU_Cycles = 2
	func (cpu *CPU) X89NOP_Immediate(){
		cpu.ReadFromPC()
		// CPU_Cycles = 2
	func (cpu *CPU) XC2NOP_Immediate(){
		cpu.ReadFromPC()
		// CPU_Cycles = 2
	func (cpu *CPU) XE2NOP_Immediate(){
		cpu.ReadFromPC()
		// CPU_Cycles = 2

	func (cpu *CPU) X04NOP_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		// CPU_Cycles = 3
	func (cpu *CPU) X44NOP_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		// CPU_Cycles = 3
	func (cpu *CPU) X64NOP_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		// CPU_Cycles = 3

	func (cpu *CPU) X14NOP_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		// CPU_Cycles = 4
	func (cpu *CPU) X34NOP_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		// CPU_Cycles = 4
	func (cpu *CPU) X54NOP_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		// CPU_Cycles = 4
	func (cpu *CPU) X74NOP_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		// CPU_Cycles = 4
	func (cpu *CPU) XD4NOP_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		// CPU_Cycles = 4
	func (cpu *CPU) XF4NOP_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		// CPU_Cycles = 4

	func (cpu *CPU) X0CNOP_Absolute(){
		cpu.GetAddress_Absolute()
		// CPU_Cycles = 4

	func (cpu *CPU) X1CNOP_Absolute_X(){
		// CPU_Cycles = 4
		cpu.GetAddress_AbsoluteIndX(true)
	func (cpu *CPU) X3CNOP_Absolute_X(){
		// CPU_Cycles = 4
		cpu.GetAddress_AbsoluteIndX(true)
	func (cpu *CPU) X5CNOP_Absolute_X(){
		// CPU_Cycles = 4
		cpu.GetAddress_AbsoluteIndX(true)
	func (cpu *CPU) X7CNOP_Absolute_X(){
		// CPU_Cycles = 4
		cpu.GetAddress_AbsoluteIndX(true)
	func (cpu *CPU) XDCNOP_Absolute_X(){
		// CPU_Cycles = 4
		cpu.GetAddress_AbsoluteIndX(true)
	func (cpu *CPU) XFCNOP_Absolute_X(){
		// CPU_Cycles = 4
		cpu.GetAddress_AbsoluteIndX(true)

	//SAX: A AND X -> M

	func (cpu *CPU) X87SAX_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		bus.Write(AddressBus, (A & X))
		// CPU_Cycles = 3
	func (cpu *CPU) X97SAX_ZeroPage_Y(){
		cpu.GetAddress_ZeroPageY()
		bus.Write(AddressBus, (A & X))
		// CPU_Cycles = 4
	func (cpu *CPU) X8FSAX_Absolute(){
		cpu.GetAddress_Absolute()
		bus.Write(AddressBus, (A & X))
		// CPU_Cycles = 4
	func (cpu *CPU) X83SAX_Indirect_X(){
		cpu.GetAddress_IndirectX()
		bus.Write(AddressBus, (A & X))
		// CPU_Cycles = 6

	//LAX: LDA + LDX

	func (cpu *CPU) XA7LAX_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		cpu.A = cpu.ReadFromAB()
		X = A
		cpu.SetZNFlags(cpu.X)
		// CPU_Cycles = 3
	func (cpu *CPU) XB7LAX_ZeroPage_Y(){
		cpu.GetAddress_ZeroPageY()
		cpu.A = cpu.ReadFromAB()
		X = A
		cpu.SetZNFlags(cpu.X)
		// CPU_Cycles = 4
	func (cpu *CPU) XAFLAX_Absolute(){
		cpu.GetAddress_Absolute()
		cpu.A = cpu.ReadFromAB()
		X = A
		cpu.SetZNFlags(cpu.X)
		// CPU_Cycles = 4
	func (cpu *CPU) XBFLAX_Absolute_Y(){
		// CPU_Cycles = 4
		cpu.GetAddress_AbsoluteIndY(true)
		cpu.A = cpu.ReadFromAB()
		X = A
		cpu.SetZNFlags(cpu.X)
	func (cpu *CPU) XA3LAX_Indirect_X(){
		cpu.GetAddress_IndirectX()
		cpu.A = cpu.ReadFromAB()
		X = A
		cpu.SetZNFlags(cpu.X)
		// CPU_Cycles = 6
	func (cpu *CPU) XB3LAX_Indirect_Y(){
		// CPU_Cycles = 5
		cpu.GetAddress_IndirectY(true)
		cpu.A = cpu.ReadFromAB()
		X = A
		cpu.SetZNFlags(cpu.X)

	//SLO: ASL + ORA

	func (cpu *CPU) X07 SLO_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		Op_ASL(AddressBus)
		Op_ORA(cpu.ReadFromAB())
		// CPU_Cycles = 5
	func (cpu *CPU) X17SLO_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		Op_ASL(AddressBus)
		Op_ORA(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) X0FSLO_Absolute(){
		cpu.GetAddress_Absolute()
		Op_ASL(AddressBus)
		Op_ORA(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) X1FSlo_Absolute_X(){
		cpu.GetAddress_AbsoluteIndX(false)
		Op_ASL(AddressBus)
		Op_ORA(cpu.ReadFromAB())
		// CPU_Cycles = 7
	func (cpu *CPU) X1BSlo_Absolute_Y(){
		cpu.GetAddress_AbsoluteIndY(false)
		Op_ASL(AddressBus)
		Op_ORA(cpu.ReadFromAB())
		// CPU_Cycles = 7
	func (cpu *CPU) X03Slo_Indirect_X(){
		cpu.GetAddress_IndirectX()
		Op_ASL(AddressBus)
		Op_ORA(cpu.ReadFromAB())
		// CPU_Cycles = 8
	func (cpu *CPU) X13Slo_Indirect_Y(){
		cpu.GetAddress_IndirectY(false)
		Op_ASL(AddressBus)
		Op_ORA(cpu.ReadFromAB())
		// CPU_Cycles = 8

	//DCP: DEC + CMP

	func (cpu *CPU) XC7DCP_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		Op_CMP(cpu.ReadFromAB())
		// CPU_Cycles = 5
	func (cpu *CPU) XD7DCP_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		Op_CMP(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) XCFDCP_Absolute(){
		cpu.GetAddress_Absolute()
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		Op_CMP(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) XDFDCP_Absolute_X(){
		cpu.GetAddress_AbsoluteIndX(false)
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		Op_CMP(cpu.ReadFromAB())
		// CPU_Cycles = 7
	func (cpu *CPU) XDBDCP_Absolute_Y(){
		cpu.GetAddress_AbsoluteIndY(false)
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		Op_CMP(cpu.ReadFromAB())
		// CPU_Cycles = 7
	func (cpu *CPU) XC3DCP_Indirect_X(){
		cpu.GetAddress_IndirectX()
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		Op_CMP(cpu.ReadFromAB())
		// CPU_Cycles = 8
	func (cpu *CPU) XD3DCP_Indirect_Y(){
		cpu.GetAddress_IndirectY(false)
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		Op_CMP(cpu.ReadFromAB())
		// CPU_Cycles = 8

	//SHA: Stores A AND X AND (high-byte of addr. + 1) at addr.

	func (cpu *CPU) X9FSHA_Absolute_Y(){
		cpu.GetAddress_AbsoluteIndY(false)
		HiByte := byte(AddressBus >> 8)
		bus.Write(AddressBus, A&X&(HiByte+1))
		// CPU_Cycles = 5
	func (cpu *CPU) X93SHA_Indirect_Y(){
		cpu.GetAddress_IndirectY(false)
		bus.Write(AddressBus, A&X&byte((AddressBus>>8)+1))
		// CPU_Cycles = 6

	//SHX: Stores X AND (high-byte of addr. + 1) at addr.

	func (cpu *CPU) X9ESHX_Absolute_Y(){
		cpu.GetAddress_AbsoluteIndY(false)
		val := (X & byte(((AddressBus&0xFF00)>>8)+1))
		bus.Write(AddressBus, val)
		// CPU_Cycles = 5

	//SHY: Stores Y AND (high-byte of addr. + 1) at addr.

	func (cpu *CPU) X9CSHY_Absolute_X(){
		cpu.GetAddress_AbsoluteIndX(false)
		val := (Y & byte(((AddressBus&0xFF00)>>8)+1))
		bus.Write(AddressBus, val)
		// CPU_Cycles = 5

	//RLA: ROL + AND

	func (cpu *CPU) X27RLA_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		Op_ROL(AddressBus)
		Op_AND(cpu.ReadFromAB())
		// CPU_Cycles = 5
	func (cpu *CPU) X37RLA_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		Op_ROL(AddressBus)
		Op_AND(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) X2FRLA_Absolute(){
		cpu.GetAddress_Absolute()
		Op_ROL(AddressBus)
		Op_AND(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) X3FRLA_Absolute_X(){
		cpu.GetAddress_AbsoluteIndX(false)
		Op_ROL(AddressBus)
		Op_AND(cpu.ReadFromAB())
		// CPU_Cycles = 7
	func (cpu *CPU) X3BRLA_Absolute_Y(){
		cpu.GetAddress_AbsoluteIndY(false)
		Op_ROL(AddressBus)
		Op_AND(cpu.ReadFromAB())
		// CPU_Cycles = 7
	func (cpu *CPU) X23RLA_Indirect_X(){
		cpu.GetAddress_IndirectX()
		Op_ROL(AddressBus)
		Op_AND(cpu.ReadFromAB())
		// CPU_Cycles = 8
	func (cpu *CPU) X33RLA_Indirect_Y(){
		cpu.GetAddress_IndirectY(false)
		Op_ROL(AddressBus)
		Op_AND(cpu.ReadFromAB())
		// CPU_Cycles = 8

	//SRE: LSR + EOR

	func (cpu *CPU) X47SRE_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		Op_LSR(AddressBus)
		Op_EOR(cpu.ReadFromAB())
		// CPU_Cycles = 5
	func (cpu *CPU) X57SRE_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		Op_LSR(AddressBus)
		Op_EOR(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) X4FSRE_Absolute(){
		cpu.GetAddress_Absolute()
		Op_LSR(AddressBus)
		Op_EOR(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) X5FSRE_Absolute_X(){
		cpu.GetAddress_AbsoluteIndX(false)
		Op_LSR(AddressBus)
		Op_EOR(cpu.ReadFromAB())
		// CPU_Cycles = 7
	func (cpu *CPU) X5BSRE_Absolute_Y(){
		cpu.GetAddress_AbsoluteIndY(false)
		Op_LSR(AddressBus)
		Op_EOR(cpu.ReadFromAB())
		// CPU_Cycles = 7
	func (cpu *CPU) X43SRE_Indirect_X(){
		cpu.GetAddress_IndirectX()
		Op_LSR(AddressBus)
		Op_EOR(cpu.ReadFromAB())
		// CPU_Cycles = 8
	func (cpu *CPU) X53SRE_Indirect_Y(){
		cpu.GetAddress_IndirectY(false)
		Op_LSR(AddressBus)
		Op_EOR(cpu.ReadFromAB())
		// CPU_Cycles = 8

	//RRA: ROR + ADC

	func (cpu *CPU) X67RRA_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		Op_ROR(AddressBus)
		Op_ADC(cpu.ReadFromAB())
		// CPU_Cycles = 5
	func (cpu *CPU) X77RRA_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		Op_ROR(AddressBus)
		Op_ADC(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) X6FRRA_Absolute(){
		cpu.GetAddress_Absolute()
		Op_ROR(AddressBus)
		Op_ADC(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) X7FRRA_Absolute_X(){
		cpu.GetAddress_AbsoluteIndX(false)
		Op_ROR(AddressBus)
		Op_ADC(cpu.ReadFromAB())
		// CPU_Cycles = 7
	func (cpu *CPU) X7BRRA_Absolute_Y(){
		cpu.GetAddress_AbsoluteIndY(false)
		Op_ROR(AddressBus)
		Op_ADC(cpu.ReadFromAB())
		// CPU_Cycles = 7
	func (cpu *CPU) X63RRA_Indirect_X(){
		cpu.GetAddress_IndirectX()
		Op_ROR(AddressBus)
		Op_ADC(cpu.ReadFromAB())
		// CPU_Cycles = 8
	func (cpu *CPU) X73RRA_Indirect_Y(){
		cpu.GetAddress_IndirectY(false)
		Op_ROR(AddressBus)
		Op_ADC(cpu.ReadFromAB())
		// CPU_Cycles = 8

	//ISC: INC + SBC

	func (cpu *CPU) XE7ISC_ZeroPage(){
		cpu.GetAddress_ZeroPage()
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		// CPU_Cycles = 5
	func (cpu *CPU) XF7ISC_ZeroPage_X(){
		cpu.GetAddress_ZeroPageX()
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) XEFISC_Absolute(){
		cpu.GetAddress_Absolute()
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		// CPU_Cycles = 6
	func (cpu *CPU) XFFISC_Absolute_X(){
		cpu.GetAddress_AbsoluteIndX(false)
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		// CPU_Cycles = 7
	func (cpu *CPU) XFBISC_Absolute_Y(){
		cpu.GetAddress_AbsoluteIndY(false)
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		// CPU_Cycles = 7
	func (cpu *CPU) XE3ISC_Indirect_X(){
		cpu.GetAddress_IndirectX()
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		// CPU_Cycles = 8
	func (cpu *CPU) XF3ISC_Indirect_Y(){
		cpu.GetAddress_IndirectY(false)
		cpu.Op_INC(cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		// CPU_Cycles = 8

	//Immediates (unofficial)

	func (cpu *CPU) X0BANC_Immediate(){
		//AND + Set Carry as ASL
		Op_AND(cpu.ReadFromPC())
		flag_Carry = flag_Negative
		// CPU_Cycles = 2
	func (cpu *CPU) X2BANC2_Immediate(){
		//AND + Set Carry as ROL
		//Same as $0B
		Op_AND(cpu.ReadFromPC())
		flag_Carry = flag_Negative
		// CPU_Cycles = 2
	func (cpu *CPU) X4BALR_Immediate(){
		//AND + LSR
		Op_AND(cpu.ReadFromPC())
		flag_Carry = (A & 1) != 0
		A >>= 1
		cpu.SetZNFlags(cpu.A)
		// CPU_Cycles = 2
	func (cpu *CPU) X6BARR_Immediate(){
		//AND + ROR
		Op_AND(cpu.ReadFromPC())
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
		X = A
		cpu.SetZNFlags(cpu.A)
		// CPU_Cycles = 2
	func (cpu *CPU) XCBSBX_Immediate(){
		//(A AND X) - oper -> X
		X = (A & X) - cpu.ReadFromPC()
		Op_CMP(cpu.X)
		cpu.SetZNFlags(cpu.X)
		// CPU_Cycles = 2
	func (cpu *CPU) XEBSBC_Immediate(){
		cpu.Op_SBC(cpu.ReadFromPC())
		// CPU_Cycles = 2

	func (cpu *CPU) X9BTAS (SHS)_Absolute_Y(){
		//A AND X -> SP, A AND X AND (H+1) -> M
		temp := A & X
		Push(temp)
		cpu.GetAddress_AbsoluteIndY(false)
		bus.Write(AddressBus, (A&X)&(byte(AddressBus&0xFF00)+1))
		// CPU_Cycles = 5

	func (cpu *CPU) XBBLAS (LAE)_Absolute_Y(){
		//LDA/TSX oper
		//M AND SP -> A, X, SP
		// CPU_Cycles = 4
		cpu.GetAddress_AbsoluteIndY(true)
		Value := (cpu.ReadFromAB() & SP)
		A = Value
		X = Value
		SP = Value
		cpu.SetZNFlags(cpu.A)

	default:
		fmt.Println("Unknown Opcode: " + fmt.Sprintf("%02X", opcode))

*/
