package cpu

import (
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
	H  byte   // High byte of address (Used by some unnoficial ops, and for page crossing checks)
	DL byte   // Data Latch, holds data between instructions
	SB byte   // Special Bus, used in certain instructions (halg-cycle stuff)

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
	MagicConstant    byte //Might be needed for some of the illegal opcodes
	BreakSource      BreakType
	NMILevelDetector bool
	RunningInterrupt bool

	AddressBus  uint16
	TempAddress uint16
	Pointer     uint16
	Target      uint16

	DelayCounter int
}

var CPU_Halted bool

type BreakType int

const (
	Break_None BreakType = iota
	Break_Software
	Break_NMI
	Break_IRQ
	Break_Reset
)

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
	cpu.opcode = 0
	cpu.operands = nil
	cpu.subCycle = 0
	cpu.MagicConstant = 0xFF
	cpu.BreakSource = Break_Reset
	cpu.NMILevelDetector, cpu.RunningInterrupt = false, false

	cpu.flag_Carry = false
	cpu.flag_Zero = false
	cpu.flag_InterruptDisable = true
	cpu.flag_Decimal = false
	cpu.flag_Overflow = false
	cpu.flag_Negative = false
	cpu.flag_B = false

	cpu.DelayCounter = 0

	CPU_Halted = false
}

