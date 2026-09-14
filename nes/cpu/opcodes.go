package cpu

import "mtt/timenes/nes/cartridge/mappers"

//----------------------------------------
//	Access
//----------------------------------------

//	STA: Store Accumulator in Memory
//	A -> M
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) X85_STA_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.WriteToAB(c.A)
		c.CompleteInstruction()
	}
}
func (c *CPU) X95_STA_ZeroPage_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.PollInterrupts()
		c.WriteToAB(c.A)
		c.CompleteInstruction()
	}
}
func (c *CPU) X8D_STA_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.WriteToAB(c.A)
		c.CompleteInstruction()
	}
}
func (c *CPU) X9D_STA_Absolute_X() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(false)
	case 4:
		c.PollInterrupts()
		c.WriteToAB(c.A)
		c.CompleteInstruction()
	}
}
func (c *CPU) X99_STA_Absolute_Y() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(false)
	case 4:
		c.PollInterrupts()
		c.WriteToAB(c.A)
		c.CompleteInstruction()
	}
}
func (c *CPU) X81_STA_Indirect_X() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectX()
	case 5:
		c.PollInterrupts()
		c.WriteToAB(c.A)
		c.CompleteInstruction()
	}
}
func (c *CPU) X91_STA_Indirect_Y() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectY(false)
	case 5:
		c.PollInterrupts()
		c.WriteToAB(c.A)
		c.CompleteInstruction()
	}
}

//	LDA: Load Accumulator with Memory
//	M -> A
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) XA9_LDA_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.A = c.ReadFromPC()
	c.SetZNFlags(c.A)
	c.CompleteInstruction()
}
func (c *CPU) XA5_LDA_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.A = c.ReadFromAB()
		c.SetZNFlags(c.A)
		c.CompleteInstruction()
	}
}
func (c *CPU) XB5_LDA_ZeroPage_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.PollInterrupts()
		c.A = c.ReadFromAB()
		c.SetZNFlags(c.A)
		c.CompleteInstruction()
	}
}
func (c *CPU) XAD_LDA_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.A = c.ReadFromAB()
		c.SetZNFlags(c.A)
		c.CompleteInstruction()
	}
}
func (c *CPU) XBD_LDA_Absolute_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(true)
	case 4:
		c.PollInterrupts()
		c.A = c.ReadFromAB()
		c.SetZNFlags(c.A)
		c.CompleteInstruction()
	}
}
func (c *CPU) XB9_LDA_Absolute_Y() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(true)
	case 4:
		c.PollInterrupts()
		c.A = c.ReadFromAB()
		c.SetZNFlags(c.A)
		c.CompleteInstruction()
	}
}
func (c *CPU) XA1_LDA_Indirect_X() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectX()
	case 5:
		c.PollInterrupts()
		c.A = c.ReadFromAB()
		c.SetZNFlags(c.A)
		c.CompleteInstruction()
	}
}
func (c *CPU) XB1_LDA_Indirect_Y() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectY(true)
	case 5:
		c.PollInterrupts()
		c.A = c.ReadFromAB()
		c.SetZNFlags(c.A)
		c.CompleteInstruction()
	}
}

//	STX: Store Index X in Memory
//	X -> M
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) X86_STX_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.WriteToAB(c.X)
		c.CompleteInstruction()
	}
}
func (c *CPU) X96_STX_ZeroPage_Y() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageY()
	case 3:
		c.PollInterrupts()
		c.WriteToAB(c.X)
		c.CompleteInstruction()
	}
}
func (c *CPU) X8E_STX_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.WriteToAB(c.X)
		c.CompleteInstruction()
	}
}

