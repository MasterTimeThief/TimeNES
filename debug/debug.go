package debug

import (
	"fmt"
	"image"
	"image/color"
	"mtt/timenes/common"
	"mtt/timenes/nes/apu"
	"mtt/timenes/nes/bus"
	"mtt/timenes/nes/cartridge"
	"mtt/timenes/nes/cpu"
	"mtt/timenes/nes/ppu"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Header
//var (
//	dbHeaderType   int // 0: iNES 1: NES.2.0
//	dbPRGSize      uint32
//	dbCHRSize      uint32
//	dbNametable    bool // 0: Horizontal 1: Vertical
//	dbNametableAlt bool
//	dbHasTrainer   bool
//	dbHasSRAM      bool
//	dbMapperID     byte
//	dbSubmapperID  byte
//)

var SizeText = map[uint32]string{
	0x0:     "Has CHR-ROM",
	0x4000:  "16 K",
	0x8000:  "32 K",
	0x10000: "64 K",
	0x20000: "128 K",
	0x40000: "256 K",
	0x80000: "512 K",
}

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
var ShowFPS bool = false
var ShowDebugWindow, ShowPatternTableWindow bool = false, false
var PatternTables = ebiten.NewImage(256, 128)

var PTValues [256][128]byte

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

var CartRamLastString string

/*func CartRAMLogger() {
	CartTestStatus := cartridge.CartRAM[0]
	CartTestValid := cartridge.CartRAM[1] == 0xDE && cartridge.CartRAM[2] == 0xB0 && cartridge.CartRAM[3] == 0x61

	newString := string(cartridge.CartRAM[4:])

	if CartTestValid && CartTestStatus == 0x80 && (newString != CartRamLastString) {
		//Valid Test line
		CartRamLastString = newString
		fmt.Print(newString)
	}
}*/

func DisplayDebugging(screen *ebiten.Image) {

	if ShowFPS {
		vector.FillRect(screen, float32(256*common.ScreenScale)-135, float32(240*common.ScreenScale)-17, 135, 17, color.Black, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f\tFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()), (256*common.ScreenScale)-132, (240*common.ScreenScale)-15)
	}
	if bus.OutsideCodeRead > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Attempting to Read at $: %04X", bus.OutsideCodeRead), 5, (240*common.ScreenScale)-30)
	}
	if bus.OutsideCodeWrite > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Attempting to Write at $: %04X", bus.OutsideCodeWrite), 5, (240*common.ScreenScale)-40)
	}
}

func Update(ui *debugui.DebugUI) {
	if common.PendingFPS {
		ShowFPS = !ShowFPS
		common.PendingFPS = false
	}
	if common.PendingDebug {
		ShowDebugWindow = !ShowDebugWindow
		common.PendingDebug = false
	}
	if common.PendingPT {
		ShowPatternTableWindow = !ShowPatternTableWindow
		common.PendingPT = false
	}
	if common.PendingPatternUpdate {
		UpdatePatternTables()
		common.PendingPatternUpdate = false
	}
	ui.Update(func(ctx *debugui.Context) error {
		DebugWindow(ctx)
		PatternTableWindow(ctx)
		return nil
	})
}

