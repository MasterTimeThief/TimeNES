// toolbar.go
//
// Toolbar struct and related functions.
//

package ui

import (
	"fmt"
	goimage "image"
	"image/color"
	"mtt/timenes/common"
	"mtt/timenes/nes"
	"os"
	"strings"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/utilities/constantutil"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
)

var FC_Red color.RGBA = color.RGBA{R: 160, G: 30, B: 37, A: 255}
var FC_Gold color.RGBA = color.RGBA{R: 185, G: 159, B: 119, A: 255}
var FC_White color.RGBA = color.RGBA{R: 232, G: 228, B: 224, A: 255}
var FC_Black color.RGBA = color.RGBA{R: 10, G: 10, B: 10, A: 255}

var currDir string = "./roms"
var fullDirectory []os.DirEntry

// NOTE: It's not strictly necessary to store references to all the buttons in the toolbar struct, but this example does
// so for completeness' sake. When you keep a reference to buttons in the struct, you can later configure them to respond
// to certain events in your application, and keep your program's logic outside the toolbar.
type toolbar struct {
	container       *widget.Container
	fileMenu        *widget.Button
	editMenu        *widget.Button
	helpButton      *widget.Button
	quitButton      *widget.Button
	smbButton       *widget.Button
	nestestButton   *widget.Button
	selectROMButton *widget.Button
	FPSButton       *widget.Button
}

type ListEntry struct {
	id    int
	isDir bool
	Name  string
}

func newToolbar(ui *ebitenui.UI, res *resources) *toolbar {
	// Create a root container for the toolbar.
	root := widget.NewContainer(
		// Use black background for the toolbar.
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.RGBA{R: 0, G: 0, B: 0, A: 100})),

		// Toolbar components must be aligned horizontally.
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			),
		),

		widget.ContainerOpts.WidgetOpts(
			// Make the toolbar fill the whole horizontal space of the screen.
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{StretchHorizontal: true}),
		),
	)

	//
	// "File" menu
	//
	file := newToolbarButton(res, "File")
	var (
		selectROM = newToolbarMenuEntry(res, "Select ROM File")
		smb       = newToolbarMenuEntry(res, "Super Mario Bros.")
		nestest   = newToolbarMenuEntry(res, "AccuracyCoin")
		quit      = newToolbarMenuEntry(res, "Quit")
	)

	// Make the toolbar entry open a menu with our "save" and "load" entries  when the user clicks it.
	file.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		openToolbarMenu(args.Button.GetWidget(), ui, selectROM, smb, nestest, quit)
	}))
	root.AddChild(file)

	//
	// "Debug" menu
	//
	edit := newToolbarButton(res, "Debug")
	var (
		fps = newToolbarMenuEntry(res, "Toggle FPS Display")
	)
	edit.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		openToolbarMenu(args.Button.GetWidget(), ui, fps)
	}))
	root.AddChild(edit)

	//
	// "Help" button
	// Unlike the "File" and "Edit" menu, this is just a regular button on the toolbar - it does not open a menu.
	// You can configure it to do something else when it's pressed, like opening a "Help" window.
	//
	help := newToolbarButton(res, "Help")
	root.AddChild(help)

	return &toolbar{
		container:       root,
		fileMenu:        file,
		editMenu:        edit,
		helpButton:      help,
		quitButton:      quit,
		smbButton:       smb,
		nestestButton:   nestest,
		selectROMButton: selectROM,
		FPSButton:       fps,
	}
}

