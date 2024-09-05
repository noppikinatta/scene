package scene

import "github.com/hajimehoshi/ebiten/v2"

type Wait struct {
	target Scene
}

func NewWait(target Scene) *Wait {
	return &Wait{target: target}
}

func (w *Wait) Init() {
}

func (w *Wait) Update() error {
	return nil
}

func (w *Wait) Draw(screen *ebiten.Image) {
}

func (w *Wait) Ended() bool {
	return w.target.Ended()
}

func (w *Wait) Dispose() {
}
