package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Init()
	Update() error
	Draw(screen *ebiten.Image)
	Ended() bool
	Dispose()
}
