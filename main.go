package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"golang.org/x/image/draw"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//var filepath string = "roms/games/smb.nes"

//var filepath string = "roms/nes-test-roms-master/apu_test/rom_singles/4-jitter.nes"

//var filepath string = "roms/nes-test-roms-master/blargg_apu_2005.07.30/04.clock_jitter.nes"

//var filepath string = "roms/AccuracyCoin.nes"

var filepath string = ""
var ROMExists bool = false
var ROMLoaded bool = false
var ShowFPS bool = true

//type App struct{ Clicks int }

//var SelectedROM string = ""

const (
	screenWidth  = 256
	screenHeight = 240
	screenScale  = 2
)

type Game struct {
	gameScreen *image.RGBA
	keys       []ebiten.Key
	ui         *ebitenui.UI
	exit       bool
}

var RAM [0x800]byte
var VRAM [0x800]byte
var PaletteRAM [0x20]byte
var ROM [0x8000]byte

// var CHRData [0x2000]byte
var CHRROM [0x2000]byte
var Header [0x10]byte
var CartRAM [0x2000]byte

// var cpuClock int = 0 // 2A03
// var ppuClock int = 0 // 2C02
//var MasterClock int = 0

//var screenScale float32 = 3
//var nametable *image.RGBA = image.NewRGBA(image.Rect(0, 0, ScreenWidth, ScreenHeight))

func (g *Game) Update() error {

	//For when CPU instruction cycles are more accurately emulated

	//Clock the 2A03 and run CPU / APU
	//CPU runs every 6 ticks
	//PPU runs every 2 ticks
	//APU runs every 12
	//DMA runs every 6
	//

	/*MasterClock++
	switch MasterClock {
	case 1:
		Emulate_CPU(g)
		Emulate_PPU(g)
		Emulate_APU(g)
		DMA_Get()
	case 2:
	case 3:
		Emulate_PPU(g)
	case 4:
	case 5:
		Emulate_PPU(g)
	case 6:
	case 7:
		Emulate_CPU(g)
		Emulate_PPU(g)
		DMA_Put()
	case 8:
	case 9:
		Emulate_PPU(g)
	case 10:
	case 11:
		Emulate_PPU(g)
	case 12:
		MasterClock = 0

	}*/

	/*
		if cpuClock >= 12 {
			cpuClock = 0
			Emulate_CPU(g)
			Emulate_APU(g)
		}
		if cpuClock == 6 {
			Emulate_APU(g)
		}

		//Clock the 2C02 and run PPU
		ppuClock++
		if ppuClock >= 4 {
			ppuClock = 0
			Emulate_PPU(g)
		}
	*/

	if filepath != "" && !ROMLoaded {
		Reset()
		if !ROMLoaded { //Invalid file
			ROMExists = false
			filepath = ""
		}
	}

	for ROMLoaded {
		UpdateControllers(g)
		Emulate_CPU(g)
		if DrawNewFrame {
			DrawNewFrame = false
			break
			//return nil
		}
		if CPU_Halted || g.exit {
			return ebiten.Termination
		}
	}
	// Update the UI
	g.ui.Update()

	//fmt.Println("CPU Paused, Drawing Window")

	//DrawWindow(ops, window)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// This graphics context is used for managing the rendering state.

	if !ROMExists {
		ebitenutil.DebugPrintAt(screen, "No ROM File loaded!", 5, 460) // Draw the UI onto the screen
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
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f\tFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()), 380, 460) // Draw the UI onto the screen
	}
	g.ui.Draw(screen)
	//ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*screenScale, screenHeight*screenScale)
	ebiten.SetWindowTitle("TimeNES (MasterTimeThief made this :3)")
	//ebiten.SetTPS(ebiten.SyncWithFPS)

	res, err := loadResources()
	if err != nil {
		log.Fatal(err)
	}

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

	g := &Game{
		gameScreen: image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
		ui:         &ui,
	}

	// Event handling
	//
	// Example 1: Configure the "Help" button to display a message in console when it's pressed.
	//
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
		SelectROM("roms/nestest.nes")
	}))

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ternary(Condition bool, ValT, ValF uint16) uint16 {
	if Condition {
		return ValT
	}
	return ValF
}

func Reset() {
	//var HeaderedROM []byte := os.ReadFile()
	flag_InterruptDisable = true

	//Reset APU variables
	apuDMAGetCycle = true
	// All the bits of $4015 are cleared
	apuDMCInterrupt = false
	apuFrameInterrupt = false
	apuDMCDelayed = false
	apuEnableDMC = false
	apuEnableNoise = false
	apuEnableTriangle = false
	apuEnablePulse2 = false
	apuEnablePulse1 = false
	apuDMC.Length = 0
	apuNoise.LengthCounter = 0
	apuTriangle.LengthCounter = 0
	apuPulse2.LengthCounter = 0
	apuPulse1.LengthCounter = 0
	apuFrameCounter = 0 // reset the frame counter

	// PPU registers
	ppuEnableNMI = false
	ppuVRAMInc32Mode = false
	ppuUse8x16Sprites = false
	ppuSpritePatternTable = false
	ppuBGPatternTable = false
	TransferAddress = 0

	//PPU_Mask_Greyscale = false
	ppuMask_EmphasisRed = false
	ppuMask_EmphasisGreen = false
	ppuMask_EmphasisBlue = false
	ppuMask_8pxMaskBG = false
	ppuMask_8pxMaskSprites = false
	ppuMask_RenderBG = false
	ppuMask_RenderSprites = false

	//PPU_Update2005Delay = 0;
	ppuScrollFineX = 0

	PPUReadBuffer = 0
	//PPU_OddFrame = false;

	ppuDot = 0
	ppuScanline = 0

	//DoDMCDMA = false;
	//DoOAMDMA = false;

	//Reset CPU cycle count
	CPU_Cycles = 0

	//Reset ROM Data
	Header = [0x10]byte{}
	ROM = [0x8000]byte{}
	CHRROM = [0x2000]byte{}

	//Reset RAM (Or don't, if I wanna do weird stuff...?)
	RAM = [0x800]byte{}
	VRAM = [0x800]byte{}
	PaletteRAM = [0x20]byte{}
	CartRAM = [0x2000]byte{}

	LoadROM()

	//copy(CHRData[:], HeaderedROM[0x8010:])

	PCL := Read(0xFFFC)
	PCH := Read(0xFFFD)
	ProgramCounter = BuildAddress(PCL, PCH)
	//ProgramCounter = 0xC000
	StackPointer -= 3
	//fmt.Printf("%#x", ProgramCounter)
	ROMLoaded = true
}

func BoolToInt(Flag bool) int {
	if Flag {
		return 1
	}
	return 0
}

func LoadROM() {
	HeaderedROM, err := os.ReadFile(filepath)
	check(err)

	copy(Header[:], HeaderedROM[0x0:])
	size := uint16(Header[4])
	copy(ROM[:], HeaderedROM[0x10:])
	if Header[5] != 0 {
		copy(CHRROM[:], HeaderedROM[(0x0010+(0x4000*size)):])
	}

}

func SelectROM(fp string) {
	filepath = fp
	ROMExists = true

	ROMLoaded = false

}