//	LDX: Load Index X with Memory
//	M -> X
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) XA2_LDX_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.X = c.ReadFromPC()
	c.SetZNFlags(c.X)
	c.CompleteInstruction()
}
func (c *CPU) XA6_LDX_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.X = c.ReadFromAB()
		c.SetZNFlags(c.X)
		c.CompleteInstruction()
	}
}
func (c *CPU) XAE_LDX_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.X = c.ReadFromAB()
		c.SetZNFlags(c.X)
		c.CompleteInstruction()
	}
}
func (c *CPU) XB6_LDX_ZeroPage_Y() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageY()
	case 3:
		c.PollInterrupts()
		c.X = c.ReadFromAB()
		c.SetZNFlags(c.X)
		c.CompleteInstruction()
	}
}
func (c *CPU) XBE_LDX_Absolute_Y() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(true)
	case 4:
		c.PollInterrupts()
		c.X = c.ReadFromAB()
		c.SetZNFlags(c.X)
		c.CompleteInstruction()
	}
}

//	STY: Store Index Y in Memory
//	Y -> M
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) X84_STY_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.WriteToAB(c.Y)
		c.CompleteInstruction()
	}
}
func (c *CPU) X94_STY_ZeroPage_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.PollInterrupts()
		c.WriteToAB(c.Y)
		c.CompleteInstruction()
	}
}
func (c *CPU) X8C_STY_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.WriteToAB(c.Y)
		c.CompleteInstruction()
	}
}

//	LDY: Load Index Y with Memory
//	M -> Y
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) XA0_LDY_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Y = c.ReadFromPC()
	c.SetZNFlags(c.Y)
	c.CompleteInstruction()
}
func (c *CPU) XA4_LDY_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.Y = c.ReadFromAB()
		c.SetZNFlags(c.Y)
		c.CompleteInstruction()
	}
}
func (c *CPU) XAC_LDY_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.Y = c.ReadFromAB()
		c.SetZNFlags(c.Y)
		c.CompleteInstruction()
	}
}
func (c *CPU) XB4_LDY_ZeroPage_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.PollInterrupts()
		c.Y = c.ReadFromAB()
		c.SetZNFlags(c.Y)
		c.CompleteInstruction()
	}
}
func (c *CPU) XBC_LDY_Absolute_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(true)
	case 4:
		c.PollInterrupts()
		c.Y = c.ReadFromAB()
		c.SetZNFlags(c.Y)
		c.CompleteInstruction()
	}
}

//----------------------------------------
//	Transfer
//----------------------------------------

//	TAX: Transfer Accumulator to Index X
//	A -> X
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) XAA_TAX() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.X = c.A
	c.SetZNFlags(c.X)
	c.CompleteInstruction()
}

//	TAY: Transfer Accumulator to Index Y
//	A -> Y
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) XA8_TAY() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Y = c.A
	c.SetZNFlags(c.Y)
	c.CompleteInstruction()
}

//	TXA: Transfer Index X to Accumulator
//	X -> A
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) X8A_TXA() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.A = c.X
	c.SetZNFlags(c.A)
	c.CompleteInstruction()
}

//	TYA: Transfer Index Y to Accumulator
//	Y -> A
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) X98_TYA() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.A = c.Y
	c.SetZNFlags(c.A)
	c.CompleteInstruction()
}

//----------------------------------------
//	Arithmetic
//----------------------------------------

//	ADC: Add Memory to Accumulator with Carry
//	A + M + C -> A, C
//	N	Z	C	I	D	V
//	+	+	+	-	-	+

func (c *CPU) X69_ADC_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Op_ADC(c.ReadFromPC())
	c.CompleteInstruction()
}
func (c *CPU) X65_ADC_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.Op_ADC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X75_ADC_ZeroPage_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.PollInterrupts()
		c.Op_ADC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X6D_ADC_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.Op_ADC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X7D_ADC_Absolute_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(true)
	case 4:
		c.PollInterrupts()
		c.Op_ADC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X79_ADC_Absolute_Y() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(true)
	case 4:
		c.PollInterrupts()
		c.Op_ADC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X61_ADC_Indirect_X() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectX()
	case 5:
		c.PollInterrupts()
		c.Op_ADC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X71_ADC_Indirect_Y() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectY(true)
	case 5:
		c.PollInterrupts()
		c.Op_ADC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}

