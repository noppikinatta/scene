package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Initer
	Updater
	Drawer
	Ended() bool
	Disposer
}

type Initer interface {
	Init()
}

type Updater interface {
	Update() error
}

type Drawer interface {
	Draw(screen *ebiten.Image)
}

type Disposer interface {
	Dispose()
}
