package common

import (
	"embed"
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/sqweek/dialog"
	"golang.org/x/image/draw"
)

var Filepath string = ""
var ROMExists bool = false
var ROMLoaded bool = false

//go:embed assets
var embeddedAssets embed.FS
var FontUI, FontMenu text.Face

func SelectROM(fp string) {
	Filepath = fp
	ROMExists = true
	ROMLoaded = false
}

func FileSelectDialog() {
	file, err := dialog.File().Title("Select NES ROM File").Filter("NES ROM files (*.nes)", "nes").Load()
	if err == nil {
		SelectROM(file)
	}
}

func SaveScreenshot(screen image.Image) {
	filename := "screenshot_" + time.Now().Format("20060102030405") + ".png"

	f, err := os.Create("screenshots/" + filename)
	if err != nil {
		log.Fatal(err)
	}

	imageScaled := image.NewRGBA(image.Rect(0, 0, ScreenWidth*ScreenScale, ScreenHeight*ScreenScale))
	draw.NearestNeighbor.Scale(imageScaled, imageScaled.Rect, screen, screen.Bounds(), draw.Over, nil)

	if err := png.Encode(f, imageScaled); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	SetUIMessage("Saved " + filename)
}

func InitMenuFont() {
	fnt, err := LoadFont("assets/fonts/PixelOperator-Bold.ttf", FontSize)
	Check(err)
	FontMenu = fnt
}

func InitUIFont() {
	fnt, err := LoadFont("assets/fonts/BetterVCR_25_09.ttf", FontSize)
	Check(err)
	FontUI = fnt
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
