package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Barrier struct {
	targetFn func() bool
}

func NewWait(targetFn func() bool) *Barrier {
	return &Barrier{targetFn: targetFn}
}

func (b *Barrier) Init() {
}

func (b *Barrier) Update() error {
	return nil
}

func (b *Barrier) Draw(screen *ebiten.Image) {
}

func (b *Barrier) Ended() bool {
	return b.targetFn()
}

func (b *Barrier) Dispose() {
}
