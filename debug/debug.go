package debug

import (
	"fmt"
	"mtt/timenes/common"
	"mtt/timenes/nes/cartridge"
	"mtt/timenes/nes/cpu"
	"mtt/timenes/nes/ppu"
)

var LoggingCPU = false
var LoggingPPU = false
var LogCount = 1000
var InstructionCount int = 0
var frame int
var tracePC, traceVRAM uint16
var traceA, traceX, tracyY, traceSP byte
var traceFlagC, traceFlagZ, traceFlagI, traceFlagD, traceFlagV, traceFlagN bool
var traceCycles int

var cycleTest string

// Sets everything outside of operands for the tracelogger to run later
func PrepTraceLogger() {
	tracePC = cpu.PC
	traceVRAM = ppu.VRAMAddress

	traceA = cpu.A
	traceX = cpu.X
	tracyY = cpu.Y
	traceSP = cpu.SP

	/*traceFlagC = nes.flag_Carry
	traceFlagZ = nes.flag_Zero
	traceFlagI = nes.flag_InterruptDisable
	traceFlagD = nes.flag_Decimal
	traceFlagV = nes.flag_Overflow
	traceFlagN = nes.flag_Negative*/

	traceCycles = common.CPU_TotalCycles
}

/*
func TraceLogger() {
	OpCodeNames := [...]string{
		"BRK", "ORA", "HLT", "SLO", "NOP", "ORA", "ASL", "SLO", "PHP", "ORA",
		"ASL", "ANC", "NOP", "ORA", "ASL", "SLO", "BPL", "ORA", "HLT", "SLO",
		"NOP", "ORA", "ASL", "SLO", "CLC", "ORA", "NOP", "SLO", "NOP", "ORA",
		"ASL", "SLO", "JSR", "AND", "HLT", "RLA", "BIT", "AND", "ROL", "RLA",
		"PLP", "AND", "ROL", "ANC", "BIT", "AND", "ROL", "RLA", "BMI", "AND",
		"HLT", "RLA", "NOP", "AND", "ROL", "RLA", "SEC", "AND", "NOP", "RLA",
		"NOP", "AND", "ROL", "RLA", "RTI", "EOR", "HLT", "SRE", "NOP", "EOR",
		"LSR", "SRE", "PHA", "EOR", "LSR", "ALR", "JMP", "EOR", "LSR", "SRE",
		"BVC", "EOR", "HLT", "SRE", "NOP", "EOR", "LSR", "SRE", "CLI", "EOR",
		"NOP", "SRE", "NOP", "EOR", "LSR", "SRE", "RTS", "ADC", "HLT", "RRA",
		"NOP", "ADC", "ROR", "RRA", "PLA", "ADC", "ROR", "ARR", "JMP", "ADC",
		"ROR", "RRA", "BVS", "ADC", "HLT", "RRA", "NOP", "ADC", "ROR", "RRA",
		"SEI", "ADC", "NOP", "RRA", "NOP", "ADC", "ROR", "RRA", "NOP", "STA",
		"NOP", "SAX", "STY", "STA", "STX", "SAX", "DEY", "NOP", "TXA", "ANE",
		"STY", "STA", "STX", "SAX", "BCC", "STA", "HLT", "SHA", "STY", "STA",
		"STX", "SAX", "TYA", "STA", "TXS", "SHS", "SHY", "STA", "SHX", "SHA",
		"LDY", "LDA", "LDX", "LAX", "LDY", "LDA", "LDX", "LAX", "TAY", "LDA",
		"TAX", "LXA", "LDY", "LDA", "LDX", "LAX", "BCS", "LDA", "HLT", "LAX",
		"LDY", "LDA", "LDX", "LAX", "CLV", "LDA", "TSX", "LAE", "LDY", "LDA",
		"LDX", "LAX", "CPY", "CMP", "NOP", "DCP", "CPY", "CMP", "DEC", "DCP",
		"INY", "CMP", "DEX", "AXS", "CPY", "CMP", "DEC", "DCP", "BNE", "CMP",
		"HLT", "DCP", "NOP", "CMP", "DEC", "DPC", "CLD", "CMP", "NOP", "DCP",
		"NOP", "CMP", "DEC", "DCP", "CPX", "SBC", "NOP", "ISC", "CPX", "SBC",
		"INC", "ISC", "INX", "SBC", "NOP", "SBC", "CPX", "SBC", "INC", "ISC",
		"BEQ", "SBC", "HLT", "ISC", "NOP", "SBC", "INC", "ISC", "SED", "SBC",
		"NOP", "ISC", "NOP", "SBC", "INC", "ISC",
	}
	if LoggingCPU /*&& InstructionCount > 35000* / {
		Traceline := "$" + fmt.Sprintf("%04X", tracePC)

		Traceline += "\t" + fmt.Sprintf("%02X ", opcode)
		opLength := len(operands)
		for i := range operands {
			Traceline += fmt.Sprintf("%02X ", operands[i])
		}

		if opLength < 2 {
			Traceline += "\t"
		}

		Traceline += "\t" + OpCodeNames[opcode]

		//Add operand / addresses after opcode

		if opcode == 0x99 {
			Traceline += " $" + fmt.Sprintf("%02X", operands[1]) + fmt.Sprintf("%02X", operands[0])
			Traceline += ", Y -> $"
			Traceline += fmt.Sprintf("%04X", BuildAddress(operands[0], operands[1])+uint16(Y))
		}

		Traceline += "\tA:" + fmt.Sprintf("%02X", traceA)
		Traceline += "\tX:" + fmt.Sprintf("%02X", traceX)
		Traceline += "\tY:" + fmt.Sprintf("%02X", tracyY)
		Traceline += "\tSP:" + fmt.Sprintf("%02X", traceSP)
		Traceline += "\t"

		//Processor Flags
		if traceFlagN {
			Traceline += "N"
		} else {
			Traceline += "n"
		}

		if traceFlagV {
			Traceline += "V"
		} else {
			Traceline += "v"
		}

		Traceline += "--"

		if traceFlagD {
			Traceline += "D"
		} else {
			Traceline += "d"
		}

		if traceFlagI {
			Traceline += "I"
		} else {
			Traceline += "i"
		}

		if traceFlagZ {
			Traceline += "Z"
		} else {
			Traceline += "z"
		}

		if traceFlagC {
			Traceline += "C"
		} else {
			Traceline += "c"
		}

		Traceline += "\tVRAM: " + fmt.Sprintf("%04X", traceVRAM)
		Traceline += "\tCycle: " + fmt.Sprintf("%d", traceCycles)
		Traceline += "\tInstructionCount: " + fmt.Sprintf("%d", InstructionCount)

		fmt.Println(Traceline)
		/*LogCount--
		if LogCount < 0 {
			CPU_Halted = true
		}* /
	}
}*/

