package ui

import (
	"embed"
	"log"
	"mtt/timenes/common"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed assets
var embeddedAssets embed.FS

type resources struct {
	font text.Face
}

const (
	uiBGColor = "000000"
)

func LoadResources() (*resources, error) {
	fnt, err := LoadFont("assets/fonts/ClearSans-Regular.ttf", 16)
	if err != nil {
		return nil, err
	}
	return &resources{
		font: fnt,
	}, nil
}

func LoadFont(path string, size float64) (text.Face, error) {
	fontFile, err := embeddedAssets.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := text.NewGoTextFaceSource(fontFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}

func InitUI() *ebitenui.UI {
	res, err := LoadResources()
	common.Check(err)

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
