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
	subCycle         int
	Magic            byte //Magic constant, for some of the more "unstable" illegal opcodes
	BreakSource      BreakType
	NMILevelDetector bool
	RunningInterrupt bool

	AddressBus  uint16
	TempAddress uint16
	Pointer     uint16
	Target      uint16

	DelayCounter int
	PendingNMI   bool
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
	cpu.subCycle = 0
	cpu.Magic = 0xFD
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
	cpu.PendingNMI = false
}

func (cpu *CPU) CPU_Cycle() {
	if cpu.DelayCounter == 0 {
		if cpu.subCycle == 0 {
			if cpu.BreakSource == Break_NMI {
				print("")
			}
			// Suppress NMI if the read was on the same cycle as VBlank being set
			if ppu.SuppressNMI && cpu.BreakSource == Break_NMI {
				cpu.BreakSource = Break_None
				cpu.PollInterrupts() // Check for an IRQ just in case?
			}

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
			cpu.SendToDebug()

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
	case 0x2:
		cpu.Kill()
	case 0x3:
		cpu.X03_SLO_Indirect_X()
	case 0x4:
		cpu.X04_NOP_ZeroPage()
	case 0x5:
		cpu.X05_ORA_ZeroPage()
	case 0x6:
		cpu.X06_ASL_ZeroPage()
	case 0x7:
		cpu.X07_SLO_ZeroPage()
	case 0x8:
		cpu.X08_PHP()
	case 0x9:
		cpu.X09_ORA_Immediate()
	case 0xA:
		cpu.X0A_ASL()
	case 0xB:
		cpu.X0B_ANC_Immediate()
	case 0xC:
		cpu.X0C_NOP_Absolute()
	case 0xD:
		cpu.X0D_ORA_Absolute()
	case 0xE:
		cpu.X0E_ASL_Absolute()
	case 0xF:
		cpu.X0F_SLO_Absolute()
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
	case 0x2:
		cpu.Kill()
	case 0x3:
		cpu.X13_SLO_Indirect_Y()
	case 0x4:
		cpu.X14_NOP_ZeroPage_X()
	case 0x5:
		cpu.X15_ORA_ZeroPage_X()
	case 0x6:
		cpu.X16_ASL_ZeroPage_X()
	case 0x7:
		cpu.X17_SLO_ZeroPage_X()
	case 0x8:
		cpu.X18_CLC()
	case 0x9:
		cpu.X19_ORA_Absolute_Y()
	case 0xA:
		cpu.X1A_NOP_Implied()
	case 0xB:
		cpu.X1B_SLO_Absolute_Y()
	case 0xC:
		cpu.X1C_NOP_Absolute_X()
	case 0xD:
		cpu.X1D_ORA_Absolute_X()
	case 0xE:
		cpu.X1E_ASL_Absolute_X()
	case 0xF:
		cpu.X1F_SLO_Absolute_X()
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
	case 0x2:
		cpu.Kill()
	case 0x3:
		cpu.X23_RLA_Indirect_X()
	case 0x4:
		cpu.X24_BIT_ZeroPage()
	case 0x5:
		cpu.X25_AND_ZeroPage()
	case 0x6:
		cpu.X26_ROL_ZeroPage()
	case 0x7:
		cpu.X27_RLA_ZeroPage()
	case 0x8:
		cpu.X28_PLP()
	case 0x9:
		cpu.X29_AND_Immediate()
	case 0xA:
		cpu.X2A_ROL()
	case 0xB:
		cpu.X2B_ANC_Immediate()
	case 0xC:
		cpu.X2C_BIT_Absolute()
	case 0xD:
		cpu.X2D_AND_Absolute()
	case 0xE:
		cpu.X2E_ROL_Absolute()
	case 0xF:
		cpu.X2F_RLA_Absolute()
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
	case 0x2:
		cpu.Kill()
	case 0x3:
		cpu.X33_RLA_Indirect_Y()
	case 0x4:
		cpu.X34_NOP_ZeroPage_X()
	case 0x5:
		cpu.X35_AND_ZeroPage_X()
	case 0x6:
		cpu.X36_ROL_ZeroPage_X()
	case 0x7:
		cpu.X37_RLA_ZeroPage_X()
	case 0x8:
		cpu.X38_SEC()
	case 0x9:
		cpu.X39_AND_Absolute_Y()
	case 0xA:
		cpu.X3A_NOP_Implied()
	case 0xB:
		cpu.X3B_RLA_Absolute_Y()
	case 0xC:
		cpu.X3C_NOP_Absolute_X()
	case 0xD:
		cpu.X3D_AND_Absolute_X()
	case 0xE:
		cpu.X3E_ROL_Absolute_X()
	case 0xF:
		cpu.X3F_RLA_Absolute_X()
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
	case 0x2:
		cpu.Kill()
	case 0x3:
		cpu.X43_SRE_Indirect_X()
	case 0x4:
		cpu.X44_NOP_ZeroPage()
	case 0x5:
		cpu.X45_EOR_ZeroPage()
	case 0x6:
		cpu.X46_LSR_ZeroPage()
	case 0x7:
		cpu.X47_SRE_ZeroPage()
	case 0x8:
		cpu.X48_PHA()
	case 0x9:
		cpu.X49_EOR_Immediate()
	case 0xA:
		cpu.X4A_LSR()
	case 0xB:
		cpu.X4B_ALR_Immediate()
	case 0xC:
		cpu.X4C_JMP()
	case 0xD:
		cpu.X4D_EOR_Absolute()
	case 0xE:
		cpu.X4E_LSR_Absolute()
	case 0xF:
		cpu.X4F_SRE_Absolute()
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
	case 0x2:
		cpu.Kill()
	case 0x3:
		cpu.X53_SRE_Indirect_Y()
	case 0x4:
		cpu.X54_NOP_ZeroPage_X()
	case 0x5:
		cpu.X55_EOR_ZeroPage_X()
	case 0x6:
		cpu.X56_LSR_ZeroPage_X()
	case 0x7:
		cpu.X57_SRE_ZeroPage_X()
	case 0x8:
		cpu.X58_CLI()
	case 0x9:
		cpu.X59_EOR_Absolute_Y()
	case 0xA:
		cpu.X5A_NOP_Implied()
	case 0xB:
		cpu.X5B_SRE_Absolute_Y()
	case 0xC:
		cpu.X5C_NOP_Absolute_X()
	case 0xD:
		cpu.X5D_EOR_Absolute_X()
	case 0xE:
		cpu.X5E_LSR_Absolute_X()
	case 0xF:
		cpu.X5F_SRE_Absolute_X()
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
	case 0x2:
		cpu.Kill()
	case 0x3:
		cpu.X63_RRA_Indirect_X()
	case 0x4:
		cpu.X64_NOP_ZeroPage()
	case 0x5:
		cpu.X65_ADC_ZeroPage()
	case 0x6:
		cpu.X66_ROR_ZeroPage()
	case 0x7:
		cpu.X67_RRA_ZeroPage()
	case 0x8:
		cpu.X68_PLA()
	case 0x9:
		cpu.X69_ADC_Immediate()
	case 0xA:
		cpu.X6A_ROR()
	case 0xB:
		cpu.X6B_ARR_Immediate()
	case 0xC:
		cpu.X6C_JMP_Indirect()
	case 0xD:
		cpu.X6D_ADC_Absolute()
	case 0xE:
		cpu.X6E_ROR_Absolute()
	case 0xF:
		cpu.X6F_RRA_Absolute()
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
	case 0x2:
		cpu.Kill()
	case 0x3:
		cpu.X73_RRA_Indirect_Y()
	case 0x4:
		cpu.X74_NOP_ZeroPage_X()
	case 0x5:
		cpu.X75_ADC_ZeroPage_X()
	case 0x6:
		cpu.X76_ROR_ZeroPage_X()
	case 0x7:
		cpu.X77_RRA_ZeroPage_X()
	case 0x8:
		cpu.X78_SEI()
	case 0x9:
		cpu.X79_ADC_Absolute_Y()
	case 0xA:
		cpu.X7A_NOP_Implied()
	case 0xB:
		cpu.X7B_RRA_Absolute_Y()
	case 0xC:
		cpu.X7C_NOP_Absolute_X()
	case 0xD:
		cpu.X7D_ADC_Absolute_X()
	case 0xE:
		cpu.X7E_ROR_Absolute_X()
	case 0xF:
		cpu.X7F_RRA_Absolute_X()
	default:
		cpu.UnknownOpcode()
	}
}
func (cpu *CPU) Opcode8X() {
	switch cpu.opcode & 0xF {
	case 0x0:
		cpu.X80_NOP_Immediate()
	case 0x1:
		cpu.X81_STA_Indirect_X()
	case 0x2:
		cpu.X82_NOP_Immediate()
	case 0x3:
		cpu.X83_SAX_Indirect_X()
	case 0x4:
		cpu.X84_STY_ZeroPage()
	case 0x5:
		cpu.X85_STA_ZeroPage()
	case 0x6:
		cpu.X86_STX_ZeroPage()
	case 0x7:
		cpu.X87_SAX_ZeroPage()
	case 0x8:
		cpu.X88_DEY()
	case 0x9:
		cpu.X89_NOP_Immediate()
	case 0xA:
		cpu.X8A_TXA()
	case 0xB:
		cpu.X8B_ANE_Immediate()
	case 0xC:
		cpu.X8C_STY_Absolute()
	case 0xD:
		cpu.X8D_STA_Absolute()
	case 0xE:
		cpu.X8E_STX_Absolute()
	case 0xF:
		cpu.X8F_SAX_Absolute()
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
	case 0x2:
		cpu.Kill()
	case 0x3:
		cpu.X93_SHA_Indirect_Y()
	case 0x4:
		cpu.X94_STY_ZeroPage_X()
	case 0x5:
		cpu.X95_STA_ZeroPage_X()
	case 0x6:
		cpu.X96_STX_ZeroPage_Y()
	case 0x7:
		cpu.X97_SAX_ZeroPage_Y()
	case 0x8:
		cpu.X98_TYA()
	case 0x9:
		cpu.X99_STA_Absolute_Y()
	case 0xA:
		cpu.X9A_TXS()
	case 0xB:
		cpu.X9B_TAS_Absolute_Y()
	case 0xC:
		cpu.X9C_SHY_Absolute_X()
	case 0xD:
		cpu.X9D_STA_Absolute_X()
	case 0xE:
		cpu.X9E_SHX_Absolute_Y()
	case 0xF:
		cpu.X9F_SHA_Absolute_Y()
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
	case 0x3:
		cpu.XA3_LAX_Indirect_X()
	case 0x4:
		cpu.XA4_LDY_ZeroPage()
	case 0x5:
		cpu.XA5_LDA_ZeroPage()
	case 0x6:
		cpu.XA6_LDX_ZeroPage()
	case 0x7:
		cpu.XA7_LAX_ZeroPage()
	case 0x8:
		cpu.XA8_TAY()
	case 0x9:
		cpu.XA9_LDA_Immediate()
	case 0xA:
		cpu.XAA_TAX()
	case 0xB:
		cpu.XAB_LXA_Immediate()
	case 0xC:
		cpu.XAC_LDY_Absolute()
	case 0xD:
		cpu.XAD_LDA_Absolute()
	case 0xE:
		cpu.XAE_LDX_Absolute()
	case 0xF:
		cpu.XAF_LAX_Absolute()
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
	case 0x2:
		cpu.Kill()
	case 0x3:
		cpu.XB3_LAX_Indirect_Y()
	case 0x4:
		cpu.XB4_LDY_ZeroPage_X()
	case 0x5:
		cpu.XB5_LDA_ZeroPage_X()
	case 0x6:
		cpu.XB6_LDX_ZeroPage_Y()
	case 0x7:
		cpu.XB7_LAX_ZeroPage_Y()
	case 0x8:
		cpu.XB8_CLV()
	case 0x9:
		cpu.XB9_LDA_Absolute_Y()
	case 0xA:
		cpu.XBA_TSX()
	case 0xB:
		cpu.XBB_LAS_Absolute_Y()
	case 0xC:
		cpu.XBC_LDY_Absolute_X()
	case 0xD:
		cpu.XBD_LDA_Absolute_X()
	case 0xE:
		cpu.XBE_LDX_Absolute_Y()
	case 0xF:
		cpu.XBF_LAX_Absolute_Y()
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
	case 0x2:
		cpu.XC2_NOP_Immediate()
	case 0x3:
		cpu.XC3_DCP_Indirect_X()
	case 0x4:
		cpu.XC4_CPY_ZeroPage()
	case 0x5:
		cpu.XC5_CMP_ZeroPage()
	case 0x6:
		cpu.XC6_DEC_ZeroPage()
	case 0x7:
		cpu.XC7_DCP_ZeroPage()
	case 0x8:
		cpu.XC8_INY()
	case 0x9:
		cpu.XC9_CMP_Immediate()
	case 0xA:
		cpu.XCA_DEX()
	case 0xB:
		cpu.XCB_SBX_Immediate()
	case 0xC:
		cpu.XCC_CPY_Absolute()
	case 0xD:
		cpu.XCD_CMP_Absolute()
	case 0xE:
		cpu.XCE_DEC_Absolute()
	case 0xF:
		cpu.XCF_DCP_Absolute()
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
	case 0x2:
		cpu.Kill()
	case 0x3:
		cpu.XD3_DCP_Indirect_Y()
	case 0x4:
		cpu.XD4_NOP_ZeroPage_X()
	case 0x5:
		cpu.XD5_CMP_ZeroPage_X()
	case 0x6:
		cpu.XD6_DEC_ZeroPage_X()
	case 0x7:
		cpu.XD7_DCP_ZeroPage_X()
	case 0x8:
		cpu.XD8_CLD()
	case 0x9:
		cpu.XD9_CMP_Absolute_Y()
	case 0xA:
		cpu.XDA_NOP_Implied()
	case 0xB:
		cpu.XDB_DCP_Absolute_Y()
	case 0xC:
		cpu.XDC_NOP_Absolute_X()
	case 0xD:
		cpu.XDD_CMP_Absolute_X()
	case 0xE:
		cpu.XDE_DEC_Absolute_X()
	case 0xF:
		cpu.XDF_DCP_Absolute_X()
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
	case 0x2:
		cpu.XE2_NOP_Immediate()
	case 0x3:
		cpu.XE3_ISC_Indirect_X()
	case 0x4:
		cpu.XE4_CPX_ZeroPage()
	case 0x5:
		cpu.XE5_SBC_ZeroPage()
	case 0x6:
		cpu.XE6_INC_ZeroPage()
	case 0x7:
		cpu.XE7_ISC_ZeroPage()
	case 0x8:
		cpu.XE8_INX()
	case 0x9:
		cpu.XE9_SBC_Immediate()
	case 0xA:
		cpu.XEA_NOP()
	case 0xB:
		cpu.XEB_SBC_Immediate()
	case 0xC:
		cpu.XEC_CPX_Absolute()
	case 0xD:
		cpu.XED_SBC_Absolute()
	case 0xE:
		cpu.XEE_INC_Absolute()
	case 0xF:
		cpu.XEF_ISC_Absolute()
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
	case 0x2:
		cpu.Kill()
	case 0x3:
		cpu.XF3_ISC_Indirect_Y()
	case 0x4:
		cpu.XF4_NOP_ZeroPage_X()
	case 0x5:
		cpu.XF5_SBC_ZeroPage_X()
	case 0x6:
		cpu.XF6_INC_ZeroPage_X()
	case 0x7:
		cpu.XF7_ISC_ZeroPage_X()
	case 0x8:
		cpu.XF8_SED()
	case 0x9:
		cpu.XF9_SBC_Absolute_Y()
	case 0xA:
		cpu.XFA_NOP_Implied()
	case 0xB:
		cpu.XFB_ISC_Absolute_Y()
	case 0xC:
		cpu.XFC_NOP_Absolute_X()
	case 0xD:
		cpu.XFD_SBC_Absolute_X()
	case 0xE:
		cpu.XFE_INC_Absolute_X()
	case 0xF:
		cpu.XFF_ISC_Absolute_X()
	default:
		cpu.UnknownOpcode()
	}
}

func (cpu *CPU) PollInterrupts() {
	if cpu.RunningInterrupt {
		return
	}

	if cpu.PendingNMI {
		cpu.PendingNMI = false
		cpu.BreakSource = Break_NMI
	} else if cpu.PollIRQ() {
		cpu.BreakSource = Break_IRQ
	}
}

func (cpu *CPU) PollInterrupts_CantDisableIRQ() {
	if cpu.RunningInterrupt {
		return
	}

	if cpu.PendingNMI {
		cpu.PendingNMI = false
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
	return !prevNMILevelDetector && cpu.NMILevelDetector && !ppu.SuppressNMI
}

func (cpu *CPU) PollIRQ() bool {
	return (apu.APUDMCInterrupt || apu.APUFrameInterrupt || mappers.MMC3_DoIRQ) && !cpu.flag_InterruptDisable
}

func (cpu *CPU) DisableNMI() {
	if cpu.BreakSource == Break_NMI {
		cpu.BreakSource = Break_None
	}
}

func (cpu *CPU) SetOpcode(code byte) {
	cpu.opcode = code
}

func (cpu *CPU) CompleteInstruction() {
	//cpu.PollInterrupts()
	cpu.subCycle = -1
	cpu.AddressBus = cpu.PC
}

func (cpu *CPU) SendToDebug() {
	status := byte(0)
	status += byte(common.Ternary(cpu.flag_Carry, 0x01, 0x00))
	status += byte(common.Ternary(cpu.flag_Zero, 0x02, 0x00))
	status += byte(common.Ternary(cpu.flag_InterruptDisable, 0x04, 0x00))
	status += byte(common.Ternary(cpu.flag_Decimal, 0x08, 0x00))
	status += byte(common.Ternary(cpu.flag_B, 0x10, 0x00)) //B Flag
	status += 0x20
	status += byte(common.Ternary(cpu.flag_Overflow, 0x40, 0x00))
	status += byte(common.Ternary(cpu.flag_Negative, 0x80, 0x00))
	debug.SetCPUData(cpu.opcode, cpu.A, cpu.X, cpu.Y, cpu.SP, status, cpu.PC, cpu.AddressBus)
	if debug.LoggingCPU {
		debug.TraceLogger(cpu.opcode, cpu.A, cpu.X, cpu.Y, cpu.SP, status, cpu.PC)
	}
}
