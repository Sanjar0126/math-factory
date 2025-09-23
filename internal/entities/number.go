package entities

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
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

func (n *Number) Update() {
	if n.IsMoving {
		n.X += n.VelocityX
		n.Y += n.VelocityY

		n.VelocityX *= 0.95
		n.VelocityY *= 0.95

		if math.Abs(n.VelocityX) < 0.01 && math.Abs(n.VelocityY) < 0.01 {
			n.VelocityX = 0
			n.VelocityY = 0
			n.IsMoving = false
		}
	}
}

func (n *Number) Draw(screen *ebiten.Image, offsetX, offsetY float64, zoom float64) {
	screenX := float32((n.X + offsetX) * zoom)
	screenY := float32((n.Y + offsetY) * zoom)
	size := float32(n.Width * zoom)

	if size < 2 {
		return
	}

	vector.DrawFilledCircle(screen, screenX+size/2, screenY+size/2, size/2, n.Color, false)

	borderColor := color.RGBA{255, 255, 255, 100}
	vector.StrokeCircle(screen, screenX+size/2, screenY+size/2, size/2, 1, borderColor, false)

	if zoom > 0.8 {
		text.Draw(screen, fmt.Sprintf("%d", n.Value), basicfont.Face7x13,
			int(screenX+2), int(screenY+12), color.White)
	}
}

func (n *Number) MoveTo(targetX, targetY float64, speed float64) {
	dx := targetX - n.X
	dy := targetY - n.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance > 0 {
		n.VelocityX = (dx / distance) * speed
		n.VelocityY = (dy / distance) * speed
		n.IsMoving = true
	}
}

func determineNumberType(value int) NumberType {
	if value < 2 {
		return TypeBasic
	}

	if isPrime(value) {
		return TypePrime
	}

	return TypeComposite
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func getNumberColor(value int) color.RGBA {
	switch determineNumberType(value) {
	case TypePrime:
		return color.RGBA{100, 255, 100, 255}
	case TypeComposite:
		return color.RGBA{255, 150, 100, 255}
	default:
		return color.RGBA{150, 150, 255, 255}
	}
}
