package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var filepath string = "roms/smb.nes"

//var filepath string = "roms/nestest.nes"

//var filepath string = "roms/AccuracyCoin.nes"

//type App struct{ Clicks int }

//var SelectedROM string = ""

const (
	screenWidth  = 256
	screenHeight = 240
	screenScale  = 3
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
var PRGROMSize uint16 = 0x4000

// var CHRData [0x2000]byte
var CHRROM [0x2000]byte
var Header [0x10]byte
var CartRAM [0x2000]byte

//var screenScale float32 = 3
//var nametable *image.RGBA = image.NewRGBA(image.Rect(0, 0, ScreenWidth, ScreenHeight))

func (g *Game) Update() error {

	for /*SelectedROM != ""*/ {
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

	//if SelectedROM != "" {
	screen.WritePixels(g.gameScreen.Pix)
	//DrawNewFrame = false
	//}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f\tFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()), 0, 225) // Draw the UI onto the screen
	g.ui.Draw(screen)
	//ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
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
	/*toolbar.smbButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		SelectedROM = "roms/smb.nes"
		//Reset the console
		Reset()
	}))
	toolbar.nestestButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		SelectedROM = "roms/nestest.nes"
		//Reset the console
		Reset()
	}))*/

	Reset()

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
	HeaderedROM, err := os.ReadFile(filepath)
	check(err)

	copy(Header[:], HeaderedROM[0x0:])
	size := uint16(Header[4])
	copy(ROM[:], HeaderedROM[0x10:])
	if Header[5] != 0 {
		copy(CHRROM[:], HeaderedROM[(0x0010+(0x4000*size)):])
	}
	//copy(CHRData[:], HeaderedROM[0x8010:])

	PCL := Read(0xFFFC)
	PCH := Read(0xFFFD)
	ProgramCounter = BuildAddress(PCL, PCH)
	//ProgramCounter = 0xC000
	StackPointer -= 3
	//fmt.Printf("%#x", ProgramCounter)
	//Run()
}

func BoolToInt(Flag bool) int {
	if Flag {
		return 1
	}
	return 0
}
