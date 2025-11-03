package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Core struct {
	Position        GridPosition
	StoredNumbers   []int
	InputPositions  []GridPosition
	ProcessingQueue []*Number
}

func NewCore(gridX, gridY int) *Core {
	pos := GridPosition{X: gridX, Y: gridY}

	// Core is 2x2, so it accepts inputs from all sides
	inputPositions := []GridPosition{
		{X: gridX - 1, Y: gridY},     // Left side, top
		{X: gridX - 1, Y: gridY + 1}, // Left side, bottom
		{X: gridX + 2, Y: gridY},     // Right side, top
		{X: gridX + 2, Y: gridY + 1}, // Right side, bottom
		{X: gridX, Y: gridY - 1},     // Top side, left
		{X: gridX + 1, Y: gridY - 1}, // Top side, right
		{X: gridX, Y: gridY + 2},     // Bottom side, left
		{X: gridX + 1, Y: gridY + 2}, // Bottom side, right
	}

	return &Core{
		Position:        pos,
		StoredNumbers:   make([]int, 0),
		InputPositions:  inputPositions,
		ProcessingQueue: make([]*Number, 0),
	}
}

func (c *Core) Update() {
	for i := len(c.ProcessingQueue) - 1; i >= 0; i-- {
		number := c.ProcessingQueue[i]

		coreWorldX, coreWorldY := c.Position.ToWorldPos()
		centerX := coreWorldX + TileSize
		centerY := coreWorldY + TileSize

		number.MoveTo(centerX, centerY, 3.0)
		number.Update()

		dx := number.X - centerX
		dy := number.Y - centerY
		distance := math.Sqrt(dx*dx + dy*dy)

		if distance < 10 {
			c.StoredNumbers = append(c.StoredNumbers, number.Value)
			c.ProcessingQueue = append(c.ProcessingQueue[:i], c.ProcessingQueue[i+1:]...)
		}
	}
}

// Update just the Draw method signature:
func (c *Core) Draw(screen *ebiten.Image, camera CameraInterface) {
    worldX, worldY := c.Position.ToWorldPos()
    screenX, screenY := camera.WorldToScreen(worldX, worldY)
    zoom := camera.GetZoom()
    size := float32(TileSize*2) * float32(zoom) // 2x2 core

    // Rest remains the same...
    if size < 8 {
        return
    }

    // Draw main core structure
    coreColor := color.RGBA{60, 100, 180, 255}
    vector.DrawFilledRect(screen, float32(screenX), float32(screenY), 
        size, size, coreColor, false)

    // Draw inner core
    innerSize := size * 0.7
    innerOffset := (size - innerSize) / 2
    innerColor := color.RGBA{100, 140, 220, 255}
    vector.DrawFilledRect(screen, float32(screenX)+innerOffset, float32(screenY)+innerOffset, 
        innerSize, innerSize, innerColor, false)

    // Draw energy animation
    time := float32(ebiten.TPS()) * 0.05
    pulse := float32(math.Sin(float64(time))) * 0.1 + 1.0
    energySize := innerSize * pulse
    energyOffset := (size - energySize) / 2
    energyColor := color.RGBA{150, 180, 255, 100}
    vector.StrokeRect(screen, float32(screenX)+energyOffset, float32(screenY)+energyOffset, 
        energySize, energySize, 2, energyColor, false)

    // Draw border
    borderColor := color.RGBA{100, 140, 220, 255}
    vector.StrokeRect(screen, float32(screenX), float32(screenY), 
        size, size, 3, borderColor, false)

    // Draw processing numbers
    for _, number := range c.ProcessingQueue {
        number.Draw(screen, camera)
    }
}

func (c *Core) CanAcceptInput(fromPos GridPosition) bool {
	for _, inputPos := range c.InputPositions {
		if inputPos.X == fromPos.X && inputPos.Y == fromPos.Y {
			return true
		}
	}
	return false
}

func (c *Core) AcceptNumber(number *Number) {
	c.ProcessingQueue = append(c.ProcessingQueue, number)
}

func (c *Core) GetGridPosition() GridPosition {
	return c.Position
}

func (c *Core) GetSize() (int, int) {
	return 2, 2
}

func (c *Core) GetStoredCount() int {
	return len(c.StoredNumbers)
}

func (c *Core) OccupiesPosition(pos GridPosition) bool {
	return pos.X >= c.Position.X && pos.X < c.Position.X+2 &&
		pos.Y >= c.Position.Y && pos.Y < c.Position.Y+2
}
