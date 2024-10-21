package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Scene represents a single scene.
type Scene interface {
	// Update updates the state of the Scene.
	// Update of ebiten.Game created by ToGame function returns the return value of it.
	Update() error
	// Draw draws a Scene.
	// Draw of ebiten.Game created by ToGame function calls this.
	// However, if Game's Update has never been called, or if the Scene is finished, it will not be called.
	Draw(screen *ebiten.Image)
	// CanEnd returns True if the Scene is finished.
	// ebiten.Game created by ToGame function, returns ebiten.Termination when it returns True.
	CanEnd() bool
}

type OnSceneStarter interface {
	OnSceneStart()
}

type OnSceneEnder interface {
	OnSceneEnd()
}

type OnTransitionStarter interface {
	OnTransitionStart()
}

type OnTransitionEnder interface {
	OnTransitionEnd()
}

func callIfImpl[T any](s Scene, fn func(t T)) {
	if t, ok := s.(T); ok {
		fn(t)
	}
}
