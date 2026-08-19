package ui

import (
	"mtt/timenes/common"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type resources struct {
	font text.Face
}

func LoadResources() (*resources, error) {
	common.InitMenuFont()
	common.InitUIFont()
	return &resources{
		font: common.FontMenu,
	}, nil
}
