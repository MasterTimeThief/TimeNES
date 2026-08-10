package nes

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"mtt/timenes/common"
	"mtt/timenes/mappers"
	"mtt/timenes/nes/apu"
	"os"

	"github.com/ebitengine/debugui"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/draw"
)

// Header Variables

var PRGROM_Size uint32 // Size of PRG ROM
var CHRROM_Size uint32 // Size of CHR ROM (value 0 means the board uses CHR RAM)
var IsNametableHorizontal bool
var HasBatteryRAM bool
var AltNametableLayout bool
var NES2_Header bool //Is the header in NES 2.0 format, rather than iNES

var MapperChipID byte

var RAM [0x800]byte
var VRAM [0x800]byte
var CartVRAM [0x1000]byte //For mapper chips
var PaletteRAM [0x20]byte
var ROM [0x80000]byte

var CHRROM [0x80000]byte
var isCHRRAM bool
var Header [0x10]byte
var CartRAM [0x2000]byte

// var cpuClock int = 0 // 2A03
// var ppuClock int = 0 // 2C02
var MasterClock int = 0

// Debugging
var OutsideCodeRead, OutsideCodeWrite uint16 = 0, 0
var ShowFPS bool = false
var pauseEmulation bool = false

type Game struct {
	gameScreen *image.RGBA
	keys       []ebiten.Key
	UI         *ebitenui.UI
	Exit       bool
	debugui    debugui.DebugUI
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
		UpdateControllers(g)
		Emulate_CPU(g)
		if DrawNewFrame {
			DrawNewFrame = false
			break
			//return nil
		}
	}

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
	} else if CPU_Halted {
		ebitenutil.DebugPrintAt(screen, "Game Crashed!", 5, (240*common.ScreenScale)-20)
	} else if common.ROMLoaded {
		screenScaled := image.NewRGBA(image.Rect(0, 0, common.ScreenWidth*common.ScreenScale, common.ScreenHeight*common.ScreenScale))
		draw.NearestNeighbor.Scale(screenScaled, screenScaled.Rect, g.gameScreen, g.gameScreen.Bounds(), draw.Over, nil)
		screen.WritePixels(screenScaled.Pix)
	}

	if ShowFPS {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f\tFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()), (256*common.ScreenScale)-132, (240*common.ScreenScale)-20) // Draw the UI onto the screen
	}
	if OutsideCodeRead > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Attempting to Read at $: %04X", OutsideCodeRead), 5, (240*common.ScreenScale)-30)
	}
	if OutsideCodeWrite > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Attempting to Write at $: %04X", OutsideCodeWrite), 5, (240*common.ScreenScale)-40)
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
	apu.InitAudioOutput()

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
	ResetCPU()
	ResetPPU()
	apu.ResetAPU()

	//Reset ROM Data
	Header = [0x10]byte{}
	ROM = [0x80000]byte{}
	CHRROM = [0x80000]byte{}
	isCHRRAM = false

	//Reset RAM (Or don't, if I wanna do weird stuff...?)
	RAM = [0x800]byte{}
	VRAM = [0x800]byte{}
	PaletteRAM = [0x20]byte{}
	CartRAM = [0x2000]byte{}

	MapperChipID = 0

	PRGROM_Size = 0
	CHRROM_Size = 0
	IsNametableHorizontal = false
	HasBatteryRAM = false
	AltNametableLayout = false

	LoadROM()

	//copy(CHRData[:], HeaderedROM[0x8010:])

	PCL := Read(0xFFFC)
	PCH := Read(0xFFFD)
	ProgramCounter = BuildAddress(PCL, PCH)
	//fmt.Printf("%#x", ProgramCounter)

	OutsideCodeRead = 0
	OutsideCodeWrite = 0
	common.ROMLoaded = true
}

func LoadROM() {
	HeaderedROM, err := os.ReadFile(common.Filepath)
	common.Check(err)

	//Header info
	copy(Header[:], HeaderedROM[0x0:])
	PRGROM_Size = uint32(Header[4]) * uint32(0x4000)
	CHRROM_Size = uint32(Header[5]) * uint32(0x2000)
	IsNametableHorizontal = (Header[6] & 1) == 0
	HasBatteryRAM = (Header[6] & 0x02) != 0
	AltNametableLayout = (Header[6] & 0x08) != 0

	NES2_Header = ((Header[7] & 0xC) >> 2) == 2

	MapperChipID = (Header[6] >> 4) | (Header[7] & 0xF0)

	//size := uint16(Header[4])
	ROM_Endpoint := uint32(0x10 + (PRGROM_Size))
	CHR_Endpoint := uint32(ROM_Endpoint + uint32(CHRROM_Size))

	copy(ROM[:], HeaderedROM[0x10:ROM_Endpoint])
	if CHRROM_Size != 0 {
		copy(CHRROM[:], HeaderedROM[ROM_Endpoint:CHR_Endpoint])
	} else {

	}

	//Initialize any PRG-RAM from mapper chips
	if HasBatteryRAM {
		switch MapperChipID {
		case 1: //MMC1
			copy(mappers.MMC1_PRGRAM[:], HeaderedROM[0x10:])
		case 2: //UxROM
		case 3: //CNROM
			//Add support for Hayauchi Super Igo?
		case 4: //MMC3
		}
	}

}

func RenderPixel(color color.RGBA) {
	pixIndex := uint64((((ppuScanline) * common.ScreenWidth) + (ppuDot - 1)) * 4)

	//pixIndex &= 0x3BFFF
	Emulator.gameScreen.Pix[pixIndex] = color.R
	Emulator.gameScreen.Pix[pixIndex+1] = color.G
	Emulator.gameScreen.Pix[pixIndex+2] = color.B
	Emulator.gameScreen.Pix[pixIndex+3] = color.A
	/*pixR, pixG, pixB := DecodeNTSC(screenWidth, int(math.Mod(float64(ppuDot*8)+3.9, 12.0)))
	g.gameScreen.Pix[pixIndex] = pixR
	g.gameScreen.Pix[pixIndex+1] = pixG
	g.gameScreen.Pix[pixIndex+2] = pixB
	g.gameScreen.Pix[pixIndex+3] = 0xFF*/
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
	CPU_Cycles_New++
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
