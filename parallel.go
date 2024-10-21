package scene

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

// Parallel runs multiple Scenes in parallel.
type Parallel struct {
	scenes []Scene
	errs   []error // errs is error cache for Update()
}

// NewParallel creates a new Parallel instance.
func NewParallel(scenes ...Scene) *Parallel {
	return &Parallel{scenes: scenes}
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

func (p *Parallel) CanEnd() bool {
	for i := range p.scenes {
		if !p.scenes[i].CanEnd() {
			return false
		}
	}

	return true
}

func (p *Parallel) OnSceneStart() {
	for _, s := range p.scenes {
		callIfImpl(s, func(o OnSceneStarter) { o.OnSceneStart() })
	}
}

func (p *Parallel) OnSceneEnd() {
	for _, s := range p.scenes {
		callIfImpl(s, func(o OnSceneEnder) { o.OnSceneEnd() })
	}
}

func (p *Parallel) OnTransitionStart() {
	for _, s := range p.scenes {
		callIfImpl(s, func(o OnTransitionStarter) { o.OnTransitionStart() })
	}
}

func (p *Parallel) OnTransitionEnd() {
	for _, s := range p.scenes {
		callIfImpl(s, func(o OnTransitionEnder) { o.OnTransitionEnd() })
	}
}
