package game

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/Sanjar0126/math-factory/internal/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	TileSize = 32
)

// BuildingType represents different types of buildings
type BuildingType int

const (
	BuildingMiner BuildingType = iota
	BuildingConveyor
	BuildingProcessor
)

// World represents the game world with grid-based entities
type World struct {
	// Grid-based storage
	Grid     map[entities.GridPosition]entities.Entity
	Deposits map[entities.GridPosition]*entities.NumberDeposit
	Core     *entities.Core
	Miners   []*entities.Miner
	Numbers  []*entities.Number

	// Building placement
	SelectedBuilding BuildingType
	BuildMode        bool
	PreviewPosition  entities.GridPosition

	// World generation
	GeneratedChunks map[ChunkPosition]bool
}

// ChunkPosition represents a chunk of the world (for generation)
type ChunkPosition struct {
	X, Y int
}

const ChunkSize = 16 // 16x16 tiles per chunk

// NewWorld creates a new grid-based world
func NewWorld() *World {
	world := &World{
		Grid:             make(map[entities.GridPosition]entities.Entity),
		Deposits:         make(map[entities.GridPosition]*entities.NumberDeposit),
		Miners:           make([]*entities.Miner, 0),
		Numbers:          make([]*entities.Number, 0),
		SelectedBuilding: BuildingMiner,
		BuildMode:        false,
		GeneratedChunks:  make(map[ChunkPosition]bool),
	}

	// Create core at origin (0,0) - it's 2x2 so occupies (0,0), (1,0), (0,1), (1,1)
	world.Core = entities.NewCore(0, 0)
	world.placeEntity(world.Core)

	// Generate initial area around spawn
	world.generateArea(-20, -20, 40, 40)

	return world
}

// Update updates the world state
func (w *World) Update() {
	// Update all entities in grid
	for _, entity := range w.Grid {
		entity.Update()
	}

	// Update core
	w.Core.Update()

	// Check miners for output and handle number movement
	for _, miner := range w.Miners {
		if miner.HasOutputReady() {
			// For now, output directly as floating numbers
			// Later this will be handled by conveyor system
			if number := miner.TryOutputNumber(); number != nil {
				// Position number at miner's output position
				outputPos := miner.GetOutputPosition()
				outputWorldX, outputWorldY := outputPos.ToWorldPos()
				number.X = outputWorldX + TileSize/2
				number.Y = outputWorldY + TileSize/2
				w.Numbers = append(w.Numbers, number)
			}
		}
	}

	// Update floating numbers and check core collection
	for i := len(w.Numbers) - 1; i >= 0; i-- {
		number := w.Numbers[i]
		number.Update()

		// Simple attraction to core for now
		coreWorldX, coreWorldY := w.Core.GetGridPosition().ToWorldPos()
		centerX := coreWorldX + TileSize
		centerY := coreWorldY + TileSize

		dx := number.X - centerX
		dy := number.Y - centerY
		distance := dx*dx + dy*dy

		// If close enough to core, attract it
		if distance < (TileSize*4)*(TileSize*4) {
			number.MoveTo(centerX, centerY, 2.0)

			// If very close, collect it
			if distance < 400 {
				w.Core.AcceptNumber(number)
				w.Numbers = append(w.Numbers[:i], w.Numbers[i+1:]...)
			}
		}
	}
}

// HandleInput processes world-related input
func (w *World) HandleInput(input *InputManager, camera *Camera) {
	// Toggle build mode
	if input.IsKeyJustPressed(ebiten.KeyB) {
		w.BuildMode = !w.BuildMode
	}

	// Cycle building types in build mode
	if w.BuildMode {
		if input.IsKeyJustPressed(ebiten.Key1) {
			w.SelectedBuilding = BuildingMiner
		}
		if input.IsKeyJustPressed(ebiten.Key2) {
			w.SelectedBuilding = BuildingConveyor
		}
	}

	// Handle building placement
	if w.BuildMode && input.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := input.GetMousePosition()
		worldX, worldY := camera.ScreenToWorld(mouseX, mouseY)
		gridPos := entities.WorldPosToGrid(worldX, worldY)

		w.tryPlaceBuilding(gridPos)
	}

	// Update preview position
	if w.BuildMode {
		mouseX, mouseY := input.GetMousePosition()
		worldX, worldY := camera.ScreenToWorld(mouseX, mouseY)
		w.PreviewPosition = entities.WorldPosToGrid(worldX, worldY)
	}

	// Generate world chunks as needed
	w.generateAroundCamera(camera)
}

// tryPlaceBuilding attempts to place a building at the given position
func (w *World) tryPlaceBuilding(pos entities.GridPosition) {
	switch w.SelectedBuilding {
	case BuildingMiner:
		w.tryPlaceMiner(pos)
	}
}

