package nes

import (
	"fmt"
	"mtt/timenes/common"
	"mtt/timenes/nes/apu"
	"mtt/timenes/nes/ppu"
)

var ProgramCounter uint16
var StackPointer byte
var A, X, Y byte
var opcode byte
var operands []byte
var CPU_Cycles, CPU_Cycles_New int
var CPU_Halted = false
var flag_Carry, flag_Zero, flag_InterruptDisable, flag_Decimal, flag_Overflow, flag_Negative, flag_B bool
var MagicConstant byte = 0xFF

var apuRun bool
var AddressBus uint16

func SetZNFlags(Value byte) {
	flag_Zero = (Value == 0x00)
	flag_Negative = (Value >= 0x80)
}

//TODO: Check the ReadOperandAbsolute functions to add a CPU cycle if the high byte was altered

func ReadOperands_AbsoluteAddressed(isJMP bool) uint16 {
	AddressBus := uint16(ReadFromPC())
	//MasterClockTick("Abs add")
	AddressBus = (uint16(ReadFromPC())<<8 | AddressBus)
	if !isJMP {
		//MasterClockTick("abs add, not jmp")
	}
	return AddressBus
}

func ReadOperands_IndirectAddressed() uint16 {
	/*Addr := ReadFromPC()
	TempAddress := Addr
	Addr = Read(uint16(TempAddress)) //Low byte of new address
	TempAddress++
	AddressBus = (uint16(Read(uint16(TempAddress)))<<8 | uint16(Addr)) //High byte
	return AddressBus*/
	AddressBus := uint16(ReadFromPC())
	AddressBus = (uint16(ReadFromPC())<<8 | AddressBus)
	//Now read from HERE
	indL := Read(AddressBus)
	var indH byte
	if AddressBus&0x00FF == 0xFF {
		//Original NMOS Bug
		indH = Read(AddressBus & 0xFF00)
	} else {
		indH = Read(AddressBus + 1)
	}
	//MasterClockTick("ind add")
	return BuildAddress(indL, indH)
}

func ReadOperands_AbsoluteAddressed_XIndexed(pbCheck bool) uint16 {
	low := ReadFromPC()
	high := ReadFromPC()
	AddressBus = BuildAddress(low, high)
	AddressBus += uint16(X)
	if pbCheck && byte((AddressBus&0xFF00)>>4) != high {
		////MasterClockTick("abs add x")
		CPU_Cycles++ //Extra cycle for crossing page boundary
	}

	return AddressBus
}

func ReadOperands_AbsoluteAddressed_YIndexed(pbCheck bool) uint16 {
	low := ReadFromPC()
	high := ReadFromPC()
	AddressBus = BuildAddress(low, high)
	AddressBus += uint16(Y)
	if pbCheck && byte((AddressBus&0xFF00)>>4) != high {
		////MasterClockTick("abs add y")
		CPU_Cycles++ //Extra cycle for crossing page boundary
	}
	return AddressBus
}

func ReadOperands_IndirectAddressed_XIndexed() uint16 {
	Addr := ReadFromPC() + X
	////MasterClockTick("ind x")
	TempAddress := Addr
	Addr = Read(uint16(TempAddress)) //Low byte of new address
	TempAddress++
	AddressBus = (uint16(Read(uint16(TempAddress)))<<8 | uint16(Addr)) //High byte
	return AddressBus
}

func ReadOperands_IndirectAddressed_YIndexed(pbCheck bool) uint16 {
	Addr := ReadFromPC()
	TempAddress := Addr
	Addr = Read(uint16(TempAddress)) //Low byte of new address
	TempAddress++
	AddressBus = (uint16(Read(uint16(TempAddress)))<<8 | uint16(Addr)) //High byte
	AddressBus += uint16(Y)
	if pbCheck && Read(uint16(TempAddress)) != byte((AddressBus&0xFF00)>>4) {
		////MasterClockTick("ind add y")
		CPU_Cycles++
	}
	return AddressBus
}

func ReadOperands_ZeroPageAddressed() uint16 {
	AddressBus := ReadFromPC()
	return BuildAddress(AddressBus, 0x00)
}

func ReadOperands_ZeroPageAddressed_XIndexed() uint16 {
	AddressBus := ReadFromPC()
	////MasterClockTick("Zp x")
	return BuildAddress(AddressBus+X, 0x00)
}

func ReadOperands_ZeroPageAddressed_YIndexed() uint16 {
	AddressBus := ReadFromPC()
	////MasterClockTick("ZP Y")
	return BuildAddress(AddressBus+Y, 0x00)
}

