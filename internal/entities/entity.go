package entities

import "github.com/hajimehoshi/ebiten/v2"

const (
	TileSize = 32 // Each grid cell is 32x32 pixels
)

// GridPosition represents a position in the grid
type GridPosition struct {
	X, Y int
}

// ToWorldPos converts grid position to world coordinates
func (gp GridPosition) ToWorldPos() (float64, float64) {
	return float64(gp.X * TileSize), float64(gp.Y * TileSize)
}

// WorldPosToGrid converts world coordinates to grid position
func WorldPosToGrid(worldX, worldY float64) GridPosition {
	return GridPosition{
		X: int(worldX / TileSize),
		Y: int(worldY / TileSize),
	}
}

// CameraInterface defines what we need from camera
type CameraInterface interface {
	WorldToScreen(worldX, worldY float64) (float64, float64)
	GetZoom() float64
}

type Entity interface {
	Update()
	Draw(screen *ebiten.Image, camera CameraInterface)
	GetGridPosition() GridPosition
	GetSize() (int, int) // Size in grid cells
}
