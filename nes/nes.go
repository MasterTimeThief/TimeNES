package nes

import (
	"flag"
	"image"
	"log"
	"mtt/timenes/common"
	"mtt/timenes/debug"
	"mtt/timenes/nes/apu"
	"mtt/timenes/nes/bus"
	"mtt/timenes/nes/cartridge"
	"mtt/timenes/nes/cpu"
	"mtt/timenes/nes/ppu"

	"github.com/ebitengine/debugui"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/draw"
)

// Header Variables

// var cpuClock int = 0 // 2A03
// var ppuClock int = 0 // 2C02
var MasterClock int = 0

var FullscreenMode bool = false
var PauseEmulation bool = false
var MenuBarSelected bool = false

type Game struct {
	cpu *cpu.CPU
	ppu *ppu.PPU
	apu *apu.APU
	bus *bus.BUS

	gameScreen *image.RGBA
	UI         *ebitenui.UI
	Exit       bool
	debugui    debugui.DebugUI
}

var Emulator *Game

func InitGame(newUI *ebitenui.UI) {
	arguments := flag.Args()
	if len(arguments) > 0 && arguments[0] != "" {
		common.Filepath = arguments[0]
	}
	if common.Filepath != "" {
		common.ROMExists = true
	}

	ebiten.SetWindowSize(common.ScreenWidth*common.ScreenScale, common.ScreenHeight*common.ScreenScale)
	ebiten.SetWindowTitle("TimeNES")
	ebiten.SetFullscreen(FullscreenMode)
	common.SetWindowIcon()
	//ebiten.SetTPS(ebiten.SyncWithFPS)

	Emulator = &Game{
		gameScreen: image.NewRGBA(image.Rect(0, 0, common.ScreenWidth, common.ScreenHeight)),
		UI:         newUI,
		cpu:        cpu.NewCPU(),
		ppu:        ppu.NewPPU(),
		apu:        apu.NewAPU(),
		bus:        bus.NewBUS(),
	}
	Emulator.cpu.SetBUS(Emulator.bus)
	Emulator.bus.SetCPU(Emulator.cpu)
	Emulator.bus.SetAPU(Emulator.apu)
	Emulator.apu.SetCPU(Emulator.cpu)

	debug.InitDEBUG(Emulator.apu)

	if err := ebiten.RunGame(Emulator); err != nil {
		log.Fatal(err)
	}
}

func Reset() {
	//var HeaderedROM []byte := os.ReadFile()
	Emulator.cpu.ResetCPU()
	ppu.ResetPPU()
	Emulator.apu.ResetAPU()

	//Reset common variables
	common.Reset()

	//Reset ROM Data
	cartridge.ResetCartridge()

	cartridge.LoadCartridge()
	debug.ResetPatternTables()

	//copy(CHRData[:], HeaderedROM[0x8010:])

	//PCL := Emulator.bus.Read(0xFFFC)
	//PCH := Emulator.bus.Read(0xFFFD)
	//Emulator.cpu.PC = Emulator.cpu.BuildAddress(PCL, PCH)

	bus.OutsideCodeRead = 0
	bus.OutsideCodeWrite = 0
}

