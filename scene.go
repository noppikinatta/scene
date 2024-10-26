package scene

import "github.com/hajimehoshi/ebiten/v2"

type OnStarter interface {
	OnStart()
}

type OnEnder interface {
	OnEnd()
}

type OnDeparturer interface {
	OnDeparture()
}

type OnArrivaler interface {
	OnArrival()
}

func callIfImpl[T any](g ebiten.Game, fn func(t T)) {
	if t, ok := g.(T); ok {
		fn(t)
	}
}
