package entities

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Miner struct {
	Position       GridPosition
	Deposit        *NumberDeposit
	MiningTimer    int
	MiningInterval int
	OutputDir      Direction
	OutputBuffer   []*Number
	MaxBuffer      int
}

type Direction int

const (
	DirectionUp Direction = iota
	DirectionRight
	DirectionDown
	DirectionLeft
)

func NewMiner(gridX, gridY int, deposit *NumberDeposit, outputDir Direction) *Miner {
	return &Miner{
		Position:       GridPosition{X: gridX, Y: gridY},
		Deposit:        deposit,
		MiningTimer:    0,
		MiningInterval: 120, // 2 seconds at 60 FPS
		OutputDir:      outputDir,
		OutputBuffer:   make([]*Number, 0),
		MaxBuffer:      5,
	}
}

func (m *Miner) Update() {
	if m.Deposit == nil || !m.Deposit.CanBeMined() {
		return
	}

	if len(m.OutputBuffer) >= m.MaxBuffer {
		return
	}

	m.MiningTimer++
	if m.MiningTimer >= m.MiningInterval {
		m.MiningTimer = 0
		m.mine()
	}
}

func (m *Miner) Draw(screen *ebiten.Image, camera CameraInterface) {
	worldX, worldY := m.Position.ToWorldPos()
	screenX, screenY := camera.WorldToScreen(worldX, worldY)
	zoom := camera.GetZoom()
	size := float32(TileSize) * float32(zoom)

	// Rest remains the same...
	if size < 4 {
		return
	}

	// Draw miner base
	baseColor := color.RGBA{120, 80, 40, 255}
	vector.DrawFilledRect(screen, float32(screenX), float32(screenY),
		size, size, baseColor, false)

	// Draw drill in center
	drillSize := size * 0.4
	drillX := float32(screenX) + (size-drillSize)/2
	drillY := float32(screenY) + (size-drillSize)/2
	drillColor := color.RGBA{80, 80, 80, 255}
	vector.DrawFilledRect(screen, drillX, drillY, drillSize, drillSize, drillColor, false)

	// Draw output direction indicator
	m.drawOutputIndicator(screen, float32(screenX), float32(screenY), size)

	// Draw mining progress
	progress := float32(m.MiningTimer) / float32(m.MiningInterval)
	if progress > 0 {
		progressColor := color.RGBA{255, 255, 100, 200}
		progressHeight := size * 0.1
		progressWidth := size * progress
		vector.DrawFilledRect(screen, float32(screenX), float32(screenY),
			progressWidth, progressHeight, progressColor, false)
	}

	// Draw border
	borderColor := color.RGBA{160, 120, 80, 255}
	vector.StrokeRect(screen, float32(screenX), float32(screenY),
		size, size, 2, borderColor, false)
}

func (m *Miner) drawOutputIndicator(screen *ebiten.Image, x, y, size float32) {
	centerX := x + size/2
	centerY := y + size/2
	arrowSize := size * 0.15

	indicatorColor := color.RGBA{255, 200, 100, 255}

	switch m.OutputDir {
	case DirectionUp:
		vector.DrawFilledRect(screen, centerX-arrowSize/2, y, arrowSize, arrowSize, indicatorColor, false)
	case DirectionRight:
		vector.DrawFilledRect(screen, x+size-arrowSize, centerY-arrowSize/2, arrowSize, arrowSize, indicatorColor, false)
	case DirectionDown:
		vector.DrawFilledRect(screen, centerX-arrowSize/2, y+size-arrowSize, arrowSize, arrowSize, indicatorColor, false)
	case DirectionLeft:
		vector.DrawFilledRect(screen, x, centerY-arrowSize/2, arrowSize, arrowSize, indicatorColor, false)
	}
}

func (m *Miner) mine() {
	if value, success := m.Deposit.Mine(); success {
		worldX, worldY := m.Position.ToWorldPos()
		number := NewNumber(worldX+TileSize/2, worldY+TileSize/2, value)
		m.OutputBuffer = append(m.OutputBuffer, number)
	}
}

func (m *Miner) GetOutputPosition() GridPosition {
	switch m.OutputDir {
	case DirectionUp:
		return GridPosition{X: m.Position.X, Y: m.Position.Y - 1}
	case DirectionRight:
		return GridPosition{X: m.Position.X + 1, Y: m.Position.Y}
	case DirectionDown:
		return GridPosition{X: m.Position.X, Y: m.Position.Y + 1}
	case DirectionLeft:
		return GridPosition{X: m.Position.X - 1, Y: m.Position.Y}
	default:
		return m.Position
	}
}

func (m *Miner) TryOutputNumber() *Number {
	if len(m.OutputBuffer) > 0 {
		number := m.OutputBuffer[0]
		m.OutputBuffer = m.OutputBuffer[1:]
		return number
	}
	return nil
}

func (m *Miner) HasOutputReady() bool {
	return len(m.OutputBuffer) > 0
}

func (m *Miner) GetGridPosition() GridPosition {
	return m.Position
}

func (m *Miner) GetSize() (int, int) {
	return 1, 1
}
