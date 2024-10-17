package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Transitioner interface {
	Transition(scene Scene) Transition
}

type fixedTransitioner struct {
	transition Transition
}

func NewFixedTransitioner(transition Transition) Transitioner {
	return &fixedTransitioner{transition: transition}
}

func (t *fixedTransitioner) Transition(Scene) Transition {
	return t.transition
}

type Transition interface {
	Init()
	Update() error
	Draw(screen *ebiten.Image)
	Completed() bool
	ShouldSwitchScenes() bool
}

var NopTransition Transition = nopTransition{}

type nopTransition struct{}

func (t nopTransition) Init()                    {}
func (t nopTransition) Update() error            { return nil }
func (t nopTransition) Draw(*ebiten.Image)       {}
func (t nopTransition) Completed() bool          { return true }
func (t nopTransition) ShouldSwitchScenes() bool { return true }

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

func (t *LinearTransition) ShouldSwitchScenes() bool {
	return t.Progress().Halfway()
}

type transitionManager struct {
	transition      Transition
	shouldSwitchRcd *boolRecorder
	completedRcd    *boolRecorder
}

func (m *transitionManager) shouldSwitchRecorder() *boolRecorder {
	if m.shouldSwitchRcd == nil {
		m.shouldSwitchRcd = &boolRecorder{}
	}
	return m.shouldSwitchRcd
}

func (m *transitionManager) completedRecorder() *boolRecorder {
	if m.completedRcd == nil {
		m.completedRcd = &boolRecorder{}
	}
	return m.completedRcd
}

func (m *transitionManager) Start(transition Transition) {
	m.transition = transition
	m.transition.Init()
	m.shouldSwitchRecorder().Reset()
	m.completedRecorder().Reset()
}

func (m *transitionManager) End() {
	m.transition = nil
	// do not reset shouldSwitchRecorder and completedRecorder
	// ShouldSwitch() and JustCompleted() are called after transition end
}

func (m *transitionManager) IsIdle() bool {
	return m.transition == nil
}

func (m *transitionManager) Update() error {
	shouldSwitch, completed, err := m.updateTransition()
	if err != nil {
		return err
	}

	m.shouldSwitchRecorder().Update(shouldSwitch)
	m.completedRecorder().Update(completed)

	if completed {
		m.End()
	}

	return nil
}

func (m *transitionManager) updateTransition() (shouldSwitch, completed bool, err error) {
	if m.IsIdle() {
		return false, false, nil
	}

	if err := m.transition.Update(); err != nil {
		return false, false, err
	}

	shouldSwitch = m.transition.ShouldSwitchScenes()
	completed = m.transition.Completed()

	// if transition completed before returning ShouldSwitchScenes() true,
	// shouldSwitch must return true to go to next scene
	shouldSwitch = shouldSwitch || completed

	return shouldSwitch, completed, nil
}

func (m *transitionManager) Draw(screen *ebiten.Image) {
	if m.IsIdle() {
		return
	}
	m.transition.Draw(screen)
}

func (m *transitionManager) ShouldSwitchScenes() bool {
	return m.shouldSwitchRecorder().TrueInThisFrame()
}

func (m *transitionManager) JustCompleted() bool {
	return m.completedRecorder().TrueInThisFrame()
}

type boolRecorder struct {
	value      bool
	trueInPast bool
}

func (r *boolRecorder) Reset() {
	r.value = false
	r.trueInPast = false
}

func (r *boolRecorder) Update(value bool) {
	if r.value {
		r.trueInPast = true
		return
	}
	r.value = value
}

func (r *boolRecorder) TrueInThisFrame() bool {
	return r.value && !r.trueInPast
}