func newToolbarButton(res *resources, label string) *widget.Button {
	// Create a button for the toolbar.
	return widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(color.Transparent),
			Hover:   image.NewNineSliceColor(colornames.Darkgray),
			Pressed: image.NewNineSliceColor(colornames.White),
		}),
		widget.ButtonOpts.Text(label, &res.font, &widget.ButtonTextColor{
			Idle:     color.White,
			Disabled: colornames.Gray,
			Hover:    color.White,
			Pressed:  color.Black,
		}),
		widget.ButtonOpts.TextPadding(&widget.Insets{
			Top:    2,
			Left:   5,
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
			Hover:   image.NewNineSliceColor(colornames.Darkgray),
			Pressed: image.NewNineSliceColor(colornames.White),
		}),
		widget.ButtonOpts.Text(label, &res.font, &widget.ButtonTextColor{
			Idle:     color.White,
			Disabled: colornames.Gray,
			Hover:    color.White,
			Pressed:  color.Black,
		}),
		widget.ButtonOpts.TextPosition(widget.TextPositionStart, widget.TextPositionCenter),
		widget.ButtonOpts.TextPadding(&widget.Insets{Left: 5, Right: 10}),
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
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(64, 0)),
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
	)

	// Immediately add the menu to the UI.
	ui.AddWindow(window)
}