func (cpu *CPU) CPU_Cycle() {
	if cpu.DelayCounter == 0 {
		if cpu.subCycle == 0 {
			if cpu.BreakSource != Break_None {
				cpu.SetOpcode(0x00)
				if cpu.BreakSource == Break_IRQ {
					apu.APUDMCInterrupt = false
					apu.APUFrameInterrupt = false
					mappers.MMC3_DoIRQ = false
				}
			} else {
				cpu.SetOpcode(cpu.ReadFromPC())
				if cpu.opcode == 0x00 {
					cpu.BreakSource = Break_Software
				}
			}
			cpu.subCycle++
			if debug.LoggingCPU {
				cpu.SendToTracelogger()
			}
		} else {
			cpu.RunInstruction()
			cpu.subCycle++
		}
	} else {
		cpu.DelayCounter--
	}

	common.CPU_TotalCycles++
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
		cpu.X00_BRK()
	case 0x1:
		cpu.X01_ORA_Indirect_X()
	//case 0x2:
	//case 0x3:
	//case 0x4:
	case 0x5:
		cpu.X05_ORA_ZeroPage()
	case 0x6:
		cpu.X06_ASL_ZeroPage()
	//case 0x7:
	case 0x8:
		cpu.X08_PHP()
	case 0x9:
		cpu.X09_ORA_Immediate()
	case 0xA:
		cpu.X0A_ASL()
	//case 0xB:
	//case 0xC:
	case 0xD:
		cpu.X0D_ORA_Absolute()
	case 0xE:
		cpu.X0E_ASL_Absolute()
	//case 0xF:
	default:
		cpu.UnknownOpcode()
	}
}
func (cpu *CPU) Opcode1X() {
	switch cpu.opcode & 0xF {
	case 0x0:
		cpu.X10_BPL()
	case 0x1:
		cpu.X11_ORA_Indirect_Y()
	//case 0x2:
	//case 0x3:
	//case 0x4:
	case 0x5:
		cpu.X15_ORA_ZeroPage_X()
	case 0x6:
		cpu.X16_ASL_ZeroPage_X()
	//case 0x7:
	case 0x8:
		cpu.X18_CLC()
	case 0x9:
		cpu.X19_ORA_Absolute_Y()
	//case 0xA:
	//case 0xB:
	//case 0xC:
	case 0xD:
		cpu.X1D_ORA_Absolute_X()
	case 0xE:
		cpu.X1E_ASL_Absolute_X()
	//case 0xF:
	default:
		cpu.UnknownOpcode()
	}
}
func (cpu *CPU) Opcode2X() {
	switch cpu.opcode & 0xF {
	case 0x0:
		cpu.X20_JSR()
	case 0x1:
		cpu.X21_AND_Indirect_X()
	//case 0x2:
	//case 0x3:
	case 0x4:
		cpu.X24_BIT_ZeroPage()
	case 0x5:
		cpu.X25_AND_ZeroPage()
	case 0x6:
		cpu.X26_ROL_ZeroPage()
	//case 0x7:
	case 0x8:
		cpu.X28_PLP()
	case 0x9:
		cpu.X29_AND_Immediate()
	case 0xA:
		cpu.X2A_ROL()
	//case 0xB:
	case 0xC:
		cpu.X2C_BIT_Absolute()
	case 0xD:
		cpu.X2D_AND_Absolute()
	case 0xE:
		cpu.X2E_ROL_Absolute()
	//case 0xF:
	default:
		cpu.UnknownOpcode()
	}
}
func (cpu *CPU) Opcode3X() {
	switch cpu.opcode & 0xF {
	case 0x0:
		cpu.X30_BMI()
	case 0x1:
		cpu.X31_AND_Indirect_Y()
	//case 0x2:
	//case 0x3:
	//case 0x4:
	case 0x5:
		cpu.X35_AND_ZeroPage_X()
	case 0x6:
		cpu.X36_ROL_ZeroPage_X()
	//case 0x7:
	case 0x8:
		cpu.X38_SEC()
	case 0x9:
		cpu.X39_AND_Absolute_Y()
	//case 0xA:
	//case 0xB:
	//case 0xC:
	case 0xD:
		cpu.X3D_AND_Absolute_X()
	case 0xE:
		cpu.X3E_ROL_Absolute_X()
	//case 0xF:
	default:
		cpu.UnknownOpcode()
	}
}
func (cpu *CPU) Opcode4X() {
	switch cpu.opcode & 0xF {
	case 0x0:
		cpu.X40_RTI()
	case 0x1:
		cpu.X41_EOR_Indirect_X()
	//case 0x2:
	//case 0x3:
	//case 0x4:
	case 0x5:
		cpu.X45_EOR_ZeroPage()
	case 0x6:
		cpu.X46_LSR_ZeroPage()
	//case 0x7:
	case 0x8:
		cpu.X48_PHA()
	case 0x9:
		cpu.X49_EOR_Immediate()
	case 0xA:
		cpu.X4A_LSR()
	//case 0xB:
	case 0xC:
		cpu.X4C_JMP()
	case 0xD:
		cpu.X4D_EOR_Absolute()
	case 0xE:
		cpu.X4E_LSR_Absolute()
	//case 0xF:
	default:
		cpu.UnknownOpcode()
	}
}
func (cpu *CPU) Opcode5X() {
	switch cpu.opcode & 0xF {
	case 0x0:
		cpu.X50_BVC()
	case 0x1:
		cpu.X51_EOR_Indirect_Y()
	//case 0x2:
	//case 0x3:
	//case 0x4:
	case 0x5:
		cpu.X55_EOR_ZeroPage_X()
	case 0x6:
		cpu.X56_LSR_ZeroPage_X()
	//case 0x7:
	case 0x8:
		cpu.X58_CLI()
	case 0x9:
		cpu.X59_EOR_Absolute_Y()
	//case 0xA:
	//case 0xB:
	//case 0xC:
	case 0xD:
		cpu.X5D_EOR_Absolute_X()
	case 0xE:
		cpu.X5E_LSR_Absolute_X()
	//case 0xF:
	default:
		cpu.UnknownOpcode()
	}
}
func (cpu *CPU) Opcode6X() {
	switch cpu.opcode & 0xF {
	case 0x0:
		cpu.X60_RTS()
	case 0x1:
		cpu.X61_ADC_Indirect_X()
	//case 0x2:
	//case 0x3:
	//case 0x4:
	case 0x5:
		cpu.X65_ADC_ZeroPage()
	case 0x6:
		cpu.X66_ROR_ZeroPage()
	//case 0x7:
	case 0x8:
		cpu.X68_PLA()
	case 0x9:
		cpu.X69_ADC_Immediate()
	case 0xA:
		cpu.X6A_ROR()
	//case 0xB:
	case 0xC:
		cpu.X6C_JMP_Indirect()
	case 0xD:
		cpu.X6D_ADC_Absolute()
	case 0xE:
		cpu.X6E_ROR_Absolute()
	//case 0xF:
	default:
		cpu.UnknownOpcode()
	}
}
func (cpu *CPU) Opcode7X() {
	switch cpu.opcode & 0xF {
	case 0x0:
		cpu.X70_BVS()
	case 0x1:
		cpu.X71_ADC_Indirect_Y()
	//case 0x2:
	//case 0x3:
	//case 0x4:
	case 0x5:
		cpu.X75_ADC_ZeroPage_X()
	case 0x6:
		cpu.X76_ROR_ZeroPage_X()
	//case 0x7:
	case 0x8:
		cpu.X78_SEI()
	case 0x9:
		cpu.X79_ADC_Absolute_Y()
	//case 0xA:
	//case 0xB:
	//case 0xC:
	case 0xD:
		cpu.X7D_ADC_Absolute_X()
	case 0xE:
		cpu.X7E_ROR_Absolute_X()
	//case 0xF:
	default:
		cpu.UnknownOpcode()
	}
}
func (cpu *CPU) Opcode8X() {
	switch cpu.opcode & 0xF {
	//case 0x0:
	case 0x1:
		cpu.X81_STA_Indirect_X()
	//case 0x2:
	//case 0x3:
	case 0x4:
		cpu.X84_STY_ZeroPage()
	case 0x5:
		cpu.X85_STA_ZeroPage()
	case 0x6:
		cpu.X86_STX_ZeroPage()
	//case 0x7:
	case 0x8:
		cpu.X88_DEY()
	//case 0x9:
	case 0xA:
		cpu.X8A_TXA()
	//case 0xB:
	case 0xC:
		cpu.X8C_STY_Absolute()
	case 0xD:
		cpu.X8D_STA_Absolute()
	case 0xE:
		cpu.X8E_STX_Absolute()
	//case 0xF:
	default:
		cpu.UnknownOpcode()
	}
}
func (cpu *CPU) Opcode9X() {
	switch cpu.opcode & 0xF {
	case 0x0:
		cpu.X90_BCC()
	case 0x1:
		cpu.X91_STA_Indirect_Y()
	//case 0x2:
	//case 0x3:
	case 0x4:
		cpu.X94_STY_ZeroPage_X()
	case 0x5:
		cpu.X95_STA_ZeroPage_X()
	case 0x6:
		cpu.X96_STX_ZeroPage_Y()
	//case 0x7:
	case 0x8:
		cpu.X98_TYA()
	case 0x9:
		cpu.X99_STA_Absolute_Y()
	case 0xA:
		cpu.X9A_TXS()
	//case 0xB:
	//case 0xC:
	case 0xD:
		cpu.X9D_STA_Absolute_X()
	//case 0xE:
	//case 0xF:
	default:
		cpu.UnknownOpcode()
	}
}
func (cpu *CPU) OpcodeAX() {
	switch cpu.opcode & 0xF {
	case 0x0:
		cpu.XA0_LDY_Immediate()
	case 0x1:
		cpu.XA1_LDA_Indirect_X()
	case 0x2:
		cpu.XA2_LDX_Immediate()
	//case 0x3:
	case 0x4:
		cpu.XA4_LDY_ZeroPage()
	case 0x5:
		cpu.XA5_LDA_ZeroPage()
	case 0x6:
		cpu.XA6_LDX_ZeroPage()
	//case 0x7:
	case 0x8:
		cpu.XA8_TAY()
	case 0x9:
		cpu.XA9_LDA_Immediate()
	case 0xA:
		cpu.XAA_TAX()
	//case 0xB:
	case 0xC:
		cpu.XAC_LDY_Absolute()
	case 0xD:
		cpu.XAD_LDA_Absolute()
	case 0xE:
		cpu.XAE_LDX_Absolute()
	//case 0xF:
	default:
		cpu.UnknownOpcode()
	}
}
func (cpu *CPU) OpcodeBX() {
	switch cpu.opcode & 0xF {
	case 0x0:
		cpu.XB0_BCS()
	case 0x1:
		cpu.XB1_LDA_Indirect_Y()
	//case 0x2:
	//case 0x3:
	case 0x4:
		cpu.XB4_LDY_ZeroPage_X()
	case 0x5:
		cpu.XB5_LDA_ZeroPage_X()
	case 0x6:
		cpu.XB6_LDX_ZeroPage_Y()
	//case 0x7:
	case 0x8:
		cpu.XB8_CLV()
	case 0x9:
		cpu.XB9_LDA_Absolute_Y()
	case 0xA:
		cpu.XBA_TSX()
	//case 0xB:
	case 0xC:
		cpu.XBC_LDY_Absolute_X()
	case 0xD:
		cpu.XBD_LDA_Absolute_X()
	case 0xE:
		cpu.XBE_LDX_Absolute_Y()
	//case 0xF:
	default:
		cpu.UnknownOpcode()
	}
}
func (cpu *CPU) OpcodeCX() {
	switch cpu.opcode & 0xF {
	case 0x0:
		cpu.XC0_CPY_Immediate()
	case 0x1:
		cpu.XC1_CMP_Indirect_X()
	//case 0x2:
	//case 0x3:
	case 0x4:
		cpu.XC4_CPY_ZeroPage()
	case 0x5:
		cpu.XC5_CMP_ZeroPage()
	case 0x6:
		cpu.XC6_DEC_ZeroPage()
	//case 0x7:
	case 0x8:
		cpu.XC8_INY()
	case 0x9:
		cpu.XC9_CMP_Immediate()
	case 0xA:
		cpu.XCA_DEX()
	//case 0xB:
	case 0xC:
		cpu.XCC_CPY_Absolute()
	case 0xD:
		cpu.XCD_CMP_Absolute()
	case 0xE:
		cpu.XCE_DEC_Absolute()
	//case 0xF:
	default:
		cpu.UnknownOpcode()
	}
}
func (cpu *CPU) OpcodeDX() {
	switch cpu.opcode & 0xF {
	case 0x0:
		cpu.XD0_BNE()
	case 0x1:
		cpu.XD1_CMP_Indirect_Y()
	//case 0x2:
	//case 0x3:
	//case 0x4:
	case 0x5:
		cpu.XD5_CMP_ZeroPage_X()
	case 0x6:
		cpu.XD6_DEC_ZeroPage_X()
	//case 0x7:
	case 0x8:
		cpu.XD8_CLD()
	case 0x9:
		cpu.XD9_CMP_Absolute_Y()
	//case 0xA:
	//case 0xB:
	//case 0xC:
	case 0xD:
		cpu.XDD_CMP_Absolute_X()
	case 0xE:
		cpu.XDE_DEC_Absolute_X()
	//case 0xF:
	default:
		cpu.UnknownOpcode()
	}
}
func (cpu *CPU) OpcodeEX() {
	switch cpu.opcode & 0xF {
	case 0x0:
		cpu.XE0_CPX_Immediate()
	case 0x1:
		cpu.XE1_SBC_Indirect_X()
	//case 0x2:
	//case 0x3:
	case 0x4:
		cpu.XE4_CPX_ZeroPage()
	case 0x5:
		cpu.XE5_SBC_ZeroPage()
	case 0x6:
		cpu.XE6_INC_ZeroPage()
	//case 0x7:
	case 0x8:
		cpu.XE8_INX()
	case 0x9:
		cpu.XE9_SBC_Immediate()
	case 0xA:
		cpu.XEA_NOP()
	//case 0xB:
	case 0xC:
		cpu.XEC_CPX_Absolute()
	case 0xD:
		cpu.XED_SBC_Absolute()
	case 0xE:
		cpu.XEE_INC_Absolute()
	//case 0xF:
	default:
		cpu.UnknownOpcode()
	}
}
func (cpu *CPU) OpcodeFX() {
	switch cpu.opcode & 0xF {
	case 0x0:
		cpu.XF0_BEQ()
	case 0x1:
		cpu.XF1_SBC_Indirect_Y()
	//case 0x2:
	//case 0x3:
	//case 0x4:
	case 0x5:
		cpu.XF5_SBC_ZeroPage_X()
	case 0x6:
		cpu.XF6_INC_ZeroPage_X()
	//case 0x7:
	case 0x8:
		cpu.XF8_SED()
	case 0x9:
		cpu.XF9_SBC_Absolute_Y()
	//case 0xA:
	//case 0xB:
	//case 0xC:
	case 0xD:
		cpu.XFD_SBC_Absolute_X()
	case 0xE:
		cpu.XFE_INC_Absolute_X()
	//case 0xF:
	default:
		cpu.UnknownOpcode()
	}
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
	//fmt.Println("Unknown Opcode: " + fmt.Sprintf("%02X", cpu.opcode))
	cpu.CompleteInstruction()
}

func (cpu *CPU) PollNMI() bool {
	prevNMILevelDetector := cpu.NMILevelDetector
	cpu.NMILevelDetector = (ppu.PPUCTRL_EnableNMI && ppu.PPUSTATUS_VBlank)
	return !prevNMILevelDetector && cpu.NMILevelDetector
}

func (cpu *CPU) PollIRQ() bool {
	return (apu.APUDMCInterrupt || apu.APUFrameInterrupt || mappers.MMC3_DoIRQ) && !cpu.flag_InterruptDisable
}

func (cpu *CPU) SetOpcode(code byte) {
	cpu.opcode = code
}

func (cpu *CPU) CompleteInstruction() {
	cpu.PollInterrupts()
	cpu.subCycle = -1
	cpu.AddressBus = cpu.PC
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

	debug.TraceLogger(cpu.opcode, cpu.A, cpu.X, cpu.Y, cpu.SP, status, cpu.PC)
}