/*
func TraceLoggerPPU() {

	if LoggingPPU {

		Traceline := fmt.Sprintf("D: %-3d", ppuDot)
		Traceline += fmt.Sprintf("\tSL: %-3d", ppuScanline)

		//Add Shift registers

		Traceline += fmt.Sprintf("\tpL: %016b", ppuShiftRegister_patternL)
		Traceline += fmt.Sprintf("\tpH: %016b", ppuShiftRegister_patternH)
		Traceline += fmt.Sprintf("\taL: %016b", ppuShiftRegister_attributeL)
		Traceline += fmt.Sprintf("\taH: %016b", ppuShiftRegister_attributeH)

		//Print PalHi and PalLow
		col0 := byte((ppuShiftRegister_patternL >> (15 - ppuScrollFineX)) & 1)
		col1 := byte((ppuShiftRegister_patternH >> (15 - ppuScrollFineX)) & 1)

		pal0 := byte((ppuShiftRegister_attributeL >> (15 - ppuScrollFineX)) & 1)
		pal1 := byte((ppuShiftRegister_attributeH >> (15 - ppuScrollFineX)) & 1)

		Traceline += fmt.Sprintf("\tLo: %02b Hi: %02b", byte((col1<<1)|col0), byte((pal1<<1)|pal0))

		fmt.Println(Traceline)
		/*LogCount--
		if LogCount < 0 {
			CPU_Halted = true
		}* /
	}
}*/
/*
func DebugWindow( g *Game ) error { //Using ebitengine's built in debug display

	if _, err := g.debugui.Update(func(ctx *debugui.Context) error {
		ctx.Window("Debugging info", image.Rect(10, 100, 260, 300), func(layout debugui.ContainerLayout) {
			ctx.SetGridLayout([]int{100, -1}, nil)
			ctx.Text("Instruction Count:")
			ctx.Text(fmt.Sprintf("$%d", InstructionCount))
			//ctx.Text("VRAM Address:")
			//ctx.Text(fmt.Sprintf("$%04X", VRAMAddress))
			//ctx.Text("T Register:")
			//ctx.Text(fmt.Sprintf("$%04X", TransferAddress))
			//ctx.Text("Fine X Scroll:")
			//ctx.Text(fmt.Sprintf("$%02X", ppuScrollFineX))
			//ctx.Text("Fine Y Scroll:")
			//ctx.Text(fmt.Sprintf("$%02X", ppuScrollFineX))
			//ctx.Text("Nametable:")
			//ctx.Text(fmt.Sprintf("$%02X", PPUCTRL_NametableSelect))

			//ctx.Text("Pause")
			ctx.Button("Pause").On(func() {
				pauseEmulation = !pauseEmulation
			})
			ctx.Button("RESET").On(func() {
				nes.ROMLoaded = false
			})
		})
		return nil
	}); err != nil {
		return err
	}
	return nil
}
*/
var CartRamLastString string

func CartRAMLogger() {
	CartTestStatus := cartridge.CartRAM[0]
	CartTestValid := cartridge.CartRAM[1] == 0xDE && cartridge.CartRAM[2] == 0xB0 && cartridge.CartRAM[3] == 0x61

	newString := string(cartridge.CartRAM[4:])

	if CartTestValid && CartTestStatus == 0x80 && (newString != CartRamLastString) {
		//Valid Test line
		CartRamLastString = newString
		fmt.Print(newString)
	}
}
