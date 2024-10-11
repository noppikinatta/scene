package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Barrier ends when passed function returned true. It is useful with Chain and Parallel.
type Barrier struct {
	targetFn func() bool
}

// NewBarrier creates a new Barrier instance.
func NewBarrier(targetFn func() bool) *Barrier {
	return &Barrier{targetFn: targetFn}
}

func (b *Barrier) Update() error {
	return nil
}

func (b *Barrier) Draw(screen *ebiten.Image) {
}

func (b *Barrier) CanEnd() bool {
	return b.targetFn()
}
