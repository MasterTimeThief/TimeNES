package cpu

import (
	"fmt"
	"mtt/timenes/common"
	"mtt/timenes/debug"
	"mtt/timenes/nes/apu"
	"mtt/timenes/nes/cartridge/mappers"
	"mtt/timenes/nes/ppu"
)

type BUS interface {
	Read(uint16) byte
	Write(uint16, byte)
}

type CPU struct {
	bus BUS
	// CPU Registers
	PC uint16 // Program Counter
	SP byte   // Stack Pointer
	A  byte   // Accumulator
	X  byte   // X-Index
	Y  byte   // Y-Index

	// Status Register
	flag_Carry            bool // Bit 0: Carry Flag
	flag_Zero             bool // Bit 1: Zero Flag
	flag_InterruptDisable bool // Bit 2: Interrupt Disable Flag
	flag_Decimal          bool // Bit 3: Decimal Flag
	flag_B                bool // Bit 4: B Flag
	flag_Overflow         bool // Bit 6: Overflow Flag
	flag_Negative         bool // Bit 7: Negative Flag

}

var opcode byte

// var operands []byte
var CPU_Cycles, CPU_Cycles_New int
var CPU_Halted = false
var MagicConstant byte = 0xFF //Might be needed for some of the illegal opcodes
var NMILevelDetector, DoNMI bool

var AddressBus uint16
var Pointer uint16
var Target uint16

func NewCPU() *CPU {
	cpu := CPU{}
	return &cpu
}

func (cpu *CPU) SetBUS(b BUS) {
	cpu.bus = b
}

func (cpu *CPU) ResetCPU() {
	cpu.SP = 0xFD
	cpu.A, cpu.X, cpu.Y = 0, 0, 0
	opcode = 0
	//operands = nil
	CPU_Cycles, CPU_Cycles_New, common.CPU_TotalCycles = 0, 0, 0
	CPU_Halted = false
	NMILevelDetector, DoNMI = false, false

	cpu.flag_Carry = false
	cpu.flag_Zero = false
	cpu.flag_InterruptDisable = true
	cpu.flag_Decimal = false
	cpu.flag_Overflow = false
	cpu.flag_Negative = false
	cpu.flag_B = false
}