func (g *Game) Update() error {
	g.CheckCommonFunctions()
	if common.Filepath != "" && !common.ROMLoaded {
		Reset()
		if !common.ROMLoaded { //Invalid file
			common.ROMExists = false
			common.Filepath = ""
		}
	}
	if debug.LoggingCPU {
		debug.SetupTraceLogger()
	}

	for common.ROMLoaded && !PauseEmulation && !cpu.CPU_Halted {
		//Emulator.cpu.CPU_Cycle()
		MasterClockTick()
		if ppu.DrawNewFrame {
			ppu.DrawNewFrame = false
			break
			//return nil
		}
	}
	ppu.RenderFrame(Emulator.gameScreen)
	apu.TransferBuffer()

	if g.Exit {
		debug.ExportLog()
		return ebiten.Termination
	}
	// Update the UI
	debug.Update(&g.debugui)
	g.UI.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// This graphics context is used for managing the rendering state.

	if common.ROMLoaded && !cpu.CPU_Halted {
		screenRect := image.Rect(0, 0, common.ScreenWidth*common.ScreenScale, common.ScreenHeight*common.ScreenScale)
		screenScaled := image.NewRGBA(screenRect)
		draw.NearestNeighbor.Scale(screenScaled, screenScaled.Rect, g.gameScreen, g.gameScreen.Bounds(), draw.Over, nil)
		if screen.Bounds() == screenRect {
			screen.WritePixels(screenScaled.Pix)
		}
	}

	ShowMessages(screen)
	debug.DisplayDebugging(screen)

	//Only draw the UI if the mouse is close to it
	if _, my := input.CursorPosition(); my < common.MouseHeight || MenuBarSelected || !common.ROMExists {
		g.UI.Draw(screen)
	}

	//Draw debugging window, if enabled
	g.debugui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.ScreenWidth * common.ScreenScale, common.ScreenHeight * common.ScreenScale
}

func MasterClockTick() {
	//Run this everytime a CPU cycle is added, and run the PPU and APU accordingly
	//For when CPU instruction cycles are more accurately emulated

	//Clock the 2A03 and run CPU / APU
	//CPU runs every 6 ticks
	//PPU runs every 2 ticks
	//APU runs every 12
	//DMA runs every 6
	//

	MasterClock++
	//cpu.CPU_Cycles_New++
	switch MasterClock {
	case 1:
		Emulator.cpu.CPU_Cycle()
		ppu.PPU_Cycle()
		Emulator.apu.APU_Cycle()
	//case 2:
	case 3:
		ppu.PPU_Cycle()
	//case 4:
	case 5:
		ppu.PPU_Cycle()
	//case 6:
	case 7:
		Emulator.cpu.CPU_Cycle()
		ppu.PPU_Cycle()
		Emulator.apu.APU_Cycle()
	//case 8:
	case 9:
		ppu.PPU_Cycle()
	//case 10:
	case 11:
		ppu.PPU_Cycle()
	case 12:
		MasterClock = 0
	}
}

func (g *Game) CheckCommonFunctions() {
	common.ParseHotkeys()
	g.CheckForPause()
	g.CheckForReset()
	g.CheckForScreenshot()
	g.CheckForFullscreen()
	g.CheckForWindowResize()
}

func (g *Game) CheckForReset() {
	if common.PendingReset {
		common.ROMLoaded = false
		common.PendingReset = false
	}
}

func (g *Game) CheckForPause() {
	if common.PendingPause {
		PauseEmulation = !PauseEmulation
		common.PendingPause = false
	}
}

func (g *Game) CheckForScreenshot() {
	if common.PendingScreenshot {
		if !cpu.CPU_Halted && common.ROMLoaded {
			common.SaveScreenshot(g.gameScreen)
		}
		common.PendingScreenshot = false
	}
}

func (g *Game) CheckForFullscreen() {
	if common.PendingFullscreen {
		FullscreenMode = !FullscreenMode
		ebiten.SetFullscreen(FullscreenMode)
		common.PendingFullscreen = false
	}
}

func (g *Game) CheckForWindowResize() {
	if common.NewScreenScale != -1 {
		common.ScreenScale = common.NewScreenScale
		ebiten.SetWindowSize(common.ScreenWidth*common.ScreenScale, common.ScreenHeight*common.ScreenScale)
		common.NewScreenScale = -1
	}
}

func ShowMessages(screen *ebiten.Image) {
	if !common.ROMExists && common.UIMessageTimer == 0 {
		common.PrintUIMessage(screen, "No ROM file loaded")
	} /* else if cpu.CPU_Halted {
		common.PrintUIMessage(screen, "Game Crashed!")
	}*/

	if common.UIMessageTimer > 0 {
		//ebitenutil.DebugPrintAt(screen, common.UIMessage, 5, (240*common.ScreenScale)-20)
		common.PrintUIMessage(screen, common.UIMessage)
		common.UIMessageTimer--
		if common.UIMessageTimer == 0 {
			common.UIMessage = ""
		}
	}
}
