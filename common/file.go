package common

import (
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/sqweek/dialog"
	"golang.org/x/image/draw"
)

var Filepath string = ""
var ROMExists bool = false
var ROMLoaded bool = false

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

	f, err := os.Create(filename)
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

	SetUIMessage("Saved screenshot " + filename)
}