func Branch(Condition bool, Value byte) {
	if Condition {
		//fmt.Println("branch taken")
		signedVal := int(Value)
		if signedVal > 127 {
			signedVal -= 256 //range from -128 to 127
		}
		CPU_Cycles = 3
		if byte((ProgramCounter&0xFF00)>>4) != byte(((ProgramCounter+uint16(signedVal))&0xFF00)>>4) {
			CPU_Cycles++ //Extra cycle for crossing page boundary
			//MasterClockTick("branch page cross")
		}
		ProgramCounter = uint16(ProgramCounter + uint16(signedVal))
		//MasterClockTick("branch taken")
	} else {
		//fmt.Println("branch not taken")
		CPU_Cycles = 2
	}
}

func Push(Value byte) {
	//Store to the stack, and decrement the stack pointer
	Write(uint16(StackPointer)+0x100, Value)
	////MasterClockTick("push")
	StackPointer--
}

func Pull() byte {
	//Increment the stack pointer, and read from the stack
	StackPointer++
	//MasterClockTick("pull SP++")
	temp := Read(uint16(StackPointer) + 0x100)
	//MasterClockTick("pull")
	return temp
}

func PushFlags() {
	temp := byte(0)
	temp += byte(common.Ternary(flag_Carry, 0x01, 0x00))
	temp += byte(common.Ternary(flag_Zero, 0x02, 0x00))
	temp += byte(common.Ternary(flag_InterruptDisable, 0x04, 0x00))
	temp += byte(common.Ternary(flag_Decimal, 0x08, 0x00))
	temp += byte(common.Ternary(flag_B, 0x10, 0x00)) //B Flag
	temp += 0x20
	temp += byte(common.Ternary(flag_Overflow, 0x40, 0x00))
	temp += byte(common.Ternary(flag_Negative, 0x80, 0x00))
	Push(temp)

}

func PullFlags() {
	temp := Pull()
	flag_Carry = (temp & 0x01) != 0
	flag_Zero = (temp & 0x02) != 0
	flag_InterruptDisable = (temp & 0x04) != 0
	flag_Decimal = (temp & 0x08) != 0
	flag_B = (temp & 0x10) != 0
	flag_Overflow = (temp & 0x40) != 0
	flag_Negative = (temp & 0x80) != 0
	//MasterClockTick("pull flags")
}

// Performs Arithmetic Shift Left onto value at Address
func Op_ASL(Address uint16) {
	Value := Read(Address)
	flag_Carry = (Value >= 0x80)
	Value <<= 1
	Write(Address, Value)
	SetZNFlags(Value)

}

// Performs Arithmetic Shift Right onto value at Address
func Op_LSR(Address uint16) {
	Value := Read(Address)
	flag_Carry = (Value & 1) != 0
	Value >>= 1
	Write(Address, Value)
	SetZNFlags(Value)
}

// Perform Rotate Left onto value at Address
func Op_ROL(Address uint16) {
	Value := Read(Address)
	futureCarry := (Value >= 0x80)
	Value <<= 1
	if flag_Carry {
		Value |= 1
	}
	//MasterClockTick("rol")
	Write(Address, Value)
	flag_Carry = futureCarry
	SetZNFlags(Value)
}

// Perform Rotate Right onto value at Address
func Op_ROR(Address uint16) {
	Value := Read(Address)
	futureCarry := (Value & 1) != 0
	Value >>= 1
	if flag_Carry {
		Value |= 0x80
	}
	//MasterClockTick("ror")
	Write(Address, Value)
	flag_Carry = futureCarry
	SetZNFlags(Value)
}

// Increment Value, and save to Address
func Op_INC(Address uint16, Value byte) {
	Value++
	Write(Address, Value)
	SetZNFlags(Value)
}

// Decrement Value, and save to Address
func Op_DEC(Address uint16, Value byte) {
	Value--
	Write(Address, Value)
	SetZNFlags(Value)
}

// Bitwise OR with A
func Op_ORA(Value byte) {
	A |= Value
	SetZNFlags(A)
}

// Bitwise AND with A
func Op_AND(Value byte) {
	A &= Value
	SetZNFlags(A)
}

// Bitwise XOR with A
func Op_EOR(Value byte) {
	A ^= Value
	SetZNFlags(A)
}

// Add Value to A with Carry
func Op_ADC(Value byte) {
	IntSum := int(A) + int(Value) + common.BoolToInt(flag_Carry)
	flag_Overflow = (^int(A^Value) & (int(A) ^ IntSum) & 0x80) != 0
	flag_Carry = IntSum > 0xFF
	A = byte(IntSum)
	SetZNFlags(A)
}

// Subtract Value from A with Carry
func Op_SBC(Value byte) {
	IntSum := int(A) - int(Value) - common.BoolToInt(!flag_Carry)
	flag_Overflow = (int(A^Value) & (int(A) ^ IntSum) & 0x80) != 0
	flag_Carry = IntSum >= 0x00
	A = byte(IntSum)
	SetZNFlags(A)
}

// Compare Value with A
func Op_CMP(Value byte) {
	flag_Carry = Value <= A
	flag_Zero = (Value == A)
	flag_Negative = ((A - Value) >= 0x80)
}

