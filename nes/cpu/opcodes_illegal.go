package cpu

//----------------------------------------
//	NOP Codes (unofficial)
//----------------------------------------

//NOP Codes (unofficial)

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

/*
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
		cpu.GetAddress_AbsoluteX(false)
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
		cpu.GetAddress_AbsoluteX(false)
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
		cpu.GetAddress_AbsoluteX(false)
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
		cpu.GetAddress_AbsoluteX(false)
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
		cpu.GetAddress_AbsoluteX(false)
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
		cpu.GetAddress_AbsoluteX(false)
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
		cpu.GetAddress_AbsoluteX(false)
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
