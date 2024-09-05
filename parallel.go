package scene

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

type Parallel struct {
	scenes []Scene
	errs   []error // errs is error cache for Update()
}

func NewParallel(scenes ...Scene) *Parallel {
	return &Parallel{scenes: scenes}
}

func (p *Parallel) Init() {
	for i := range p.scenes {
		p.scenes[i].Init()
	}
}

func (p *Parallel) Update() error {
	if len(p.errs) < len(p.scenes) {
		p.errs = make([]error, len(p.scenes))
	}
	p.errs = p.errs[:len(p.scenes)]

	for i := range p.scenes {
		p.errs[i] = p.scenes[i].Update()
	}

	return errors.Join(p.errs...)
}

func (p *Parallel) Draw(screen *ebiten.Image) {
	for i := range p.scenes {
		p.scenes[i].Draw(screen)
	}
}

func (p *Parallel) Ended() bool {
	for i := range p.scenes {
		if !p.scenes[i].Ended() {
			return false
		}
	}

	return true
}

func (p *Parallel) Dispose() {
	for i := range p.scenes {
		p.scenes[i].Dispose()
	}
}
