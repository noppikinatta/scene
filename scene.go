package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Scene represents a single scene.
type Scene interface {
	// Init is the initialization process of the Scene.
	// Init is called in the process of the first Update of ebiten.Game created by ToGame function, before the Scene's Update.
	// Init is called when the Scene is switched inside the Chain, before the Scene's Update.
	Init()
	// Update updates the state of the Scene.
	// Update of ebiten.Game created by ToGame function returns the return value of it.
	Update() error
	// Draw draws a Scene.
	// Draw of ebiten.Game created by ToGame function calls this.
	// However, if Game's Update has never been called, or if the Scene is finished, it will not be called.
	Draw(screen *ebiten.Image)
	// Ended returns True if the Scene is finished.
	// ebiten.Game created by ToGame function, returns ebiten.Termination when it returns True.
	Ended() bool
	// Dispose releases the resources used by this Scene.
	// Game created by the ToGame function is called when the Scene's Ended returns true.
	// Chain calls the previous Scene's Dispose when the Scene to be run switches.
	Dispose()
}