//	SBC: Subtract Memory from Accumulator with Borrow
//	A - M - C̅ -> A
//	N	Z	C	I	D	V
//	+	+	+	-	-	+

func (c *CPU) XE9_SBC_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Op_SBC(c.ReadFromPC())
	c.CompleteInstruction()
}
func (c *CPU) XE5_SBC_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.Op_SBC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XED_SBC_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.Op_SBC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XF5_SBC_ZeroPage_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.PollInterrupts()
		c.Op_SBC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XFD_SBC_Absolute_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(true)
	case 4:
		c.PollInterrupts()
		c.Op_SBC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XF9_SBC_Absolute_Y() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(true)
	case 4:
		c.PollInterrupts()
		c.Op_SBC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XE1_SBC_Indirect_X() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectX()
	case 5:
		c.PollInterrupts()
		c.Op_SBC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XF1_SBC_Indirect_Y() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectY(true)
	case 5:
		c.PollInterrupts()
		c.Op_SBC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}

//	INC: Increment Memory by One
//	M + 1 -> M
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) XE6_INC_ZeroPage() {
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
		c.Op_INC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XEE_INC_Absolute() {
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
		c.Op_INC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XF6_INC_ZeroPage_X() {
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
		c.Op_INC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XFE_INC_Absolute_X() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_AbsoluteX(false)
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_INC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}

//	DEC: Decrement Memory by One
//	M - 1 -> M
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) XC6_DEC_ZeroPage() {
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
		c.Op_DEC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XCE_DEC_Absolute() {
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
		c.Op_DEC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XD6_DEC_ZeroPage_X() {
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
		c.Op_DEC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XDE_DEC_Absolute_X() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_AbsoluteX(false)
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_DEC(c.ReadFromAB())
		c.CompleteInstruction()
	}
}

//	INX: Increment Index X by One
//	X + 1 -> X
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) XE8_INX() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.ReadFromAB() // Dummy Read
	c.X++
	c.SetZNFlags(c.X)
	c.CompleteInstruction()
}

//	DEX: Decrement Index X by One
//	X - 1 -> X
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) XCA_DEX() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.ReadFromAB() // Dummy Read
	c.X--
	c.SetZNFlags(c.X)
	c.CompleteInstruction()
}

//	INY: Increment Index Y by One
//	Y + 1 -> Y
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) XC8_INY() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.ReadFromAB() // Dummy Read
	c.Y++
	c.SetZNFlags(c.Y)
	c.CompleteInstruction()
}

//	DEY: Decrement Index Y by One
//	Y - 1 -> Y
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) X88_DEY() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.ReadFromAB() // Dummy Read
	c.Y--
	c.SetZNFlags(c.Y)
	c.CompleteInstruction()
}

//----------------------------------------
//	Shift
//----------------------------------------

//	ASL: Shift Left One Bit (Memory or Accumulator)
//	C <- [76543210] <- 0
//	N	Z	C	I	D	V
//	+	+	+	-	-	-

func (c *CPU) X0A_ASL() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.flag_Carry = c.A > 127
	c.A <<= 1
	c.SetZNFlags(c.A)
	c.CompleteInstruction()
}
func (c *CPU) X06_ASL_ZeroPage() {
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
		c.Op_ASL()
		c.CompleteInstruction()
	}
}
func (c *CPU) X0E_ASL_Absolute() {
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
		c.Op_ASL()
		c.CompleteInstruction()
	}
}
func (c *CPU) X16_ASL_ZeroPage_X() {
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
		c.Op_ASL()
		c.CompleteInstruction()
	}
}
func (c *CPU) X1E_ASL_Absolute_X() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_AbsoluteX(false)
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_ASL()
		c.CompleteInstruction()
	}
}

//	LSR: Shift One Bit Right (Memory or Accumulator)
//	0 -> [76543210] -> C
//	N	Z	C	I	D	V
//	0	+	+	-	-	-