func openFileSelectWindow(res *resources, ui *ebitenui.UI) {
	var rw widget.RemoveWindowFunc
	var window *widget.Window

	titleBar := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(FC_Red)),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Stretch([]bool{true, false, false}, []bool{true}),
			widget.GridLayoutOpts.Padding(&widget.Insets{
				Left:   30,
				Right:  5,
				Top:    6,
				Bottom: 5,
			}))))

	titleBar.AddChild(widget.NewText(
		widget.TextOpts.Text("Select ROM File", &res.font, FC_White),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	))

	titleBar.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(color.Transparent),
			Hover:   image.NewNineSliceColor(colornames.Gray),
			Pressed: image.NewNineSliceColor(colornames.Black),
		}),
		widget.ButtonOpts.TextPadding(&widget.Insets{Left: 8, Right: 8, Top: 8, Bottom: 8}),
		widget.ButtonOpts.Text("X", &res.font, &widget.ButtonTextColor{
			Idle:     color.Black,
			Disabled: colornames.Darkgray,
			Hover:    color.Black,
			Pressed:  color.White,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
		}),
		widget.ButtonOpts.TabOrder(99),
	))

	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(FC_Gold)),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
				widget.GridLayoutOpts.Padding(&widget.Insets{
					Left:   13,
					Right:  13,
					Top:    7,
					Bottom: 7,
				}),
				widget.GridLayoutOpts.Spacing(0, 15),
			),
		),
	)

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("Select a folder / file, and click \"Open\".", &res.font, FC_Black),
	))

	/*tOpts := []widget.TextInputOpt{
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     image.NewNineSliceColor(color.RGBA{R: 0, G: 0, B: 100, A: 255}),
			Disabled: image.NewNineSliceColor(color.RGBA{R: 100, G: 0, B: 0, A: 255}),
		}),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:     color.White,
			Disabled: colornames.Gray,
		}),
		widget.TextInputOpts.Padding(&widget.Insets{
			Left:   13,
			Right:  13,
			Top:    7,
			Bottom: 7,
		}),
		widget.TextInputOpts.Face(&res.font),
		widget.TextInputOpts.CaretWidth(2),
	}

	t := widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.TextInputOpts.Placeholder("Enter text here"))...,
	)
	textContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout()))
	textContainer.AddChild(t)
	c.AddChild(textContainer)*/

	fileList := widget.NewList(
		// Set how wide the list should be
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(150, 0),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchVertical:    true,
				Padding:            widget.NewInsetsSimple(50),
			}),
		)),
		// Set the entries in the list
		// widget.ListOpts.Entries(entries),
		widget.ListOpts.ScrollContainerImage(&widget.ScrollContainerImage{
			Idle:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			Disabled: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			Mask:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
		}),

		widget.ListOpts.SliderParams(&widget.SliderParams{
			// Set the background images/color for the background of the slider track
			TrackImage: &widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			HandleImage: &widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255}),
				Hover:   image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{R: 255, G: 100, B: 120, A: 255}),
			},
			MinHandleSize: constantutil.ConstantToPointer(5),
			// Set how wide the track should be
			TrackPadding: widget.NewInsetsSimple(2),
		}),
		// Hide the horizontal slider
		widget.ListOpts.HideHorizontalSlider(),
		// Set the font for the list options
		widget.ListOpts.EntryFontFace(&res.font),
		// Set the colors for the list
		widget.ListOpts.EntryColor(&widget.ListEntryColor{
			Selected:                   color.NRGBA{R: 0, G: 255, B: 0, A: 255},     // Foreground color for the unfocused selected entry
			Unselected:                 color.NRGBA{R: 254, G: 255, B: 255, A: 255}, // Foreground color for the unfocused unselected entry
			SelectedBackground:         color.NRGBA{R: 130, G: 130, B: 200, A: 255}, // Background color for the unfocused selected entry
			SelectingBackground:        color.NRGBA{R: 130, G: 130, B: 130, A: 255}, // Background color for the unfocused being selected entry
			SelectingFocusedBackground: color.NRGBA{R: 130, G: 140, B: 170, A: 255}, // Background color for the focused being selected entry
			SelectedFocusedBackground:  color.NRGBA{R: 130, G: 130, B: 170, A: 255}, // Background color for the focused selected entry
			FocusedBackground:          color.NRGBA{R: 170, G: 170, B: 180, A: 255}, // Background color for the focused unselected entry
			DisabledUnselected:         color.NRGBA{R: 100, G: 100, B: 100, A: 255}, // Foreground color for the disabled unselected entry
			DisabledSelected:           color.NRGBA{R: 100, G: 100, B: 100, A: 255}, // Foreground color for the disabled selected entry
			DisabledSelectedBackground: color.NRGBA{R: 100, G: 100, B: 100, A: 255}, // Background color for the disabled selected entry
		}),
		// This required function returns the string displayed in the list.
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(ListEntry).Name
		}),
		// Padding for each entry
		widget.ListOpts.EntryTextPadding(widget.NewInsetsSimple(5)),
		// Text position for each entry
		widget.ListOpts.EntryTextPosition(widget.TextPositionStart, widget.TextPositionCenter),
		// This handler defines what function to run when a list item is selected.
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			entry := args.Entry.(ListEntry)
			/*if entry.isDir {
				currDir += entry.Name
				fileList
			} else {
				filepath = currDir + "/" + entry.Name
				ROMExists = true
				ROMLoaded = false
				rw()
			}*/
			fmt.Println("Entry Selected: ", entry)
		}),
		// This option will select the entry as it is focused
		// widget.ListOpts.SelectFocus(),
		// This option will select the entry when pressing the mouse button instead of releasing it
		// widget.ListOpts.SelectPressed(),

		// This option will disable default keys (up and down)
		// widget.ListOpts.DisableDefaultKeys(true),
	)

	fileList.SetEntries(RefreshDirectoryList())

	c.AddChild(fileList)

	bc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(15),
		)),
	)
	c.AddChild(bc)

	/*o2b := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(color.Transparent),
			Hover:   image.NewNineSliceColor(colornames.Darkgray),
			Pressed: image.NewNineSliceColor(colornames.White),
		}),
		widget.ButtonOpts.TextPadding(&widget.Insets{Left: 16, Right: 64}),
		widget.ButtonOpts.Text("Open Another", &res.font, &widget.ButtonTextColor{
			Idle:     color.White,
			Disabled: colornames.Gray,
			Hover:    color.White,
			Pressed:  color.Black,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			//openWindow2(res, ui)
		}),
	)
	bc.AddChild(o2b)*/

	openButton := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(FC_White),
			Hover:   image.NewNineSliceColor(colornames.Darkgray),
			Pressed: image.NewNineSliceColor(FC_Black),
		}),
		widget.ButtonOpts.TextPadding(&widget.Insets{Left: 32, Right: 32, Top: 5, Bottom: 5}),
		widget.ButtonOpts.Text("Open", &res.font, &widget.ButtonTextColor{
			Idle:     FC_Black,
			Disabled: colornames.Gray,
			Hover:    FC_Black,
			Pressed:  FC_White,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			entry := fileList.SelectedEntry()
			if entry != nil {
				if entry.(ListEntry).isDir {
					currDir += entry.(ListEntry).Name
					fileList.SetEntries(RefreshDirectoryList())
				} else {
					common.Filepath = currDir + "/" + entry.(ListEntry).Name
					common.ROMExists = true
					common.ROMLoaded = false
					rw()
				}
			}
		}),
	)
	bc.AddChild(openButton)

	closeButton := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(FC_White),
			Hover:   image.NewNineSliceColor(colornames.Darkgray),
			Pressed: image.NewNineSliceColor(FC_Black),
		}),
		widget.ButtonOpts.TextPadding(&widget.Insets{Left: 32, Right: 32, Top: 5, Bottom: 5}),
		widget.ButtonOpts.Text("Close", &res.font, &widget.ButtonTextColor{
			Idle:     FC_Black,
			Disabled: colornames.Gray,
			Hover:    FC_Black,
			Pressed:  FC_White,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
		}),
	)
	bc.AddChild(closeButton)

	window = widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),
		widget.WindowOpts.TitleBar(titleBar, 30),
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.Resizeable(),
		widget.WindowOpts.MinSize(100, 100),
		widget.WindowOpts.MaxSize(500, 450),
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Resize: ", args.Rect)
		}),
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Move: ", args.Rect)
		}),
	)

	promptWidth := 400
	promptHeight := 300

	r := goimage.Rect(0, 0, promptWidth, promptHeight)
	r = r.Add(goimage.Point{((common.ScreenWidth * common.ScreenScale) - promptWidth) / 2, ((common.ScreenHeight * common.ScreenScale) - promptHeight) / 2})
	window.SetLocation(r)

	rw = ui.AddWindow(window)
}

