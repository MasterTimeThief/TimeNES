package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var Controller1, Controller2 byte
var Controller1ShiftRegister, Controller2ShiftRegister uint16

// Update states of controllers
func UpdateControllers(g *Game) {
	/*
		0 - B
		1 - A
		2 - Select
		3 - Start
		4 - Up
		5 - Down
		6 - Left
		7 - Right
	*/
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	Controller1 = 0
	Controller2 = 0
	for _, k := range g.keys {
		switch k {
		case ebiten.KeyX:
			Controller1 |= 0x80
		case ebiten.KeyZ:
			Controller1 |= 0x40
		case ebiten.KeyShiftRight:
			Controller1 |= 0x20
		case ebiten.KeyEnter:
			Controller1 |= 0x10
		case ebiten.KeyArrowUp:
			Controller1 |= 0x8
		case ebiten.KeyArrowDown:
			Controller1 |= 0x4
		case ebiten.KeyArrowLeft:
			Controller1 |= 0x2
		case ebiten.KeyArrowRight:
			Controller1 |= 0x1
		}
	}

	/*Controller2 |= byte(ternary(Controller2State[0], 0x80, 0x00))
	Controller2 |= byte(ternary(Controller2State[1], 0x40, 0x00))
	Controller2 |= byte(ternary(Controller2State[2], 0x20, 0x00))
	Controller2 |= byte(ternary(Controller2State[3], 0x10, 0x00))
	Controller2 |= byte(ternary(Controller2State[4], 0x08, 0x00))
	Controller2 |= byte(ternary(Controller2State[5], 0x04, 0x00))
	Controller2 |= byte(ternary(Controller2State[6], 0x02, 0x00))
	Controller2 |= byte(ternary(Controller2State[7], 0x01, 0x00))*/
}
