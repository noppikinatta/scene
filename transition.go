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
	CurrentFrame int
	MaxFrames    int
	Drawer       LinearTransitionDrawer
}

type LinearTransitionDrawer interface {
	Draw(screen *ebiten.Image, progress float64)
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
