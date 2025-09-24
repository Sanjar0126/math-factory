package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TileSize = 32
)

type GridPosition struct {
	X, Y int
}

func (gp GridPosition) ToWorldPos() (float64, float64) {
	return float64(gp.X * TileSize), float64(gp.Y * TileSize)
}

func WorldPosToGrid(worldX, worldY float64) GridPosition {
	return GridPosition{
		X: int(worldX / TileSize),
		Y: int(worldY / TileSize),
	}
}

type Entity interface {
	Update()
	Draw(screen *ebiten.Image, camera Camera)
	GetGridPosition() GridPosition
	GetSize() (int, int)
}

type Camera interface {
	WorldToScreen(worldX, worldY float64) (float64, float64)
	GetZoom() float64
}
