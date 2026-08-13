// toolbar.go
//
// Toolbar struct and related functions.
//

package ui

import (
	goimage "image"
	"image/color"
	"mtt/timenes/common"
	"mtt/timenes/debug"
	"mtt/timenes/nes"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
)

//var FC_Red color.RGBA = color.RGBA{R: 160, G: 30, B: 37, A: 255}
//var FC_Gold color.RGBA = color.RGBA{R: 185, G: 159, B: 119, A: 255}
//var FC_White color.RGBA = color.RGBA{R: 232, G: 228, B: 224, A: 255}
//var FC_Black color.RGBA = color.RGBA{R: 10, G: 10, B: 10, A: 255}

// NOTE: It's not strictly necessary to store references to all the buttons in the toolbar struct, but this example does
// so for completeness' sake. When you keep a reference to buttons in the struct, you can later configure them to respond
// to certain events in your application, and keep your program's logic outside the toolbar.
type toolbar struct {
	container *widget.Container
	// File
	fileMenu   *widget.Button
	openButton *widget.Button
	quitButton *widget.Button
	// Game
	gameMenu    *widget.Button
	pauseButton *widget.Button
	muteButton  *widget.Button
	resetButton *widget.Button
	// Debug
	debugMenu  *widget.Button
	smbButton  *widget.Button
	testButton *widget.Button
	FPSButton  *widget.Button
}

type ListEntry struct {
	id    int
	isDir bool
	Name  string
}

func InitUI() *ebitenui.UI {
	res, err := LoadResources()
	common.Check(err)

	//Setup hitbox for menubar

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

	SetupToolbarOptions(res, toolbar)

	return &ui
}

func newToolbar(ui *ebitenui.UI, res *resources) *toolbar {
	// Create a root container for the toolbar.
	root := widget.NewContainer(
		// Use black background for the toolbar.
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.RGBA{R: 0, G: 0, B: 0, A: 200})),

		// Toolbar components must be aligned horizontally.
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			),
		),

		widget.ContainerOpts.WidgetOpts(
			// Make the toolbar fill the whole horizontal space of the screen.
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{StretchHorizontal: true}),
			widget.WidgetOpts.MinSize(0, common.MenuBarHeight),
			//widget.WidgetOpts.TrackHover(true),
		),
	)

	//
	// "File" menu
	//
	file := newToolbarButton(res, "File")
	var (
		open = newToolbarMenuEntry(res, "Open ROM")
		quit = newToolbarMenuEntry(res, "Quit")
	)

	// Make the toolbar entry open a menu with our "save" and "load" entries  when the user clicks it.
	file.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		nes.MenuBarSelected = true
		openToolbarMenu(args.Button.GetWidget(), ui, open, quit)
	}))
	root.AddChild(file)

	//
	// "Game" menu
	//
	game := newToolbarButton(res, "Game")
	var (
		mute  = newToolbarMenuEntry(res, "Toggle Audio")
		pause = newToolbarMenuEntry(res, "Pause")
		reset = newToolbarMenuEntry(res, "Reset")
	)
	game.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		nes.MenuBarSelected = true
		openToolbarMenu(args.Button.GetWidget(), ui, mute, pause, reset)
	}))
	root.AddChild(game)

	//
	// "Debug" menu
	//
	debug := newToolbarButton(res, "Debug")
	var (
		smb  = newToolbarMenuEntry(res, "Super Mario Bros.")
		test = newToolbarMenuEntry(res, "AccuracyCoin")
		fps  = newToolbarMenuEntry(res, "Toggle FPS Display")
	)
	debug.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		nes.MenuBarSelected = true
		openToolbarMenu(args.Button.GetWidget(), ui, smb, test, fps)
	}))
	root.AddChild(debug)

	return &toolbar{
		container: root,
		// File
		fileMenu:   file,
		openButton: open,
		quitButton: quit,
		// Game
		gameMenu:    game,
		muteButton:  mute,
		pauseButton: pause,
		resetButton: reset,
		// Debug
		debugMenu:  debug,
		smbButton:  smb,
		testButton: test,
		FPSButton:  fps,
	}
}

