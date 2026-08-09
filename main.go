package main

import (
	"fmt"
	"image"
	"log"
	"mtt/timenes/mappers"
	"os"

	"golang.org/x/image/draw"

	"github.com/ebitengine/debugui"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//var filepath string = "roms/games/smb.nes"

//var filepath string = "roms/nes-test-roms-master/apu_test/rom_singles/4-jitter.nes"

//var filepath string = "roms/nes-test-roms-master/blargg_apu_2005.07.30/04.clock_jitter.nes"

//var filepath string = "roms/AccuracyCoin.nes"

var filepath string = ""
var ROMExists bool = false
var ROMLoaded bool = false
var ShowFPS bool = false
var pauseEmulation bool = false

// Header Variables

var PRGROM_Size uint32 // Size of PRG ROM
var CHRROM_Size uint32 // Size of CHR ROM (value 0 means the board uses CHR RAM)
var IsNametableHorizontal bool
var HasBatteryRAM bool
var AltNametableLayout bool
var NES2_Header bool //Is the header in NES 2.0 format, rather than iNES

var MapperChipID byte

//type App struct{ Clicks int }

//var SelectedROM string = ""

const (
	screenWidth  = 256
	screenHeight = 240
	screenScale  = 2
)

type Game struct {
	gameScreen   *image.RGBA
	keys         []ebiten.Key
	ui           *ebitenui.UI
	exit         bool
	debugui      debugui.DebugUI
	directory    []os.DirEntry
	audioContext *audio.Context
	player       *audio.Player
}

var g *Game

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

//var screenScale float32 = 3
//var nametable *image.RGBA = image.NewRGBA(image.Rect(0, 0, ScreenWidth, ScreenHeight))

func main() {
	InitAPU()
	//InitAudioOutput()

	if len(os.Args) > 1 && os.Args[1] != "" {
		filepath = os.Args[1]
	}
	if filepath != "" {
		ROMExists = true
	}
	ebiten.SetWindowSize(screenWidth*screenScale, screenHeight*screenScale)
	ebiten.SetWindowTitle("TimeNES (MasterTimeThief made this :3)")
	//ebiten.SetTPS(ebiten.SyncWithFPS)

	res, err := loadResources()
	check(err)
	//soundLoop()
	// Construct a new container that serves as the root of the UI hierarchy.
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	ui := ebitenui.UI{
		Container: root,
	}

	// Create a toolbar and add it to the UI.
	toolbar := newToolbar(&ui, res)
	root.AddChild(toolbar.container)

	g = &Game{
		gameScreen: image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
		ui:         &ui,
	}

	SetupToolbarOptions(res, g, toolbar)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}

func Reset() {
	//var HeaderedROM []byte := os.ReadFile()
	ResetCPU()
	ResetPPU()
	ResetAPU()

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

	outsideCodeRead = 0
	outsideCodeWrite = 0
	ROMLoaded = true
}

func (g *Game) Update() error {

	if filepath != "" && !ROMLoaded {
		Reset()
		if !ROMLoaded { //Invalid file
			ROMExists = false
			filepath = ""
		}
	}

	for ROMLoaded && !pauseEmulation {
		UpdateControllers(g)
		Emulate_CPU(g)
		if DrawNewFrame {
			DrawNewFrame = false
			break
			//return nil
		}
	}
	if /*CPU_Halted ||*/ g.exit {
		return ebiten.Termination
	}
	// Update the UI

	/*mX, _ := ebiten.CursorPosition()
	if mX > 50 {
		g.ui.Container.SetLocation(image.Rect(-20, 0, 20, 20))
	} else {
		g.ui.Container.SetLocation(image.Rect(0, 0, 20, 20))
	}*/
	//DebugWindow(g)

	g.ui.Update()

	//fmt.Println("CPU Paused, Drawing Window")

	//DrawWindow(ops, window)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// This graphics context is used for managing the rendering state.

	if !ROMExists {
		ebitenutil.DebugPrintAt(screen, "No ROM File loaded!", 5, (240*screenScale)-20)
	} else if CPU_Halted {
		ebitenutil.DebugPrintAt(screen, "Game Crashed!", 5, (240*screenScale)-20)
	} else if ROMLoaded {
		//if SelectedROM != "" {
		//if DrawNewFrame {
		screenScaled := image.NewRGBA(image.Rect(0, 0, screenWidth*screenScale, screenHeight*screenScale))
		draw.NearestNeighbor.Scale(screenScaled, screenScaled.Rect, g.gameScreen, g.gameScreen.Bounds(), draw.Over, nil)
		screen.WritePixels(screenScaled.Pix)
		//DrawNewFrame = false
		//}
		//}
	}

	if ShowFPS {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f\tFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()), (256*screenScale)-132, (240*screenScale)-20) // Draw the UI onto the screen
	}
	if outsideCodeRead > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Attempting to Read at $: %04X", outsideCodeRead), 5, (240*screenScale)-30)
	}
	if outsideCodeWrite > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Attempting to Write at $: %04X", outsideCodeWrite), 5, (240*screenScale)-40)
	}

	g.ui.Draw(screen)
	//g.debugui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func LoadROM() {
	HeaderedROM, err := os.ReadFile(filepath)
	check(err)

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

func SelectROM(fp string) {
	filepath = fp
	ROMExists = true
	ROMLoaded = false
}

func SetupToolbarOptions(res *resources, g *Game, toolbar *toolbar) {
	// Event handling
	toolbar.helpButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		println("The help button was pressed!")
	}))

	// Example 2: Configure the "Quit" menu entry to end the program when it's pressed.
	toolbar.quitButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		g.exit = true
	}))

	//Select ROM
	toolbar.smbButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		SelectROM("roms/games/smb.nes")
	}))

	toolbar.nestestButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		SelectROM("roms/AccuracyCoin.nes")
		//SelectROM("roms/nestest.nes")
	}))

	toolbar.selectROMButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		//Add file select
		openFileSelectWindow(res, g.ui)
	}))

	toolbar.FPSButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		toggleFPS()
	}))

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
