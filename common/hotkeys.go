package common

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func ParseHotkeys() {

	var keys []ebiten.Key
	keys = inpututil.AppendJustPressedKeys(keys[:0])

	for _, k := range keys {
		switch k {
		case ebiten.KeyF1: // Pause
			PendingPause = true
		case ebiten.KeyF2: // Debug window
			PendingDebug = true
		case ebiten.KeyF3: // Reset
			PendingReset = true
		case ebiten.KeyF4:
		case ebiten.KeyF5:
		case ebiten.KeyF6:
			PendingPT = true
		case ebiten.KeyF7:
		case ebiten.KeyF8:
		case ebiten.KeyF9: // Toggle Audio
			PendingMute = true
		case ebiten.KeyF10: // Show FPS
			PendingFPS = true
		case ebiten.KeyF11: // Fullscreen
			PendingFullscreen = true
		case ebiten.KeyF12: //Screenshot
			PendingScreenshot = true
		case ebiten.KeyMinus:
			if ScreenScale > 1 {
				NewScreenScale = ScreenScale - 1
			}
		case ebiten.KeyEqual:
			if ScreenScale < 4 {
				NewScreenScale = ScreenScale + 1
			}

		}
	}
}
