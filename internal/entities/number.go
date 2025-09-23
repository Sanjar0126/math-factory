package entities

import (
	"image/color"
)

type NumberType int

const (
	TypeBasic NumberType = iota
	TypePrime
	TypeComposite
)

type Number struct {
	BaseEntity
	Value     int
	Type      NumberType
	VelocityX float64
	VelocityY float64
	Color     color.RGBA
	IsMoving  bool
}

func NewNumber(x, y float64, value int) *Number {
	num := &Number{
		BaseEntity: BaseEntity{
			X:      x,
			Y:      y,
			Width:  16,
			Height: 16,
		},
		Value:    value,
		Type:     determineNumberType(value),
		Color:    getNumberColor(value),
		IsMoving: false,
	}
	return num
}
