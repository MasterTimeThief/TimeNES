package main

import (
	"mtt/timenes/common"
	"mtt/timenes/nes"
	"mtt/timenes/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(common.ScreenWidth*common.ScreenScale, common.ScreenHeight*common.ScreenScale)
	ebiten.SetWindowTitle("TimeNES")
	//ebiten.SetTPS(ebiten.SyncWithFPS)

	ui := ui.InitUI()

	nes.InitGame(ui)

}
