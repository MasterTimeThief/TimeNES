// toolbar.go
//
// Toolbar struct and related functions.
//

package main

import (
	goimage "image"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
)

// NOTE: It's not strictly necessary to store references to all the buttons in the toolbar struct, but this example does
// so for completeness' sake. When you keep a reference to buttons in the struct, you can later configure them to respond
// to certain events in your application, and keep your program's logic outside the toolbar.
type toolbar struct {
	container     *widget.Container
	fileMenu      *widget.Button
	editMenu      *widget.Button
	helpButton    *widget.Button
	quitButton    *widget.Button
	smbButton     *widget.Button
	nestestButton *widget.Button
	toBeDecided   *widget.Button
	toBeDecided2  *widget.Button
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
		smb     = newToolbarMenuEntry(res, "Super Mario Bros.")
		nestest = newToolbarMenuEntry(res, "nestest")
		quit    = newToolbarMenuEntry(res, "Quit")
	)

	// Make the toolbar entry open a menu with our "save" and "load" entries  when the user clicks it.
	file.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		openToolbarMenu(args.Button.GetWidget(), ui, smb, nestest, quit)
	}))
	root.AddChild(file)

	//
	// "Edit" menu
	// This is the same thing as the "File" menu, just with more entries.
	//
	edit := newToolbarButton(res, "Edit")
	var (
		TBD  = newToolbarMenuEntry(res, "TBD")
		TBD2 = newToolbarMenuEntry(res, "TBD2")
	)
	edit.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		openToolbarMenu(args.Button.GetWidget(), ui, TBD, TBD2)
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
		container:     root,
		fileMenu:      file,
		editMenu:      edit,
		helpButton:    help,
		quitButton:    quit,
		smbButton:     smb,
		nestestButton: nestest,
		toBeDecided:   TBD,
		toBeDecided2:  TBD2,
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
			Top:    0,
			Left:   3,
			Right:  5,
			Bottom: 0,
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
		widget.ButtonOpts.TextPadding(&widget.Insets{Left: 16, Right: 64}),
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
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.RGBA{R: 0, G: 100, B: 0, A: 200})),

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