// tryPlaceMiner attempts to place a miner at the given position
func (w *World) tryPlaceMiner(pos entities.GridPosition) {
	// Check if position is occupied
	if w.isPositionOccupied(pos) {
		return
	}

	// Check if there's a deposit at this position
	deposit, exists := w.Deposits[pos]
	if !exists || !deposit.CanBeMined() {
		return
	}

	// Create miner with default output direction (right)
	miner := entities.NewMiner(pos.X, pos.Y, deposit, entities.DirectionRight)
	deposit.SetMined(true)

	// Add to world
	w.Miners = append(w.Miners, miner)
	w.placeEntity(miner)
}

// Draw renders the world
func (w *World) Draw(screen *ebiten.Image, camera *Camera) {
	w.drawGrid(screen, camera)
	w.drawDeposits(screen, camera)
	w.drawEntities(screen, camera)
	w.drawNumbers(screen, camera)
	w.drawBuildPreview(screen, camera)
}

// drawGrid draws the world grid
func (w *World) drawGrid(screen *ebiten.Image, camera *Camera) {
	if camera.GetZoom() < 0.5 {
		return // Don't draw grid when zoomed out too much
	}

	// Calculate visible area
	screenW, screenH := screen.Bounds().Dx(), screen.Bounds().Dy()
	startWorldX := camera.X - float64(screenW)/(2*camera.GetZoom())
	endWorldX := camera.X + float64(screenW)/(2*camera.GetZoom())
	startWorldY := camera.Y - float64(screenH)/(2*camera.GetZoom())
	endWorldY := camera.Y + float64(screenH)/(2*camera.GetZoom())

	startGridX := int(startWorldX/TileSize) - 1
	endGridX := int(endWorldX/TileSize) + 1
	startGridY := int(startWorldY/TileSize) - 1
	endGridY := int(endWorldY/TileSize) + 1

	gridColor := color.RGBA{64, 64, 64, 255}

	// Draw vertical lines
	for x := startGridX; x <= endGridX; x++ {
		worldX := float64(x * TileSize)
		screenX1, screenY1 := camera.WorldToScreen(worldX, startWorldY)
		screenX2, screenY2 := camera.WorldToScreen(worldX, endWorldY)
		vector.StrokeLine(screen, float32(screenX1), float32(screenY1), float32(screenX2), float32(screenY2), 1, gridColor, false)
	}

	// Draw horizontal lines
	for y := startGridY; y <= endGridY; y++ {
		worldY := float64(y * TileSize)
		screenX1, screenY1 := camera.WorldToScreen(startWorldX, worldY)
		screenX2, screenY2 := camera.WorldToScreen(endWorldX, worldY)
		vector.StrokeLine(screen, float32(screenX1), float32(screenY1), float32(screenX2), float32(screenY2), 1, gridColor, false)
	}
}

// drawDeposits draws all number deposits
func (w *World) drawDeposits(screen *ebiten.Image, camera *Camera) {
	for _, deposit := range w.Deposits {
		deposit.Draw(screen, camera)
	}
}

// drawEntities draws all placed entities
func (w *World) drawEntities(screen *ebiten.Image, camera *Camera) {
	for _, entity := range w.Grid {
		entity.Draw(screen, camera)
	}
}

// drawNumbers draws all floating numbers
func (w *World) drawNumbers(screen *ebiten.Image, camera *Camera) {
	for _, number := range w.Numbers {
		number.Draw(screen, camera)
	}
}

// drawBuildPreview draws building placement preview
func (w *World) drawBuildPreview(screen *ebiten.Image, camera *Camera) {
	if !w.BuildMode {
		return
	}

	worldX, worldY := w.PreviewPosition.ToWorldPos()
	screenX, screenY := camera.WorldToScreen(worldX, worldY)
	size := float32(TileSize * camera.GetZoom())

	// Choose preview color based on validity
	previewColor := color.RGBA{100, 255, 100, 100} // Green for valid
	if !w.canPlaceAt(w.PreviewPosition) {
		previewColor = color.RGBA{255, 100, 100, 100} // Red for invalid
	}

	// Draw preview rectangle
	vector.DrawFilledRect(screen, float32(screenX), float32(screenY),
		size, size, previewColor, false)
}

// generateArea generates deposits in the specified area
func (w *World) generateArea(startX, startY, width, height int) {
	for x := startX; x < startX+width; x++ {
		for y := startY; y < startY+height; y++ {
			pos := entities.GridPosition{X: x, Y: y}

			// Skip if already has deposit or is occupied by core
			if _, exists := w.Deposits[pos]; exists {
				continue
			}
			if w.Core.OccupiesPosition(pos) {
				continue
			}

			// Generate deposit based on distance from origin and randomness
			if w.shouldGenerateDeposit(x, y) {
				value := w.generateNumberForPosition(x, y)
				infinite := w.shouldBeInfinite(x, y, value)
				deposit := entities.NewNumberDeposit(x, y, value, infinite)
				w.Deposits[pos] = deposit
			}
		}
	}
}

