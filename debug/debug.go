package debug

import (
	"fmt"
	"image"
	"image/color"
	"mtt/timenes/common"
	"mtt/timenes/nes/bus"
	"mtt/timenes/nes/cartridge"

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

var MapperName = map[byte]string{
	0: "NROM",
	1: "MMC1",
	2: "UxROM",
	3: "CNROM",
	4: "MMC3",
	7: "AxROM",
}

type DEBUG struct {
	apu APU
}

var D DEBUG

type APU interface {
	GetPulse1Mute() *bool
	GetPulse2Mute() *bool
	GetTriangleMute() *bool
	GetNoiseMute() *bool
	GetDMCMute() *bool
	SetEmulatorvolume(float64)
}

var LoggingCPU = false
var LoggingPPU = false
var LogCount = 1000
var InstructionCount int = 0
var frame int

var cycleTest string
var ShowFPS bool = false
var ShowDebugWindow, ShowPatternTableWindow bool = false, false
var PatternTables = ebiten.NewImage(256, 128)

var PTValues [256][128]byte

func InitDEBUG(a APU) {
	D.apu = a
}

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
		ctx.Window("Debugging info", image.Rect(10, 100, 260, 400), func(layout debugui.ContainerLayout) {

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
				ctx.Text(fmt.Sprintf("%d (%s)", cartridge.MapperChipID, MapperName[cartridge.MapperChipID]))
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
				ctx.SetGridLayout([]int{60, -1}, nil)
				ctx.Text("Volume: ")
				ctx.GridCell(func(bounds image.Rectangle) {
					ctx.SliderF(&common.EmulatorVolume, 0, 100, 1, 0).On(func() {
						D.apu.SetEmulatorvolume(common.EmulatorVolume / 100)
						//g.writeLog(fmt.Sprintf("Selected option: %s", g.dropdownOptions1[g.selectedOption1]))
					})
				})
				ctx.GridCell(func(bounds image.Rectangle) {
					ctx.Button("Mute").On(func() {
						common.PendingMute = true
					})
				})
				ctx.GridCell(func(bounds image.Rectangle) {
					ctx.TreeNode("Toggle Channels", func() {
						ctx.Checkbox(D.apu.GetPulse1Mute(), "Square 1")
						ctx.Checkbox(D.apu.GetPulse2Mute(), "Square 2")
						ctx.Checkbox(D.apu.GetTriangleMute(), "Triangle")
						ctx.Checkbox(D.apu.GetNoiseMute(), "Noise")
						ctx.Checkbox(D.apu.GetDMCMute(), "DMC")
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
