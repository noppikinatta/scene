package scene

import "github.com/hajimehoshi/ebiten/v2"

type OnStarter interface {
	OnStart()
}

type OnEnder interface {
	OnEnd()
}

type OnArrivaler interface {
	OnArrival()
}

type OnDeparturer interface {
	OnDeparture()
}

func callIfImpl[T any](g ebiten.Game, fn func(t T)) {
	if t, ok := g.(T); ok {
		fn(t)
	}
}