func DebugWindow(ctx *debugui.Context) {
	if ShowDebugWindow {
		ctx.Window("Debugging info", image.Rect(10, 100, 260, 300), func(layout debugui.ContainerLayout) {

			screenScaleOptions := []string{"1x", "2x", "3x", "4x"}
			ctx.Header("Functions", true, func() {
				ctx.SetGridLayout([]int{-1, -1}, nil)
				ctx.Button("Pause").On(func() {
					common.PendingPause = true
				})
				ctx.Button("Reset").On(func() {
					common.PendingReset = true
				})
				ctx.Text("Window Scale:")
				screenScale := common.ScreenScale - 1
				ctx.Dropdown(&screenScale, screenScaleOptions).On(func() {
					common.NewScreenScale = screenScale + 1
					//g.writeLog(fmt.Sprintf("Selected option: %s", g.dropdownOptions1[g.selectedOption1]))
				})
			})
			//Header Info
			headerType := "iNES"
			if cartridge.NES2_Header {
				headerType = "NES 2.0"
			}
			ctx.Header("Header Info ("+headerType+")", false, func() {
				ctx.SetGridLayout([]int{-1, -1}, nil)
				ctx.Text("PRG-ROM Size:")
				ctx.Text(SizeText[cartridge.PRGROM_Size])
				ctx.Text("CHR-RAM Size:")
				ctx.Text(SizeText[cartridge.CHRROM_Size])
				ctx.Text("Mapper Chip:")
				ctx.Text(fmt.Sprintf("%d", cartridge.MapperChipID))
				ctx.Text("Submapper:")
				ctx.Text(fmt.Sprintf("%d", cartridge.SubmapperID))
				ctx.Text("Nametable Arrangement:")
				if cartridge.IsNametableHorizontal {
					ctx.Text("Horizontal")
				} else {
					ctx.Text("Vertical")
				}
				ctx.Text("Alt Nametable Arr:")
				if cartridge.AltNametableLayout {
					ctx.Text("Yes")
				} else {
					ctx.Text("No")
				}
				ctx.Text("SRAM")
				if cartridge.HasBatteryRAM {
					ctx.Text("Yes")
				} else {
					ctx.Text("No")
				}
				ctx.Text("Trainer")
				if cartridge.HasTrainer {
					ctx.Text("Yes")
				} else {
					ctx.Text("No")
				}

			})
			//CPU Info
			ctx.Header("CPU Info", false, func() {
				ctx.SetGridLayout([]int{-1, -1}, nil)
			})

			//APU Info
			ctx.Header("APU Info", false, func() {
				ctx.SetGridLayout([]int{-1, -1}, nil)
				ctx.GridCell(func(bounds image.Rectangle) {
					ctx.TreeNode("Toggle Channels", func() {
						ctx.Checkbox(&apu.Pulse1.ForceMute, "Square 1")
						ctx.Checkbox(&apu.Pulse2.ForceMute, "Square 2")
						ctx.Checkbox(&apu.Triangle.ForceMute, "Triangle")
						ctx.Checkbox(&apu.Noise.ForceMute, "Noise")
						ctx.Checkbox(&apu.DMC.ForceMute, "DMC")
					})
				})
				ctx.GridCell(func(bounds image.Rectangle) {
					ctx.Button("Toggle Audio").On(func() {
						common.PendingMute = true
					})
				})
			})

			ctx.SetGridLayout([]int{100, -1}, nil)
			//ctx.Text("Instruction Count:")
			//ctx.Text(fmt.Sprintf("$%d", InstructionCount))
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

		})
	}
}

func PatternTableWindow(ctx *debugui.Context) {
	if ShowPatternTableWindow {
		px, py := 100, 20
		ctx.Window("Pattern Tables", image.Rect(px, py, px+256+10, py+128+35), func(layout debugui.ContainerLayout) {

			ctx.GridCell(func(bounds image.Rectangle) {
				ctx.DrawOnlyWidget(func(screen *ebiten.Image) {
					op := &ebiten.DrawImageOptions{}
					scale := float64(1)
					op.GeoM.Translate(float64(bounds.Min.X)/scale, float64(bounds.Min.Y)/scale)
					op.GeoM.Scale(scale, scale)

					DrawPatternTables()
					screen.DrawImage(PatternTables, op)
				})
			})

		})
	}
}

func UpdatePatternTables() {
	for table := 0; table < 2; table++ {
		for row := 0; row < 16; row++ {
			for col := 0; col < 16; col++ {
				for y := 0; y < 8; y++ {
					lowByte := cartridge.CHRROM[y+(col*16)+(row*256)+(table*4096)]
					hiByte := cartridge.CHRROM[8+y+(col*16)+(row*256)+(table*4096)]
					for x := 0; x < 8; x++ {
						var TwoBit uint8
						if ((lowByte >> (7 - x)) & 1) == 1 {
							TwoBit += 1
						}
						if ((hiByte >> (7 - x)) & 1) == 1 {
							TwoBit += 2
						}
						PTValues[x+(col*8)+(table*128)][y+(row*8)] = TwoBit * 85
					}
				}
			}
		}
	}
}

func ResetPatternTables() {
	PatternTables.Fill(color.Black)
	common.PendingPatternUpdate = true
}

func DrawPatternTables() {
	screenPT := image.NewRGBA(image.Rect(0, 0, 256, 128))
	index := 0
	for y := 0; y < 128; y++ {
		for x := 0; x < 256; x++ {
			col := PTValues[x][y]
			screenPT.Pix[index] = col
			screenPT.Pix[index+1] = col
			screenPT.Pix[index+2] = col
			screenPT.Pix[index+3] = 0xFF

			index += 4
		}
	}
	PatternTables.WritePixels(screenPT.Pix)
}