// Compare Value with X
func Op_CPX(Value byte) {
	flag_Carry = Value <= X
	flag_Zero = (Value == X)
	flag_Negative = ((X - Value) >= 0x80)
}

// Compare Value with Y
func Op_CPY(Value byte) {
	flag_Carry = Value <= Y
	flag_Zero = (Value == Y)
	flag_Negative = ((Y - Value) >= 0x80)
}

// Uhh
func Op_BIT(Value byte) {
	//Bit Test
	flag_Zero = ((A & Value) == 0)
	flag_Negative = ((Value & 0x80) != 0)
	flag_Overflow = ((Value & 0x40) != 0)
}

func ResetCPU() {

	StackPointer = 0xFD
	A, X, Y = 0, 0, 0
	opcode = 0
	operands = nil
	CPU_Cycles, CPU_Cycles_New, common.CPU_TotalCycles = 0, 0, 0
	CPU_Halted = false

	flag_Carry = false
	flag_Zero = false
	flag_InterruptDisable = true
	flag_Decimal = false
	flag_Overflow = false
	flag_Negative = false
	flag_B = false
}

func BuildAddress(Value_Low, Value_High byte) uint16 {
	//b := []byte{Value_Low, Value_High}
	//return binary.LittleEndian.Uint16(b[0:])

	AddressBus = (uint16(Value_High)<<8 | uint16(Value_Low))
	return AddressBus
}

func ReadFromPC() byte {
	Value := Read(ProgramCounter)
	//if LoggingCPU {
	//	operands = append(operands, Value)
	//}
	ProgramCounter++
	return Value
}

