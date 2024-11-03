package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Sequence provides the ability to run multiple ebiten.Games in sequence.
type Sequence struct {
	current           ebiten.Game
	transitionUpdater *transitionUpdater
	onStartCalled     bool
}

// NewSequence creates a new Sequence instance.
func NewSequence(first ebiten.Game) *Sequence {
	return &Sequence{current: first}
}

// Update is ebiten.Game implementation.
func (s *Sequence) Update() error {
	if s.inTransition() {
		if err := s.transitionUpdater.Update(); err != nil {
			return err
		}
	}

	if !s.onStartCalled {
		s.OnStart()
		s.OnArrival()
		s.onStartCalled = true
	}

	return s.current.Update()
}

// Draw is ebiten.Game implementation.
func (s *Sequence) Draw(screen *ebiten.Image) {
	s.current.Draw(screen)
	if s.inTransition() {
		s.transitionUpdater.Draw(screen)
	}
}

// Layout is ebiten.Game implementation.
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

// DrawFinalScreen is ebiten.FinalScreenDrawer implementation.
func (s *Sequence) DrawFinalScreen(screen ebiten.FinalScreen, offScreen *ebiten.Image, geoM ebiten.GeoM) {
	if f, ok := s.current.(ebiten.FinalScreenDrawer); ok {
		f.DrawFinalScreen(screen, offScreen, geoM)
	} else {
		defaultDrawFinalScreenTemporaryImplRemoveItWhenEbitengineV290Released(screen, offScreen, geoM)
	}
}

// LayoutF is ebiten.LayoutFer implementation.
func (s *Sequence) LayoutF(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	if l, ok := s.current.(ebiten.LayoutFer); ok {
		return l.LayoutF(outsideWidth, outsideHeight)
	}

	owi := int(outsideWidth)
	ohi := int(outsideHeight)

	if owi < 1 {
		owi = 1
	}
	if ohi < 1 {
		ohi = 1
	}

	wi, hi := s.current.Layout(owi, ohi)
	return float64(wi), float64(hi)
}

// OnStart is OnStarter implementation.
func (s *Sequence) OnStart() {
	callIfImpl(s.current, func(o OnStarter) { o.OnStart() })
	s.onStartCalled = true
}

// OnEnd is OnEnder implementation.
func (s *Sequence) OnEnd() {
	callIfImpl(s.current, func(o OnEnder) { o.OnEnd() })
}

// OnArrival is OnArrivaler implementation.
func (s *Sequence) OnArrival() {
	callIfImpl(s.current, func(o OnArrivaler) { o.OnArrival() })
}

// OnDeparture is OnDeparturer implementation.
func (s *Sequence) OnDeparture() {
	callIfImpl(s.current, func(o OnDeparturer) { o.OnDeparture() })
}