func (c *CPU) X4A_LSR() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.flag_Carry = (c.A & 1) != 0
	c.A >>= 1
	c.SetZNFlags(c.A)
	c.CompleteInstruction()
}
func (c *CPU) X46_LSR_ZeroPage() {
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
		c.Op_LSR()
		c.CompleteInstruction()
	}
}
func (c *CPU) X4E_LSR_Absolute() {
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
		c.Op_LSR()
		c.CompleteInstruction()
	}
}
func (c *CPU) X56_LSR_ZeroPage_X() {
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
		c.Op_LSR()
		c.CompleteInstruction()
	}
}
func (c *CPU) X5E_LSR_Absolute_X() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_AbsoluteX(false)
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_LSR()
		c.CompleteInstruction()
	}
}

//	ROL: Rotate One Bit Left (Memory or Accumulator)
//	C <- [76543210] <- C
//	N	Z	C	I	D	V
//	+	+	+	-	-	-

func (c *CPU) X2A_ROL() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	futureCarry := (c.A >= 0x80)
	c.A <<= 1
	if c.flag_Carry {
		c.A |= 1
	}
	c.flag_Carry = futureCarry
	c.SetZNFlags(c.A)
	c.CompleteInstruction()
}
func (c *CPU) X26_ROL_ZeroPage() {
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
		c.Op_ROL()
		c.CompleteInstruction()
	}
}
func (c *CPU) X2E_ROL_Absolute() {
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
		c.Op_ROL()
		c.CompleteInstruction()
	}
}
func (c *CPU) X36_ROL_ZeroPage_X() {
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
		c.Op_ROL()
		c.CompleteInstruction()
	}
}
func (c *CPU) X3E_ROL_Absolute_X() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_AbsoluteX(false)
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_ROL()
		c.CompleteInstruction()
	}
}

//	ROR: Rotate One Bit Right (Memory or Accumulator)
//	C -> [76543210] -> C
//	N	Z	C	I	D	V
//	+	+	+	-	-	-

func (c *CPU) X6A_ROR() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	futureCarry := (c.A & 1) != 0
	c.A >>= 1
	if c.flag_Carry {
		c.A |= 0x80
	}
	c.flag_Carry = futureCarry
	c.SetZNFlags(c.A)
	c.CompleteInstruction()
}
func (c *CPU) X66_ROR_ZeroPage() {
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
		c.Op_ROR()
		c.CompleteInstruction()
	}
}
func (c *CPU) X6E_ROR_Absolute() {
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
		c.Op_ROR()
		c.CompleteInstruction()
	}
}
func (c *CPU) X76_ROR_ZeroPage_X() {
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
		c.Op_ROR()
		c.CompleteInstruction()
	}
}
func (c *CPU) X7E_ROR_Absolute_X() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_AbsoluteX(false)
	case 5:
		c.WriteToAB(c.DL) // Dummy write
	case 6:
		c.PollInterrupts()
		c.Op_ROR()
		c.CompleteInstruction()
	}
}

//----------------------------------------
//	Bitwise
//----------------------------------------

