package scene

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

// Sequence provides the ability to run multiple ebiten.Games in sequence.
type Sequence struct {
	current           ebiten.Game
	transitionUpdater *transitionUpdater
	inited            bool
}

// NewSequence creates a new Sequence instance.
func NewSequence(first ebiten.Game) *Sequence {
	return &Sequence{current: first}
}

// Update is ebiten.Game.Update implementation.
func (s *Sequence) Update() error {
	if s.inTransition() {
		if err := s.transitionUpdater.Update(); err != nil {
			return err
		}
	}

	if !s.inited {
		callIfImpl(s.current, func(o OnStarter) { o.OnStart() })
		s.inited = true
	}

	err := s.current.Update()
	if errors.Is(err, ebiten.Termination) {
		callIfImpl(s.current, func(o OnEnder) { o.OnEnd() })
	}

	return err
}

// Draw is ebiten.Game.Draw implementation.
func (s *Sequence) Draw(screen *ebiten.Image) {
	s.current.Draw(screen)
	if s.inTransition() {
		s.transitionUpdater.Draw(screen)
	}
}

// Layout is ebiten.Game.Layout implementation.
func (s *Sequence) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return s.current.Layout(outsideWidth, outsideHeight)
}

// Switch switches the ebiten.Game to run in Sequence.
func (s *Sequence) Switch(next ebiten.Game) bool {
	return s.SwitchWithTransition(next, NopTransition)
}

// SwitchWithTransition switches the ebiten.Game to run in Sequence with the Transition.
func (s *Sequence) SwitchWithTransition(next ebiten.Game, transition Transition) bool {
	if s.inTransition() {
		return false
	}
	p := newTransitionUpdater(s, next, transition)
	s.transitionUpdater = p
	transition.Reset()
	callIfImpl(s.current, func(o OnDeparturer) { o.OnDeparture() })
	return true
}

// inTransition returns true if the Transition is being processed.
func (s *Sequence) inTransition() bool {
	return s.transitionUpdater != nil
}

// switchScenes switches scenes.
func (s *Sequence) switchScenes(next ebiten.Game) {
	callIfImpl(s.current, func(o OnEnder) { o.OnEnd() })
	s.current = next
	callIfImpl(s.current, func(o OnStarter) { o.OnStart() })
}

// endTransition is called when the Transition completed.
func (s *Sequence) endTransition() {
	s.transitionUpdater = nil
	callIfImpl(s.current, func(o OnArrivaler) { o.OnArrival() })
}

// drawFinalScreenFunc is FinalScreenDrawerConvertible implementation.
func (s *Sequence) drawFinalScreenFunc() func(screen ebiten.FinalScreen, offScreen *ebiten.Image, geoM ebiten.GeoM) {
	return func(screen ebiten.FinalScreen, offScreen *ebiten.Image, geoM ebiten.GeoM) {
		if f, ok := s.current.(ebiten.FinalScreenDrawer); ok {
			f.DrawFinalScreen(screen, offScreen, geoM)
		} else {
			defaultDrawFinalScreenTemporaryImplRemoveItWhenEbitengineV290Released(screen, offScreen, geoM)
		}
	}
}

// layoutFFunc is LayoutFConvertible implementation.
func (s *Sequence) layoutFFunc() func(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	return func(outsideWidth, outsideHeight float64) (screenWidth float64, screenHeight float64) {
		if l, ok := s.current.(ebiten.LayoutFer); ok {
			return l.LayoutF(outsideWidth, outsideHeight)
		}

		wi, hi := s.current.Layout(int(outsideWidth), int(outsideHeight))
		return float64(wi), float64(hi)
	}
}
