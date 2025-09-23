package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputManager struct {
	mouseX, mouseY int
	wheelX, wheelY float64
}

func NewInputManager() *InputManager {
	return &InputManager{}
}

func (im *InputManager) Update() {
	im.mouseX, im.mouseY = ebiten.CursorPosition()
	im.wheelX, im.wheelY = ebiten.Wheel()
}

func (im *InputManager) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

func (im *InputManager) IsKeyJustPressed(key ebiten.Key) bool {
	return inpututil.IsKeyJustPressed(key)
}

func (im *InputManager) IsMouseButtonPressed(button ebiten.MouseButton) bool {
	return ebiten.IsMouseButtonPressed(button)
}

func (im *InputManager) IsMouseButtonJustPressed(button ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustPressed(button)
}

func (im *InputManager) GetMousePosition() (int, int) {
	return im.mouseX, im.mouseY
}

func (im *InputManager) GetWheelDelta() (float64, float64) {
	return im.wheelX, im.wheelY
}