//	AND: AND Memory with Accumulator
//	A AND M -> A
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) X29_AND_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Op_AND(c.ReadFromPC())
	c.CompleteInstruction()
}
func (c *CPU) X25_AND_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.Op_AND(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X2D_AND_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.Op_AND(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X35_AND_ZeroPage_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.PollInterrupts()
		c.Op_AND(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X3D_AND_Absolute_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(true)
	case 4:
		c.PollInterrupts()
		c.Op_AND(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X39_AND_Absolute_Y() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(true)
	case 4:
		c.PollInterrupts()
		c.Op_AND(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X21_AND_Indirect_X() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectX()
	case 5:
		c.PollInterrupts()
		c.Op_AND(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X31_AND_Indirect_Y() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectY(true)
	case 5:
		c.PollInterrupts()
		c.Op_AND(c.ReadFromAB())
		c.CompleteInstruction()
	}
}

//	ORA: OR Memory with Accumulator
//	A OR M -> A
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) X09_ORA_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Op_ORA(c.ReadFromPC())
	c.CompleteInstruction()
}
func (c *CPU) X05_ORA_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.Op_ORA(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X0D_ORA_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.Op_ORA(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X15_ORA_ZeroPage_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.PollInterrupts()
		c.Op_ORA(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X1D_ORA_Absolute_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(true)
	case 4:
		c.PollInterrupts()
		c.Op_ORA(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X19_ORA_Absolute_Y() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(true)
	case 4:
		c.PollInterrupts()
		c.Op_ORA(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X01_ORA_Indirect_X() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectX()
	case 5:
		c.PollInterrupts()
		c.Op_ORA(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X11_ORA_Indirect_Y() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectY(true)
	case 5:
		c.PollInterrupts()
		c.Op_ORA(c.ReadFromAB())
		c.CompleteInstruction()
	}
}

//	EOR: Exclusive-OR Memory with Accumulator
//	A EOR M -> A
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) X49_EOR_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Op_EOR(c.ReadFromPC())
	c.CompleteInstruction()
}
func (c *CPU) X45_EOR_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.Op_EOR(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X4D_EOR_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.Op_EOR(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X55_EOR_ZeroPage_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.PollInterrupts()
		c.Op_EOR(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X5D_EOR_Absolute_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(true)
	case 4:
		c.PollInterrupts()
		c.Op_EOR(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X59_EOR_Absolute_Y() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(true)
	case 4:
		c.PollInterrupts()
		c.Op_EOR(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X41_EOR_Indirect_X() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectX()
	case 5:
		c.PollInterrupts()
		c.Op_EOR(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X51_EOR_Indirect_Y() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectY(true)
	case 5:
		c.PollInterrupts()
		c.Op_EOR(c.ReadFromAB())
		c.CompleteInstruction()
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

func (c *CPU) X24_BIT_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.Op_BIT(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) X2C_BIT_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.Op_BIT(c.ReadFromAB())
		c.CompleteInstruction()
	}
}

//----------------------------------------
//	Compare
//----------------------------------------

//	CMP: Compare Memory with Accumulator
//	A - M
//	N	Z	C	I	D	V
//	+	+	+	-	-	-

func (c *CPU) XC9_CMP_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Op_CMP(c.ReadFromPC())
	c.CompleteInstruction()
}
func (c *CPU) XC5_CMP_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.Op_CMP(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XCD_CMP_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.Op_CMP(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XD5_CMP_ZeroPage_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_ZeroPageX()
	case 3:
		c.PollInterrupts()
		c.Op_CMP(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XDD_CMP_Absolute_X() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteX(true)
	case 4:
		c.PollInterrupts()
		c.Op_CMP(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XD9_CMP_Absolute_Y() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.GetAddress_AbsoluteY(true)
	case 4:
		c.PollInterrupts()
		c.Op_CMP(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XC1_CMP_Indirect_X() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectX()
	case 5:
		c.PollInterrupts()
		c.Op_CMP(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XD1_CMP_Indirect_Y() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1, 2, 3, 4:
		c.GetAddress_IndirectY(true)
	case 5:
		c.PollInterrupts()
		c.Op_CMP(c.ReadFromAB())
		c.CompleteInstruction()
	}
}

//	CPX: Compare Memory and Index X
//	X - M
//	N	Z	C	I	D	V
//	+	+	+	-	-	-

func (c *CPU) XE0_CPX_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Op_CPX(c.ReadFromPC())
	c.CompleteInstruction()
}
func (c *CPU) XE4_CPX_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.Op_CPX(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XEC_CPX_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.Op_CPX(c.ReadFromAB())
		c.CompleteInstruction()
	}
}

//	CPY: Compare Memory and Index Y
//	Y - M
//	N	Z	C	I	D	V
//	+	+	+	-	-	-

func (c *CPU) XC0_CPY_Immediate() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.Op_CPY(c.ReadFromPC())
	c.CompleteInstruction()
}
func (c *CPU) XC4_CPY_ZeroPage() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_ZeroPage()
	case 2:
		c.PollInterrupts()
		c.Op_CPY(c.ReadFromAB())
		c.CompleteInstruction()
	}
}
func (c *CPU) XCC_CPY_Absolute() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.PollInterrupts()
		c.Op_CPY(c.ReadFromAB())
		c.CompleteInstruction()
	}
}

//----------------------------------------
//	Branch
//----------------------------------------

//	BCC: Branch on Carry Clear
//	branch on C = 0
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) X90_BCC() {
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.Branch(!c.flag_Carry)
	}
}

//	BCS: Branch on Carry Set
//	branch on C = 1
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) XB0_BCS() {
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.Branch(c.flag_Carry)
	}
}

