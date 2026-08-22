package cpu2

import (
	"fmt"
	"mtt/timenes/nes/apu"
	c "mtt/timenes/nes/cpu"
	"mtt/timenes/nes/ppu"
)

// CPU Registers
type CPU struct {
	PC uint16 // Program Counter
	SP byte   // Stack Pointer
	A  byte   // Accumulator
	X  byte   // X-Index
	Y  byte   // Y-Index
	H  byte   // High byte of address (Used by some unnoficial ops, and for page crossing checks)
	DL byte   // Data Latch, holds data between instructions

	// Status Register
	flag_Carry            bool // Bit 0: Carry Flag
	flag_Zero             bool // Bit 1: Zero Flag
	flag_InterruptDisable bool // Bit 2: Interrupt Disable Flag
	flag_Decimal          bool // Bit 3: Decimal Flag
	flag_B                bool // Bit 4: B Flag
	flag_Overflow         bool // Bit 6: Overflow Flag
	flag_Negative         bool // Bit 7: Negative Flag

	opcode           byte
	operands         []byte
	subCycle         int
	CPU_Halted       bool
	MagicConstant    byte //Might be needed for some of the illegal opcodes
	BreakSource      BreakType
	NMILevelDetector bool
	RunningInterrupt bool

	AddressBus  uint16
	TempAddress uint16
	Pointer     uint16
	Target      uint16
}

type BreakType int

const (
	Break_None BreakType = iota
	Break_Software
	Break_NMI
	Break_IRQ
	Break_Reset
)

func (cpu *CPU) ResetCPU() {
	cpu.SP = 0xFD
	cpu.A, cpu.X, cpu.Y = 0, 0, 0
	cpu.opcode = 0
	cpu.operands = nil
	cpu.subCycle = 0
	cpu.CPU_Halted = false
	cpu.MagicConstant = 0xFF
	cpu.NMILevelDetector, cpu.RunningInterrupt = false, false

	cpu.flag_Carry = false
	cpu.flag_Zero = false
	cpu.flag_InterruptDisable = true
	cpu.flag_Decimal = false
	cpu.flag_Overflow = false
	cpu.flag_Negative = false
	cpu.flag_B = false
}

func (cpu *CPU) CPU_Cycle() {

	if cpu.subCycle == 0 {
		cpu.SetOpcode(c.ReadFromPC())

		if cpu.BreakSource != 0 {
			cpu.SetOpcode(0x00)
		} else if cpu.opcode == 0x00 {
			cpu.BreakSource = Break_Software
		}

	} else {
		cpu.RunInstruction()
		cpu.subCycle++
	}

}

func (cpu *CPU) RunInstruction() {
	switch cpu.opcode & 0xF0 {
	case 0x00:
		cpu.Opcode0X()
	case 0x10:
		cpu.Opcode1X()
	case 0x20:
		cpu.Opcode2X()
	case 0x30:
		cpu.Opcode3X()
	case 0x40:
		cpu.Opcode4X()
	case 0x50:
		cpu.Opcode5X()
	case 0x60:
		cpu.Opcode6X()
	case 0x70:
		cpu.Opcode7X()
	case 0x80:
		cpu.Opcode8X()
	case 0x90:
		cpu.Opcode9X()
	case 0xA0:
		cpu.OpcodeAX()
	case 0xB0:
		cpu.OpcodeBX()
	case 0xC0:
		cpu.OpcodeCX()
	case 0xD0:
		cpu.OpcodeDX()
	case 0xE0:
		cpu.OpcodeEX()
	case 0xF0:
		cpu.OpcodeFX()
	}
}

func (cpu *CPU) Opcode0X() {
	switch cpu.opcode & 0xF {
	case 0x0:
	case 0x1:
	case 0x2:
	case 0x3:
	case 0x4:
	case 0x5:
	case 0x6:
	case 0x7:
	case 0x8:
	case 0x9:
	case 0xA:
	case 0xB:
	case 0xC:
	case 0xD:
	case 0xE:
	case 0xF:
	default:
	}
}

func (cpu *CPU) Opcode1X() {}
func (cpu *CPU) Opcode2X() {}
func (cpu *CPU) Opcode3X() {}
func (cpu *CPU) Opcode4X() {}
func (cpu *CPU) Opcode5X() {}
func (cpu *CPU) Opcode6X() {}
func (cpu *CPU) Opcode7X() {}
func (cpu *CPU) Opcode8X() {}
func (cpu *CPU) Opcode9X() {}
func (cpu *CPU) OpcodeAX() {}
func (cpu *CPU) OpcodeBX() {}
func (cpu *CPU) OpcodeCX() {}
func (cpu *CPU) OpcodeDX() {}
func (cpu *CPU) OpcodeEX() {}
func (cpu *CPU) OpcodeFX() {}

func Emulate_CPU() {
	/*
		//Check for, and perform Interrupt Request (IRQ)
		if (apu.APUDMCInterrupt || apu.APUFrameInterrupt) && !ppu.DoNMI && !flag_InterruptDisable {
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
			ppu.DoNMI = false
		}
	*/
}

func (cpu *CPU) PollInterrupts() {
	if cpu.RunningInterrupt {
		return
	}

	if cpu.PollNMI() {
		cpu.BreakSource = Break_NMI
	} else if cpu.PollIRQ() {
		cpu.BreakSource = Break_IRQ
	}
}

func (cpu *CPU) PollInterrupts_CantDisableIRQ() {
	if cpu.RunningInterrupt {
		return
	}

	if cpu.PollNMI() {
		cpu.BreakSource = Break_NMI
	} else if cpu.BreakSource != Break_IRQ && cpu.PollIRQ() {
		cpu.BreakSource = Break_IRQ
	}
}

func (cpu *CPU) UnknownOpcode() {
	fmt.Println("Unknown Opcode: " + fmt.Sprintf("%02X", cpu.opcode))
}

func (cpu *CPU) PollNMI() bool {
	prevNMILevelDetector := cpu.NMILevelDetector
	cpu.NMILevelDetector = (ppu.PPUCTRL_EnableNMI && ppu.PPUSTATUS_VBlank)
	return !prevNMILevelDetector && cpu.NMILevelDetector
}

func (cpu *CPU) PollIRQ() bool {
	return (apu.APUDMCInterrupt || apu.APUFrameInterrupt) && !cpu.flag_InterruptDisable
}

func (cpu *CPU) SetOpcode(code byte) {
	cpu.opcode = code
}

func (cpu *CPU) CompleteInstruction() {
	cpu.PollInterrupts()
	cpu.subCycle = -1
	cpu.AddressBus = cpu.PC
}