func Emulate_CPU( /*g *Game*/ ) {
	//Non Maskable Interrupt check
	prevNMILevelDetector := ppu.NMILevelDetector
	ppu.NMILevelDetector = (ppu.PPUCTRL_EnableNMI && ppu.PPUSTATUS_VBlank)
	if !prevNMILevelDetector && ppu.NMILevelDetector {
		ppu.DoNMI = true
	}
	CPU_Cycles_New = 0

	//Get the opcode
	//var opcode byte
	if !ppu.DoNMI {
		//If we're not running an NMI
		opcode = Read(ProgramCounter)
		//if debug.LoggingCPU {
		//	debug.prepTraceLogger()
		//}
		ProgramCounter++
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
		Address := ReadOperands_ZeroPageAddressed()
		Write(Address, A)
		CPU_Cycles = 3
	case 0x95: //STA Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Write(Address, A)
		CPU_Cycles = 4
	case 0x8D: //STA Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Write(Address, A)
		CPU_Cycles = 4
	case 0x9D: //STA Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed(false)
		Write(Address, A)
		//MasterClockTick("STA Absolute X")
		CPU_Cycles = 5
	case 0x99: //STA Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed(false)
		Write(Address, A)
		//MasterClockTick("sta absolute y")
		CPU_Cycles = 5
	case 0x81: //STA Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Write(Address, A)
		CPU_Cycles = 6
	case 0x91: //STA Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed(false)
		Write(Address, A)
		//MasterClockTick("sta ind Y")
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
		Address := ReadOperands_ZeroPageAddressed()
		A = Read(Address)
		SetZNFlags(A)
		CPU_Cycles = 3
	case 0xB5: //LDA Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		A = Read(Address)
		SetZNFlags(A)
		CPU_Cycles = 4
	case 0xAD: //LDA Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		A = Read(Address)
		SetZNFlags(A)
		CPU_Cycles = 4
	case 0xBD: //LDA Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed(true)
		A = Read(Address)
		SetZNFlags(A)
		CPU_Cycles = 4
	case 0xB9: //LDA Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed(true)
		A = Read(Address)
		SetZNFlags(A)
		CPU_Cycles = 4
	case 0xA1: //LDA Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		A = Read(Address)
		SetZNFlags(A)
		CPU_Cycles = 6
	case 0xB1: //LDA Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed(true)
		A = Read(Address)
		SetZNFlags(A)
		CPU_Cycles = 5

	//	STX: Store Index X in Memory
	//	X -> M
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x86: //STX Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Write(Address, X)
		CPU_Cycles = 3
	case 0x96: //STX Zero Page, Y
		Address := ReadOperands_ZeroPageAddressed_YIndexed()
		Write(Address, X)
		CPU_Cycles = 4
	case 0x8E: //STX Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Write(Address, X)
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
		Address := ReadOperands_ZeroPageAddressed()
		X = Read(Address)
		SetZNFlags(X)
		CPU_Cycles = 3
	case 0xAE: //LDX Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		X = Read(Address)
		SetZNFlags(X)
		CPU_Cycles = 4
	case 0xB6: //LDX Zero Page, Y
		Address := ReadOperands_ZeroPageAddressed_YIndexed()
		X = Read(Address)
		SetZNFlags(X)
		CPU_Cycles = 4
	case 0xBE: //LDX Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed(true)
		X = Read(Address)
		SetZNFlags(X)
		CPU_Cycles = 4

	//	STY: Store Index Y in Memory
	//	Y -> M
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x84: //STY Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Write(Address, Y)
		CPU_Cycles = 3
	case 0x94: //STY Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Write(Address, Y)
		CPU_Cycles = 4
	case 0x8C: //STY Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Write(Address, Y)
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
		Address := ReadOperands_ZeroPageAddressed()
		Y = Read(Address)
		SetZNFlags(Y)
		CPU_Cycles = 3
	case 0xAC: //LDY Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Y = Read(Address)
		SetZNFlags(Y)
		CPU_Cycles = 4
	case 0xB4: //LDY Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Y = Read(Address)
		SetZNFlags(Y)
		CPU_Cycles = 4
	case 0xBC: //LDY Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed(true)
		Y = Read(Address)
		SetZNFlags(Y)
		CPU_Cycles = 4

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
		Address := ReadOperands_ZeroPageAddressed()
		Op_ADC(Read(Address))
		CPU_Cycles = 3
	case 0x75: //ADC Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_ADC(Read(Address))
		CPU_Cycles = 4
	case 0x6D: //ADC Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Op_ADC(Read(Address))
		CPU_Cycles = 4
	case 0x7D: //ADC Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed(true)
		Op_ADC(Read(Address))
		CPU_Cycles = 4
	case 0x79: //ADC Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed(true)
		Op_ADC(Read(Address))
		CPU_Cycles = 4
	case 0x61: //ADC Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Op_ADC(Read(Address))
		CPU_Cycles = 6
	case 0x71: //ADC Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed(true)
		Op_ADC(Read(Address))
		CPU_Cycles = 5

	//	SBC: Subtract Memory from Accumulator with Borrow
	//	A - M - C̅ -> A
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	+

	case 0xE9: //SBC Immediate
		Op_SBC(ReadFromPC())
		CPU_Cycles = 2
	case 0xE5: //SBC Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_SBC(Read(Address))
		CPU_Cycles = 3
	case 0xED: //SBC Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Op_SBC(Read(Address))
		CPU_Cycles = 4
	case 0xF5: //SBC Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_SBC(Read(Address))
		CPU_Cycles = 4
	case 0xFD: //SBC Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed(true)
		Op_SBC(Read(Address))
		CPU_Cycles = 4
	case 0xF9: //SBC Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed(true)
		Op_SBC(Read(Address))
		CPU_Cycles = 4
	case 0xE1: //SBC Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Op_SBC(Read(Address))
		CPU_Cycles = 6
	case 0xF1: //SBC Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed(true)
		Op_SBC(Read(Address))
		CPU_Cycles = 5

	//	INC: Increment Memory by One
	//	M + 1 -> M
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xE6: //INC Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_INC(Address, Read(Address))
		CPU_Cycles = 5
	case 0xEE: //INC Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Op_INC(Address, Read(Address))
		CPU_Cycles = 6
	case 0xF6: //INC Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_INC(Address, Read(Address))
		CPU_Cycles = 6
	case 0xFE: //INC Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_INC(Address, Read(Address))
		CPU_Cycles = 7

	//	DEC: Decrement Memory by One
	//	M - 1 -> M
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xC6: //DEC Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_DEC(Address, Read(Address))
		CPU_Cycles = 5
	case 0xCE: //DEC Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Op_DEC(Address, Read(Address))
		CPU_Cycles = 6
	case 0xD6: //DEC Zeropage, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_DEC(Address, Read(Address))
		CPU_Cycles = 6
	case 0xDE: //DEC Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_DEC(Address, Read(Address))
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
		Op_ASL(ReadOperands_ZeroPageAddressed())
		CPU_Cycles = 5
	case 0x0E: //ASL Absolute
		Op_ASL(ReadOperands_AbsoluteAddressed(false))
		CPU_Cycles = 6
	case 0x16: //ASL Zero Page, X
		Op_ASL(ReadOperands_ZeroPageAddressed_XIndexed())
		CPU_Cycles = 6
	case 0x1E: //ASL Absolute, X
		Op_ASL(ReadOperands_AbsoluteAddressed_XIndexed(false))
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
		Op_LSR(ReadOperands_ZeroPageAddressed())
		CPU_Cycles = 5
	case 0x4E: //LSR Absolute
		Op_LSR(ReadOperands_AbsoluteAddressed(false))
		CPU_Cycles = 6
	case 0x56: //LSR Zero Page, X
		Op_LSR(ReadOperands_ZeroPageAddressed_XIndexed())
		CPU_Cycles = 6
	case 0x5E: //LSR Absolute, X
		Op_LSR(ReadOperands_AbsoluteAddressed_XIndexed(false))
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
		Op_ROL(ReadOperands_ZeroPageAddressed())
		CPU_Cycles = 5
	case 0x2E: //ROL Absolute
		Op_ROL(ReadOperands_AbsoluteAddressed(false))
		CPU_Cycles = 6
	case 0x36: //ROL Zero Page, X
		Op_ROL(ReadOperands_ZeroPageAddressed_XIndexed())
		CPU_Cycles = 6
	case 0x3E: //ROL Absolute, X
		Op_ROL(ReadOperands_AbsoluteAddressed_XIndexed(false))
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
		Op_ROR(ReadOperands_ZeroPageAddressed())
		CPU_Cycles = 5
	case 0x6E: //ROR Absolute
		Op_ROR(ReadOperands_AbsoluteAddressed(false))
		CPU_Cycles = 6
	case 0x76: //ROR Zero Page, X
		Op_ROR(ReadOperands_ZeroPageAddressed_XIndexed())
		CPU_Cycles = 6
	case 0x7E: //ROR Absolute, X
		Op_ROR(ReadOperands_AbsoluteAddressed_XIndexed(false))
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
		Address := ReadOperands_ZeroPageAddressed()
		Op_AND(Read(Address))
		CPU_Cycles = 3
	case 0x2D: //AND Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Op_AND(Read(Address))
		CPU_Cycles = 4
	case 0x35: //AND Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_AND(Read(Address))
		CPU_Cycles = 4
	case 0x3D: //AND Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed(true)
		Op_AND(Read(Address))
		CPU_Cycles = 4
	case 0x39: //AND Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed(true)
		Op_AND(Read(Address))
		CPU_Cycles = 4
	case 0x21: //AND Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Op_AND(Read(Address))
		CPU_Cycles = 6
	case 0x31: //AND Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed(true)
		Op_AND(Read(Address))
		CPU_Cycles = 5

	//	ORA: OR Memory with Accumulator
	//	A OR M -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0x09: //ORA Immediate
		Op_ORA(ReadFromPC())
		CPU_Cycles = 2
	case 0x05: //ORA Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_ORA(Read(Address))
		CPU_Cycles = 3
	case 0x0D: //ORA Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Op_ORA(Read(Address))
		CPU_Cycles = 4
	case 0x15: //ORA Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_ORA(Read(Address))
		CPU_Cycles = 4
	case 0x1D: //ORA Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed(true)
		Op_ORA(Read(Address))
		CPU_Cycles = 4
	case 0x19: //ORA Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed(true)
		Op_ORA(Read(Address))
		CPU_Cycles = 4
	case 0x01: //ORA Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Op_ORA(Read(Address))
		CPU_Cycles = 6
	case 0x11: //ORA Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed(true)
		Op_ORA(Read(Address))
		CPU_Cycles = 5

	//	EOR: Exclusive-OR Memory with Accumulator
	//	A EOR M -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0x49: //EOR Immediate
		Op_EOR(ReadFromPC())
		CPU_Cycles = 2
	case 0x45: //EOR Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_EOR(Read(Address))
		CPU_Cycles = 3
	case 0x4D: //EOR Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Op_EOR(Read(Address))
		CPU_Cycles = 4
	case 0x55: //EOR Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_EOR(Read(Address))
		CPU_Cycles = 4
	case 0x5D: //EOR Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed(true)
		Op_EOR(Read(Address))
		CPU_Cycles = 4
	case 0x59: //EOR Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed(true)
		Op_EOR(Read(Address))
		CPU_Cycles = 4
	case 0x41: //EOR Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Op_EOR(Read(Address))
		CPU_Cycles = 6
	case 0x51: //EOR Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed(true)
		Op_EOR(Read(Address))
		CPU_Cycles = 5

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
		Address := ReadOperands_ZeroPageAddressed()
		Op_BIT(Read(Address))
		CPU_Cycles = 3
	case 0x2C: //BIT Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Op_BIT(Read(Address))
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
		Address := ReadOperands_ZeroPageAddressed()
		Op_CMP(Read(Address))
		CPU_Cycles = 3
	case 0xCD: //CMP Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Op_CMP(Read(Address))
		CPU_Cycles = 4
	case 0xD5: //CMP Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_CMP(Read(Address))
		CPU_Cycles = 4
	case 0xDD: //CMP Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed(true)
		Op_CMP(Read(Address))
		CPU_Cycles = 4
	case 0xD9: //CMP Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed(true)
		Op_CMP(Read(Address))
		CPU_Cycles = 4
	case 0xC1: //CMP Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Op_CMP(Read(Address))
		CPU_Cycles = 6
	case 0xD1: //CMP Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed(true)
		Op_CMP(Read(Address))
		CPU_Cycles = 5

	//	CPX: Compare Memory and Index X
	//	X - M
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	case 0xE0: //CPX Immediate
		Op_CPX(ReadFromPC())
		CPU_Cycles = 2
	case 0xE4: //CPX Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_CPX(Read(Address))
		CPU_Cycles = 3
	case 0xEC: //CPX Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Op_CPX(Read(Address))
		CPU_Cycles = 4

	//	CPY: Compare Memory and Index Y
	//	Y - M
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	case 0xC0: //CPY Immediate
		Op_CPY(ReadFromPC())
		CPU_Cycles = 2
	case 0xC4: //CPY Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_CPY(Read(Address))
		CPU_Cycles = 3
	case 0xCC: //CPY Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Op_CPY(Read(Address))
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
		Address := ReadOperands_AbsoluteAddressed(true)
		ProgramCounter = Address
		CPU_Cycles = 3

	case 0x6C: //JMP Indirect
		Address := ReadOperands_IndirectAddressed()
		ProgramCounter = Address
		CPU_Cycles = 5 //TODO: What the fuck

	//	JSR: Jump to New Location Saving Return Address
	//	push (PC+2),
	//	operand 1st byte -> PCL
	//	operand 2nd byte -> PCH
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x20: //JSR
		temp_low := ReadFromPC()
		Push(byte(ProgramCounter / 0x100))
		Push(byte(ProgramCounter))
		temp_high := ReadFromPC()
		ProgramCounter = BuildAddress(temp_low, temp_high)
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
		ProgramCounter = BuildAddress(temp_low, temp_high)
		ProgramCounter++
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
		if !ppu.DoNMI {
			ProgramCounter++
			flag_B = true
		}
		Push(byte(ProgramCounter >> 8))
		Push(byte(ProgramCounter))
		PushFlags()
		//flag_InterruptDisable = true

		PCL := Read(common.Ternary(ppu.DoNMI, 0xFFFA, 0xFFFE))
		PCH := Read(common.Ternary(ppu.DoNMI, 0xFFFB, 0xFFFF))
		ProgramCounter = uint16((uint16(PCH) * 0x100) + uint16(PCL)) //BuildAddress(PCL, PCH)
		ppu.DoNMI = false
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
		ProgramCounter = BuildAddress(temp_low, temp_high)
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
		X = StackPointer
		SetZNFlags(X)
		CPU_Cycles = 2

	//	TXS: Transfer Index X to Stack Register
	//	X -> SP
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x9A: //TXS
		StackPointer = X
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
		ReadOperands_AbsoluteAddressed_XIndexed(true)
		CPU_Cycles = 4
	case 0x3C: //NOP Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(true)
		CPU_Cycles = 4
	case 0x5C: //NOP Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(true)
		CPU_Cycles = 4
	case 0x7C: //NOP Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(true)
		CPU_Cycles = 4
	case 0xDC: //NOP Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(true)
		CPU_Cycles = 4
	case 0xFC: //NOP Absolute, X
		ReadOperands_AbsoluteAddressed_XIndexed(true)
		CPU_Cycles = 4

	//SAX: A AND X -> M

	case 0x87: //SAX Zero Page
		AddressBus = ReadOperands_ZeroPageAddressed()
		Write(AddressBus, (A & X))
		CPU_Cycles = 3
	case 0x97: //SAX Zero Page, Y
		AddressBus = ReadOperands_ZeroPageAddressed_YIndexed()
		Write(AddressBus, (A & X))
		CPU_Cycles = 4
	case 0x8F: //SAX Absolute
		AddressBus = ReadOperands_AbsoluteAddressed(false)
		Write(AddressBus, (A & X))
		CPU_Cycles = 4
	case 0x83: //SAX Indirect, X
		AddressBus = ReadOperands_IndirectAddressed_XIndexed()
		Write(AddressBus, (A & X))
		CPU_Cycles = 6

	//LAX: LDA + LDX

	case 0xA7: //LAX Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		A = Read(Address)
		X = A
		SetZNFlags(X)
		CPU_Cycles = 3
	case 0xB7: //LAX Zero Page, Y
		Address := ReadOperands_ZeroPageAddressed_YIndexed()
		A = Read(Address)
		X = A
		SetZNFlags(X)
		CPU_Cycles = 4
	case 0xAF: //LAX Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		A = Read(Address)
		X = A
		SetZNFlags(X)
		CPU_Cycles = 4
	case 0xBF: //LAX Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed(true)
		A = Read(Address)
		X = A
		SetZNFlags(X)
		CPU_Cycles = 4
	case 0xA3: //LAX Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		A = Read(Address)
		X = A
		SetZNFlags(X)
		CPU_Cycles = 6
	case 0xB3: //LAX Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed(true)
		A = Read(Address)
		X = A
		SetZNFlags(X)
		CPU_Cycles = 5

	//SLO: ASL + ORA

	case 0x07: // SLO Zero Page
		Addr := ReadOperands_ZeroPageAddressed()
		Op_ASL(Addr)
		Op_ORA(Read(Addr))
		CPU_Cycles = 5
	case 0x17: //SLO Zero Page, X
		Addr := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_ASL(Addr)
		Op_ORA(Read(Addr))
		CPU_Cycles = 6
	case 0x0F: //SLO Absolute
		Addr := ReadOperands_AbsoluteAddressed(false)
		Op_ASL(Addr)
		Op_ORA(Read(Addr))
		CPU_Cycles = 6
	case 0x1F: //Slo Absolute, X
		Addr := ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_ASL(Addr)
		Op_ORA(Read(Addr))
		CPU_Cycles = 7
	case 0x1B: //Slo Absolute, Y
		Addr := ReadOperands_AbsoluteAddressed_YIndexed(false)
		Op_ASL(Addr)
		Op_ORA(Read(Addr))
		CPU_Cycles = 7
	case 0x03: //Slo Indirect, X
		Addr := ReadOperands_IndirectAddressed_XIndexed()
		Op_ASL(Addr)
		Op_ORA(Read(Addr))
		CPU_Cycles = 8
	case 0x13: //Slo Indirect, Y
		Addr := ReadOperands_IndirectAddressed_YIndexed(false)
		Op_ASL(Addr)
		Op_ORA(Read(Addr))
		CPU_Cycles = 8

	//DCP: DEC + CMP

	case 0xC7: //DCP Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_DEC(Address, Read(Address))
		Op_CMP(Read(Address))
		CPU_Cycles = 5
	case 0xD7: //DCP Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_DEC(Address, Read(Address))
		Op_CMP(Read(Address))
		CPU_Cycles = 6
	case 0xCF: //DCP Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Op_DEC(Address, Read(Address))
		Op_CMP(Read(Address))
		CPU_Cycles = 6
	case 0xDF: //DCP Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_DEC(Address, Read(Address))
		Op_CMP(Read(Address))
		CPU_Cycles = 7
	case 0xDB: //DCP Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed(false)
		Op_DEC(Address, Read(Address))
		Op_CMP(Read(Address))
		CPU_Cycles = 7
	case 0xC3: //DCP Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Op_DEC(Address, Read(Address))
		Op_CMP(Read(Address))
		CPU_Cycles = 8
	case 0xD3: //DCP Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed(false)
		Op_DEC(Address, Read(Address))
		Op_CMP(Read(Address))
		CPU_Cycles = 8

	//SHA: Stores A AND X AND (high-byte of addr. + 1) at addr.

	case 0x9F: //SHA Absolute, Y
		Addr := ReadOperands_AbsoluteAddressed_YIndexed(false)
		HiByte := byte(Addr >> 8)
		Write(Addr, A&X&(HiByte+1))
		CPU_Cycles = 5
	case 0x93: //SHA Indirect, Y
		Addr := ReadOperands_IndirectAddressed_YIndexed(false)
		Write(Addr, A&X&byte((Addr>>8)+1))
		CPU_Cycles = 6

	//SHX: Stores X AND (high-byte of addr. + 1) at addr.

	case 0x9E: //SHX Absolute, Y
		Addr := ReadOperands_AbsoluteAddressed_YIndexed(false)
		val := (X & byte(((Addr&0xFF00)>>8)+1))
		Write(Addr, val)
		CPU_Cycles = 5

	//SHY: Stores Y AND (high-byte of addr. + 1) at addr.

	case 0x9C: //SHY Absolute, X
		Addr := ReadOperands_AbsoluteAddressed_XIndexed(false)
		val := (Y & byte(((Addr&0xFF00)>>8)+1))
		Write(Addr, val)
		CPU_Cycles = 5

	//RLA: ROL + AND

	case 0x27: //RLA Zero Page
		Addr := ReadOperands_ZeroPageAddressed()
		Op_ROL(Addr)
		Op_AND(Read(Addr))
		CPU_Cycles = 5
	case 0x37: //RLA Zero Page, X
		Addr := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_ROL(Addr)
		Op_AND(Read(Addr))
		CPU_Cycles = 6
	case 0x2F: //RLA Absolute
		Addr := ReadOperands_AbsoluteAddressed(false)
		Op_ROL(Addr)
		Op_AND(Read(Addr))
		CPU_Cycles = 6
	case 0x3F: //RLA Absolute, X
		Addr := ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_ROL(Addr)
		Op_AND(Read(Addr))
		CPU_Cycles = 7
	case 0x3B: //RLA Absolute, Y
		Addr := ReadOperands_AbsoluteAddressed_YIndexed(false)
		Op_ROL(Addr)
		Op_AND(Read(Addr))
		CPU_Cycles = 7
	case 0x23: //RLA Indirect, X
		Addr := ReadOperands_IndirectAddressed_XIndexed()
		Op_ROL(Addr)
		Op_AND(Read(Addr))
		CPU_Cycles = 8
	case 0x33: //RLA Indirect, Y
		Addr := ReadOperands_IndirectAddressed_YIndexed(false)
		Op_ROL(Addr)
		Op_AND(Read(Addr))
		CPU_Cycles = 8

	//SRE: LSR + EOR

	case 0x47: //SRE Zero Page
		AddressBus = ReadOperands_ZeroPageAddressed()
		Op_LSR(AddressBus)
		Op_EOR(Read(AddressBus))
		CPU_Cycles = 5
	case 0x57: //SRE Zero Page, X
		AddressBus = ReadOperands_ZeroPageAddressed_XIndexed()
		Op_LSR(AddressBus)
		Op_EOR(Read(AddressBus))
		CPU_Cycles = 6
	case 0x4F: //SRE Absolute
		AddressBus = ReadOperands_AbsoluteAddressed(false)
		Op_LSR(AddressBus)
		Op_EOR(Read(AddressBus))
		CPU_Cycles = 6
	case 0x5F: //SRE Absolute, X
		AddressBus = ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_LSR(AddressBus)
		Op_EOR(Read(AddressBus))
		CPU_Cycles = 7
	case 0x5B: //SRE Absolute, Y
		AddressBus = ReadOperands_AbsoluteAddressed_YIndexed(false)
		Op_LSR(AddressBus)
		Op_EOR(Read(AddressBus))
		CPU_Cycles = 7
	case 0x43: //SRE Indirect, X
		AddressBus = ReadOperands_IndirectAddressed_XIndexed()
		Op_LSR(AddressBus)
		Op_EOR(Read(AddressBus))
		CPU_Cycles = 8
	case 0x53: //SRE Indirect, Y
		AddressBus = ReadOperands_IndirectAddressed_YIndexed(false)
		Op_LSR(AddressBus)
		Op_EOR(Read(AddressBus))
		CPU_Cycles = 8

	//RRA: ROR + ADC

	case 0x67: //RRA Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_ROR(Address)
		Op_ADC(Read(Address))
		CPU_Cycles = 5
	case 0x77: //RRA Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_ROR(Address)
		Op_ADC(Read(Address))
		CPU_Cycles = 6
	case 0x6F: //RRA Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Op_ROR(Address)
		Op_ADC(Read(Address))
		CPU_Cycles = 6
	case 0x7F: //RRA Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_ROR(Address)
		Op_ADC(Read(Address))
		CPU_Cycles = 7
	case 0x7B: //RRA Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed(false)
		Op_ROR(Address)
		Op_ADC(Read(Address))
		CPU_Cycles = 7
	case 0x63: //RRA Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Op_ROR(Address)
		Op_ADC(Read(Address))
		CPU_Cycles = 8
	case 0x73: //RRA Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed(false)
		Op_ROR(Address)
		Op_ADC(Read(Address))
		CPU_Cycles = 8

	//ISC: INC + SBC

	case 0xE7: //ISC Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_INC(Address, Read(Address))
		Op_SBC(Read(Address))
		CPU_Cycles = 5
	case 0xF7: //ISC Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_INC(Address, Read(Address))
		Op_SBC(Read(Address))
		CPU_Cycles = 6
	case 0xEF: //ISC Absolute
		Address := ReadOperands_AbsoluteAddressed(false)
		Op_INC(Address, Read(Address))
		Op_SBC(Read(Address))
		CPU_Cycles = 6
	case 0xFF: //ISC Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed(false)
		Op_INC(Address, Read(Address))
		Op_SBC(Read(Address))
		CPU_Cycles = 7
	case 0xFB: //ISC Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed(false)
		Op_INC(Address, Read(Address))
		Op_SBC(Read(Address))
		CPU_Cycles = 7
	case 0xE3: //ISC Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Op_INC(Address, Read(Address))
		Op_SBC(Read(Address))
		CPU_Cycles = 8
	case 0xF3: //ISC Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed(false)
		Op_INC(Address, Read(Address))
		Op_SBC(Read(Address))
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
		Address := ReadOperands_AbsoluteAddressed_YIndexed(false)
		Write(Address, (A&X)&(byte(Address&0xFF00)+1))
		CPU_Cycles = 5

	case 0xBB: //LAS (LAE) Absolute, Y
		//LDA/TSX oper
		//M AND SP -> A, X, SP
		Value := (Read(ReadOperands_AbsoluteAddressed_YIndexed(true)) & StackPointer)
		A = Value
		X = Value
		StackPointer = Value
		SetZNFlags(A)
		CPU_Cycles = 4

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
	if (apu.APUDMCInterrupt || apu.APUFrameInterrupt) && !ppu.DoNMI && !flag_InterruptDisable {
		flag_B = false
		Push(byte(ProgramCounter >> 8))
		Push(byte(ProgramCounter))
		PushFlags()

		PCL := Read(0xFFFE)
		PCH := Read(0xFFFF)
		ProgramCounter = BuildAddress(PCL, PCH)

		//Disable interrupts
		apu.APUDMCInterrupt = false
		apu.APUFrameInterrupt = false
		ppu.DoNMI = false
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