//	BEQ: Branch on Result Zero
//	branch on Z = 1
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) XF0_BEQ() {
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.Branch(c.flag_Zero)
	}
}

//	BNE: Branch on Result not Zero
//	branch on Z = 0
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) XD0_BNE() {
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.Branch(!c.flag_Zero)
	}
}

//	BPL: Branch on Result Plus
//	branch on N = 0
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) X10_BPL() {
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.Branch(!c.flag_Negative)
	}
}

//	BMI: Branch on Result Minus
//	branch on N = 1
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) X30_BMI() {
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.Branch(c.flag_Negative)
	}
}

//	BVC: Branch on Overflow Clear
//	branch on V = 0
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) X50_BVC() {
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.Branch(!c.flag_Overflow)
	}
}

//	BVS: Branch on Overflow Set
//	branch on V = 1
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) X70_BVS() {
	switch c.InstructionCycle {
	case 1, 2, 3:
		c.Branch(c.flag_Overflow)
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

func (c *CPU) X4C_JMP() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_Absolute()
	case 2:
		c.PollInterrupts()
		c.GetAddress_Absolute()
		c.PC = c.AddressBus
		c.CompleteInstruction()
	}
}

func (c *CPU) X6C_JMP_Indirect() {
	// CPU_Cycles = 5
	switch c.InstructionCycle {
	case 1, 2:
		c.GetAddress_Absolute()
	case 3:
		c.SB = c.ReadFromAB()
	case 4:
		c.PollInterrupts()
		c.DL = c.bus.Read((c.AddressBus & 0xFF00) | ((c.AddressBus + 1) & 0xFF))
		c.PC = c.BuildAddress(c.SB, c.DL)
		c.CompleteInstruction()
	}
}

//	JSR: Jump to New Location Saving Return Address
//	push (PC+2),
//	operand 1st byte -> PCL
//	operand 2nd byte -> PCH
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) X20_JSR() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1: // fetch the byte that will be PC low
		c.GetAddress_Immediate()
	case 2: // transfer stack pointer to address bus, and alu to stack pointer. I'm just reusing `dl` here, but this instruction actually uses the Arithmetic Logic Unit for this.
		c.AddressBus = uint16(c.SP) | 0x100
		//c.SP = c.DL
		c.ReadFromAB() // Dummy Read
	case 3: // push PC high to stack via address bus
		c.Push(byte(c.PC / 0x100))
	case 4: // push PC low to stack via address bus
		c.Push(byte(c.PC))
	case 5: // fetch PC High, transfer stack pointer to PC low, address bus to stack pointer.
		c.PollInterrupts()
		c.PC = c.BuildAddress(c.DL, c.ReadFromPC())
		c.CompleteInstruction()
	}
}

//	RTS: Return from Subroutine
//	pull PC, PC+1 -> PC
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) X60_RTS() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_Immediate() // Dummy read to advance PC
	case 2:
		c.SP--
		c.Pull() // Dummy Read (Pull)
	case 3: // Target byte low
		c.DL = c.Pull()
	case 4: // Target byte high
		c.PC = c.BuildAddress(c.DL, c.Pull())
	case 5:
		c.PollInterrupts()
		c.ReadFromPC() // Dummy read to advance PC
		c.CompleteInstruction()
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

