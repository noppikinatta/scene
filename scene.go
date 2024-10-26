package scene

import "github.com/hajimehoshi/ebiten/v2"

// OnStarter is an interface that executes processing at the start of a scene.
type OnStarter interface {
	// OnStart　is called immediately after the start of the scene.
	OnStart()
}

// OnEnder is an interface that executes processing at the end of a scene.
type OnEnder interface {
	// OnEnd is called just before the end of the scene.
	OnEnd()
}

// OnArrivaler is an interface that executes processing when a scene begins and the beginning transition is completed.
type OnArrivaler interface {
	OnArrival()
}

// OnDeparturer is an interface that executes processing when a scene ends and the ending transition begins.
type OnDeparturer interface {
	OnDeparture()
}

func callIfImpl[T any](g ebiten.Game, fn func(t T)) {
	if t, ok := g.(T); ok {
		fn(t)
	}
}
