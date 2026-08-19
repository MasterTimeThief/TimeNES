package common

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func BoolToInt(Flag bool) int {
	if Flag {
		return 1
	}
	return 0
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func Ternary(Condition bool, ValT, ValF uint16) uint16 {
	if Condition {
		return ValT
	}
	return ValF
}

func FirstN(str string, n int) string {
	v := []rune(str)
	if n >= len(v) {
		return str
	}
	return string(v[:n])
}

func Print(s string) func() {
	return func() {
		fmt.Println(s)
	}
}

func Combine2Bytes(Lo, Hi byte) uint16 {
	return uint16(Hi)<<8 | uint16(Lo)
}

func SetUIMessage(message string) {
	UIMessage = message
	UIMessageTimer = UIMessageTimerValue
}

func PrintUIMessage(screen *ebiten.Image, message string) {
	ebitenutil.DebugPrintAt(screen, message, 5, (240*ScreenScale)-20)
}
