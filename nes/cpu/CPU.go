package cpu

import (
	"fmt"
	"mtt/timenes/common"
	"mtt/timenes/nes/apu"
	"mtt/timenes/nes/bus"
	"mtt/timenes/nes/ppu"
)

type CPU struct {
}

// CPU Registers
var (
	PC uint16 // Program Counter
	SP byte   // Stack Pointer
	A  byte   // Accumulator
	X  byte   // X-Index
	Y  byte   // Y-Index
)

// Status Register
var (
	flag_Carry            bool // Bit 0: Carry Flag
	flag_Zero             bool // Bit 1: Zero Flag
	flag_InterruptDisable bool // Bit 2: Interrupt Disable Flag
	flag_Decimal          bool // Bit 3: Decimal Flag
	flag_B                bool // Bit 4: B Flag
	flag_Overflow         bool // Bit 6: Overflow Flag
	flag_Negative         bool // Bit 7: Negative Flag
)

var opcode byte
var operands []byte
var CPU_Cycles, CPU_Cycles_New int
var CPU_Halted = false
var MagicConstant byte = 0xFF //Might be needed for some of the illegal opcodes
var NMILevelDetector, DoNMI bool

var apuRun bool
var AddressBus uint16
var Pointer uint16
var Target uint16

func ResetCPU() {
	SP = 0xFD
	A, X, Y = 0, 0, 0
	opcode = 0
	operands = nil
	CPU_Cycles, CPU_Cycles_New, common.CPU_TotalCycles = 0, 0, 0
	CPU_Halted = false
	NMILevelDetector, DoNMI = false, false

	flag_Carry = false
	flag_Zero = false
	flag_InterruptDisable = true
	flag_Decimal = false
	flag_Overflow = false
	flag_Negative = false
	flag_B = false
}

