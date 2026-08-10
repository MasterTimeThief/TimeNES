package common

import (
	"image"

	"github.com/ebitengine/debugui"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type Game struct {
	gameScreen *image.RGBA
	keys       []ebiten.Key
	ui         *ebitenui.UI
	exit       bool
	debugui    debugui.DebugUI
	//directory    []os.DirEntry
	audioContext *audio.Context
	player       *audio.Player
}

var G *Game
