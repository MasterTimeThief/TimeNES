package main

import (
	"flag"
	"log"
	"mtt/timenes/common"
	"mtt/timenes/nes"
	"mtt/timenes/ui"

	"github.com/hajimehoshi/ebiten/v2"

	_ "net/http/pprof"
)

func main() {

	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()
	log.SetFlags(0)
	flag.IntVar(&common.ScreenScale, "scale", 2, "Screen scale")
	flag.Parse()

	ebiten.SetWindowSize(common.ScreenWidth*common.ScreenScale, common.ScreenHeight*common.ScreenScale)
	ebiten.SetWindowTitle("TimeNES")
	//ebiten.SetTPS(ebiten.SyncWithFPS)

	ui := ui.InitUI()

	nes.InitGame(ui)

}