func (c *CPU) X00_BRK() {
	// CPU_Cycles = 7
	switch c.InstructionCycle {
	case 1:
		if c.BreakSource == Break_Software {
			c.GetAddress_Immediate() // Dummy read that increments PC
		} else {
			c.ReadFromAB() // Dummy read that does not increment PC
		}
	case 2:
		if c.BreakSource != Break_Reset {
			c.Push(byte(c.PC >> 8))
		} else {
			c.ResetReadPush()
		}
	case 3:
		if c.BreakSource != Break_Reset {
			c.Push(byte(c.PC))
		} else {
			c.ResetReadPush()
		}
	case 4:
		if c.BreakSource != Break_Reset {
			c.flag_B = false
			if c.BreakSource == Break_Software {
				c.flag_B = true
			}
			c.PushFlags()
		} else {
			c.ResetReadPush()
		}
		c.PollInterrupts()
	case 5:
		if c.BreakSource == Break_NMI {
			c.PC = (c.PC & 0xFF00) | uint16(c.bus.Read(0xFFFA))
		} else if c.BreakSource == Break_Reset {
			c.PC = (c.PC & 0xFF00) | uint16(c.bus.Read(0xFFFC))
		} else {
			c.PC = (c.PC & 0xFF00) | uint16(c.bus.Read(0xFFFE))
		}
	case 6:
		if c.BreakSource == Break_NMI {
			c.PC = (c.PC & 0xFF) | (uint16(c.bus.Read(0xFFFB)) << 8)
		} else if c.BreakSource == Break_Reset {
			c.PC = (c.PC & 0xFF) | (uint16(c.bus.Read(0xFFFD)) << 8)
		} else {
			c.PC = (c.PC & 0xFF) | (uint16(c.bus.Read(0xFFFF)) << 8)
		}
		c.BreakSource = Break_None
		c.NMIPending = false
		c.IRQPending = false
		c.IRQLine = false
		mappers.MMC3_IRQPending = false
		c.flag_InterruptDisable = true
		c.CompleteInstruction()
	}
}

//	RTI: Return from Interrupt
//	The status register is pulled with the break flag
//	and bit 5 ignored. Then PC is pulled from the stack.
//
//	pull SR, pull PC
//	N	Z	C	I	D	V
//	from stack

func (c *CPU) X40_RTI() {
	// CPU_Cycles = 6
	switch c.InstructionCycle {
	case 1:
		c.GetAddress_Immediate()
	case 2:
		c.AddressBus = uint16(c.SP) | 0x100
		c.ReadFromAB() //Dummy read
	case 3:
		status := c.Pull()
		c.flag_Carry = (status & 0x01) != 0
		c.flag_Zero = (status & 0x02) != 0
		c.flag_InterruptDisable = (status & 0x04) != 0
		c.flag_Decimal = (status & 0x08) != 0
		c.flag_Overflow = (status & 0x40) != 0
		c.flag_Negative = (status & 0x80) != 0
	case 4:
		c.DL = c.Pull()
	case 5:
		c.PollInterrupts()
		c.PC = c.BuildAddress(c.DL, c.Pull())
		c.CompleteInstruction()
	}
}

//----------------------------------------
//	Stack
//----------------------------------------

//	PHA: Push Accumulator on Stack
//	push A
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) X48_PHA() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.DL = c.ReadFromAB() // Dummy read
	case 2:
		c.PollInterrupts()
		c.Push(c.A)
		c.CompleteInstruction()
	}
}

//	PLA: Pull Accumulator from Stack
//	pull A
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) X68_PLA() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1:
		c.AddressBus = c.PC
		c.ReadFromAB() // Dummy Read
	case 2:
		c.AddressBus = uint16(c.SP) | 0x100
		c.ReadFromAB() //Dummy read
	case 3:
		c.PollInterrupts()
		c.A = c.Pull()
		c.SetZNFlags(c.A)
		c.CompleteInstruction()
	}
}

