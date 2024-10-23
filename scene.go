package scene

import "github.com/hajimehoshi/ebiten/v2"

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

func callIfImpl[T any](g ebiten.Game, fn func(t T)) {
	if t, ok := g.(T); ok {
		fn(t)
	}
}
