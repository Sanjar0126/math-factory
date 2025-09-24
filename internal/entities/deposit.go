package entities

import (
	"fmt"
	"image/color"

	"github.com/Sanjar0126/math-factory/internal/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type NumberDeposit struct {
	Position     GridPosition
	NumberValue  int
	IsInfinite   bool
	RemainingOre int
	DepositType  NumberType
	IsMined      bool // Has a miner placed on it
}

func NewNumberDeposit(gridX, gridY int, value int, infinite bool) *NumberDeposit {
	remaining := 1000
	if infinite {
		remaining = -1 // -1 means infinite
	}

	return &NumberDeposit{
		Position:     GridPosition{X: gridX, Y: gridY},
		NumberValue:  value,
		IsInfinite:   infinite,
		RemainingOre: remaining,
		DepositType:  determineNumberType(value),
		IsMined:      false,
	}
}

func (d *NumberDeposit) Update() {}

func (d *NumberDeposit) Draw(screen *ebiten.Image, camera Camera) {
	worldX, worldY := d.Position.ToWorldPos()
	screenX, screenY := camera.WorldToScreen(worldX, worldY)
	zoom := camera.GetZoom()
	size := float32(TileSize) * float32(zoom)

	if size < 4 {
		return
	}

	var bgColor, borderColor color.RGBA
	if d.IsMined {
		bgColor = color.RGBA{100, 100, 100, 200} // Darker when mined
		borderColor = color.RGBA{150, 150, 150, 255}
	} else {
		switch d.DepositType {
		case TypePrime:
			bgColor = color.RGBA{50, 150, 50, 255} // Dark green
			borderColor = color.RGBA{100, 255, 100, 255}
		case TypeComposite:
			bgColor = color.RGBA{150, 80, 50, 255} // Dark orange
			borderColor = color.RGBA{255, 150, 100, 255}
		default:
			bgColor = color.RGBA{50, 50, 150, 255} // Dark blue
			borderColor = color.RGBA{150, 150, 255, 255}
		}
	}

	vector.DrawFilledRect(screen, float32(screenX), float32(screenY),
		size, size, bgColor, false)

	vector.StrokeRect(screen, float32(screenX), float32(screenY),
		size, size, 2, borderColor, false)

	if zoom > 0.6 {
		textColor := color.RGBA{255, 255, 255, 1}
		if d.IsMined {
			textColor = color.RGBA{200, 200, 200, 255}
		}

		opts := &text.DrawOptions{}
		opts.GeoM.Translate(float64(screenX+4), float64(screenY+20))
		opts.ColorScale.ScaleWithColor(textColor)
		text.Draw(screen, fmt.Sprintf("%d", d.NumberValue), fonts.MplusNormalFont, opts)
	}

	if d.IsInfinite && zoom > 0.8 {
		opts := &text.DrawOptions{}
		opts.GeoM.Translate(float64(float32(screenX)+size-12), float64(screenY+12))
		opts.ColorScale.ScaleWithColor(color.RGBA{255, 255, 0, 255})
		text.Draw(screen, "∞", fonts.MplusNormalFont, opts)
	}
}

func (d *NumberDeposit) GetGridPosition() GridPosition {
	return d.Position
}

func (d *NumberDeposit) GetSize() (int, int) {
	return 1, 1
}

func (d *NumberDeposit) CanBeMined() bool {
	return !d.IsMined && (d.IsInfinite || d.RemainingOre > 0)
}

func (d *NumberDeposit) SetMined(mined bool) {
	d.IsMined = mined
}

func (d *NumberDeposit) Mine() (int, bool) {
	if !d.CanBeMined() {
		return 0, false
	}

	if !d.IsInfinite {
		d.RemainingOre--
	}

	return d.NumberValue, true
}
