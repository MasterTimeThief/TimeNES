package common

import (
	"embed"
	"image"
	"image/png"
	"io/fs"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
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
var icons []image.Image

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

	if err := os.MkdirAll("screenshots", 0777); err != nil {
		log.Fatal(err)
	}
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
	fnt, err := LoadFont("assets/fonts/PixelOperatorMono-Bold.ttf", FontSize)
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
	defer func(f fs.File) {
		_ = f.Close()
	}(fontFile)

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

func LoadIcon(path string) image.Image {
	iconFile, err := embeddedAssets.Open(path)
	if err != nil {
		return nil
	}
	defer func(f fs.File) {
		_ = f.Close()
	}(iconFile)

	icon, err := png.Decode(iconFile)
	if err != nil {
		return nil
	}

	return icon
}

func SetWindowIcon() {
	//var iconSlice []image.Image
	icons = append(icons, LoadIcon("assets/icons/16.png"))
	icons = append(icons, LoadIcon("assets/icons/32.png"))
	icons = append(icons, LoadIcon("assets/icons/48.png"))
	icons = append(icons, LoadIcon("assets/icons/64.png"))
	icons = append(icons, LoadIcon("assets/icons/128.png"))
	icons = append(icons, LoadIcon("assets/icons/256.png"))
	ebiten.SetWindowIcon(icons)
}
