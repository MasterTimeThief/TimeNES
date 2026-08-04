// toolbar.go
//
// Toolbar struct and related functions.
//

package main

import (
	"fmt"
	goimage "image"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
)

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
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.RGBA{R: 0, G: 100, B: 0, A: 255})),
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
		widget.TextOpts.Text("Modal Window", &res.font, color.RGBA{R: 52, G: 0, B: 78, A: 255}),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	))

	titleBar.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(color.Transparent),
			Hover:   image.NewNineSliceColor(colornames.Darkgray),
			Pressed: image.NewNineSliceColor(colornames.White),
		}),
		widget.ButtonOpts.TextPadding(&widget.Insets{Left: 16, Right: 64}),
		widget.ButtonOpts.Text("X", &res.font, &widget.ButtonTextColor{
			Idle:     color.White,
			Disabled: colornames.Gray,
			Hover:    color.White,
			Pressed:  color.Black,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
		}),
		widget.ButtonOpts.TabOrder(99),
	))

	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.RGBA{R: 0, G: 10, B: 10, A: 255})),
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
		widget.TextOpts.Text("This window blocks all input to widgets below it.", &res.font, color.RGBA{R: 234, G: 194, B: 12, A: 255}),
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

	// construct a new container that serves as the root of the UI hierarchy
	/*rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an grid layout to layout its ScrollableContainer and Slider
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Spacing(2, 0),
			//widget.GridLayoutOpts.Stretch([]bool{true, false}, []bool{true}),
		)),
	)

	// Create the container with the content that should be scrolled
	content := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewRowLayout(
		widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		widget.RowLayoutOpts.Spacing(20),
		widget.RowLayoutOpts.Padding(&widget.Insets{Top: 10, Bottom: 10}),
	)))

	// Add 20 buttons to the scrollable content container
	for x := 0; x < 20; x++ {
		// Capture x for use in callback
		x := x
		// construct a button
		button := widget.NewButton(
			// set general widget options
			widget.ButtonOpts.WidgetOpts(
				// instruct the container's anchor layout to center the button both horizontally and vertically
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
			),

			// specify the images to use
			widget.ButtonOpts.Image(&widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.Transparent),
				Hover:   image.NewNineSliceColor(colornames.Darkgray),
				Pressed: image.NewNineSliceColor(colornames.White),
			}),

			// specify the button's text, the font face, and the color
			widget.ButtonOpts.Text(fmt.Sprintf("Hello, World! - %d", x), &res.font, &widget.ButtonTextColor{
				Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
			}),

			// specify that the button's text needs some padding for correct display
			widget.ButtonOpts.TextPadding(&widget.Insets{
				Left:   30,
				Right:  30,
				Top:    5,
				Bottom: 5,
			}),

			// add a handler that reacts to clicking the button
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				println(fmt.Sprintf("Button %d Clicked!", x))
			}),
		)

		// add the button as a child of the container
		content.AddChild(button)
	}

	// Create the new ScrollContainer object
	scrollContainer := widget.NewScrollContainer(
		// Set the content that will be scrolled
		widget.ScrollContainerOpts.Content(content),
		// Tell the container to stretch the content width to match available space
		widget.ScrollContainerOpts.StretchContentWidth(),
		// Set the background images for the scrollable container
		widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle: image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff}),
			Mask: image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff}),
		}),
	)
	// Add the scrollable container to the left side of the window
	rootContainer.AddChild(scrollContainer)

	// Create a function to return the page size used by the slider
	pageSizeFunc := func() int {
		return int(math.Round(float64(scrollContainer.ViewRect().Dy())/float64(content.GetWidget().Rect.Dy())*1000) / 3)
	}
	// Create a vertical Slider bar to control the ScrollableContainer
	vSlider := widget.NewSlider(
		widget.SliderOpts.Direction(widget.DirectionVertical),
		widget.SliderOpts.MinMax(0, 1000),
		widget.SliderOpts.PageSizeFunc(pageSizeFunc),
		// On change update scroll location based on the Slider's value
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			scrollContainer.ScrollTop = float64(args.Slider.Current) / 1000
		}),
		widget.SliderOpts.Images(
			// Set the track images
			&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			// Set the handle images
			&widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Hover:   image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
			},
		),
	)
	// Set the slider's position if the scrollContainer is scrolled by other means than the slider
	scrollContainer.GetWidget().ScrolledEvent.AddHandler(func(args interface{}) {
		if a, ok := args.(*widget.WidgetScrolledEventArgs); ok {
			vSlider.Current -= int(math.Round(a.Y * float64(pageSizeFunc())))
		}
	})

	// Add the slider to the second slot in the root container
	rootContainer.AddChild(vSlider)

	c.AddChild(rootContainer)
	*/

	cb := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(color.Transparent),
			Hover:   image.NewNineSliceColor(colornames.Darkgray),
			Pressed: image.NewNineSliceColor(colornames.White),
		}),
		widget.ButtonOpts.TextPadding(&widget.Insets{Left: 16, Right: 64}),
		widget.ButtonOpts.Text("Close", &res.font, &widget.ButtonTextColor{
			Idle:     color.White,
			Disabled: colornames.Gray,
			Hover:    color.White,
			Pressed:  color.Black,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
		}),
	)
	bc.AddChild(cb)

	window = widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),
		widget.WindowOpts.TitleBar(titleBar, 30),
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.Resizeable(),
		widget.WindowOpts.MinSize(500, 200),
		widget.WindowOpts.MaxSize(700, 400),
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Resize: ", args.Rect)
		}),
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Move: ", args.Rect)
		}),
	)

	promptWidth := 200
	promptHeight := 250

	windowSize := input.GetWindowSize()
	r := goimage.Rect(0, 0, promptWidth, promptHeight)
	r = r.Add(goimage.Point{(screenWidth - promptWidth) / 2, windowSize.Y * 2 / 3 / 2})
	window.SetLocation(r)

	rw = ui.AddWindow(window)
}
