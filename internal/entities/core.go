package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Core struct {
	BaseEntity
	StoredNumbers   []int
	CollectionRange float64
	Numbers         []*Number
}

func NewCore(x, y float64) *Core {
	return &Core{
		BaseEntity: BaseEntity{
			X:      x,
			Y:      y,
			Width:  64,
			Height: 64,
		},
		StoredNumbers:   make([]int, 0),
		CollectionRange: 80,
		Numbers:         make([]*Number, 0),
	}
}

func (c *Core) Update() {}

func (c *Core) Draw(screen *ebiten.Image, offsetX, offsetY float64, zoom float64) {
	screenX := float32((c.X + offsetX) * zoom)
	screenY := float32((c.Y + offsetY) * zoom)
	size := float32(c.Width * zoom)

	if size < 4 {
		return
	}

	centerX := screenX + size/2
	centerY := screenY + size/2
	radius := size / 2

	coreColor := color.RGBA{80, 120, 200, 255}
	vector.DrawFilledCircle(screen, centerX, centerY, radius, coreColor, false)

	innerColor := color.RGBA{120, 160, 255, 255}
	vector.DrawFilledCircle(screen, centerX, centerY, radius*0.6, innerColor, false)

	if zoom > 0.5 {
		rangeColor := color.RGBA{80, 120, 200, 50}
		rangeRadius := float32(c.CollectionRange * zoom)
		vector.StrokeCircle(screen, centerX, centerY, rangeRadius, 2, rangeColor, false)
	}

	pulse := float32(math.Sin(float64(ebiten.TPS())*0.1))*0.2 + 1.0
	energyColor := color.RGBA{200, 200, 255, 100}
	vector.StrokeCircle(screen, centerX, centerY, radius*pulse, 1, energyColor, false)
}

func (c *Core) CanCollect(number *Number) bool {
	dx := number.X - (c.X + c.Width/2)
	dy := number.Y - (c.Y + c.Height/2)
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance <= c.CollectionRange
}

func (c *Core) CollectNumber(number *Number) {
	c.StoredNumbers = append(c.StoredNumbers, number.Value)
	c.Numbers = append(c.Numbers, number)
}

func (c *Core) GetStoredCount() int {
	return len(c.StoredNumbers)
}