//	PHP: Push Processor Status on Stack
//	The status register will be pushed with the break
//	flag and bit 5 set to 1.
//	push SR
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) X08_PHP() {
	// CPU_Cycles = 3
	switch c.InstructionCycle {
	case 1:
		c.ReadFromAB() // Dummy read
	case 2:
		c.PollInterrupts()
		c.flag_B = true
		c.PushFlags()
		c.CompleteInstruction()
	}
}

//	PLP: Pull Processor Status from Stack
//	The status register will be pulled with the break
//	flag and bit 5 ignored.
//	pull SR
//	N	Z	C	I	D	V
//	from stack

func (c *CPU) X28_PLP() {
	// CPU_Cycles = 4
	switch c.InstructionCycle {
	case 1:
		c.AddressBus = c.PC
		c.ReadFromAB() // Dummy Read
	case 2:
		c.AddressBus = uint16(c.SP) | 0x100
		c.ReadFromAB() //Dummy read
	case 3:
		c.PollInterrupts()
		c.PullFlags()
		c.CompleteInstruction()
	}
}

//	TSX: Transfer Stack Pointer to Index X
//	SP -> X
//	N	Z	C	I	D	V
//	+	+	-	-	-	-

func (c *CPU) XBA_TSX() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.X = c.SP
	c.ReadFromAB() // Dummy read
	c.SetZNFlags(c.X)
	c.CompleteInstruction()
}

//	TXS: Transfer Index X to Stack Register
//	X -> SP
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) X9A_TXS() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.SP = c.X
	c.ReadFromAB() // Dummy read
	c.CompleteInstruction()
}

//----------------------------------------
//	Flags
//----------------------------------------

//	CLC: Clear Carry Flag
//	0 -> C
//	N	Z	C	I	D	V
//	-	-	0	-	-	-

func (c *CPU) X18_CLC() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.flag_Carry = false
	c.ReadFromAB() // Dummy read
	c.CompleteInstruction()
}

//	SEC: Set Carry Flag
//	1 -> C
//	N	Z	C	I	D	V
//	-	-	1	-	-	-

func (c *CPU) X38_SEC() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.flag_Carry = true
	c.ReadFromAB() // Dummy read
	c.CompleteInstruction()
}

//	CLI: Clear Interrupt Disable Bit
//	0 -> I
//	N	Z	C	I	D	V
//	-	-	-	0	-	-

func (c *CPU) X58_CLI() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.flag_InterruptDisable = false
	c.ReadFromAB() // Dummy read
	c.CompleteInstruction()
}

//	SEI: Set Interrupt Disable Status
//	1 -> I
//	N	Z	C	I	D	V
//	-	-	-	1	-	-

func (c *CPU) X78_SEI() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.flag_InterruptDisable = true
	c.ReadFromAB() // Dummy read
	c.CompleteInstruction()
}

//	CLD: Clear Decimal Mode
//	0 -> D
//	N	Z	C	I	D	V
//	-	-	-	-	0	-

func (c *CPU) XD8_CLD() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.flag_Decimal = false
	c.ReadFromAB() // Dummy read
	c.CompleteInstruction()
}

//	SED: Set Decimal Flag
//	1 -> D
//	N	Z	C	I	D	V
//	-	-	-	-	1	-

func (c *CPU) XF8_SED() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.flag_Decimal = true
	c.ReadFromAB() // Dummy read
	c.CompleteInstruction()
}

//	CLV: Clear Overflow Flag
//	0 -> V
//	N	Z	C	I	D	V
//	-	-	-	-	-	0

func (c *CPU) XB8_CLV() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.flag_Overflow = false
	c.ReadFromAB() // Dummy read
	c.CompleteInstruction()
}

//----------------------------------------
//	Other
//----------------------------------------

//	NOP: No Operation
//	N	Z	C	I	D	V
//	-	-	-	-	-	-

func (c *CPU) XEA_NOP() {
	// CPU_Cycles = 2
	c.PollInterrupts()
	c.ReadFromAB() // Dummy read
	c.CompleteInstruction()
}