// shouldGenerateDeposit determines if a deposit should be generated at this position
func (w *World) shouldGenerateDeposit(x, y int) bool {
	// Don't generate too close to origin
	if x >= -1 && x <= 2 && y >= -1 && y <= 2 {
		return false
	}

	// Use better random seed based on position
	// Create a proper pseudo-random number for this position
	seed := int64(x*1000000 + y*1000 + 12345) // Better seed combination
	rng := rand.New(rand.NewSource(seed))

	// 15% chance for deposits
	return rng.Float64() < 0.15
}

// generateNumberForPosition generates an appropriate number for the given position
func (w *World) generateNumberForPosition(x, y int) int {
	// Calculate distance from origin (0,0)
	distanceFromOrigin := math.Sqrt(float64(x*x + y*y))

	// Use the same position-based random seed for consistency
	seed := int64(x*1000000 + y*1000 + 12345)
	rng := rand.New(rand.NewSource(seed))

	if distanceFromOrigin < 3 {
		// Very close to origin: tiny numbers (1-5)
		return rng.Intn(5) + 1
	} else if distanceFromOrigin < 8 {
		// Close to origin: small numbers (1-20)
		return rng.Intn(20) + 1
	} else if distanceFromOrigin < 15 {
		// Medium distance: medium numbers (5-100)
		return rng.Intn(96) + 5
	} else if distanceFromOrigin < 25 {
		// Far distance: larger numbers (20-500)
		return rng.Intn(481) + 20
	} else {
		// Very far: huge numbers and more primes
		maxValue := int(distanceFromOrigin * 50) // Numbers scale with distance

		return generatePrimeInRange(50, maxValue, rng)
	}
}

// shouldBeInfinite determines if a deposit should be infinite
func (w *World) shouldBeInfinite(x, y int, value int) bool {
	// Small chance for infinite deposits, higher for primes
	baseChance := 0.05
	if entities.IsPrime(value) {
		baseChance = 0.15
	}
	return rand.Float64() < baseChance
}

// generateAroundCamera generates world chunks around the camera
func (w *World) generateAroundCamera(camera *Camera) {
	chunkX := int(camera.X) / (ChunkSize * TileSize)
	chunkY := int(camera.Y) / (ChunkSize * TileSize)

	// Generate 3x3 chunks around camera
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			chunk := ChunkPosition{X: chunkX + dx, Y: chunkY + dy}
			if !w.GeneratedChunks[chunk] {
				w.generateChunk(chunk)
				w.GeneratedChunks[chunk] = true
			}
		}
	}
}

// generateChunk generates a single chunk
func (w *World) generateChunk(chunk ChunkPosition) {
	startX := chunk.X * ChunkSize
	startY := chunk.Y * ChunkSize
	w.generateArea(startX, startY, ChunkSize, ChunkSize)
}

// Utility methods
func (w *World) isPositionOccupied(pos entities.GridPosition) bool {
	_, occupied := w.Grid[pos]
	return occupied || w.Core.OccupiesPosition(pos)
}

func (w *World) canPlaceAt(pos entities.GridPosition) bool {
	switch w.SelectedBuilding {
	case BuildingMiner:
		return !w.isPositionOccupied(pos) && w.hasDepositAt(pos)
	default:
		return !w.isPositionOccupied(pos)
	}
}

func (w *World) hasDepositAt(pos entities.GridPosition) bool {
	deposit, exists := w.Deposits[pos]
	return exists && deposit.CanBeMined()
}

func (w *World) placeEntity(entity entities.Entity) {
	pos := entity.GetGridPosition()
	sizeX, sizeY := entity.GetSize()

	// Place entity in all positions it occupies
	for dx := 0; dx < sizeX; dx++ {
		for dy := 0; dy < sizeY; dy++ {
			occupiedPos := entities.GridPosition{X: pos.X + dx, Y: pos.Y + dy}
			w.Grid[occupiedPos] = entity
		}
	}
}

// Update this function to accept an RNG parameter
func generatePrimeInRange(min, max int, rng *rand.Rand) int {
	if min > max || min < 2 {
		return 2
	}

	// Try to find a prime in the range
	for attempts := 0; attempts < 100; attempts++ {
		candidate := rng.Intn(max-min+1) + min
		if entities.IsPrime(candidate) {
			return candidate
		}
	}

	// Fallback: find next prime after min
	for candidate := min; candidate <= max; candidate++ {
		if entities.IsPrime(candidate) {
			return candidate
		}
	}
	return 2 // Ultimate fallback
}

// GetStats returns world statistics
func (w *World) GetStats() (int, int, int, int) {
	return len(w.Numbers), w.Core.GetStoredCount(), len(w.Miners), len(w.Deposits)
}
