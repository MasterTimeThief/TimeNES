package main

import (
	"embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed assets
var embeddedAssets embed.FS

type resources struct {
	font text.Face
}

func loadResources() (*resources, error) {
	fnt, err := loadFont("assets/fonts/MinecraftStandard.otf", 6)
	if err != nil {
		return nil, err
	}
	return &resources{
		font: fnt,
	}, nil
}

func loadFont(path string, size float64) (text.Face, error) {
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
