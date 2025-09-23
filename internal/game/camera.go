package game

import (
    "github.com/hajimehoshi/ebiten/v2"
)

const (
    cameraSpeed = 1.0
    minZoom     = 0.25
    maxZoom     = 4.0
    zoomStep    = 0.1
)

type Camera struct {
    X, Y        float64
    Zoom        float64
    screenWidth int
    screenHeight int
}

func NewCamera(screenWidth, screenHeight int) *Camera {
    return &Camera{
        X:            0,
        Y:            0,
        Zoom:         1.0,
        screenWidth:  screenWidth,
        screenHeight: screenHeight,
    }
}

func (c *Camera) HandleInput(input *InputManager) {
    if input.IsKeyPressed(ebiten.KeyW) || input.IsKeyPressed(ebiten.KeyArrowUp) {
        c.Y -= cameraSpeed / c.Zoom
    }
    if input.IsKeyPressed(ebiten.KeyS) || input.IsKeyPressed(ebiten.KeyArrowDown) {
        c.Y += cameraSpeed / c.Zoom
    }
    if input.IsKeyPressed(ebiten.KeyA) || input.IsKeyPressed(ebiten.KeyArrowLeft) {
        c.X -= cameraSpeed / c.Zoom
    }
    if input.IsKeyPressed(ebiten.KeyD) || input.IsKeyPressed(ebiten.KeyArrowRight) {
        c.X += cameraSpeed / c.Zoom
    }

    // Zoom
    _, wheelY := input.GetWheelDelta()
    if wheelY > 0 && c.Zoom < maxZoom {
        c.Zoom += zoomStep
    } else if wheelY < 0 && c.Zoom > minZoom {
        c.Zoom -= zoomStep
    }
}

func (c *Camera) WorldToScreen(worldX, worldY float64) (screenX, screenY float64) {
    screenX = (worldX-c.X)*c.Zoom + float64(c.screenWidth)/2
    screenY = (worldY-c.Y)*c.Zoom + float64(c.screenHeight)/2
    return
}

func (c *Camera) ScreenToWorld(screenX, screenY int) (worldX, worldY float64) {
    worldX = (float64(screenX)-float64(c.screenWidth)/2)/c.Zoom + c.X
    worldY = (float64(screenY)-float64(c.screenHeight)/2)/c.Zoom + c.Y
    return
}

func (c *Camera) GetTransform() ebiten.GeoM {
    var transform ebiten.GeoM
    transform.Scale(c.Zoom, c.Zoom)
    transform.Translate(-c.X*c.Zoom+float64(c.screenWidth)/2, 
                       -c.Y*c.Zoom+float64(c.screenHeight)/2)
    return transform
}