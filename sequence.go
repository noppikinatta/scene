package scene

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sequence struct {
	current           ebiten.Game
	transitionUpdater *transitionUpdater
	inited            bool
}

func NewSequence(first ebiten.Game) *Sequence {
	return &Sequence{current: first}
}

func (s *Sequence) Update() error {
	if s.inTransition() {
		if err := s.transitionUpdater.Update(); err != nil {
			return err
		}
	}

	if !s.inited {
		callIfImpl(s.current, func(o OnSceneStarter) { o.OnSceneStart() })
		s.inited = true
	}

	err := s.current.Update()
	if errors.Is(err, ebiten.Termination) {
		callIfImpl(s.current, func(o OnSceneEnder) { o.OnSceneEnd() })
	}

	return err
}

func (s *Sequence) Draw(screen *ebiten.Image) {
	s.current.Draw(screen)
	if s.inTransition() {
		s.transitionUpdater.Draw(screen)
	}
}

func (s *Sequence) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return s.current.Layout(outsideWidth, outsideHeight)
}

func (s *Sequence) Switch(next ebiten.Game) bool {
	return s.SwitchWithTransition(next, NopTransition)
}

func (s *Sequence) SwitchWithTransition(next ebiten.Game, transition Transition) bool {
	if s.inTransition() {
		return false
	}
	p := newTransitionUpdater(s, next, transition)
	s.transitionUpdater = p
	transition.Reset()
	callIfImpl(s.current, func(o OnTransitionStarter) { o.OnTransitionStart() })
	return true
}

func (s *Sequence) inTransition() bool {
	return s.transitionUpdater != nil
}

func (s *Sequence) switchScenes(next ebiten.Game) {
	callIfImpl(s.current, func(o OnSceneEnder) { o.OnSceneEnd() })
	s.current = next
	callIfImpl(s.current, func(o OnSceneStarter) { o.OnSceneStart() })
}

func (s *Sequence) endTransition() {
	s.transitionUpdater = nil
	callIfImpl(s.current, func(o OnTransitionEnder) { o.OnTransitionEnd() })
}

func (p *Sequence) drawFinalScreenFunc() func(screen ebiten.FinalScreen, offScreen *ebiten.Image, geoM ebiten.GeoM) {
	return func(screen ebiten.FinalScreen, offScreen *ebiten.Image, geoM ebiten.GeoM) {
		if f, ok := p.current.(ebiten.FinalScreenDrawer); ok {
			f.DrawFinalScreen(screen, offScreen, geoM)
		} else {
			defaultDrawFinalScreenTemporaryImplRemoveItWhenEbitengineV290Released(screen, offScreen, geoM)
		}
	}
}

func (p *Sequence) layoutFFunc() func(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	return func(outsideWidth, outsideHeight float64) (screenWidth float64, screenHeight float64) {
		if l, ok := p.current.(ebiten.LayoutFer); ok {
			return l.LayoutF(outsideWidth, outsideHeight)
		}

		wi, hi := p.current.Layout(int(outsideWidth), int(outsideHeight))
		return float64(wi), float64(hi)
	}
}
