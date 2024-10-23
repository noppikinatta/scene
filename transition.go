package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Transition interface {
	Init()
	Update() error
	Draw(screen *ebiten.Image)
	Completed() bool
	CanSwitchScenes() bool
}

var NopTransition Transition = nopTransition{}

type nopTransition struct{}

func (t nopTransition) Init()                 {}
func (t nopTransition) Update() error         { return nil }
func (t nopTransition) Draw(*ebiten.Image)    {}
func (t nopTransition) Completed() bool       { return true }
func (t nopTransition) CanSwitchScenes() bool { return true }

type LinearTransition struct {
	currentFrame int
	maxFrames    int
	drawer       LinearTransitionDrawer
}

func NewLinearTransition(maxFrames int, drawer LinearTransitionDrawer) *LinearTransition {
	return &LinearTransition{maxFrames: maxFrames, drawer: drawer}
}

type LinearTransitionDrawer interface {
	Draw(screen *ebiten.Image, progress LinearTransitionProgress)
}

type LinearTransitionProgress struct {
	CurrentFrame int
	MaxFrames    int
}

func (p LinearTransitionProgress) Rate() float64 {
	return float64(p.CurrentFrame) / float64(p.MaxFrames)
}

func (p LinearTransitionProgress) Halfway() bool {
	return p.CurrentFrame == (p.MaxFrames/2 + 1)
}

func (t *LinearTransition) Init() {
	t.currentFrame = 0
}

func (t *LinearTransition) Update() error {
	if t.Completed() {
		return nil
	}
	t.currentFrame++
	return nil
}

func (t *LinearTransition) Draw(screen *ebiten.Image) {
	p := t.Progress()
	t.drawer.Draw(screen, p)
}

func (t *LinearTransition) Progress() LinearTransitionProgress {
	return LinearTransitionProgress{
		CurrentFrame: t.currentFrame,
		MaxFrames:    t.maxFrames,
	}
}

func (t *LinearTransition) Completed() bool {
	return t.currentFrame > t.maxFrames
}

func (t *LinearTransition) CanSwitchScenes() bool {
	return t.Progress().Halfway()
}

type transitionProgress struct {
	flow       *Sequence
	next       ebiten.Game
	transition Transition
	switched   bool
}

func newTransitionProgress(flow *Sequence, next ebiten.Game, transition Transition) *transitionProgress {
	return &transitionProgress{
		flow:       flow,
		next:       next,
		transition: transition,
		switched:   false,
	}
}

func (t *transitionProgress) Update() error {
	if err := t.transition.Update(); err != nil {
		return err
	}

	if t.transition.CanSwitchScenes() {
		t.switchOnce()
	}
	if t.transition.Completed() {
		t.switchOnce()
		t.flow.endTransition()
	}

	return nil
}

func (t *transitionProgress) switchOnce() {
	if t.switched {
		return
	}
	t.switched = true
	t.flow.switchScenes(t.next)
}

func (t *transitionProgress) Draw(screen *ebiten.Image) {
	t.transition.Draw(screen)
}
