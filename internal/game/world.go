package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	tileSize = 32
	gridSize = 64
)

type World struct {
}

func NewWorld() *World {
	return &World{}
}

func (w *World) Update() {

}

func (w *World) Draw(screen *ebiten.Image, camera *Camera) {
	w.drawGrid(screen, camera)
	w.drawOrigin(screen, camera)
}

func (w *World) drawGrid(screen *ebiten.Image, camera *Camera) {

	startX := int(camera.X-float64(camera.screenWidth)/(2*camera.Zoom*tileSize)) - 1
	endX := int(camera.X+float64(camera.screenWidth)/(2*camera.Zoom*tileSize)) + 1
	startY := int(camera.Y-float64(camera.screenHeight)/(2*camera.Zoom*tileSize)) - 1
	endY := int(camera.Y+float64(camera.screenHeight)/(2*camera.Zoom*tileSize)) + 1

	for x := startX; x <= endX; x++ {
		worldX := float64(x * tileSize)
		screenX1, screenY1 := camera.WorldToScreen(worldX, float64(startY*tileSize))
		screenX2, screenY2 := camera.WorldToScreen(worldX, float64(endY*tileSize))

		if camera.Zoom > 0.5 {
			ebitenutil.DrawLine(screen, screenX1, screenY1, screenX2, screenY2,
				color.RGBA{64, 64, 64, 255})
		}
	}

	for y := startY; y <= endY; y++ {
		worldY := float64(y * tileSize)
		screenX1, screenY1 := camera.WorldToScreen(float64(startX*tileSize), worldY)
		screenX2, screenY2 := camera.WorldToScreen(float64(endX*tileSize), worldY)

		if camera.Zoom > 0.5 {
			ebitenutil.DrawLine(screen, screenX1, screenY1, screenX2, screenY2,
				color.RGBA{64, 64, 64, 255})
		}
	}
}

func (w *World) drawOrigin(screen *ebiten.Image, camera *Camera) {
	screenX, screenY := camera.WorldToScreen(0, 0)

	ebitenutil.DrawLine(screen, screenX-10, screenY, screenX+10, screenY,
		color.RGBA{255, 0, 0, 255})
	ebitenutil.DrawLine(screen, screenX, screenY-10, screenX, screenY+10,
		color.RGBA{255, 0, 0, 255})
}
