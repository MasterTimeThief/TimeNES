package ui

import (
	"embed"
	"log"
	"mtt/timenes/common"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed assets
var embeddedAssets embed.FS

type resources struct {
	font text.Face
}

func LoadResources() (*resources, error) {
	fnt, err := LoadFont("assets/fonts/ClearSans-Regular.ttf", common.FontSize)
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
