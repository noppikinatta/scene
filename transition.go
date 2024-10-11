package scene

import "github.com/hajimehoshi/ebiten/v2"

type Transitioner interface {
	Transition(scene Scene) Transition
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
	currentFrame       int
	maxFrames          int
	shouldSwitchScenes bool
	drawer             LinearTransitionDrawer
}

func NewLinearTransition(maxFrames int, drawer LinearTransitionDrawer) *LinearTransition {
	return &LinearTransition{maxFrames: maxFrames, drawer: drawer}
}

type LinearTransitionDrawer interface {
	Draw(screen *ebiten.Image, progress float64) bool
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
	t.shouldSwitchScenes = t.drawer.Draw(screen, p)
}

func (t *LinearTransition) Progress() float64 {
	p := float64(t.currentFrame) / float64(t.maxFrames)
	if p > 1 {
		p = 1
	}
	return p
}

func (t *LinearTransition) Completed() bool {
	return t.currentFrame > t.maxFrames
}

func (t *LinearTransition) ShouldSwitchScenes() bool {
	return t.shouldSwitchScenes
}

type transitionManager struct {
	transition             Transition
	shouldSwitchSceneCount int
	completedCount         int
}

func (m *transitionManager) Start(transition Transition) {
	m.transition = transition
	m.transition.Init()
	m.shouldSwitchSceneCount = 0
	m.completedCount = 0
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
	if m.transition.ShouldSwitchScenes() {
		m.shouldSwitchSceneCount++
	}
	if m.transition.Completed() {
		m.completedCount++
	}

	return nil
}

func (m *transitionManager) Draw(screen *ebiten.Image) {
	if m.IsIdle() {
		return
	}
	m.Draw(screen)
}

func (m *transitionManager) ShouldSwitchScenes() bool {
	if m.IsIdle() {
		return false
	}
	return m.shouldSwitchSceneCount == 1
}

func (m *transitionManager) JustCompleted() bool {
	if m.IsIdle() {
		return false
	}
	return m.completedCount == 1
}

func (m *transitionManager) IsIdle() bool {
	return m.transition == nil
}
