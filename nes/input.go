package nes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var Controller1, Controller2 byte
var Controller1ShiftRegister, Controller2ShiftRegister uint16

//var enableZapper bool = false

// Update states of controllers
func UpdateControllers(g *Game) {

	//	0 - B
	//	1 - A
	//	2 - Select
	//	3 - Start
	//	4 - Up
	//	5 - Down
	//	6 - Left
	//	7 - Right

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	Controller1 = 0
	Controller2 = 0
	for _, k := range g.keys {
		switch k {
		//Player 1
		case ebiten.KeyX:
			Controller1 |= 0x80
		case ebiten.KeyZ:
			Controller1 |= 0x40
		case ebiten.KeyShiftRight:
			Controller1 |= 0x20
		case ebiten.KeyEnter:
			//if !enableZapper {
			Controller1 |= 0x10
			//}
		case ebiten.KeyArrowUp:
			//if !enableZapper {
			Controller1 |= 0x8
			//}
		case ebiten.KeyArrowDown:
			Controller1 |= 0x4
		case ebiten.KeyArrowLeft:
			Controller1 |= 0x2
		case ebiten.KeyArrowRight:
			Controller1 |= 0x1

		//Player 2
		case ebiten.KeyK:
			Controller2 |= 0x80
		case ebiten.KeyL:
			Controller2 |= 0x40
		case ebiten.KeyI:
			Controller2 |= 0x20
		case ebiten.KeyO:
			Controller2 |= 0x10
		case ebiten.KeyW:
			Controller2 |= 0x8
		case ebiten.KeyS:
			Controller2 |= 0x4
		case ebiten.KeyA:
			Controller2 |= 0x2
		case ebiten.KeyD:
			Controller2 |= 0x1
		}
	}

	//if enableZapper {
	//	Controller1 |= byte(ternary(inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0), 0x0, 0x10))
	//	Controller1 |= 0x8
	//}

}