func (cpu *CPU) CPU_Cycle() {
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
		opcode = cpu.bus.Read(cpu.PC)
		//if debug.LoggingCPU {
		//	debug.prepTraceLogger()
		//}
		cpu.PC++
		//MasterClockTick("OPCODE")
	} else {
		//If we're running an NMI, force opcode $00
		opcode = 0x00
	}

	if debug.LoggingCPU {
		cpu.SendToTracelogger()
	}

	CPU_Cycles = 0
	switch opcode {

	//Access

	//	STA: Store Accumulator in Memory
	//	A -> M
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x85: //STA Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.WriteToAB(cpu.A)
		CPU_Cycles = 3
	case 0x95: //STA Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.WriteToAB(cpu.A)
		CPU_Cycles = 4
	case 0x8D: //STA Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.WriteToAB(cpu.A)
		CPU_Cycles = 4
	case 0x9D: //STA Absolute, X
		CPU_Cycles = 5
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(false)
		cpu.WriteToAB(cpu.A)
	case 0x99: //STA Absolute, Y
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(false)
		cpu.WriteToAB(cpu.A)
		CPU_Cycles = 5
	case 0x81: //STA Indirect, X
		cpu.ReadOperands_IndirectAddressed_XIndexed()
		cpu.WriteToAB(cpu.A)
		CPU_Cycles = 6
	case 0x91: //STA Indirect, Y
		cpu.ReadOperands_IndirectAddressed_YIndexed(false)
		cpu.WriteToAB(cpu.A)
		CPU_Cycles = 6

	//	LDA: Load Accumulator with Memory
	//	M -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xA9: //LDA Immediate
		cpu.A = cpu.ReadFromPC()
		cpu.SetZNFlags(cpu.A)
		CPU_Cycles = 2
	case 0xA5: //LDA Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.A = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.A)
		CPU_Cycles = 3
	case 0xB5: //LDA Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.A = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.A)
		CPU_Cycles = 4
	case 0xAD: //LDA Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.A = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.A)
		CPU_Cycles = 4
	case 0xBD: //LDA Absolute, X
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(true)
		cpu.A = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.A)
	case 0xB9: //LDA Absolute, Y
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(true)
		cpu.A = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.A)
	case 0xA1: //LDA Indirect, X
		cpu.ReadOperands_IndirectAddressed_XIndexed()
		cpu.A = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.A)
		CPU_Cycles = 6
	case 0xB1: //LDA Indirect, Y
		CPU_Cycles = 5
		cpu.ReadOperands_IndirectAddressed_YIndexed(true)
		cpu.A = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.A)

	//	STX: Store Index X in Memory
	//	X -> M
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x86: //STX Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.WriteToAB(cpu.X)
		CPU_Cycles = 3
	case 0x96: //STX Zero Page, Y
		cpu.ReadOperands_ZeroPageAddressed_YIndexed()
		cpu.WriteToAB(cpu.X)
		CPU_Cycles = 4
	case 0x8E: //STX Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.WriteToAB(cpu.X)
		CPU_Cycles = 4

	//	LDX: Load Index X with Memory
	//	M -> X
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xA2: //LDX Immediate
		cpu.X = cpu.ReadFromPC()
		cpu.SetZNFlags(cpu.X)
		CPU_Cycles = 2
	case 0xA6: //LDX Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.X = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.X)
		CPU_Cycles = 3
	case 0xAE: //LDX Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.X = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.X)
		CPU_Cycles = 4
	case 0xB6: //LDX Zero Page, Y
		cpu.ReadOperands_ZeroPageAddressed_YIndexed()
		cpu.X = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.X)
		CPU_Cycles = 4
	case 0xBE: //LDX Absolute, Y
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(true)
		cpu.X = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.X)

	//	STY: Store Index Y in Memory
	//	Y -> M
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x84: //STY Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.WriteToAB(cpu.Y)
		CPU_Cycles = 3
	case 0x94: //STY Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.WriteToAB(cpu.Y)
		CPU_Cycles = 4
	case 0x8C: //STY Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.WriteToAB(cpu.Y)
		CPU_Cycles = 4

	//	LDY: Load Index Y with Memory
	//	M -> Y
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xA0: //LDY Immediate
		cpu.Y = cpu.ReadFromPC()
		cpu.SetZNFlags(cpu.Y)
		CPU_Cycles = 2
	case 0xA4: //LDY Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Y = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.Y)
		CPU_Cycles = 3
	case 0xAC: //LDY Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Y = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.Y)
		CPU_Cycles = 4
	case 0xB4: //LDY Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Y = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.Y)
		CPU_Cycles = 4
	case 0xBC: //LDY Absolute, X
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(true)
		cpu.Y = cpu.ReadFromAB()
		cpu.SetZNFlags(cpu.Y)

	//Transfer

	//	TAX: Transfer Accumulator to Index X
	//	A -> X
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xAA: //TAX
		cpu.X = cpu.A
		cpu.SetZNFlags(cpu.X)
		CPU_Cycles = 2

	//	TAY: Transfer Accumulator to Index Y
	//	A -> Y
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xA8: //TAY
		cpu.Y = cpu.A
		cpu.SetZNFlags(cpu.Y)
		CPU_Cycles = 2

	//	TXA: Transfer Index X to Accumulator
	//	X -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0x8A: //TXA
		cpu.A = cpu.X
		cpu.SetZNFlags(cpu.A)
		CPU_Cycles = 2

	//	TYA: Transfer Index Y to Accumulator
	//	Y -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0x98: //TYA
		cpu.A = cpu.Y
		cpu.SetZNFlags(cpu.A)
		CPU_Cycles = 2

	//Arithmetic

	//	ADC: Add Memory to Accumulator with Carry
	//	A + M + C -> A, C
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	+

	case 0x69: //ADC Immediate
		cpu.Op_ADC(cpu.ReadFromPC())
		CPU_Cycles = 2
	case 0x65: //ADC Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_ADC(cpu.ReadFromAB())
		CPU_Cycles = 3
	case 0x75: //ADC Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_ADC(cpu.ReadFromAB())
		CPU_Cycles = 4
	case 0x6D: //ADC Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_ADC(cpu.ReadFromAB())
		CPU_Cycles = 4
	case 0x7D: //ADC Absolute, X
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(true)
		cpu.Op_ADC(cpu.ReadFromAB())
	case 0x79: //ADC Absolute, Y
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(true)
		cpu.Op_ADC(cpu.ReadFromAB())
	case 0x61: //ADC Indirect, X
		cpu.ReadOperands_IndirectAddressed_XIndexed()
		cpu.Op_ADC(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0x71: //ADC Indirect, Y
		CPU_Cycles = 5
		cpu.ReadOperands_IndirectAddressed_YIndexed(true)
		cpu.Op_ADC(cpu.ReadFromAB())

	//	SBC: Subtract Memory from Accumulator with Borrow
	//	A - M - C̅ -> A
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	+

	case 0xE9: //SBC Immediate
		cpu.Op_SBC(cpu.ReadFromPC())
		CPU_Cycles = 2
	case 0xE5: //SBC Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_SBC(cpu.ReadFromAB())
		CPU_Cycles = 3
	case 0xED: //SBC Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_SBC(cpu.ReadFromAB())
		CPU_Cycles = 4
	case 0xF5: //SBC Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_SBC(cpu.ReadFromAB())
		CPU_Cycles = 4
	case 0xFD: //SBC Absolute, X
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(true)
		cpu.Op_SBC(cpu.ReadFromAB())
	case 0xF9: //SBC Absolute, Y
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(true)
		cpu.Op_SBC(cpu.ReadFromAB())
	case 0xE1: //SBC Indirect, X
		cpu.ReadOperands_IndirectAddressed_XIndexed()
		cpu.Op_SBC(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0xF1: //SBC Indirect, Y
		CPU_Cycles = 5
		cpu.ReadOperands_IndirectAddressed_YIndexed(true)
		cpu.Op_SBC(cpu.ReadFromAB())

	//	INC: Increment Memory by One
	//	M + 1 -> M
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xE6: //INC Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_INC(AddressBus, cpu.ReadFromAB())
		CPU_Cycles = 5
	case 0xEE: //INC Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_INC(AddressBus, cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0xF6: //INC Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_INC(AddressBus, cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0xFE: //INC Absolute, X
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(false)
		cpu.Op_INC(AddressBus, cpu.ReadFromAB())
		CPU_Cycles = 7

	//	DEC: Decrement Memory by One
	//	M - 1 -> M
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xC6: //DEC Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		CPU_Cycles = 5
	case 0xCE: //DEC Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0xD6: //DEC Zeropage, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0xDE: //DEC Absolute, X
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(false)
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		CPU_Cycles = 7

	//	INX: Increment Index X by One
	//	X + 1 -> X
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xE8: //INX
		cpu.X++
		CPU_Cycles = 2
		//MasterClockTick("inx")
		cpu.SetZNFlags(cpu.X)

	//	DEX: Decrement Index X by One
	//	X - 1 -> X
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xCA: //DEX
		cpu.X--
		CPU_Cycles = 2
		//MasterClockTick("dex")
		cpu.SetZNFlags(cpu.X)

	//	INY: Increment Index Y by One
	//	Y + 1 -> Y
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xC8: //INY
		cpu.Y++
		CPU_Cycles = 2
		//MasterClockTick("iny")
		cpu.SetZNFlags(cpu.Y)

	//	DEY: Decrement Index Y by One
	//	Y - 1 -> Y
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0x88: //DEY
		cpu.Y--
		CPU_Cycles = 2
		//MasterClockTick("dey")
		cpu.SetZNFlags(cpu.Y)

	//Shift

	//	ASL: Shift Left One Bit (Memory or Accumulator)
	//	C <- [76543210] <- 0
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	case 0x0A: //ASL A
		cpu.flag_Carry = cpu.A > 127
		cpu.A <<= 1
		cpu.SetZNFlags(cpu.A)
		CPU_Cycles = 2
	case 0x06: //ASL Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_ASL(AddressBus)
		CPU_Cycles = 5
	case 0x0E: //ASL Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_ASL(AddressBus)
		CPU_Cycles = 6
	case 0x16: //ASL Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_ASL(AddressBus)
		CPU_Cycles = 6
	case 0x1E: //ASL Absolute, X
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(false)
		cpu.Op_ASL(AddressBus)
		CPU_Cycles = 7

	//	LSR: Shift One Bit Right (Memory or Accumulator)
	//	0 -> [76543210] -> C
	//	N	Z	C	I	D	V
	//	0	+	+	-	-	-

	case 0x4A: //LSR A
		cpu.flag_Carry = (cpu.A & 1) != 0
		cpu.A >>= 1
		cpu.SetZNFlags(cpu.A)
		CPU_Cycles = 2
	case 0x46: //LSR Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_LSR(AddressBus)
		CPU_Cycles = 5
	case 0x4E: //LSR Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_LSR(AddressBus)
		CPU_Cycles = 6
	case 0x56: //LSR Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_LSR(AddressBus)
		CPU_Cycles = 6
	case 0x5E: //LSR Absolute, X
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(false)
		cpu.Op_LSR(AddressBus)
		CPU_Cycles = 7

	//	ROL: Rotate One Bit Left (Memory or Accumulator)
	//	C <- [76543210] <- C
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	case 0x2A: //ROL A
		futureCarry := (cpu.A >= 0x80)
		cpu.A <<= 1
		if cpu.flag_Carry {
			cpu.A |= 1
		}
		cpu.flag_Carry = futureCarry
		cpu.SetZNFlags(cpu.A)
		CPU_Cycles = 2
	case 0x26: //ROL Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_ROL(AddressBus)
		CPU_Cycles = 5
	case 0x2E: //ROL Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_ROL(AddressBus)
		CPU_Cycles = 6
	case 0x36: //ROL Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_ROL(AddressBus)
		CPU_Cycles = 6
	case 0x3E: //ROL Absolute, X
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(false)
		cpu.Op_ROL(AddressBus)
		CPU_Cycles = 7

	//	ROR: Rotate One Bit Right (Memory or Accumulator)
	//	C -> [76543210] -> C
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	case 0x6A: //ROR A
		futureCarry := (cpu.A & 1) != 0
		cpu.A >>= 1
		if cpu.flag_Carry {
			cpu.A |= 0x80
		}
		cpu.flag_Carry = futureCarry
		cpu.SetZNFlags(cpu.A)
		CPU_Cycles = 2
	case 0x66: //ROR Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_ROR(AddressBus)
		CPU_Cycles = 5
	case 0x6E: //ROR Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_ROR(AddressBus)
		CPU_Cycles = 6
	case 0x76: //ROR Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_ROR(AddressBus)
		CPU_Cycles = 6
	case 0x7E: //ROR Absolute, X
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(false)
		cpu.Op_ROR(AddressBus)
		CPU_Cycles = 7

	//Bitwise

	//	AND: AND Memory with Accumulator
	//	A AND M -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0x29: //AND Immediate
		cpu.Op_AND(cpu.ReadFromPC())
		CPU_Cycles = 2
	case 0x25: //AND Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_AND(cpu.ReadFromAB())
		CPU_Cycles = 3
	case 0x2D: //AND Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_AND(cpu.ReadFromAB())
		CPU_Cycles = 4
	case 0x35: //AND Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_AND(cpu.ReadFromAB())
		CPU_Cycles = 4
	case 0x3D: //AND Absolute, X
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(true)
		cpu.Op_AND(cpu.ReadFromAB())
	case 0x39: //AND Absolute, Y
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(true)
		cpu.Op_AND(cpu.ReadFromAB())
	case 0x21: //AND Indirect, X
		cpu.ReadOperands_IndirectAddressed_XIndexed()
		cpu.Op_AND(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0x31: //AND Indirect, Y
		CPU_Cycles = 5
		cpu.ReadOperands_IndirectAddressed_YIndexed(true)
		cpu.Op_AND(cpu.ReadFromAB())

	//	ORA: OR Memory with Accumulator
	//	A OR M -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0x09: //ORA Immediate
		cpu.Op_ORA(cpu.ReadFromPC())
		CPU_Cycles = 2
	case 0x05: //ORA Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_ORA(cpu.ReadFromAB())
		CPU_Cycles = 3
	case 0x0D: //ORA Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_ORA(cpu.ReadFromAB())
		CPU_Cycles = 4
	case 0x15: //ORA Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_ORA(cpu.ReadFromAB())
		CPU_Cycles = 4
	case 0x1D: //ORA Absolute, X
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(true)
		cpu.Op_ORA(cpu.ReadFromAB())
	case 0x19: //ORA Absolute, Y
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(true)
		cpu.Op_ORA(cpu.ReadFromAB())
	case 0x01: //ORA Indirect, X
		cpu.ReadOperands_IndirectAddressed_XIndexed()
		cpu.Op_ORA(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0x11: //ORA Indirect, Y
		CPU_Cycles = 5
		cpu.ReadOperands_IndirectAddressed_YIndexed(true)
		cpu.Op_ORA(cpu.ReadFromAB())

	//	EOR: Exclusive-OR Memory with Accumulator
	//	A EOR M -> A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0x49: //EOR Immediate
		cpu.Op_EOR(cpu.ReadFromPC())
		CPU_Cycles = 2
	case 0x45: //EOR Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_EOR(cpu.ReadFromAB())
		CPU_Cycles = 3
	case 0x4D: //EOR Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_EOR(cpu.ReadFromAB())
		CPU_Cycles = 4
	case 0x55: //EOR Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_EOR(cpu.ReadFromAB())
		CPU_Cycles = 4
	case 0x5D: //EOR Absolute, X
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(true)
		cpu.Op_EOR(cpu.ReadFromAB())
	case 0x59: //EOR Absolute, Y
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(true)
		cpu.Op_EOR(cpu.ReadFromAB())
	case 0x41: //EOR Indirect, X
		cpu.ReadOperands_IndirectAddressed_XIndexed()
		cpu.Op_EOR(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0x51: //EOR Indirect, Y
		CPU_Cycles = 5
		cpu.ReadOperands_IndirectAddressed_YIndexed(true)
		cpu.Op_EOR(cpu.ReadFromAB())

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

	case 0x24: //BIT Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_BIT(cpu.ReadFromAB())
		CPU_Cycles = 3
	case 0x2C: //BIT Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_BIT(cpu.ReadFromAB())
		CPU_Cycles = 4

	//Compare

	//	CMP: Compare Memory with Accumulator
	//	A - M
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	case 0xC9: //CMP Immediate
		cpu.Op_CMP(cpu.ReadFromPC())
		CPU_Cycles = 2
	case 0xC5: //CMP Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_CMP(cpu.ReadFromAB())
		CPU_Cycles = 3
	case 0xCD: //CMP Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_CMP(cpu.ReadFromAB())
		CPU_Cycles = 4
	case 0xD5: //CMP Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_CMP(cpu.ReadFromAB())
		CPU_Cycles = 4
	case 0xDD: //CMP Absolute, X
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(true)
		cpu.Op_CMP(cpu.ReadFromAB())
	case 0xD9: //CMP Absolute, Y
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(true)
		cpu.Op_CMP(cpu.ReadFromAB())
	case 0xC1: //CMP Indirect, X
		cpu.ReadOperands_IndirectAddressed_XIndexed()
		cpu.Op_CMP(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0xD1: //CMP Indirect, Y
		CPU_Cycles = 5
		cpu.ReadOperands_IndirectAddressed_YIndexed(true)
		cpu.Op_CMP(cpu.ReadFromAB())

	//	CPX: Compare Memory and Index X
	//	X - M
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	case 0xE0: //CPX Immediate
		cpu.Op_CPX(cpu.ReadFromPC())
		CPU_Cycles = 2
	case 0xE4: //CPX Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_CPX(cpu.ReadFromAB())
		CPU_Cycles = 3
	case 0xEC: //CPX Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_CPX(cpu.ReadFromAB())
		CPU_Cycles = 4

	//	CPY: Compare Memory and Index Y
	//	Y - M
	//	N	Z	C	I	D	V
	//	+	+	+	-	-	-

	case 0xC0: //CPY Immediate
		cpu.Op_CPY(cpu.ReadFromPC())
		CPU_Cycles = 2
	case 0xC4: //CPY Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_CPY(cpu.ReadFromAB())
		CPU_Cycles = 3
	case 0xCC: //CPY Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_CPY(cpu.ReadFromAB())
		CPU_Cycles = 4

	//Branch

	//	BCC: Branch on Carry Clear
	//	branch on C = 0
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x90: //BCC (Branch on Carry Clear)
		temp := cpu.ReadFromPC()
		cpu.Branch(!cpu.flag_Carry, temp)

	//	BCS: Branch on Carry Set
	//	branch on C = 1
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0xB0: //BCS (Branch on Carry Set)
		temp := cpu.ReadFromPC()
		cpu.Branch(cpu.flag_Carry, temp)

	//	BEQ: Branch on Result Zero
	//	branch on Z = 1
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0xF0: //BEQ (Branch on Equal)
		temp := cpu.ReadFromPC()
		cpu.Branch(cpu.flag_Zero, temp)

	//	BNE: Branch on Result not Zero
	//	branch on Z = 0
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0xD0: //BNE (Branch on Not Equal)
		temp := cpu.ReadFromPC()
		cpu.Branch(!cpu.flag_Zero, temp)

	//	BPL: Branch on Result Plus
	//	branch on N = 0
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x10: //BPL (Branch on Plus)
		temp := cpu.ReadFromPC()
		cpu.Branch(!cpu.flag_Negative, temp)

	//	BMI: Branch on Result Minus
	//	branch on N = 1
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x30: //BMI (Branch on Minus)
		temp := cpu.ReadFromPC()
		cpu.Branch(cpu.flag_Negative, temp)

	//	BVC: Branch on Overflow Clear
	//	branch on V = 0
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x50: //BVC (Branch on Overflow Clear)
		temp := cpu.ReadFromPC()
		cpu.Branch(!cpu.flag_Overflow, temp)

	//	BVS: Branch on Overflow Set
	//	branch on V = 1
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x70: //BVS (Branch on Overflow Set)
		temp := cpu.ReadFromPC()
		cpu.Branch(cpu.flag_Overflow, temp)

	//Jump

	//	JMP: Jump to New Location
	//	operand 1st byte -> PCL
	//	operand 2nd byte -> PCH
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x4C: //JMP
		cpu.ReadOperands_AbsoluteAddressed(true)
		cpu.PC = AddressBus
		CPU_Cycles = 3

	case 0x6C: //JMP Indirect
		cpu.ReadOperands_IndirectAddressed()
		cpu.PC = AddressBus
		CPU_Cycles = 5

	//	JSR: Jump to New Location Saving Return Address
	//	push (PC+2),
	//	operand 1st byte -> PCL
	//	operand 2nd byte -> PCH
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x20: //JSR
		temp_low := cpu.ReadFromPC()
		cpu.Push(byte(cpu.PC / 0x100))
		cpu.Push(byte(cpu.PC))
		temp_high := cpu.ReadFromPC()
		cpu.PC = cpu.BuildAddress(temp_low, temp_high)
		//MasterClockTick("jsr")
		CPU_Cycles = 6

	//	RTS: Return from Subroutine
	//	pull PC, PC+1 -> PC
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x60: //RTS
		temp_low := cpu.Pull()
		temp_high := cpu.Pull()
		//MasterClockTick("rts Pull1")
		//MasterClockTick("rts Pull2")
		cpu.PC = cpu.BuildAddress(temp_low, temp_high)
		cpu.PC++
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
		cpu.flag_B = false
		if !DoNMI {
			cpu.PC++
			cpu.flag_B = true
		}
		cpu.Push(byte(cpu.PC >> 8))
		cpu.Push(byte(cpu.PC))
		cpu.PushFlags()
		//flag_InterruptDisable = true

		PCL := cpu.bus.Read(common.Ternary(DoNMI, 0xFFFA, 0xFFFE))
		PCH := cpu.bus.Read(common.Ternary(DoNMI, 0xFFFB, 0xFFFF))
		cpu.PC = uint16((uint16(PCH) * 0x100) + uint16(PCL)) //cpu.BuildAddress(PCL, PCH)
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
		temp := cpu.Pull()
		cpu.flag_Carry = (temp & 1) != 0
		cpu.flag_Zero = (temp & 2) != 0
		cpu.flag_InterruptDisable = (temp & 4) != 0
		cpu.flag_Decimal = (temp & 8) != 0
		cpu.flag_Overflow = (temp & 64) != 0
		cpu.flag_Negative = (temp & 128) != 0
		//MasterClockTick("rti Pull1")
		//MasterClockTick("rti Pull2")
		temp_low := cpu.Pull()
		temp_high := cpu.Pull()
		cpu.PC = cpu.BuildAddress(temp_low, temp_high)
		CPU_Cycles = 6

	//Stack

	//	PHA: Push Accumulator on Stack
	//	push A
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x48: //PHA
		cpu.Push(cpu.A)
		CPU_Cycles = 3

	//	PLA: Pull Accumulator from Stack
	//	pull A
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0x68: //PLA
		cpu.A = cpu.Pull()
		//MasterClockTick("pla")
		cpu.SetZNFlags(cpu.A)
		CPU_Cycles = 4

	//	PHP: Push Processor Status on Stack
	//	The status register will be pushed with the break
	//	flag and bit 5 set to 1.
	//	push SR
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-
	case 0x08: //PHP
		cpu.flag_B = true
		cpu.PushFlags()
		CPU_Cycles = 3

	//	PLP: Pull Processor Status from Stack
	//	The status register will be pulled with the break
	//	flag and bit 5 ignored.
	//	pull SR
	//	N	Z	C	I	D	V
	//	from stack

	case 0x28: //PLP
		cpu.PullFlags()
		CPU_Cycles = 4

	//	TSX: Transfer Stack Pointer to Index X
	//	SP -> X
	//	N	Z	C	I	D	V
	//	+	+	-	-	-	-

	case 0xBA: //TSX
		cpu.X = cpu.SP
		cpu.SetZNFlags(cpu.X)
		CPU_Cycles = 2

	//	TXS: Transfer Index X to Stack Register
	//	X -> SP
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	-

	case 0x9A: //TXS
		cpu.SP = cpu.X
		CPU_Cycles = 2

	//Flags

	//	CLC: Clear Carry Flag
	//	0 -> C
	//	N	Z	C	I	D	V
	//	-	-	0	-	-	-

	case 0x18: //CLC
		cpu.flag_Carry = false
		//MasterClockTick("clc")
		CPU_Cycles = 2

	//	SEC: Set Carry Flag
	//	1 -> C
	//	N	Z	C	I	D	V
	//	-	-	1	-	-	-

	case 0x38: //SEC
		cpu.flag_Carry = true
		//MasterClockTick("sec")
		CPU_Cycles = 2

	//	CLI: Clear Interrupt Disable Bit
	//	0 -> I
	//	N	Z	C	I	D	V
	//	-	-	-	0	-	-

	case 0x58: //CLI
		cpu.flag_InterruptDisable = false
		//MasterClockTick("cli")
		CPU_Cycles = 2

	//	SEI: Set Interrupt Disable Status
	//	1 -> I
	//	N	Z	C	I	D	V
	//	-	-	-	1	-	-

	case 0x78: //SEI
		cpu.flag_InterruptDisable = true
		//MasterClockTick("sei")
		CPU_Cycles = 2

	//	CLD: Clear Decimal Mode
	//	0 -> D
	//	N	Z	C	I	D	V
	//	-	-	-	-	0	-

	case 0xD8: //CLD
		cpu.flag_Decimal = false
		//MasterClockTick("cld")
		CPU_Cycles = 2

	//	SED: Set Decimal Flag
	//	1 -> D
	//	N	Z	C	I	D	V
	//	-	-	-	-	1	-

	case 0xF8: //SED
		cpu.flag_Decimal = true
		//MasterClockTick("sed")
		CPU_Cycles = 2

	//	CLV: Clear Overflow Flag
	//	0 -> V
	//	N	Z	C	I	D	V
	//	-	-	-	-	-	0

	case 0xB8: //CLV
		cpu.flag_Overflow = false
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
		cpu.Kill()
	case 0x12: //HLT
		cpu.Kill()
	case 0x22: //HLT
		cpu.Kill()
	case 0x32: //HLT
		cpu.Kill()
	case 0x42: //HLT
		cpu.Kill()
	case 0x52: //HLT
		cpu.Kill()
	case 0x62: //HLT
		cpu.Kill()
	case 0x72: //HLT
		cpu.Kill()
	case 0x92: //HLT
		cpu.Kill()
	case 0xB2: //HLT
		cpu.Kill()
	case 0xD2: //HLT
		cpu.Kill()
	case 0xF2: //HLT
		cpu.Kill()

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
		cpu.ReadFromPC()
		CPU_Cycles = 2
	case 0x82: //NOP Immediate
		cpu.ReadFromPC()
		CPU_Cycles = 2
	case 0x89: //NOP Immediate
		cpu.ReadFromPC()
		CPU_Cycles = 2
	case 0xC2: //NOP Immediate
		cpu.ReadFromPC()
		CPU_Cycles = 2
	case 0xE2: //NOP Immediate
		cpu.ReadFromPC()
		CPU_Cycles = 2

	case 0x04: //NOP Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		CPU_Cycles = 3
	case 0x44: //NOP Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		CPU_Cycles = 3
	case 0x64: //NOP Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		CPU_Cycles = 3

	case 0x14: //NOP Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		CPU_Cycles = 4
	case 0x34: //NOP Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		CPU_Cycles = 4
	case 0x54: //NOP Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		CPU_Cycles = 4
	case 0x74: //NOP Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		CPU_Cycles = 4
	case 0xD4: //NOP Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		CPU_Cycles = 4
	case 0xF4: //NOP Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		CPU_Cycles = 4

	case 0x0C: //NOP Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		CPU_Cycles = 4

	case 0x1C: //NOP Absolute, X
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(true)
	case 0x3C: //NOP Absolute, X
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(true)
	case 0x5C: //NOP Absolute, X
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(true)
	case 0x7C: //NOP Absolute, X
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(true)
	case 0xDC: //NOP Absolute, X
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(true)
	case 0xFC: //NOP Absolute, X
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(true)

	//SAX: A AND X -> M

	case 0x87: //SAX Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.bus.Write(AddressBus, (cpu.A & cpu.X))
		CPU_Cycles = 3
	case 0x97: //SAX Zero Page, Y
		cpu.ReadOperands_ZeroPageAddressed_YIndexed()
		cpu.bus.Write(AddressBus, (cpu.A & cpu.X))
		CPU_Cycles = 4
	case 0x8F: //SAX Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.bus.Write(AddressBus, (cpu.A & cpu.X))
		CPU_Cycles = 4
	case 0x83: //SAX Indirect, X
		cpu.ReadOperands_IndirectAddressed_XIndexed()
		cpu.bus.Write(AddressBus, (cpu.A & cpu.X))
		CPU_Cycles = 6

	//LAX: LDA + LDX

	case 0xA7: //LAX Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.A = cpu.ReadFromAB()
		cpu.X = cpu.A
		cpu.SetZNFlags(cpu.X)
		CPU_Cycles = 3
	case 0xB7: //LAX Zero Page, Y
		cpu.ReadOperands_ZeroPageAddressed_YIndexed()
		cpu.A = cpu.ReadFromAB()
		cpu.X = cpu.A
		cpu.SetZNFlags(cpu.X)
		CPU_Cycles = 4
	case 0xAF: //LAX Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.A = cpu.ReadFromAB()
		cpu.X = cpu.A
		cpu.SetZNFlags(cpu.X)
		CPU_Cycles = 4
	case 0xBF: //LAX Absolute, Y
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(true)
		cpu.A = cpu.ReadFromAB()
		cpu.X = cpu.A
		cpu.SetZNFlags(cpu.X)
	case 0xA3: //LAX Indirect, X
		cpu.ReadOperands_IndirectAddressed_XIndexed()
		cpu.A = cpu.ReadFromAB()
		cpu.X = cpu.A
		cpu.SetZNFlags(cpu.X)
		CPU_Cycles = 6
	case 0xB3: //LAX Indirect, Y
		CPU_Cycles = 5
		cpu.ReadOperands_IndirectAddressed_YIndexed(true)
		cpu.A = cpu.ReadFromAB()
		cpu.X = cpu.A
		cpu.SetZNFlags(cpu.X)

	//SLO: ASL + ORA

	case 0x07: // SLO Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_ASL(AddressBus)
		cpu.Op_ORA(cpu.ReadFromAB())
		CPU_Cycles = 5
	case 0x17: //SLO Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_ASL(AddressBus)
		cpu.Op_ORA(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0x0F: //SLO Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_ASL(AddressBus)
		cpu.Op_ORA(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0x1F: //Slo Absolute, X
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(false)
		cpu.Op_ASL(AddressBus)
		cpu.Op_ORA(cpu.ReadFromAB())
		CPU_Cycles = 7
	case 0x1B: //Slo Absolute, Y
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(false)
		cpu.Op_ASL(AddressBus)
		cpu.Op_ORA(cpu.ReadFromAB())
		CPU_Cycles = 7
	case 0x03: //Slo Indirect, X
		cpu.ReadOperands_IndirectAddressed_XIndexed()
		cpu.Op_ASL(AddressBus)
		cpu.Op_ORA(cpu.ReadFromAB())
		CPU_Cycles = 8
	case 0x13: //Slo Indirect, Y
		cpu.ReadOperands_IndirectAddressed_YIndexed(false)
		cpu.Op_ASL(AddressBus)
		cpu.Op_ORA(cpu.ReadFromAB())
		CPU_Cycles = 8

	//DCP: DEC + CMP

	case 0xC7: //DCP Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		cpu.Op_CMP(cpu.ReadFromAB())
		CPU_Cycles = 5
	case 0xD7: //DCP Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		cpu.Op_CMP(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0xCF: //DCP Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		cpu.Op_CMP(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0xDF: //DCP Absolute, X
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(false)
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		cpu.Op_CMP(cpu.ReadFromAB())
		CPU_Cycles = 7
	case 0xDB: //DCP Absolute, Y
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(false)
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		cpu.Op_CMP(cpu.ReadFromAB())
		CPU_Cycles = 7
	case 0xC3: //DCP Indirect, X
		cpu.ReadOperands_IndirectAddressed_XIndexed()
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		cpu.Op_CMP(cpu.ReadFromAB())
		CPU_Cycles = 8
	case 0xD3: //DCP Indirect, Y
		cpu.ReadOperands_IndirectAddressed_YIndexed(false)
		cpu.Op_DEC(AddressBus, cpu.ReadFromAB())
		cpu.Op_CMP(cpu.ReadFromAB())
		CPU_Cycles = 8

	//SHA: Stores A AND X AND (high-byte of addr. + 1) at addr.

	case 0x9F: //SHA Absolute, Y
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(false)
		HiByte := byte(AddressBus >> 8)
		cpu.bus.Write(AddressBus, cpu.A&cpu.X&(HiByte+1))
		CPU_Cycles = 5
	case 0x93: //SHA Indirect, Y
		cpu.ReadOperands_IndirectAddressed_YIndexed(false)
		cpu.bus.Write(AddressBus, cpu.A&cpu.X&byte((AddressBus>>8)+1))
		CPU_Cycles = 6

	//SHX: Stores X AND (high-byte of addr. + 1) at addr.

	case 0x9E: //SHX Absolute, Y
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(false)
		val := (cpu.X & byte(((AddressBus&0xFF00)>>8)+1))
		cpu.bus.Write(AddressBus, val)
		CPU_Cycles = 5

	//SHY: Stores Y AND (high-byte of addr. + 1) at addr.

	case 0x9C: //SHY Absolute, X
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(false)
		val := (cpu.Y & byte(((AddressBus&0xFF00)>>8)+1))
		cpu.bus.Write(AddressBus, val)
		CPU_Cycles = 5

	//RLA: ROL + AND

	case 0x27: //RLA Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_ROL(AddressBus)
		cpu.Op_AND(cpu.ReadFromAB())
		CPU_Cycles = 5
	case 0x37: //RLA Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_ROL(AddressBus)
		cpu.Op_AND(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0x2F: //RLA Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_ROL(AddressBus)
		cpu.Op_AND(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0x3F: //RLA Absolute, X
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(false)
		cpu.Op_ROL(AddressBus)
		cpu.Op_AND(cpu.ReadFromAB())
		CPU_Cycles = 7
	case 0x3B: //RLA Absolute, Y
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(false)
		cpu.Op_ROL(AddressBus)
		cpu.Op_AND(cpu.ReadFromAB())
		CPU_Cycles = 7
	case 0x23: //RLA Indirect, X
		cpu.ReadOperands_IndirectAddressed_XIndexed()
		cpu.Op_ROL(AddressBus)
		cpu.Op_AND(cpu.ReadFromAB())
		CPU_Cycles = 8
	case 0x33: //RLA Indirect, Y
		cpu.ReadOperands_IndirectAddressed_YIndexed(false)
		cpu.Op_ROL(AddressBus)
		cpu.Op_AND(cpu.ReadFromAB())
		CPU_Cycles = 8

	//SRE: LSR + EOR

	case 0x47: //SRE Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_LSR(AddressBus)
		cpu.Op_EOR(cpu.ReadFromAB())
		CPU_Cycles = 5
	case 0x57: //SRE Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_LSR(AddressBus)
		cpu.Op_EOR(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0x4F: //SRE Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_LSR(AddressBus)
		cpu.Op_EOR(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0x5F: //SRE Absolute, X
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(false)
		cpu.Op_LSR(AddressBus)
		cpu.Op_EOR(cpu.ReadFromAB())
		CPU_Cycles = 7
	case 0x5B: //SRE Absolute, Y
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(false)
		cpu.Op_LSR(AddressBus)
		cpu.Op_EOR(cpu.ReadFromAB())
		CPU_Cycles = 7
	case 0x43: //SRE Indirect, X
		cpu.ReadOperands_IndirectAddressed_XIndexed()
		cpu.Op_LSR(AddressBus)
		cpu.Op_EOR(cpu.ReadFromAB())
		CPU_Cycles = 8
	case 0x53: //SRE Indirect, Y
		cpu.ReadOperands_IndirectAddressed_YIndexed(false)
		cpu.Op_LSR(AddressBus)
		cpu.Op_EOR(cpu.ReadFromAB())
		CPU_Cycles = 8

	//RRA: ROR + ADC

	case 0x67: //RRA Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_ROR(AddressBus)
		cpu.Op_ADC(cpu.ReadFromAB())
		CPU_Cycles = 5
	case 0x77: //RRA Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_ROR(AddressBus)
		cpu.Op_ADC(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0x6F: //RRA Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_ROR(AddressBus)
		cpu.Op_ADC(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0x7F: //RRA Absolute, X
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(false)
		cpu.Op_ROR(AddressBus)
		cpu.Op_ADC(cpu.ReadFromAB())
		CPU_Cycles = 7
	case 0x7B: //RRA Absolute, Y
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(false)
		cpu.Op_ROR(AddressBus)
		cpu.Op_ADC(cpu.ReadFromAB())
		CPU_Cycles = 7
	case 0x63: //RRA Indirect, X
		cpu.ReadOperands_IndirectAddressed_XIndexed()
		cpu.Op_ROR(AddressBus)
		cpu.Op_ADC(cpu.ReadFromAB())
		CPU_Cycles = 8
	case 0x73: //RRA Indirect, Y
		cpu.ReadOperands_IndirectAddressed_YIndexed(false)
		cpu.Op_ROR(AddressBus)
		cpu.Op_ADC(cpu.ReadFromAB())
		CPU_Cycles = 8

	//ISC: INC + SBC

	case 0xE7: //ISC Zero Page
		cpu.ReadOperands_ZeroPageAddressed()
		cpu.Op_INC(AddressBus, cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		CPU_Cycles = 5
	case 0xF7: //ISC Zero Page, X
		cpu.ReadOperands_ZeroPageAddressed_XIndexed()
		cpu.Op_INC(AddressBus, cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0xEF: //ISC Absolute
		cpu.ReadOperands_AbsoluteAddressed(false)
		cpu.Op_INC(AddressBus, cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		CPU_Cycles = 6
	case 0xFF: //ISC Absolute, X
		cpu.ReadOperands_AbsoluteAddressed_XIndexed(false)
		cpu.Op_INC(AddressBus, cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		CPU_Cycles = 7
	case 0xFB: //ISC Absolute, Y
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(false)
		cpu.Op_INC(AddressBus, cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		CPU_Cycles = 7
	case 0xE3: //ISC Indirect, X
		cpu.ReadOperands_IndirectAddressed_XIndexed()
		cpu.Op_INC(AddressBus, cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		CPU_Cycles = 8
	case 0xF3: //ISC Indirect, Y
		cpu.ReadOperands_IndirectAddressed_YIndexed(false)
		cpu.Op_INC(AddressBus, cpu.ReadFromAB())
		cpu.Op_SBC(cpu.ReadFromAB())
		CPU_Cycles = 8

	//Immediates (unofficial)

	case 0x0B: //ANC Immediate
		//AND + Set Carry as ASL
		cpu.Op_AND(cpu.ReadFromPC())
		cpu.flag_Carry = cpu.flag_Negative
		CPU_Cycles = 2
	case 0x2B: //ANC2 Immediate
		//AND + Set Carry as ROL
		//Same as $0B
		cpu.Op_AND(cpu.ReadFromPC())
		cpu.flag_Carry = cpu.flag_Negative
		CPU_Cycles = 2
	case 0x4B: //ALR Immediate
		//AND + LSR
		cpu.Op_AND(cpu.ReadFromPC())
		cpu.flag_Carry = (cpu.A & 1) != 0
		cpu.A >>= 1
		cpu.SetZNFlags(cpu.A)
		CPU_Cycles = 2
	case 0x6B: //ARR Immediate
		//AND + ROR
		cpu.Op_AND(cpu.ReadFromPC())
		cpu.flag_Overflow = cpu.A == 0

		cpu.A >>= 1
		if cpu.flag_Carry {
			cpu.A |= 0x80
		}
		cpu.flag_Carry = ((cpu.A & 0x40) >> 6) == 1
		cpu.flag_Overflow = (((cpu.A & 0x20) >> 5) ^ ((cpu.A & 0x40) >> 6)) == 1

		cpu.SetZNFlags(cpu.A)
		CPU_Cycles = 2
	case 0x8B: //ANE Immediate
		//Highly unstable
		cpu.A = cpu.X & cpu.ReadFromPC()
		cpu.SetZNFlags(cpu.A)
		CPU_Cycles = 2
	case 0xAB: //LXA Immediate
		//Highly unstable
		cpu.A = cpu.ReadFromPC()
		cpu.X = cpu.A
		cpu.SetZNFlags(cpu.A)
		CPU_Cycles = 2
	case 0xCB: //SBX Immediate
		//(A AND X) - oper -> X
		cpu.X = (cpu.A & cpu.X) - cpu.ReadFromPC()
		cpu.Op_CMP(cpu.X)
		cpu.SetZNFlags(cpu.X)
		CPU_Cycles = 2
	case 0xEB: //SBC Immediate
		cpu.Op_SBC(cpu.ReadFromPC())
		CPU_Cycles = 2

	case 0x9B: //TAS (SHS) Absolute, Y
		//A AND X -> SP, A AND X AND (H+1) -> M
		temp := cpu.A & cpu.X
		cpu.Push(temp)
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(false)
		cpu.bus.Write(AddressBus, (cpu.A&cpu.X)&(byte(AddressBus&0xFF00)+1))
		CPU_Cycles = 5

	case 0xBB: //LAS (LAE) Absolute, Y
		//LDA/TSX oper
		//M AND SP -> A, X, SP
		CPU_Cycles = 4
		cpu.ReadOperands_AbsoluteAddressed_YIndexed(true)
		Value := (cpu.ReadFromAB() & cpu.SP)
		cpu.A = Value
		cpu.X = Value
		cpu.SP = Value
		cpu.SetZNFlags(cpu.A)

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
	if (apu.APUDMCInterrupt || apu.APUFrameInterrupt || mappers.MMC3_DoIRQ) && !DoNMI && !cpu.flag_InterruptDisable {
		cpu.flag_B = false
		cpu.Push(byte(cpu.PC >> 8))
		cpu.Push(byte(cpu.PC))
		cpu.PushFlags()

		PCL := cpu.bus.Read(0xFFFE)
		PCH := cpu.bus.Read(0xFFFF)
		cpu.PC = cpu.BuildAddress(PCL, PCH)

		//Disable interrupts
		apu.APUDMCInterrupt = false
		apu.APUFrameInterrupt = false
		mappers.MMC3_DoIRQ = false
		DoNMI = false
	}

	//CartRAMLogger()
	//operands = nil

	common.CPU_TotalCycles += CPU_Cycles
	for CPU_Cycles > 0 {
		CPU_Cycles--
		ppu.PPU_Cycle()
		ppu.PPU_Cycle()
		ppu.PPU_Cycle()

		apu.APU_Cycle()
	}

	//Force Stop, just in case
	//if TotalCycles > 900000 {

	//if InstructionCount > 50000 {
	//fmt.Println("Too many cycles, end")
	//CPU_Halted = true
	//}

	//InstructionCount++
}

func (cpu *CPU) SendToTracelogger() {
	status := byte(0)
	status += byte(common.Ternary(cpu.flag_Carry, 0x01, 0x00))
	status += byte(common.Ternary(cpu.flag_Zero, 0x02, 0x00))
	status += byte(common.Ternary(cpu.flag_InterruptDisable, 0x04, 0x00))
	status += byte(common.Ternary(cpu.flag_Decimal, 0x08, 0x00))
	status += byte(common.Ternary(cpu.flag_B, 0x10, 0x00)) //B Flag
	status += 0x20
	status += byte(common.Ternary(cpu.flag_Overflow, 0x40, 0x00))
	status += byte(common.Ternary(cpu.flag_Negative, 0x80, 0x00))

	debug.TraceLogger(opcode, cpu.A, cpu.X, cpu.Y, cpu.SP, status, cpu.PC)
}
