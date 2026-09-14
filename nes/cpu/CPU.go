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

type APU interface {
	RunDMCDMA()
	GetFrameInterrupt() bool
	SetFrameInterrupt(bool)
}

type CPU struct {
	bus BUS
	apu APU

	// CPU Registers
	PC uint16 // Program Counter
	SP byte   // Stack Pointer
	A  byte   // Accumulator
	X  byte   // X-Index
	Y  byte   // Y-Index
	H  byte   // High byte of address (Used by some unnoficial ops, and for page crossing checks)
	DL byte   // Data Latch, holds data between instructions
	SB byte   // Special Bus, used in certain instructions (half-cycle stuff)

	// Status Register
	flag_Carry            bool // Bit 0: Carry Flag
	flag_Zero             bool // Bit 1: Zero Flag
	flag_InterruptDisable bool // Bit 2: Interrupt Disable Flag
	flag_Decimal          bool // Bit 3: Decimal Flag
	flag_B                bool // Bit 4: B Flag
	flag_Overflow         bool // Bit 6: Overflow Flag
	flag_Negative         bool // Bit 7: Negative Flag

	opcode           byte
	InstructionCycle int
	Magic            byte //Magic constant, for some of the more "unstable" illegal opcodes
	BreakSource      BreakType
	NMILevelDetector bool
	RunningInterrupt bool

	AddressBus  uint16
	TempAddress uint16
	Pointer     uint16
	Target      uint16

	DelayCounter int
	NMILine      bool
	IRQLine      bool
	NMIPending   bool
	IRQPending   bool
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

func (c *CPU) SetBUS(b BUS) {
	c.bus = b
}

func (c *CPU) SetAPU(a APU) {
	c.apu = a
}

func (c *CPU) ResetCPU() {
	c.SP = 0xFD
	c.A, c.X, c.Y = 0, 0, 0
	c.opcode = 0
	c.InstructionCycle = 0
	c.Magic = 0xFD
	c.BreakSource = Break_Reset
	c.NMILevelDetector, c.RunningInterrupt = false, false

	c.flag_Carry = false
	c.flag_Zero = false
	c.flag_InterruptDisable = true
	c.flag_Decimal = false
	c.flag_Overflow = false
	c.flag_Negative = false
	c.flag_B = false

	c.DelayCounter = 0

	CPU_Halted = false
	c.NMIPending = false
}

func (c *CPU) CPU_Cycle() {
	if c.DelayCounter == 0 {
		if c.InstructionCycle == 0 {
			if c.NMIPending {
				print("")
			}
			// Suppress NMI if the read was on the same cycle as VBlank being set
			if ppu.SuppressNMI && c.NMIPending {
				c.DisableNMI()
				//c.PollInterrupts() // Check for an IRQ just in case?
			}

			if c.NMIPending {
				c.SetOpcode(0x00)
				c.DisableNMI()
				c.BreakSource = Break_NMI
			} else if c.IRQPending {
				c.SetOpcode(0x00)
				c.BreakSource = Break_IRQ
				c.DisableIRQFlags()
			} else if c.BreakSource == Break_Reset {
				c.SetOpcode(0x00)
			} else {
				c.SetOpcode(c.ReadFromPC())
				if c.opcode == 0x00 {
					c.BreakSource = Break_Software
				}
			}
			c.SendToDebug()
		} else {
			c.RunInstruction()
		}
		c.InstructionCycle++
	} else {
		c.DelayCounter--
	}

	common.CPU_TotalCycles++
}

func (c *CPU) RunInstruction() {
	switch c.opcode & 0xF0 {
	case 0x00:
		c.Opcode0X()
	case 0x10:
		c.Opcode1X()
	case 0x20:
		c.Opcode2X()
	case 0x30:
		c.Opcode3X()
	case 0x40:
		c.Opcode4X()
	case 0x50:
		c.Opcode5X()
	case 0x60:
		c.Opcode6X()
	case 0x70:
		c.Opcode7X()
	case 0x80:
		c.Opcode8X()
	case 0x90:
		c.Opcode9X()
	case 0xA0:
		c.OpcodeAX()
	case 0xB0:
		c.OpcodeBX()
	case 0xC0:
		c.OpcodeCX()
	case 0xD0:
		c.OpcodeDX()
	case 0xE0:
		c.OpcodeEX()
	case 0xF0:
		c.OpcodeFX()
	}
}

func (c *CPU) Opcode0X() {
	switch c.opcode & 0xF {
	case 0x0:
		c.X00_BRK()
	case 0x1:
		c.X01_ORA_Indirect_X()
	case 0x2:
		c.Kill()
	case 0x3:
		c.X03_SLO_Indirect_X()
	case 0x4:
		c.X04_NOP_ZeroPage()
	case 0x5:
		c.X05_ORA_ZeroPage()
	case 0x6:
		c.X06_ASL_ZeroPage()
	case 0x7:
		c.X07_SLO_ZeroPage()
	case 0x8:
		c.X08_PHP()
	case 0x9:
		c.X09_ORA_Immediate()
	case 0xA:
		c.X0A_ASL()
	case 0xB:
		c.X0B_ANC_Immediate()
	case 0xC:
		c.X0C_NOP_Absolute()
	case 0xD:
		c.X0D_ORA_Absolute()
	case 0xE:
		c.X0E_ASL_Absolute()
	case 0xF:
		c.X0F_SLO_Absolute()
	default:
		c.UnknownOpcode()
	}
}
func (c *CPU) Opcode1X() {
	switch c.opcode & 0xF {
	case 0x0:
		c.X10_BPL()
	case 0x1:
		c.X11_ORA_Indirect_Y()
	case 0x2:
		c.Kill()
	case 0x3:
		c.X13_SLO_Indirect_Y()
	case 0x4:
		c.X14_NOP_ZeroPage_X()
	case 0x5:
		c.X15_ORA_ZeroPage_X()
	case 0x6:
		c.X16_ASL_ZeroPage_X()
	case 0x7:
		c.X17_SLO_ZeroPage_X()
	case 0x8:
		c.X18_CLC()
	case 0x9:
		c.X19_ORA_Absolute_Y()
	case 0xA:
		c.X1A_NOP_Implied()
	case 0xB:
		c.X1B_SLO_Absolute_Y()
	case 0xC:
		c.X1C_NOP_Absolute_X()
	case 0xD:
		c.X1D_ORA_Absolute_X()
	case 0xE:
		c.X1E_ASL_Absolute_X()
	case 0xF:
		c.X1F_SLO_Absolute_X()
	default:
		c.UnknownOpcode()
	}
}
func (c *CPU) Opcode2X() {
	switch c.opcode & 0xF {
	case 0x0:
		c.X20_JSR()
	case 0x1:
		c.X21_AND_Indirect_X()
	case 0x2:
		c.Kill()
	case 0x3:
		c.X23_RLA_Indirect_X()
	case 0x4:
		c.X24_BIT_ZeroPage()
	case 0x5:
		c.X25_AND_ZeroPage()
	case 0x6:
		c.X26_ROL_ZeroPage()
	case 0x7:
		c.X27_RLA_ZeroPage()
	case 0x8:
		c.X28_PLP()
	case 0x9:
		c.X29_AND_Immediate()
	case 0xA:
		c.X2A_ROL()
	case 0xB:
		c.X2B_ANC_Immediate()
	case 0xC:
		c.X2C_BIT_Absolute()
	case 0xD:
		c.X2D_AND_Absolute()
	case 0xE:
		c.X2E_ROL_Absolute()
	case 0xF:
		c.X2F_RLA_Absolute()
	default:
		c.UnknownOpcode()
	}
}
func (c *CPU) Opcode3X() {
	switch c.opcode & 0xF {
	case 0x0:
		c.X30_BMI()
	case 0x1:
		c.X31_AND_Indirect_Y()
	case 0x2:
		c.Kill()
	case 0x3:
		c.X33_RLA_Indirect_Y()
	case 0x4:
		c.X34_NOP_ZeroPage_X()
	case 0x5:
		c.X35_AND_ZeroPage_X()
	case 0x6:
		c.X36_ROL_ZeroPage_X()
	case 0x7:
		c.X37_RLA_ZeroPage_X()
	case 0x8:
		c.X38_SEC()
	case 0x9:
		c.X39_AND_Absolute_Y()
	case 0xA:
		c.X3A_NOP_Implied()
	case 0xB:
		c.X3B_RLA_Absolute_Y()
	case 0xC:
		c.X3C_NOP_Absolute_X()
	case 0xD:
		c.X3D_AND_Absolute_X()
	case 0xE:
		c.X3E_ROL_Absolute_X()
	case 0xF:
		c.X3F_RLA_Absolute_X()
	default:
		c.UnknownOpcode()
	}
}
func (c *CPU) Opcode4X() {
	switch c.opcode & 0xF {
	case 0x0:
		c.X40_RTI()
	case 0x1:
		c.X41_EOR_Indirect_X()
	case 0x2:
		c.Kill()
	case 0x3:
		c.X43_SRE_Indirect_X()
	case 0x4:
		c.X44_NOP_ZeroPage()
	case 0x5:
		c.X45_EOR_ZeroPage()
	case 0x6:
		c.X46_LSR_ZeroPage()
	case 0x7:
		c.X47_SRE_ZeroPage()
	case 0x8:
		c.X48_PHA()
	case 0x9:
		c.X49_EOR_Immediate()
	case 0xA:
		c.X4A_LSR()
	case 0xB:
		c.X4B_ALR_Immediate()
	case 0xC:
		c.X4C_JMP()
	case 0xD:
		c.X4D_EOR_Absolute()
	case 0xE:
		c.X4E_LSR_Absolute()
	case 0xF:
		c.X4F_SRE_Absolute()
	default:
		c.UnknownOpcode()
	}
}
func (c *CPU) Opcode5X() {
	switch c.opcode & 0xF {
	case 0x0:
		c.X50_BVC()
	case 0x1:
		c.X51_EOR_Indirect_Y()
	case 0x2:
		c.Kill()
	case 0x3:
		c.X53_SRE_Indirect_Y()
	case 0x4:
		c.X54_NOP_ZeroPage_X()
	case 0x5:
		c.X55_EOR_ZeroPage_X()
	case 0x6:
		c.X56_LSR_ZeroPage_X()
	case 0x7:
		c.X57_SRE_ZeroPage_X()
	case 0x8:
		c.X58_CLI()
	case 0x9:
		c.X59_EOR_Absolute_Y()
	case 0xA:
		c.X5A_NOP_Implied()
	case 0xB:
		c.X5B_SRE_Absolute_Y()
	case 0xC:
		c.X5C_NOP_Absolute_X()
	case 0xD:
		c.X5D_EOR_Absolute_X()
	case 0xE:
		c.X5E_LSR_Absolute_X()
	case 0xF:
		c.X5F_SRE_Absolute_X()
	default:
		c.UnknownOpcode()
	}
}
func (c *CPU) Opcode6X() {
	switch c.opcode & 0xF {
	case 0x0:
		c.X60_RTS()
	case 0x1:
		c.X61_ADC_Indirect_X()
	case 0x2:
		c.Kill()
	case 0x3:
		c.X63_RRA_Indirect_X()
	case 0x4:
		c.X64_NOP_ZeroPage()
	case 0x5:
		c.X65_ADC_ZeroPage()
	case 0x6:
		c.X66_ROR_ZeroPage()
	case 0x7:
		c.X67_RRA_ZeroPage()
	case 0x8:
		c.X68_PLA()
	case 0x9:
		c.X69_ADC_Immediate()
	case 0xA:
		c.X6A_ROR()
	case 0xB:
		c.X6B_ARR_Immediate()
	case 0xC:
		c.X6C_JMP_Indirect()
	case 0xD:
		c.X6D_ADC_Absolute()
	case 0xE:
		c.X6E_ROR_Absolute()
	case 0xF:
		c.X6F_RRA_Absolute()
	default:
		c.UnknownOpcode()
	}
}
func (c *CPU) Opcode7X() {
	switch c.opcode & 0xF {
	case 0x0:
		c.X70_BVS()
	case 0x1:
		c.X71_ADC_Indirect_Y()
	case 0x2:
		c.Kill()
	case 0x3:
		c.X73_RRA_Indirect_Y()
	case 0x4:
		c.X74_NOP_ZeroPage_X()
	case 0x5:
		c.X75_ADC_ZeroPage_X()
	case 0x6:
		c.X76_ROR_ZeroPage_X()
	case 0x7:
		c.X77_RRA_ZeroPage_X()
	case 0x8:
		c.X78_SEI()
	case 0x9:
		c.X79_ADC_Absolute_Y()
	case 0xA:
		c.X7A_NOP_Implied()
	case 0xB:
		c.X7B_RRA_Absolute_Y()
	case 0xC:
		c.X7C_NOP_Absolute_X()
	case 0xD:
		c.X7D_ADC_Absolute_X()
	case 0xE:
		c.X7E_ROR_Absolute_X()
	case 0xF:
		c.X7F_RRA_Absolute_X()
	default:
		c.UnknownOpcode()
	}
}
func (c *CPU) Opcode8X() {
	switch c.opcode & 0xF {
	case 0x0:
		c.X80_NOP_Immediate()
	case 0x1:
		c.X81_STA_Indirect_X()
	case 0x2:
		c.X82_NOP_Immediate()
	case 0x3:
		c.X83_SAX_Indirect_X()
	case 0x4:
		c.X84_STY_ZeroPage()
	case 0x5:
		c.X85_STA_ZeroPage()
	case 0x6:
		c.X86_STX_ZeroPage()
	case 0x7:
		c.X87_SAX_ZeroPage()
	case 0x8:
		c.X88_DEY()
	case 0x9:
		c.X89_NOP_Immediate()
	case 0xA:
		c.X8A_TXA()
	case 0xB:
		c.X8B_ANE_Immediate()
	case 0xC:
		c.X8C_STY_Absolute()
	case 0xD:
		c.X8D_STA_Absolute()
	case 0xE:
		c.X8E_STX_Absolute()
	case 0xF:
		c.X8F_SAX_Absolute()
	default:
		c.UnknownOpcode()
	}
}
func (c *CPU) Opcode9X() {
	switch c.opcode & 0xF {
	case 0x0:
		c.X90_BCC()
	case 0x1:
		c.X91_STA_Indirect_Y()
	case 0x2:
		c.Kill()
	case 0x3:
		c.X93_SHA_Indirect_Y()
	case 0x4:
		c.X94_STY_ZeroPage_X()
	case 0x5:
		c.X95_STA_ZeroPage_X()
	case 0x6:
		c.X96_STX_ZeroPage_Y()
	case 0x7:
		c.X97_SAX_ZeroPage_Y()
	case 0x8:
		c.X98_TYA()
	case 0x9:
		c.X99_STA_Absolute_Y()
	case 0xA:
		c.X9A_TXS()
	case 0xB:
		c.X9B_TAS_Absolute_Y()
	case 0xC:
		c.X9C_SHY_Absolute_X()
	case 0xD:
		c.X9D_STA_Absolute_X()
	case 0xE:
		c.X9E_SHX_Absolute_Y()
	case 0xF:
		c.X9F_SHA_Absolute_Y()
	default:
		c.UnknownOpcode()
	}
}
func (c *CPU) OpcodeAX() {
	switch c.opcode & 0xF {
	case 0x0:
		c.XA0_LDY_Immediate()
	case 0x1:
		c.XA1_LDA_Indirect_X()
	case 0x2:
		c.XA2_LDX_Immediate()
	case 0x3:
		c.XA3_LAX_Indirect_X()
	case 0x4:
		c.XA4_LDY_ZeroPage()
	case 0x5:
		c.XA5_LDA_ZeroPage()
	case 0x6:
		c.XA6_LDX_ZeroPage()
	case 0x7:
		c.XA7_LAX_ZeroPage()
	case 0x8:
		c.XA8_TAY()
	case 0x9:
		c.XA9_LDA_Immediate()
	case 0xA:
		c.XAA_TAX()
	case 0xB:
		c.XAB_LXA_Immediate()
	case 0xC:
		c.XAC_LDY_Absolute()
	case 0xD:
		c.XAD_LDA_Absolute()
	case 0xE:
		c.XAE_LDX_Absolute()
	case 0xF:
		c.XAF_LAX_Absolute()
	default:
		c.UnknownOpcode()
	}
}
func (c *CPU) OpcodeBX() {
	switch c.opcode & 0xF {
	case 0x0:
		c.XB0_BCS()
	case 0x1:
		c.XB1_LDA_Indirect_Y()
	case 0x2:
		c.Kill()
	case 0x3:
		c.XB3_LAX_Indirect_Y()
	case 0x4:
		c.XB4_LDY_ZeroPage_X()
	case 0x5:
		c.XB5_LDA_ZeroPage_X()
	case 0x6:
		c.XB6_LDX_ZeroPage_Y()
	case 0x7:
		c.XB7_LAX_ZeroPage_Y()
	case 0x8:
		c.XB8_CLV()
	case 0x9:
		c.XB9_LDA_Absolute_Y()
	case 0xA:
		c.XBA_TSX()
	case 0xB:
		c.XBB_LAS_Absolute_Y()
	case 0xC:
		c.XBC_LDY_Absolute_X()
	case 0xD:
		c.XBD_LDA_Absolute_X()
	case 0xE:
		c.XBE_LDX_Absolute_Y()
	case 0xF:
		c.XBF_LAX_Absolute_Y()
	default:
		c.UnknownOpcode()
	}
}
func (c *CPU) OpcodeCX() {
	switch c.opcode & 0xF {
	case 0x0:
		c.XC0_CPY_Immediate()
	case 0x1:
		c.XC1_CMP_Indirect_X()
	case 0x2:
		c.XC2_NOP_Immediate()
	case 0x3:
		c.XC3_DCP_Indirect_X()
	case 0x4:
		c.XC4_CPY_ZeroPage()
	case 0x5:
		c.XC5_CMP_ZeroPage()
	case 0x6:
		c.XC6_DEC_ZeroPage()
	case 0x7:
		c.XC7_DCP_ZeroPage()
	case 0x8:
		c.XC8_INY()
	case 0x9:
		c.XC9_CMP_Immediate()
	case 0xA:
		c.XCA_DEX()
	case 0xB:
		c.XCB_SBX_Immediate()
	case 0xC:
		c.XCC_CPY_Absolute()
	case 0xD:
		c.XCD_CMP_Absolute()
	case 0xE:
		c.XCE_DEC_Absolute()
	case 0xF:
		c.XCF_DCP_Absolute()
	default:
		c.UnknownOpcode()
	}
}
func (c *CPU) OpcodeDX() {
	switch c.opcode & 0xF {
	case 0x0:
		c.XD0_BNE()
	case 0x1:
		c.XD1_CMP_Indirect_Y()
	case 0x2:
		c.Kill()
	case 0x3:
		c.XD3_DCP_Indirect_Y()
	case 0x4:
		c.XD4_NOP_ZeroPage_X()
	case 0x5:
		c.XD5_CMP_ZeroPage_X()
	case 0x6:
		c.XD6_DEC_ZeroPage_X()
	case 0x7:
		c.XD7_DCP_ZeroPage_X()
	case 0x8:
		c.XD8_CLD()
	case 0x9:
		c.XD9_CMP_Absolute_Y()
	case 0xA:
		c.XDA_NOP_Implied()
	case 0xB:
		c.XDB_DCP_Absolute_Y()
	case 0xC:
		c.XDC_NOP_Absolute_X()
	case 0xD:
		c.XDD_CMP_Absolute_X()
	case 0xE:
		c.XDE_DEC_Absolute_X()
	case 0xF:
		c.XDF_DCP_Absolute_X()
	default:
		c.UnknownOpcode()
	}
}
func (c *CPU) OpcodeEX() {
	switch c.opcode & 0xF {
	case 0x0:
		c.XE0_CPX_Immediate()
	case 0x1:
		c.XE1_SBC_Indirect_X()
	case 0x2:
		c.XE2_NOP_Immediate()
	case 0x3:
		c.XE3_ISC_Indirect_X()
	case 0x4:
		c.XE4_CPX_ZeroPage()
	case 0x5:
		c.XE5_SBC_ZeroPage()
	case 0x6:
		c.XE6_INC_ZeroPage()
	case 0x7:
		c.XE7_ISC_ZeroPage()
	case 0x8:
		c.XE8_INX()
	case 0x9:
		c.XE9_SBC_Immediate()
	case 0xA:
		c.XEA_NOP()
	case 0xB:
		c.XEB_SBC_Immediate()
	case 0xC:
		c.XEC_CPX_Absolute()
	case 0xD:
		c.XED_SBC_Absolute()
	case 0xE:
		c.XEE_INC_Absolute()
	case 0xF:
		c.XEF_ISC_Absolute()
	default:
		c.UnknownOpcode()
	}
}
func (c *CPU) OpcodeFX() {
	switch c.opcode & 0xF {
	case 0x0:
		c.XF0_BEQ()
	case 0x1:
		c.XF1_SBC_Indirect_Y()
	case 0x2:
		c.Kill()
	case 0x3:
		c.XF3_ISC_Indirect_Y()
	case 0x4:
		c.XF4_NOP_ZeroPage_X()
	case 0x5:
		c.XF5_SBC_ZeroPage_X()
	case 0x6:
		c.XF6_INC_ZeroPage_X()
	case 0x7:
		c.XF7_ISC_ZeroPage_X()
	case 0x8:
		c.XF8_SED()
	case 0x9:
		c.XF9_SBC_Absolute_Y()
	case 0xA:
		c.XFA_NOP_Implied()
	case 0xB:
		c.XFB_ISC_Absolute_Y()
	case 0xC:
		c.XFC_NOP_Absolute_X()
	case 0xD:
		c.XFD_SBC_Absolute_X()
	case 0xE:
		c.XFE_INC_Absolute_X()
	case 0xF:
		c.XFF_ISC_Absolute_X()
	default:
		c.UnknownOpcode()
	}
}

func (c *CPU) PollInterrupts() {
	if c.RunningInterrupt {
		return
	}

	if c.PollNMI() {
		c.NMIPending = true
	}
	c.IRQPending = c.IRQLine && !c.flag_InterruptDisable
}

func (c *CPU) PollInterrupts_CantDisableIRQ() {
	if c.RunningInterrupt {
		return
	}

	if c.PollNMI() {
		c.NMIPending = true
	}
	if !c.IRQPending {
		c.IRQPending = c.IRQLine && !c.flag_InterruptDisable
	}
}

func (c *CPU) UnknownOpcode() {
	//fmt.Println("Unknown Opcode: " + fmt.Sprintf("%02X", c.opcode))
	c.CompleteInstruction()
}

func (c *CPU) PollNMI() bool {
	prevNMILevelDetector := c.NMILevelDetector
	c.NMILevelDetector = c.NMILine
	return !prevNMILevelDetector && c.NMILevelDetector && !ppu.SuppressNMI
}

func (c *CPU) SetNMILine() {
	if !c.NMILine {
		c.NMILine = ppu.PPUCTRL_EnableNMI && ppu.PPUSTATUS_VBlank
	}
}

func (c *CPU) SetIRQLine() {
	c.IRQLine = apu.IRQLevelDetector || mappers.MMC3_IRQPending
}

func (c *CPU) PollIRQ() bool {
	//if c.flag_InterruptDisable {
	//	c.IRQPending = false
	//}
	return (c.IRQPending || c.apu.GetFrameInterrupt() || mappers.MMC3_IRQPending) && !c.flag_InterruptDisable
}

func (c *CPU) DisableNMI() {
	c.NMIPending = false
}

func (c *CPU) DisableIRQFlags() {
	c.IRQPending = false
	apu.APUDMCInterrupt = false
	c.apu.SetFrameInterrupt(false)
	mappers.MMC3_IRQPending = false
}

func (c *CPU) SetOpcode(code byte) {
	c.opcode = code
}

func (c *CPU) CompleteInstruction() {
	//c.PollInterrupts()
	c.InstructionCycle = -1
	c.AddressBus = c.PC
}

func (c *CPU) SendToDebug() {
	status := byte(0)
	status += byte(common.Ternary(c.flag_Carry, 0x01, 0x00))
	status += byte(common.Ternary(c.flag_Zero, 0x02, 0x00))
	status += byte(common.Ternary(c.flag_InterruptDisable, 0x04, 0x00))
	status += byte(common.Ternary(c.flag_Decimal, 0x08, 0x00))
	status += byte(common.Ternary(c.flag_B, 0x10, 0x00)) //B Flag
	status += 0x20
	status += byte(common.Ternary(c.flag_Overflow, 0x40, 0x00))
	status += byte(common.Ternary(c.flag_Negative, 0x80, 0x00))
	debug.SetCPUData(c.opcode, c.A, c.X, c.Y, c.SP, status, c.PC, c.AddressBus)
	if debug.LoggingCPU {
		debug.TraceLogger(c.opcode, c.A, c.X, c.Y, c.SP, status, c.PC)
	}
}