func Emulate_CPU() {
	//Non Maskable Interrupt check
	prevNMILevelDetector := NMILevelDetector
	NMILevelDetector = (ppu.PPUCTRL_EnableNMI && ppu.PPUSTATUS_VBlank)
	if !prevNMILevelDetector && NMILevelDetector {
		DoNMI = true
	}
	CPU_Cycles_New = 0

	//Get the opcode
	//var opcode byte
	if !DoNMI {
		//If we're not running an NMI
		opcode = bus.Read(PC)
		//if debug.LoggingCPU {
		//	debug.prepTraceLogger()
		//}
		PC++
		//MasterClockTick("OPCODE")
	} else {
		//If we're running an NMI, force opcode $00
		opcode = 0x00
	}

	CPU_Cycles = 0
	switch opcode {

	//Access

	//	STA: Store Accumulator in Memory
	//	A -> M
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x85: //STA Zero Page
		ReadOperands_ZeroPageAddressed()
		WriteToAB(A)
		CPU_Cycles = 3
	case 0x95: //STA Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		WriteToAB(A)
		CPU_Cycles = 4
	case 0x8D: //STA Absolute
		ReadOperands_AbsoluteAddressed(false)
		WriteToAB(A)
		CPU_Cycles = 4
	case 0x9D: //STA Absolute, X
		CPU_Cycles = 5
		ReadOperands_AbsoluteAddressed_XIndexed(false)
		WriteToAB(A)
	case 0x99: //STA Absolute, Y
		ReadOperands_AbsoluteAddressed_YIndexed(false)
		WriteToAB(A)
		CPU_Cycles = 5
	case 0x81: //STA Indirect, X
		ReadOperands_IndirectAddressed_XIndexed()
		WriteToAB(A)
		CPU_Cycles = 6
	case 0x91: //STA Indirect, Y
		ReadOperands_IndirectAddressed_YIndexed(false)
		WriteToAB(A)
		CPU_Cycles = 6

	//	LDA: Load Accumulator with Memory
	//	M -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xA9: //LDA Immediate
		A = ReadFromPC()
		SetZNFlags(A)
		CPU_Cycles = 2
	case 0xA5: //LDA Zero Page
		ReadOperands_ZeroPageAddressed()
		A = ReadFromAB()
		SetZNFlags(A)
		CPU_Cycles = 3
	case 0xB5: //LDA Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		A = ReadFromAB()
		SetZNFlags(A)
		CPU_Cycles = 4
	case 0xAD: //LDA Absolute
		ReadOperands_AbsoluteAddressed(false)
		A = ReadFromAB()
		SetZNFlags(A)
		CPU_Cycles = 4
	case 0xBD: //LDA Absolute, X
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_XIndexed(true)
		A = ReadFromAB()
		SetZNFlags(A)
	case 0xB9: //LDA Absolute, Y
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_YIndexed(true)
		A = ReadFromAB()
		SetZNFlags(A)
	case 0xA1: //LDA Indirect, X
		ReadOperands_IndirectAddressed_XIndexed()
		A = ReadFromAB()
		SetZNFlags(A)
		CPU_Cycles = 6
	case 0xB1: //LDA Indirect, Y
		CPU_Cycles = 5
		ReadOperands_IndirectAddressed_YIndexed(true)
		A = ReadFromAB()
		SetZNFlags(A)

	//	STX: Store Index X in Memory
	//	X -> M
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x86: //STX Zero Page
		ReadOperands_ZeroPageAddressed()
		WriteToAB(X)
		CPU_Cycles = 3
	case 0x96: //STX Zero Page, Y
		ReadOperands_ZeroPageAddressed_YIndexed()
		WriteToAB(X)
		CPU_Cycles = 4
	case 0x8E: //STX Absolute
		ReadOperands_AbsoluteAddressed(false)
		WriteToAB(X)
		CPU_Cycles = 4

	//	LDX: Load Index X with Memory
	//	M -> X
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xA2: //LDX Immediate
		X = ReadFromPC()
		SetZNFlags(X)
		CPU_Cycles = 2
	case 0xA6: //LDX Zero Page
		ReadOperands_ZeroPageAddressed()
		X = ReadFromAB()
		SetZNFlags(X)
		CPU_Cycles = 3
	case 0xAE: //LDX Absolute
		ReadOperands_AbsoluteAddressed(false)
		X = ReadFromAB()
		SetZNFlags(X)
		CPU_Cycles = 4
	case 0xB6: //LDX Zero Page, Y
		ReadOperands_ZeroPageAddressed_YIndexed()
		X = ReadFromAB()
		SetZNFlags(X)
		CPU_Cycles = 4
	case 0xBE: //LDX Absolute, Y
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_YIndexed(true)
		X = ReadFromAB()
		SetZNFlags(X)

	//	STY: Store Index Y in Memory
	//	Y -> M
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x84: //STY Zero Page
		ReadOperands_ZeroPageAddressed()
		WriteToAB(Y)
		CPU_Cycles = 3
	case 0x94: //STY Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		WriteToAB(Y)
		CPU_Cycles = 4
	case 0x8C: //STY Absolute
		ReadOperands_AbsoluteAddressed(false)
		WriteToAB(Y)
		CPU_Cycles = 4

	//	LDY: Load Index Y with Memory
	//	M -> Y
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xA0: //LDY Immediate
		Y = ReadFromPC()
		SetZNFlags(Y)
		CPU_Cycles = 2
	case 0xA4: //LDY Zero Page
		ReadOperands_ZeroPageAddressed()
		Y = ReadFromAB()
		SetZNFlags(Y)
		CPU_Cycles = 3
	case 0xAC: //LDY Absolute
		ReadOperands_AbsoluteAddressed(false)
		Y = ReadFromAB()
		SetZNFlags(Y)
		CPU_Cycles = 4
	case 0xB4: //LDY Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Y = ReadFromAB()
		SetZNFlags(Y)
		CPU_Cycles = 4
	case 0xBC: //LDY Absolute, X
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_XIndexed(true)
		Y = ReadFromAB()
		SetZNFlags(Y)

	//Transfer

	//	TAX: Transfer Accumulator to Index X
	//	A -> X
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xAA: //TAX
		X = A
		SetZNFlags(X)
		CPU_Cycles = 2

	//	TAY: Transfer Accumulator to Index Y
	//	A -> Y
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xA8: //TAY
		Y = A
		SetZNFlags(Y)
		CPU_Cycles = 2

	//	TXA: Transfer Index X to Accumulator
	//	X -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0x8A: //TXA
		A = X
		SetZNFlags(A)
		CPU_Cycles = 2

	//	TYA: Transfer Index Y to Accumulator
	//	Y -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0x98: //TYA
		A = Y
		SetZNFlags(A)
		CPU_Cycles = 2

	//Arithmetic

	//	ADC: Add Memory to Accumulator with Carry
	//	A + M + C -> A, C
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	+

	case 0x69: //ADC Immediate
		Op_ADC(ReadFromPC())
		CPU_Cycles = 2
	case 0x65: //ADC Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_ADC(ReadFromAB())
		CPU_Cycles = 3
	case 0x75: //ADC Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_ADC(ReadFromAB())
		CPU_Cycles = 4
	case 0x6D: //ADC Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_ADC(ReadFromAB())
		CPU_Cycles = 4
	case 0x7D: //ADC Absolute, X
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_XIndexed(true)
		Op_ADC(ReadFromAB())
	case 0x79: //ADC Absolute, Y
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_YIndexed(true)
		Op_ADC(ReadFromAB())
	case 0x61: //ADC Indirect, X
		ReadOperands_IndirectAddressed_XIndexed()
		Op_ADC(ReadFromAB())
		CPU_Cycles = 6
	case 0x71: //ADC Indirect, Y
		CPU_Cycles = 5
		ReadOperands_IndirectAddressed_YIndexed(true)
		Op_ADC(ReadFromAB())

	//	SBC: Subtract Memory from Accumulator with Borrow
	//	A - M - C̅ -> A
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	+

	case 0xE9: //SBC Immediate
		Op_SBC(ReadFromPC())
		CPU_Cycles = 2
	case 0xE5: //SBC Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_SBC(ReadFromAB())
		CPU_Cycles = 3
	case 0xED: //SBC Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_SBC(ReadFromAB())
		CPU_Cycles = 4
	case 0xF5: //SBC Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_SBC(ReadFromAB())
		CPU_Cycles = 4
	case 0xFD: //SBC Absolute, X
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_XIndexed(true)
		Op_SBC(ReadFromAB())
	case 0xF9: //SBC Absolute, Y
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_YIndexed(true)
		Op_SBC(ReadFromAB())
	case 0xE1: //SBC Indirect, X
		ReadOperands_IndirectAddressed_XIndexed()
		Op_SBC(ReadFromAB())
		CPU_Cycles = 6
	case 0xF1: //SBC Indirect, Y
		CPU_Cycles = 5
		ReadOperands_IndirectAddressed_YIndexed(true)
		Op_SBC(ReadFromAB())

	//	INC: Increment Memory by One
	//	M + 1 -> M
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xE6: //INC Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_INC(AddressBus, ReadFromAB())
		CPU_Cycles = 5
	case 0xEE: //INC Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_INC(AddressBus, ReadFromAB())
		CPU_Cycles = 6
	case 0xF6: //INC Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_INC(AddressBus, ReadFromAB())
		CPU_Cycles = 6
	case 0xFE: //INC Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_INC(AddressBus, ReadFromAB())
		CPU_Cycles = 7

	//	DEC: Decrement Memory by One
	//	M - 1 -> M
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xC6: //DEC Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_DEC(AddressBus, ReadFromAB())
		CPU_Cycles = 5
	case 0xCE: //DEC Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_DEC(AddressBus, ReadFromAB())
		CPU_Cycles = 6
	case 0xD6: //DEC Zeropage, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_DEC(AddressBus, ReadFromAB())
		CPU_Cycles = 6
	case 0xDE: //DEC Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_DEC(AddressBus, ReadFromAB())
		CPU_Cycles = 7

	//	INX: Increment Index X by One
	//	X + 1 -> X
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xE8: //INX
		X++
		CPU_Cycles = 2
		//MasterClockTick("inx")
		SetZNFlags(X)

	//	DEX: Decrement Index X by One
	//	X - 1 -> X
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xCA: //DEX
		X--
		CPU_Cycles = 2
		//MasterClockTick("dex")
		SetZNFlags(X)

	//	INY: Increment Index Y by One
	//	Y + 1 -> Y
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xC8: //INY
		Y++
		CPU_Cycles = 2
		//MasterClockTick("iny")
		SetZNFlags(Y)

	//	DEY: Decrement Index Y by One
	//	Y - 1 -> Y
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0x88: //DEY
		Y--
		CPU_Cycles = 2
		//MasterClockTick("dey")
		SetZNFlags(Y)

	//Shift

	//	ASL: Shift Left One Bit (Memory or Accumulator)
	//	C <- [76543210] <- 0
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	case 0x0A: //ASL A
		flag_Carry = A > 127
		A <<= 1
		SetZNFlags(A)
		CPU_Cycles = 2
	case 0x06: //ASL Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_ASL(AddressBus)
		CPU_Cycles = 5
	case 0x0E: //ASL Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_ASL(AddressBus)
		CPU_Cycles = 6
	case 0x16: //ASL Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_ASL(AddressBus)
		CPU_Cycles = 6
	case 0x1E: //ASL Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_ASL(AddressBus)
		CPU_Cycles = 7

	//	LSR: Shift One Bit Right (Memory or Accumulator)
	//	0 -> [76543210] -> C
	//	N	Z	C	I	D	V
	//	0	+	+	-	-	-

	case 0x4A: //LSR A
		flag_Carry = (A & 1) != 0
		A >>= 1
		SetZNFlags(A)
		CPU_Cycles = 2
	case 0x46: //LSR Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_LSR(AddressBus)
		CPU_Cycles = 5
	case 0x4E: //LSR Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_LSR(AddressBus)
		CPU_Cycles = 6
	case 0x56: //LSR Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_LSR(AddressBus)
		CPU_Cycles = 6
	case 0x5E: //LSR Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_LSR(AddressBus)
		CPU_Cycles = 7

	//	ROL: Rotate One Bit Left (Memory or Accumulator)
	//	C <- [76543210] <- C
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	case 0x2A: //ROL A
		futureCarry := (A >= 0x80)
		A <<= 1
		if flag_Carry {
			A |= 1
		}
		flag_Carry = futureCarry
		SetZNFlags(A)
		CPU_Cycles = 2
	case 0x26: //ROL Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_ROL(AddressBus)
		CPU_Cycles = 5
	case 0x2E: //ROL Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_ROL(AddressBus)
		CPU_Cycles = 6
	case 0x36: //ROL Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_ROL(AddressBus)
		CPU_Cycles = 6
	case 0x3E: //ROL Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_ROL(AddressBus)
		CPU_Cycles = 7

	//	ROR: Rotate One Bit Right (Memory or Accumulator)
	//	C -> [76543210] -> C
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	case 0x6A: //ROR A
		futureCarry := (A & 1) != 0
		A >>= 1
		if flag_Carry {
			A |= 0x80
		}
		flag_Carry = futureCarry
		SetZNFlags(A)
		CPU_Cycles = 2
	case 0x66: //ROR Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_ROR(AddressBus)
		CPU_Cycles = 5
	case 0x6E: //ROR Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_ROR(AddressBus)
		CPU_Cycles = 6
	case 0x76: //ROR Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_ROR(AddressBus)
		CPU_Cycles = 6
	case 0x7E: //ROR Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_ROR(AddressBus)
		CPU_Cycles = 7

	//Bitwise

	//	AND: AND Memory with Accumulator
	//	A AND M -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0x29: //AND Immediate
		Op_AND(ReadFromPC())
		CPU_Cycles = 2
	case 0x25: //AND Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_AND(ReadFromAB())
		CPU_Cycles = 3
	case 0x2D: //AND Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_AND(ReadFromAB())
		CPU_Cycles = 4
	case 0x35: //AND Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_AND(ReadFromAB())
		CPU_Cycles = 4
	case 0x3D: //AND Absolute, X
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_XIndexed(true)
		Op_AND(ReadFromAB())
	case 0x39: //AND Absolute, Y
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_YIndexed(true)
		Op_AND(ReadFromAB())
	case 0x21: //AND Indirect, X
		ReadOperands_IndirectAddressed_XIndexed()
		Op_AND(ReadFromAB())
		CPU_Cycles = 6
	case 0x31: //AND Indirect, Y
		CPU_Cycles = 5
		ReadOperands_IndirectAddressed_YIndexed(true)
		Op_AND(ReadFromAB())

	//	ORA: OR Memory with Accumulator
	//	A OR M -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0x09: //ORA Immediate
		Op_ORA(ReadFromPC())
		CPU_Cycles = 2
	case 0x05: //ORA Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_ORA(ReadFromAB())
		CPU_Cycles = 3
	case 0x0D: //ORA Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_ORA(ReadFromAB())
		CPU_Cycles = 4
	case 0x15: //ORA Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_ORA(ReadFromAB())
		CPU_Cycles = 4
	case 0x1D: //ORA Absolute, X
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_XIndexed(true)
		Op_ORA(ReadFromAB())
	case 0x19: //ORA Absolute, Y
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_YIndexed(true)
		Op_ORA(ReadFromAB())
	case 0x01: //ORA Indirect, X
		ReadOperands_IndirectAddressed_XIndexed()
		Op_ORA(ReadFromAB())
		CPU_Cycles = 6
	case 0x11: //ORA Indirect, Y
		CPU_Cycles = 5
		ReadOperands_IndirectAddressed_YIndexed(true)
		Op_ORA(ReadFromAB())

	//	EOR: Exclusive-OR Memory with Accumulator
	//	A EOR M -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0x49: //EOR Immediate
		Op_EOR(ReadFromPC())
		CPU_Cycles = 2
	case 0x45: //EOR Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_EOR(ReadFromAB())
		CPU_Cycles = 3
	case 0x4D: //EOR Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_EOR(ReadFromAB())
		CPU_Cycles = 4
	case 0x55: //EOR Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_EOR(ReadFromAB())
		CPU_Cycles = 4
	case 0x5D: //EOR Absolute, X
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_XIndexed(true)
		Op_EOR(ReadFromAB())
	case 0x59: //EOR Absolute, Y
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_YIndexed(true)
		Op_EOR(ReadFromAB())
	case 0x41: //EOR Indirect, X
		ReadOperands_IndirectAddressed_XIndexed()
		Op_EOR(ReadFromAB())
		CPU_Cycles = 6
	case 0x51: //EOR Indirect, Y
		CPU_Cycles = 5
		ReadOperands_IndirectAddressed_YIndexed(true)
		Op_EOR(ReadFromAB())

	//	BIT: //Test Bits in Memory with Accumulator
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

	case 0x24: //BIT Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_BIT(ReadFromAB())
		CPU_Cycles = 3
	case 0x2C: //BIT Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_BIT(ReadFromAB())
		CPU_Cycles = 4

	//Compare

	//	CMP: Compare Memory with Accumulator
	//	A - M
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	case 0xC9: //CMP Immediate
		Op_CMP(ReadFromPC())
		CPU_Cycles = 2
	case 0xC5: //CMP Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_CMP(ReadFromAB())
		CPU_Cycles = 3
	case 0xCD: //CMP Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_CMP(ReadFromAB())
		CPU_Cycles = 4
	case 0xD5: //CMP Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_CMP(ReadFromAB())
		CPU_Cycles = 4
	case 0xDD: //CMP Absolute, X
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_XIndexed(true)
		Op_CMP(ReadFromAB())
	case 0xD9: //CMP Absolute, Y
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_YIndexed(true)
		Op_CMP(ReadFromAB())
	case 0xC1: //CMP Indirect, X
		ReadOperands_IndirectAddressed_XIndexed()
		Op_CMP(ReadFromAB())
		CPU_Cycles = 6
	case 0xD1: //CMP Indirect, Y
		CPU_Cycles = 5
		ReadOperands_IndirectAddressed_YIndexed(true)
		Op_CMP(ReadFromAB())

	//	CPX: Compare Memory and Index X
	//	X - M
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	case 0xE0: //CPX Immediate
		Op_CPX(ReadFromPC())
		CPU_Cycles = 2
	case 0xE4: //CPX Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_CPX(ReadFromAB())
		CPU_Cycles = 3
	case 0xEC: //CPX Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_CPX(ReadFromAB())
		CPU_Cycles = 4

	//	CPY: Compare Memory and Index Y
	//	Y - M
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	case 0xC0: //CPY Immediate
		Op_CPY(ReadFromPC())
		CPU_Cycles = 2
	case 0xC4: //CPY Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_CPY(ReadFromAB())
		CPU_Cycles = 3
	case 0xCC: //CPY Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_CPY(ReadFromAB())
		CPU_Cycles = 4

	//Branch

	//	BCC: Branch on Carry Clear
	//	branch on C = 0
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x90: //BCC (Branch on Carry Clear)
		temp := ReadFromPC()
		Branch(!flag_Carry, temp)

	//	BCS: Branch on Carry Set
	//	branch on C = 1
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0xB0: //BCS (Branch on Carry Set)
		temp := ReadFromPC()
		Branch(flag_Carry, temp)

	//	BEQ: Branch on Result Zero
	//	branch on Z = 1
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0xF0: //BEQ (Branch on Equal)
		temp := ReadFromPC()
		Branch(flag_Zero, temp)

	//	BNE: Branch on Result not Zero
	//	branch on Z = 0
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0xD0: //BNE (Branch on Not Equal)
		temp := ReadFromPC()
		Branch(!flag_Zero, temp)

	//	BPL: Branch on Result Plus
	//	branch on N = 0
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x10: //BPL (Branch on Plus)
		temp := ReadFromPC()
		Branch(!flag_Negative, temp)

	//	BMI: Branch on Result Minus
	//	branch on N = 1
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x30: //BMI (Branch on Minus)
		temp := ReadFromPC()
		Branch(flag_Negative, temp)

	//	BVC: Branch on Overflow Clear
	//	branch on V = 0
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x50: //BVC (Branch on Overflow Clear)
		temp := ReadFromPC()
		Branch(!flag_Overflow, temp)

	//	BVS: Branch on Overflow Set
	//	branch on V = 1
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x70: //BVS (Branch on Overflow Set)
		temp := ReadFromPC()
		Branch(flag_Overflow, temp)

	//Jump

	//	JMP: Jump to New Location
	//	operand 1st byte -> PCL
	//	operand 2nd byte -> PCH
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x4C: //JMP
		ReadOperands_AbsoluteAddressed(true)
		PC = AddressBus
		CPU_Cycles = 3

	case 0x6C: //JMP Indirect
		ReadOperands_IndirectAddressed()
		PC = AddressBus
		CPU_Cycles = 5 //TODO: What the fuck

	//	JSR: Jump to New Location Saving Return Address
	//	push (PC+2),
	//	operand 1st byte -> PCL
	//	operand 2nd byte -> PCH
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x20: //JSR
		temp_low := ReadFromPC()
		Push(byte(PC / 0x100))
		Push(byte(PC))
		temp_high := ReadFromPC()
		PC = BuildAddress(temp_low, temp_high)
		//MasterClockTick("jsr")
		CPU_Cycles = 6

	//	RTS: Return from Subroutine
	//	pull PC, PC+1 -> PC
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x60: //RTS
		temp_low := Pull()
		temp_high := Pull()
		//MasterClockTick("rts Pull1")
		//MasterClockTick("rts Pull2")
		PC = BuildAddress(temp_low, temp_high)
		PC++
		//MasterClockTick("rts pc++")
		CPU_Cycles = 6

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

	case 0x00: //BRK
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
		CPU_Cycles = 7

	//	RTI: Return from Interrupt
	//	The status register is pulled with the break flag
	//	and bit 5 ignored. Then PC is pulled from the stack.
	//
	//	pull SR, pull PC
	//	N	Z	C	I	D	V
	//	from stack

	case 0x40: //RTI
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
		CPU_Cycles = 6

	//Stack

	//	PHA: Push Accumulator on Stack
	//	push A
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x48: //PHA
		Push(A)
		CPU_Cycles = 3

	//	PLA: Pull Accumulator from Stack
	//	pull A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0x68: //PLA
		A = Pull()
		//MasterClockTick("pla")
		SetZNFlags(A)
		CPU_Cycles = 4

	//	PHP: Push Processor Status on Stack
	//	The status register will be pushed with the break
	//	flag and bit 5 set to 1.
	//	push SR
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-
	case 0x08: //PHP
		flag_B = true
		PushFlags()
		CPU_Cycles = 3

	//	PLP: Pull Processor Status from Stack
	//	The status register will be pulled with the break
	//	flag and bit 5 ignored.
	//	pull SR
	//	N	Z	C	I	D	V
	//	from stack

	case 0x28: //PLP
		PullFlags()
		CPU_Cycles = 4

	//	TSX: Transfer Stack Pointer to Index X
	//	SP -> X
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xBA: //TSX
		X = SP
		SetZNFlags(X)
		CPU_Cycles = 2

	//	TXS: Transfer Index X to Stack Register
	//	X -> SP
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x9A: //TXS
		SP = X
		CPU_Cycles = 2

	//Flags

	//	CLC: Clear Carry Flag
	//	0 -> C
	//	N	Z	C	I	D	V
	//	-	-	0	-	-	-

	case 0x18: //CLC
		flag_Carry = false
		//MasterClockTick("clc")
		CPU_Cycles = 2

	//	SEC: Set Carry Flag
	//	1 -> C
	//	N	Z	C	I	D	V
	//	-	-	1	-	-	-

	case 0x38: //SEC
		flag_Carry = true
		//MasterClockTick("sec")
		CPU_Cycles = 2

	//	CLI: Clear Interrupt Disable Bit
	//	0 -> I
	//	N	Z	C	I	D	V
	//	-	-	-	0	-	-

	case 0x58: //CLI
		flag_InterruptDisable = false
		//MasterClockTick("cli")
		CPU_Cycles = 2

	//	SEI: Set Interrupt Disable Status
	//	1 -> I
	//	N	Z	C	I	D	V
	//	-	-	-	1	-	-

	case 0x78: //SEI
		flag_InterruptDisable = true
		//MasterClockTick("sei")
		CPU_Cycles = 2

	//	CLD: Clear Decimal Mode
	//	0 -> D
	//	N	Z	C	I	D	V
	//	-	-	-	-	0	-

	case 0xD8: //CLD
		flag_Decimal = false
		//MasterClockTick("cld")
		CPU_Cycles = 2

	//	SED: Set Decimal Flag
	//	1 -> D
	//	N	Z	C	I	D	V
	//	-	-	-	-	1	-

	case 0xF8: //SED
		flag_Decimal = true
		//MasterClockTick("sed")
		CPU_Cycles = 2

	//	CLV: Clear Overflow Flag
	//	0 -> V
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	0

	case 0xB8: //CLV
		flag_Overflow = false
		//MasterClockTick("clv")
		CPU_Cycles = 2

	//Other

	//	NOP: No Operation
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0xEA: //NOP
		CPU_Cycles = 2

	//----------------------------------------------------------------------
	//Unofficial Opcodes
	//----------------------------------------------------------------------

	//HLT Codes

	case 0x02: //HLT
		CPU_Halted = true
	case 0x12: //HLT
		CPU_Halted = true
	case 0x22: //HLT
		CPU_Halted = true
	case 0x32: //HLT
		CPU_Halted = true
	case 0x42: //HLT
		CPU_Halted = true
	case 0x52: //HLT
		CPU_Halted = true
	case 0x62: //HLT
		CPU_Halted = true
	case 0x72: //HLT
		CPU_Halted = true
	case 0x92: //HLT
		CPU_Halted = true
	case 0xB2: //HLT
		CPU_Halted = true
	case 0xD2: //HLT
		CPU_Halted = true
	case 0xF2: //HLT
		CPU_Halted = true

	//NOP Codes (unofficial)

	case 0x1A: //NOP Implied
		CPU_Cycles = 2
	case 0x3A: //NOP Implied
		CPU_Cycles = 2
	case 0x5A: //NOP Implied
		CPU_Cycles = 2
	case 0x7A: //NOP Implied
		CPU_Cycles = 2
	case 0xDA: //NOP Implied
		CPU_Cycles = 2
	case 0xFA: //NOP Implied
		CPU_Cycles = 2

	case 0x80: //NOP Immediate
		ReadFromPC()
		CPU_Cycles = 2
	case 0x82: //NOP Immediate
		ReadFromPC()
		CPU_Cycles = 2
	case 0x89: //NOP Immediate
		ReadFromPC()
		CPU_Cycles = 2
	case 0xC2: //NOP Immediate
		ReadFromPC()
		CPU_Cycles = 2
	case 0xE2: //NOP Immediate
		ReadFromPC()
		CPU_Cycles = 2

	case 0x04: //NOP Zero Page
		ReadOperands_ZeroPageAddressed()
		CPU_Cycles = 3
	case 0x44: //NOP Zero Page
		ReadOperands_ZeroPageAddressed()
		CPU_Cycles = 3
	case 0x64: //NOP Zero Page
		ReadOperands_ZeroPageAddressed()
		CPU_Cycles = 3

	case 0x14: //NOP Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		CPU_Cycles = 4
	case 0x34: //NOP Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		CPU_Cycles = 4
	case 0x54: //NOP Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		CPU_Cycles = 4
	case 0x74: //NOP Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		CPU_Cycles = 4
	case 0xD4: //NOP Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		CPU_Cycles = 4
	case 0xF4: //NOP Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		CPU_Cycles = 4

	case 0x0C: //NOP Absolute
		ReadOperands_AbsoluteAddressed(false)
		CPU_Cycles = 4

	case 0x1C: //NOP Absolute, X
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_XIndexed(true)
	case 0x3C: //NOP Absolute, X
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_XIndexed(true)
	case 0x5C: //NOP Absolute, X
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_XIndexed(true)
	case 0x7C: //NOP Absolute, X
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_XIndexed(true)
	case 0xDC: //NOP Absolute, X
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_XIndexed(true)
	case 0xFC: //NOP Absolute, X
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_XIndexed(true)

	//SAX: A AND X -> M

	case 0x87: //SAX Zero Page
		ReadOperands_ZeroPageAddressed()
		bus.Write(AddressBus, (A & X))
		CPU_Cycles = 3
	case 0x97: //SAX Zero Page, Y
		ReadOperands_ZeroPageAddressed_YIndexed()
		bus.Write(AddressBus, (A & X))
		CPU_Cycles = 4
	case 0x8F: //SAX Absolute
		ReadOperands_AbsoluteAddressed(false)
		bus.Write(AddressBus, (A & X))
		CPU_Cycles = 4
	case 0x83: //SAX Indirect, X
		ReadOperands_IndirectAddressed_XIndexed()
		bus.Write(AddressBus, (A & X))
		CPU_Cycles = 6

	//LAX: LDA + LDX

	case 0xA7: //LAX Zero Page
		ReadOperands_ZeroPageAddressed()
		A = ReadFromAB()
		X = A
		SetZNFlags(X)
		CPU_Cycles = 3
	case 0xB7: //LAX Zero Page, Y
		ReadOperands_ZeroPageAddressed_YIndexed()
		A = ReadFromAB()
		X = A
		SetZNFlags(X)
		CPU_Cycles = 4
	case 0xAF: //LAX Absolute
		ReadOperands_AbsoluteAddressed(false)
		A = ReadFromAB()
		X = A
		SetZNFlags(X)
		CPU_Cycles = 4
	case 0xBF: //LAX Absolute, Y
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_YIndexed(true)
		A = ReadFromAB()
		X = A
		SetZNFlags(X)
	case 0xA3: //LAX Indirect, X
		ReadOperands_IndirectAddressed_XIndexed()
		A = ReadFromAB()
		X = A
		SetZNFlags(X)
		CPU_Cycles = 6
	case 0xB3: //LAX Indirect, Y
		CPU_Cycles = 5
		ReadOperands_IndirectAddressed_YIndexed(true)
		A = ReadFromAB()
		X = A
		SetZNFlags(X)

	//SLO: ASL + ORA

	case 0x07: // SLO Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_ASL(AddressBus)
		Op_ORA(ReadFromAB())
		CPU_Cycles = 5
	case 0x17: //SLO Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_ASL(AddressBus)
		Op_ORA(ReadFromAB())
		CPU_Cycles = 6
	case 0x0F: //SLO Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_ASL(AddressBus)
		Op_ORA(ReadFromAB())
		CPU_Cycles = 6
	case 0x1F: //Slo Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_ASL(AddressBus)
		Op_ORA(ReadFromAB())
		CPU_Cycles = 7
	case 0x1B: //Slo Absolute, Y
		ReadOperands_AbsoluteAddressed_YIndexed(false)
		Op_ASL(AddressBus)
		Op_ORA(ReadFromAB())
		CPU_Cycles = 7
	case 0x03: //Slo Indirect, X
		ReadOperands_IndirectAddressed_XIndexed()
		Op_ASL(AddressBus)
		Op_ORA(ReadFromAB())
		CPU_Cycles = 8
	case 0x13: //Slo Indirect, Y
		ReadOperands_IndirectAddressed_YIndexed(false)
		Op_ASL(AddressBus)
		Op_ORA(ReadFromAB())
		CPU_Cycles = 8

	//DCP: DEC + CMP

	case 0xC7: //DCP Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_DEC(AddressBus, ReadFromAB())
		Op_CMP(ReadFromAB())
		CPU_Cycles = 5
	case 0xD7: //DCP Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_DEC(AddressBus, ReadFromAB())
		Op_CMP(ReadFromAB())
		CPU_Cycles = 6
	case 0xCF: //DCP Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_DEC(AddressBus, ReadFromAB())
		Op_CMP(ReadFromAB())
		CPU_Cycles = 6
	case 0xDF: //DCP Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_DEC(AddressBus, ReadFromAB())
		Op_CMP(ReadFromAB())
		CPU_Cycles = 7
	case 0xDB: //DCP Absolute, Y
		ReadOperands_AbsoluteAddressed_YIndexed(false)
		Op_DEC(AddressBus, ReadFromAB())
		Op_CMP(ReadFromAB())
		CPU_Cycles = 7
	case 0xC3: //DCP Indirect, X
		ReadOperands_IndirectAddressed_XIndexed()
		Op_DEC(AddressBus, ReadFromAB())
		Op_CMP(ReadFromAB())
		CPU_Cycles = 8
	case 0xD3: //DCP Indirect, Y
		ReadOperands_IndirectAddressed_YIndexed(false)
		Op_DEC(AddressBus, ReadFromAB())
		Op_CMP(ReadFromAB())
		CPU_Cycles = 8

	//SHA: Stores A AND X AND (high-byte of addr. + 1) at addr.

	case 0x9F: //SHA Absolute, Y
		ReadOperands_AbsoluteAddressed_YIndexed(false)
		HiByte := byte(AddressBus >> 8)
		bus.Write(AddressBus, A&X&(HiByte+1))
		CPU_Cycles = 5
	case 0x93: //SHA Indirect, Y
		ReadOperands_IndirectAddressed_YIndexed(false)
		bus.Write(AddressBus, A&X&byte((AddressBus>>8)+1))
		CPU_Cycles = 6

	//SHX: Stores X AND (high-byte of addr. + 1) at addr.

	case 0x9E: //SHX Absolute, Y
		ReadOperands_AbsoluteAddressed_YIndexed(false)
		val := (X & byte(((AddressBus&0xFF00)>>8)+1))
		bus.Write(AddressBus, val)
		CPU_Cycles = 5

	//SHY: Stores Y AND (high-byte of addr. + 1) at addr.

	case 0x9C: //SHY Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(false)
		val := (Y & byte(((AddressBus&0xFF00)>>8)+1))
		bus.Write(AddressBus, val)
		CPU_Cycles = 5

	//RLA: ROL + AND

	case 0x27: //RLA Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_ROL(AddressBus)
		Op_AND(ReadFromAB())
		CPU_Cycles = 5
	case 0x37: //RLA Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_ROL(AddressBus)
		Op_AND(ReadFromAB())
		CPU_Cycles = 6
	case 0x2F: //RLA Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_ROL(AddressBus)
		Op_AND(ReadFromAB())
		CPU_Cycles = 6
	case 0x3F: //RLA Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_ROL(AddressBus)
		Op_AND(ReadFromAB())
		CPU_Cycles = 7
	case 0x3B: //RLA Absolute, Y
		ReadOperands_AbsoluteAddressed_YIndexed(false)
		Op_ROL(AddressBus)
		Op_AND(ReadFromAB())
		CPU_Cycles = 7
	case 0x23: //RLA Indirect, X
		ReadOperands_IndirectAddressed_XIndexed()
		Op_ROL(AddressBus)
		Op_AND(ReadFromAB())
		CPU_Cycles = 8
	case 0x33: //RLA Indirect, Y
		ReadOperands_IndirectAddressed_YIndexed(false)
		Op_ROL(AddressBus)
		Op_AND(ReadFromAB())
		CPU_Cycles = 8

	//SRE: LSR + EOR

	case 0x47: //SRE Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_LSR(AddressBus)
		Op_EOR(ReadFromAB())
		CPU_Cycles = 5
	case 0x57: //SRE Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_LSR(AddressBus)
		Op_EOR(ReadFromAB())
		CPU_Cycles = 6
	case 0x4F: //SRE Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_LSR(AddressBus)
		Op_EOR(ReadFromAB())
		CPU_Cycles = 6
	case 0x5F: //SRE Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_LSR(AddressBus)
		Op_EOR(ReadFromAB())
		CPU_Cycles = 7
	case 0x5B: //SRE Absolute, Y
		ReadOperands_AbsoluteAddressed_YIndexed(false)
		Op_LSR(AddressBus)
		Op_EOR(ReadFromAB())
		CPU_Cycles = 7
	case 0x43: //SRE Indirect, X
		ReadOperands_IndirectAddressed_XIndexed()
		Op_LSR(AddressBus)
		Op_EOR(ReadFromAB())
		CPU_Cycles = 8
	case 0x53: //SRE Indirect, Y
		ReadOperands_IndirectAddressed_YIndexed(false)
		Op_LSR(AddressBus)
		Op_EOR(ReadFromAB())
		CPU_Cycles = 8

	//RRA: ROR + ADC

	case 0x67: //RRA Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_ROR(AddressBus)
		Op_ADC(ReadFromAB())
		CPU_Cycles = 5
	case 0x77: //RRA Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_ROR(AddressBus)
		Op_ADC(ReadFromAB())
		CPU_Cycles = 6
	case 0x6F: //RRA Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_ROR(AddressBus)
		Op_ADC(ReadFromAB())
		CPU_Cycles = 6
	case 0x7F: //RRA Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_ROR(AddressBus)
		Op_ADC(ReadFromAB())
		CPU_Cycles = 7
	case 0x7B: //RRA Absolute, Y
		ReadOperands_AbsoluteAddressed_YIndexed(false)
		Op_ROR(AddressBus)
		Op_ADC(ReadFromAB())
		CPU_Cycles = 7
	case 0x63: //RRA Indirect, X
		ReadOperands_IndirectAddressed_XIndexed()
		Op_ROR(AddressBus)
		Op_ADC(ReadFromAB())
		CPU_Cycles = 8
	case 0x73: //RRA Indirect, Y
		ReadOperands_IndirectAddressed_YIndexed(false)
		Op_ROR(AddressBus)
		Op_ADC(ReadFromAB())
		CPU_Cycles = 8

	//ISC: INC + SBC

	case 0xE7: //ISC Zero Page
		ReadOperands_ZeroPageAddressed()
		Op_INC(AddressBus, ReadFromAB())
		Op_SBC(ReadFromAB())
		CPU_Cycles = 5
	case 0xF7: //ISC Zero Page, X
		ReadOperands_ZeroPageAddressed_XIndexed()
		Op_INC(AddressBus, ReadFromAB())
		Op_SBC(ReadFromAB())
		CPU_Cycles = 6
	case 0xEF: //ISC Absolute
		ReadOperands_AbsoluteAddressed(false)
		Op_INC(AddressBus, ReadFromAB())
		Op_SBC(ReadFromAB())
		CPU_Cycles = 6
	case 0xFF: //ISC Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_INC(AddressBus, ReadFromAB())
		Op_SBC(ReadFromAB())
		CPU_Cycles = 7
	case 0xFB: //ISC Absolute, Y
		ReadOperands_AbsoluteAddressed_YIndexed(false)
		Op_INC(AddressBus, ReadFromAB())
		Op_SBC(ReadFromAB())
		CPU_Cycles = 7
	case 0xE3: //ISC Indirect, X
		ReadOperands_IndirectAddressed_XIndexed()
		Op_INC(AddressBus, ReadFromAB())
		Op_SBC(ReadFromAB())
		CPU_Cycles = 8
	case 0xF3: //ISC Indirect, Y
		ReadOperands_IndirectAddressed_YIndexed(false)
		Op_INC(AddressBus, ReadFromAB())
		Op_SBC(ReadFromAB())
		CPU_Cycles = 8

	//Immediates (unofficial)

	case 0x0B: //ANC Immediate
		//AND + Set Carry as ASL
		Op_AND(ReadFromPC())
		flag_Carry = flag_Negative
		CPU_Cycles = 2
	case 0x2B: //ANC2 Immediate
		//AND + Set Carry as ROL
		//Same as $0B
		Op_AND(ReadFromPC())
		flag_Carry = flag_Negative
		CPU_Cycles = 2
	case 0x4B: //ALR Immediate
		//AND + LSR
		Op_AND(ReadFromPC())
		flag_Carry = (A & 1) != 0
		A >>= 1
		SetZNFlags(A)
		CPU_Cycles = 2
	case 0x6B: //ARR Immediate
		//AND + ROR
		Op_AND(ReadFromPC())
		flag_Overflow = A == 0

		A >>= 1
		if flag_Carry {
			A |= 0x80
		}
		flag_Carry = ((A & 0x40) >> 6) == 1
		flag_Overflow = (((A & 0x20) >> 5) ^ ((A & 0x40) >> 6)) == 1

		SetZNFlags(A)
		CPU_Cycles = 2
	case 0x8B: //ANE Immediate
		//Highly unstable
		A = X & ReadFromPC()
		SetZNFlags(A)
		CPU_Cycles = 2
	case 0xAB: //LXA Immediate
		//Highly unstable
		A = ReadFromPC()
		X = A
		SetZNFlags(A)
		CPU_Cycles = 2
	case 0xCB: //SBX Immediate
		//(A AND X) - oper -> X
		X = (A & X) - ReadFromPC()
		Op_CMP(X)
		SetZNFlags(X)
		CPU_Cycles = 2
	case 0xEB: //SBC Immediate
		Op_SBC(ReadFromPC())
		CPU_Cycles = 2

	case 0x9B: //TAS (SHS) Absolute, Y
		//A AND X -> SP, A AND X AND (H+1) -> M
		temp := A & X
		Push(temp)
		ReadOperands_AbsoluteAddressed_YIndexed(false)
		bus.Write(AddressBus, (A&X)&(byte(AddressBus&0xFF00)+1))
		CPU_Cycles = 5

	case 0xBB: //LAS (LAE) Absolute, Y
		//LDA/TSX oper
		//M AND SP -> A, X, SP
		CPU_Cycles = 4
		ReadOperands_AbsoluteAddressed_YIndexed(true)
		Value := (ReadFromAB() & SP)
		A = Value
		X = Value
		SP = Value
		SetZNFlags(A)

	default:
		fmt.Println("Unknown Opcode: " + fmt.Sprintf("%02X", opcode))
	}

	// Every CPU instruction takes at least 2 cycles, guaranteed
	// The first is from fetching the opcode
	// This will be the second
	//MasterClockTick(g)

	//if LoggingCPU {
	//	TraceLogger()
	//}

	if CPU_Cycles != CPU_Cycles_New {
		//fmt.Print(cycleTest)
		//fmt.Printf("Cycle Mismatch! [%02X] %d/%d \n", opcode, CPU_Cycles, CPU_Cycles_New)
	}
	//cycleTest = ""

	//Check for, and perform Interrupt Request (IRQ)
	if (apu.APUDMCInterrupt || apu.APUFrameInterrupt) && !DoNMI && !flag_InterruptDisable {
		flag_B = false
		Push(byte(PC >> 8))
		Push(byte(PC))
		PushFlags()

		PCL := bus.Read(0xFFFE)
		PCH := bus.Read(0xFFFF)
		PC = BuildAddress(PCL, PCH)

		//Disable interrupts
		apu.APUDMCInterrupt = false
		apu.APUFrameInterrupt = false
		DoNMI = false
	}

	//CartRAMLogger()
	operands = nil

	common.CPU_TotalCycles += CPU_Cycles
	for CPU_Cycles > 0 {
		CPU_Cycles--
		ppu.Emulate_PPU()
		ppu.Emulate_PPU()
		ppu.Emulate_PPU()

		//Run the APU
		//apuRun = !apuRun
		//if apuRun {
		apu.Emulate_APU()
		//}
	}

	//Force Stop, just in case
	//if TotalCycles > 900000 {

	//if InstructionCount > 50000 {
	//fmt.Println("Too many cycles, end")
	//CPU_Halted = true
	//}

	//InstructionCount++
}
