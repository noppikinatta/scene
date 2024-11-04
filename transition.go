package bamenn

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Transition is an interface for drawing during scene transitions.
type Transition interface {
	// Reset is called at the start of a scene transition.
	// If the same instance is to be used multiple times, Reset should be used to initialize the Transition state.
	Reset()
	// Update is called in ebiten.Game.Update to update the state of Transition.
	Update() error
	// Draw draws during scene transitions.
	// The argument screen reflects the result of drawing by ebiten.Game.Draw.
	Draw(screen *ebiten.Image)
	// Completed returns true when the scene transition is complete.
	Completed() bool
	// CanSwitchScenes returns true when ebiten.Game, which represents a scene, can be switched.
	// Even if true is returned at multiple frames during a single scene transition, the scene is switched only once.
	// If CanSwitchScenes does not return true, the scene switch will occur when Completed returns true.
	CanSwitchScenes() bool
}

// NopTransition is a Transition that does not draw anything.
var NopTransition Transition = nopTransition{}

type nopTransition struct{}

func (t nopTransition) Reset()                {}
func (t nopTransition) Update() error         { return nil }
func (t nopTransition) Draw(*ebiten.Image)    {}
func (t nopTransition) Completed() bool       { return true }
func (t nopTransition) CanSwitchScenes() bool { return true }

// LinearTransition is a Transition that transitions linearly for a specified number of frames.
type LinearTransition struct {
	currentFrame  int
	frameToSwitch int
	maxFrames     int
	drawer        LinearTransitionDrawer
}

// NewLinearTransition returns a new LinearTransition.
func NewLinearTransition(frameToSwitch, maxFrames int, drawer LinearTransitionDrawer) *LinearTransition {
	return &LinearTransition{frameToSwitch: frameToSwitch, maxFrames: maxFrames, drawer: drawer}
}

// LinearTransitionDrawer is an interface that draws as the LinearTransition progresses.
type LinearTransitionDrawer interface {
	// Draw draws as the LinearTransition progresses.
	Draw(screen *ebiten.Image, progress LinearTransitionProgress)
}

// LinearTransitionProgress represents the progress of LinearTransition.
type LinearTransitionProgress struct {
	CurrentFrame  int  // CurrentFrame is the current frame.
	MaxFrames     int  // MaxFrames is the maximum number of frames.
	FrameToSwitch bool // FrameToSwitch returns true if a scene change occurs in this frame.
}

// Rate returns the progress rate of LinearTransition in the range of 0.0~1.0.
func (p LinearTransitionProgress) Rate() float64 {
	return float64(p.CurrentFrame) / float64(p.MaxFrames)
}

// Reset is called at the start of a scene transition.
// If the same instance is to be used multiple times, Reset should be used to initialize the Transition state.
func (t *LinearTransition) Reset() {
	t.currentFrame = 0
}

// Update is called in ebiten.Game.Update to update the state of Transition.
func (t *LinearTransition) Update() error {
	if t.Completed() {
		return nil
	}
	t.currentFrame++
	return nil
}

// Draw draws during scene transitions.
// The argument screen reflects the result of drawing by ebiten.Game.Draw.
func (t *LinearTransition) Draw(screen *ebiten.Image) {
	p := t.Progress()
	t.drawer.Draw(screen, p)
}

// Progress returns the progress of it.
func (t *LinearTransition) Progress() LinearTransitionProgress {
	return LinearTransitionProgress{
		FrameToSwitch: t.CanSwitchScenes(),
		CurrentFrame:  t.currentFrame,
		MaxFrames:     t.maxFrames,
	}
}

// Completed returns true when the scene transition is complete.
func (t *LinearTransition) Completed() bool {
	return t.currentFrame > t.maxFrames
}

// CanSwitchScenes returns true when ebiten.Game, which represents a scene, can be switched.
// Even if true is returned at multiple frames during a single scene transition, the scene is switched only once.
// If CanSwitchScenes does not return true, the scene switch will occur when Completed returns true.
func (t *LinearTransition) CanSwitchScenes() bool {
	return t.currentFrame == t.frameToSwitch
}

type transitionUpdater struct {
	seq        *Sequence
	next       ebiten.Game
	transition Transition
	switched   bool
}

func newTransitionUpdater(seq *Sequence, next ebiten.Game, transition Transition) *transitionUpdater {
	return &transitionUpdater{
		seq:        seq,
		next:       next,
		transition: transition,
		switched:   false,
	}
}

func (t *transitionUpdater) Update() error {
	if err := t.transition.Update(); err != nil {
		return err
	}

	if t.transition.CanSwitchScenes() {
		t.switchOnce()
	}
	if t.transition.Completed() {
		t.switchOnce()
		t.seq.endTransition()
	}

	return nil
}

func (t *transitionUpdater) switchOnce() {
	if t.switched {
		return
	}
	t.switched = true
	t.seq.switchScenes(t.next)
}

func (t *transitionUpdater) Draw(screen *ebiten.Image) {
	t.transition.Draw(screen)
}
