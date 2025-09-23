package game

import (
	"fmt"
	"image/color"

	"github.com/Sanjar0126/math-factory/internal/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	screenWidth  int
	screenHeight int
	world        *World
	camera       *Camera
	input        *InputManager
}

func NewGame(screenWidth, screenHeight int) *Game {
	world := NewWorld()
	camera := NewCamera(screenWidth, screenHeight)
	input := NewInputManager()
	fonts.InitFonts()

	return &Game{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		world:        world,
		camera:       camera,
		input:        input,
	}
}

func (g *Game) Update() error {
	g.input.Update()

	g.camera.HandleInput(g.input)

	g.world.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{32, 32, 32, 255})

	g.world.Draw(screen, g.camera)

	g.drawUI(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

func (g *Game) drawUI(screen *ebiten.Image) {
	debugStr := fmt.Sprintf(
		"Math Factory v0.1\n"+
			"WASD: Move camera\n"+
			"Camera: (%.1f, %.1f)\n"+
			"Zoom: %.2f",
		g.camera.X, g.camera.Y, g.camera.Zoom)
	ebitenutil.DebugPrint(screen, debugStr)
}
