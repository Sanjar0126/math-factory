package game

import (
	"fmt"
	"image/color"

	"github.com/Sanjar0126/math-factory/internal/fonts"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game represents the main game state
type Game struct {
	screenWidth  int
	screenHeight int
	world        *World
	camera       *Camera
	input        *InputManager
}

// NewGame creates a new game instance
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

// Update updates the game state
func (g *Game) Update() error {
	// Update input
	g.input.Update()

	// Handle camera movement
	g.camera.HandleInput(g.input)

	// Handle world input (building placement, etc.)
	g.world.HandleInput(g.input, g.camera)

	// Update world
	g.world.Update()

	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{32, 32, 32, 255})

	// Draw world
	g.world.Draw(screen, g.camera)

	// Draw UI
	g.drawUI(screen)
}

// Layout returns the screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

// drawUI draws the user interface
func (g *Game) drawUI(screen *ebiten.Image) {
	numbersInWorld, numbersStored, minerCount, depositCount := g.world.GetStats()

	uiText := fmt.Sprintf("Math Factory v0.3 - Grid System\n"+
		"WASD: Move camera, Mouse wheel: Zoom\n"+
		"B: Toggle build mode, 1: Miner, 2: Conveyor\n"+
		"Camera: (%.1f, %.1f) Zoom: %.2f\n"+
		"Numbers in world: %d, Stored: %d\n"+
		"Miners: %d, Deposits: %d",
		g.camera.X, g.camera.Y, g.camera.Zoom,
		numbersInWorld, numbersStored,
		minerCount, depositCount)

	if g.world.BuildMode {
		buildingName := "Unknown"
		switch g.world.SelectedBuilding {
		case BuildingMiner:
			buildingName = "Miner"
		case BuildingConveyor:
			buildingName = "Conveyor"
		}
		uiText += fmt.Sprintf("\nBUILD MODE: %s", buildingName)
	}

	ebitenutil.DebugPrintAt(screen, uiText, 10, 10)
}
