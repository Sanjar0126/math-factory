package entities

import (
	"fmt"
	"image/color"
	"math"

	"github.com/Sanjar0126/math-factory/internal/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type NumberType int

const (
	TypeBasic NumberType = iota
	TypePrime
	TypeComposite
)

type Number struct {
	X, Y      float64
	Value     int
	Type      NumberType
	VelocityX float64
	VelocityY float64
	Color     color.RGBA
	IsMoving  bool
	Size      float64
}

func NewNumber(x, y float64, value int) *Number {
	return &Number{
		X:        x,
		Y:        y,
		Value:    value,
		Type:     determineNumberType(value),
		Color:    getNumberColor(value),
		IsMoving: false,
		Size:     12,
	}
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

// Update just the Draw method signature:
func (n *Number) Draw(screen *ebiten.Image, camera CameraInterface) {
	screenX, screenY := camera.WorldToScreen(n.X, n.Y)
	zoom := camera.GetZoom()
	size := float32(n.Size * zoom)

	if size < 2 {
		return
	}

	// Draw circle
	vector.DrawFilledCircle(screen, float32(screenX), float32(screenY), size/2, n.Color, false)

	// Draw border
	borderColor := color.RGBA{255, 255, 255, 150}
	vector.StrokeCircle(screen, float32(screenX), float32(screenY), size/2, 1, borderColor, false)

	// Draw number if zoom is sufficient
	if zoom > 0.7 {
		opts := &text.DrawOptions{}
		opts.GeoM.Translate(screenX-6, screenY+4)
		opts.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, fmt.Sprintf("%d", n.Value), fonts.MplusNormalFont, opts)
	}
}

func (n *Number) MoveTo(targetX, targetY, speed float64) {
	dx := targetX - n.X
	dy := targetY - n.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance > 0 {
		n.VelocityX = (dx / distance) * speed
		n.VelocityY = (dy / distance) * speed
		n.IsMoving = true
	}
}

func IsPrime(n int) bool {
	return isPrime(n)
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