func newToolbarButton(res *resources, label string) *widget.Button {
	// Create a button for the toolbar.
	return widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(color.Transparent),
			Hover:   image.NewNineSliceColor(colornames.Dimgrey),
			Pressed: image.NewNineSliceColor(colornames.Darkgrey),
		}),
		widget.ButtonOpts.Text(label, &res.font, &widget.ButtonTextColor{
			Idle:     color.White,
			Disabled: colornames.Gray,
			Hover:    color.White,
			Pressed:  color.Black,
		}),
		widget.ButtonOpts.TextPadding(&widget.Insets{
			Top:    2,
			Left:   10,
			Right:  10,
			Bottom: 2,
		}),
	)
}

func newToolbarMenuEntry(res *resources, label string) *widget.Button {
	// Create a button for a menu entry.
	return widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(color.Transparent),
			Hover:   image.NewNineSliceColor(colornames.Dimgrey),
			Pressed: image.NewNineSliceColor(colornames.Darkgrey),
		}),
		widget.ButtonOpts.Text(label, &res.font, &widget.ButtonTextColor{
			Idle:     color.White,
			Disabled: colornames.Gray,
			Hover:    color.White,
			Pressed:  color.Black,
		}),
		widget.ButtonOpts.TextPosition(widget.TextPositionStart, widget.TextPositionCenter),
		widget.ButtonOpts.TextPadding(&widget.Insets{Left: 20, Right: 10, Top: 2, Bottom: 2}),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
		),
	)
}

func openToolbarMenu(opener *widget.Widget, ui *ebitenui.UI, entries ...*widget.Button) {
	c := widget.NewContainer(
		// Set the background to a translucent black.
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.RGBA{R: 0, G: 0, B: 0, A: 200})),

		// Menu entries should be arranged vertically.
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(4),
				widget.RowLayoutOpts.Padding(&widget.Insets{Top: 1, Bottom: 1}),
			),
		),

		// Set the minimum size for the menu.
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(128, 0)),
	)

	for _, entry := range entries {
		c.AddChild(entry)
	}

	w, h := c.PreferredSize()

	window := widget.NewWindow(
		// Set the menu to be a modal. This makes it block UI interactions to anything ese.
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),

		// Close the menu if the user clicks outside of it.
		widget.WindowOpts.CloseMode(widget.CLICK),

		// Position the menu below the menu button that it belongs to.
		widget.WindowOpts.Location(
			goimage.Rect(
				opener.Rect.Min.X,
				opener.Rect.Min.Y+opener.Rect.Max.Y,
				opener.Rect.Min.X+w,
				opener.Rect.Min.Y+opener.Rect.Max.Y+opener.Rect.Min.Y+h,
			),
		),

		widget.WindowOpts.ClosedHandler(func(args *widget.WindowClosedEventArgs) {
			nes.MenuBarSelected = false
		}),
	)

	// Immediately add the menu to the UI.
	ui.AddWindow(window)
}

func SetupToolbarOptions(res *resources, toolbar *toolbar) {
	//
	// File
	//

	// Open ROM
	toolbar.openButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		//Add file select
		//openFileSelectWindow(res, nes.Emulator.UI)
		common.FileSelectDialog()
	}))

	// Quit Emulator
	toolbar.quitButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		nes.Emulator.Exit = true
	}))

	//
	// Game
	//

	// Mute Emulator
	toolbar.muteButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		common.MuteEmulator = !common.MuteEmulator
	}))

	// Pause Emulator
	toolbar.pauseButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		nes.PauseEmulation = !nes.PauseEmulation
	}))

	// Reset Emulator
	toolbar.resetButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		common.ROMLoaded = false
	}))

	//
	// Debug
	//

	// Load Super Mario Bros.
	toolbar.smbButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		common.SelectROM("roms/smb.nes")
	}))

	// Load AccuracyCoin
	toolbar.testButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		common.SelectROM("roms/tests/AccuracyCoin.nes")
		//SelectROM("roms/tests/nestest.nes")
	}))

	toolbar.FPSButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		debug.ShowFPS = !debug.ShowFPS
	}))

}
