package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Entity interface {
	Update()
	Draw(screen *ebiten.Image, offsetX, offsetY float64, zoom float64)
	GetPosition() (float64, float64)
	GetSize() (float64, float64)
}

type BaseEntity struct {
	X, Y          float64
	Width, Height float64
}

