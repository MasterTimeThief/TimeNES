package nes

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"mtt/timenes/common"
	"mtt/timenes/nes/apu"
	"mtt/timenes/nes/bus"
	"mtt/timenes/nes/cartridge"
	"mtt/timenes/nes/cpu"
	"mtt/timenes/nes/ppu"
	"os"

	"github.com/ebitengine/debugui"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/draw"
)

// Header Variables

// var cpuClock int = 0 // 2A03
// var ppuClock int = 0 // 2C02
var MasterClock int = 0

// Debugging
var ShowFPS bool = false
var pauseEmulation bool = false

type Game struct {
	gameScreen *image.RGBA
	//keys       []ebiten.Key
	UI      *ebitenui.UI
	Exit    bool
	debugui debugui.DebugUI
	//directory    []os.DirEntry
	audioContext *audio.Context
	player       *audio.Player
}

var Emulator *Game

func (g *Game) Update() error {

	if common.Filepath != "" && !common.ROMLoaded {
		Reset()
		if !common.ROMLoaded { //Invalid file
			common.ROMExists = false
			common.Filepath = ""
		}
	}

	for common.ROMLoaded && !pauseEmulation {
		cpu.Emulate_CPU()
		if ppu.DrawNewFrame {
			ppu.DrawNewFrame = false
			break
			//return nil
		}
	}
	RenderFrame()
	apu.TransferBuffer()

	if /*CPU_Halted ||*/ g.Exit {
		return ebiten.Termination
	}
	// Update the UI
	g.UI.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// This graphics context is used for managing the rendering state.

	if !common.ROMExists {
		ebitenutil.DebugPrintAt(screen, "No ROM File loaded!", 5, (240*common.ScreenScale)-20)
	} else if cpu.CPU_Halted {
		ebitenutil.DebugPrintAt(screen, "Game Crashed!", 5, (240*common.ScreenScale)-20)
	} else if common.ROMLoaded {
		screenScaled := image.NewRGBA(image.Rect(0, 0, common.ScreenWidth*common.ScreenScale, common.ScreenHeight*common.ScreenScale))
		draw.NearestNeighbor.Scale(screenScaled, screenScaled.Rect, g.gameScreen, g.gameScreen.Bounds(), draw.Over, nil)
		screen.WritePixels(screenScaled.Pix)
	}

	if ShowFPS {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f\tFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()), (256*common.ScreenScale)-132, (240*common.ScreenScale)-20) // Draw the UI onto the screen
	}
	if bus.OutsideCodeRead > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Attempting to Read at $: %04X", bus.OutsideCodeRead), 5, (240*common.ScreenScale)-30)
	}
	if bus.OutsideCodeWrite > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Attempting to Write at $: %04X", bus.OutsideCodeWrite), 5, (240*common.ScreenScale)-40)
	}

	g.UI.Draw(screen)
	//g.debugui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func InitGame(newUI *ebitenui.UI) {
	if len(os.Args) > 1 && os.Args[1] != "" {
		common.Filepath = os.Args[1]
	}
	if common.Filepath != "" {
		common.ROMExists = true
	}
	apu.InitAPU()

	ebiten.SetWindowSize(common.ScreenWidth*common.ScreenScale, common.ScreenHeight*common.ScreenScale)
	ebiten.SetWindowTitle("TimeNES")
	//ebiten.SetTPS(ebiten.SyncWithFPS)

	Emulator = &Game{
		gameScreen: image.NewRGBA(image.Rect(0, 0, common.ScreenWidth, common.ScreenHeight)),
		UI:         newUI,
	}

	if Emulator.audioContext == nil {
		Emulator.audioContext = apu.NewAudioContext()
	}
	if Emulator.player == nil {
		Emulator.player = apu.NewAudioPlayer(Emulator.audioContext)
	}

	if err := ebiten.RunGame(Emulator); err != nil {
		log.Fatal(err)
	}
}

func Reset() {
	//var HeaderedROM []byte := os.ReadFile()
	cpu.ResetCPU()
	ppu.ResetPPU()
	apu.ResetAPU()

	//Reset ROM Data
	cartridge.ResetCartridge()

	cartridge.LoadCartridge()

	//copy(CHRData[:], HeaderedROM[0x8010:])

	PCL := bus.Read(0xFFFC)
	PCH := bus.Read(0xFFFD)
	cpu.PC = cpu.BuildAddress(PCL, PCH)
	//fmt.Printf("%#x", ProgramCounter)

	bus.OutsideCodeRead = 0
	bus.OutsideCodeWrite = 0
}

func RenderPixel(color color.RGBA) {
	pixIndex := uint64((((ppu.PPUScanline) * common.ScreenWidth) + (ppu.PPUDot - 1)) * 4)

	Emulator.gameScreen.Pix[pixIndex] = color.R
	Emulator.gameScreen.Pix[pixIndex+1] = color.G
	Emulator.gameScreen.Pix[pixIndex+2] = color.B
	Emulator.gameScreen.Pix[pixIndex+3] = color.A
}

func RenderFrame() {
	for i, color := range ppu.FrameColorBuffer {
		Emulator.gameScreen.Pix[(i * 4)] = color.R
		Emulator.gameScreen.Pix[(i*4)+1] = color.G
		Emulator.gameScreen.Pix[(i*4)+2] = color.B
		//Emulator.gameScreen.Pix[(i*4)+3] = color.A
	}
	ppu.FrameColorBufferPos = 0
}

func MasterClockTick(location string) {
	//Run this everytime a CPU cycle is added, and run the PPU and APU accordingly
	//For when CPU instruction cycles are more accurately emulated

	//Clock the 2A03 and run CPU / APU
	//CPU runs every 6 ticks
	//PPU runs every 2 ticks
	//APU runs every 12
	//DMA runs every 6
	//

	MasterClock++
	cpu.CPU_Cycles_New++
	/*switch MasterClock {
	case 1:
		Emulate_CPU(g)
		Emulate_PPU(g)
		Emulate_APU(g)
		DMA_Get()
	//case 2:
	case 3:
		Emulate_PPU(g)
	//case 4:
	case 5:
		Emulate_PPU(g)
	//case 6:
	case 7:
		Emulate_CPU(g)
		Emulate_PPU(g)
		DMA_Put()
	//case 8:
	case 9:
		Emulate_PPU(g)
	//case 10:
	case 11:
		Emulate_PPU(g)
	case 12:
		MasterClock = 0

	}*/
	//cycleTest += fmt.Sprint("Cycle: " + location + "\n")
}
