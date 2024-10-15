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
	transition           Transition
	shouldSwitchRecorder *boolRecorder
	completedRecorder    *boolRecorder
}

func (m *transitionManager) initRecordersIfNil() {
	if m.shouldSwitchRecorder == nil {
		m.shouldSwitchRecorder = &boolRecorder{}
	}
	if m.completedRecorder == nil {
		m.completedRecorder = &boolRecorder{}
	}
}

func (m *transitionManager) Start(transition Transition) {
	m.transition = transition
	m.transition.Init()
	m.initRecordersIfNil()
	m.shouldSwitchRecorder.Reset()
	m.completedRecorder.Reset()
}

func (m *transitionManager) Update() error {
	if m.IsIdle() {
		return nil
	}
	if m.transition.Completed() {
		m.transition = nil
		return nil
	}

	if err := m.transition.Update(); err != nil {
		return err
	}

	m.initRecordersIfNil()
	m.shouldSwitchRecorder.Update(m.transition.ShouldSwitchScenes())
	m.completedRecorder.Update(m.transition.Completed())

	return nil
}

func (m *transitionManager) Draw(screen *ebiten.Image) {
	if m.IsIdle() {
		return
	}
	m.transition.Draw(screen)
}

func (m *transitionManager) ShouldSwitchScenes() bool {
	if m.IsIdle() {
		return false
	}
	m.initRecordersIfNil()
	return m.shouldSwitchRecorder.TrueInThisFrame()
}

func (m *transitionManager) JustCompleted() bool {
	if m.IsIdle() {
		return false
	}
	m.initRecordersIfNil()
	return m.completedRecorder.TrueInThisFrame()
}

func (m *transitionManager) IsIdle() bool {
	return m.transition == nil
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
