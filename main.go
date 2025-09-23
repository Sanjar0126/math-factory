package main

import (
    "log"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/Sanjar0126/math-factory/internal/game"
)

const (
    screenWidth  = 1280
    screenHeight = 720
    gameTitle    = "Math Factory"
)

func main() {
    g := game.NewGame(screenWidth, screenHeight)

    ebiten.SetWindowSize(screenWidth, screenHeight)
    ebiten.SetWindowTitle(gameTitle)
    ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

    if err := ebiten.RunGame(g); err != nil {
        log.Fatal(err)
    }
}