func RefreshDirectoryList() []any {
	var list []any

	//Get Directory
	GetCurrentDirectory(currDir)

	list = append(list, ListEntry{0, true, "/.."})

	for i := range fullDirectory {
		split := strings.Split(fullDirectory[i].Name(), ".")

		//isGitFolder := fullDirectory[i].IsDir() && firstN(fullDirectory[i].Name(), 1) == "."
		isNESFile := (!fullDirectory[i].IsDir() && split[len(split)-1] == "nes")

		if common.FirstN(fullDirectory[i].Name(), 1) != "." && (fullDirectory[i].IsDir() || isNESFile) {
			itemText := fullDirectory[i].Name()
			if fullDirectory[i].IsDir() {
				itemText = "/" + itemText
			}
			list = append(list, ListEntry{i, fullDirectory[i].IsDir(), itemText})
		}
	}

	return list
}

func GetCurrentDirectory(dir string) {
	dirEntries, err := os.ReadDir(dir)
	common.Check(err)

	fullDirectory = dirEntries
}

func SetupToolbarOptions(res *resources, toolbar *toolbar) {
	// Event handling
	toolbar.helpButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		println("The help button was pressed!")
	}))

	// Example 2: Configure the "Quit" menu entry to end the program when it's pressed.
	toolbar.quitButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		nes.Emulator.Exit = true
	}))

	//Select ROM
	toolbar.smbButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		common.SelectROM("roms/smb.nes")
	}))

	toolbar.nestestButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		common.SelectROM("roms/tests/AccuracyCoin.nes")
		//SelectROM("roms/tests/nestest.nes")
	}))

	toolbar.selectROMButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		//Add file select
		openFileSelectWindow(res, nes.Emulator.UI)
	}))

	toolbar.FPSButton.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		nes.ShowFPS = !nes.ShowFPS
	}))

}
