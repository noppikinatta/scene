package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Chain runs Scenes in sequence. The next Scene when one Scene is finished is controlled using Flow.
type Chain struct {
	first             Scene
	current           Scene
	flow              Flow
	transitioner      Transitioner
	transitionManager *transitionManager
}

// NewChain creates a new Chain instance.
func NewChain(first Scene, flow Flow, transitioner Transitioner) *Chain {
	return &Chain{first: first, flow: flow, transitioner: transitioner, transitionManager: &transitionManager{}}
}

func (c *Chain) OnSceneStart() {
	c.flow.Init()
	c.current = c.first
	if o, ok := c.current.(OnSceneStarter); ok {
		o.OnSceneStart()
	}
}

func (c *Chain) Update() error {
	if c.current == nil {
		return nil
	}
	if err := c.updateTransition(); err != nil {
		return err
	}
	return c.current.Update()
}

func (c *Chain) updateTransition() error {
	if c.current.CanEnd() && c.transitionManager.IsIdle() {
		trn := c.transitioner.Transition(c.current)
		c.transitionManager.Start(trn)
		if o, ok := c.current.(OnTransitionStarter); ok {
			o.OnTransitionStart()
		}
	}

	if err := c.transitionManager.Update(); err != nil {
		return err
	}

	if c.transitionManager.ShouldSwitchScenes() {
		c.goToNext()
	}

	if c.transitionManager.JustCompleted() {
		if !c.CanEnd() {
			tryCall(c.current, func(o OnTransitionEnder) { o.OnTransitionEnd() })
		}
	}

	return nil
}

func (c *Chain) Draw(screen *ebiten.Image) {
	if c.current == nil {
		return
	}
	c.current.Draw(screen)
}

func (c *Chain) goToNext() {
	s, ok := c.nextScene()
	if !ok {
		return
	}

	if o, ok := c.current.(OnSceneEnder); ok {
		o.OnSceneEnd()
	}
	c.current = s
	if o, ok := c.current.(OnSceneStarter); ok {
		o.OnSceneStart()
	}
}

func (c *Chain) nextScene() (Scene, bool) {
	return c.flow.NextScene(c.current)
}

func (c *Chain) CanEnd() bool {
	if !c.current.CanEnd() {
		return false
	}

	// scenes are remaining
	if _, ok := c.nextScene(); ok {
		return false
	}

	return true
}

func (c *Chain) OnSceneEnd() {
	if c.current == nil {
		return
	}
	if o, ok := c.current.(OnSceneEnder); ok {
		o.OnSceneEnd()
	}
	c.current = nil
}

// Flow determines the next Scene.
type Flow interface {
	// Init initializes this object.
	Init()
	// NextScene returns the next Scene. If the next Scene does not exist, nil and false are returned.
	NextScene(current Scene) (Scene, bool)
}

// SequencialFlow is an implementation of Flow that connects Scenes in the order of the given Scene slices.
type SequencialFlow struct {
	startIdx int
	Scenes   []Scene
	Loop     bool
}

// NewSequencialFlow creates a new SequencialFlow instance.
func NewSequencialFlow(first Scene, rest ...Scene) *SequencialFlow {
	ss := make([]Scene, 0, len(rest)+1)
	ss = append(ss, first)
	ss = append(ss, rest...)
	return &SequencialFlow{Scenes: ss, Loop: false}
}

// NewSequencialLoopFlow creates a new SequencialFlow with loop.
func NewSequencialLoopFlow(first Scene, rest ...Scene) *SequencialFlow {
	f := NewSequencialFlow(first, rest...)
	f.Loop = true
	return f
}

func (s *SequencialFlow) Init() {
	s.startIdx = 0
}

func (s *SequencialFlow) NextScene(current Scene) (Scene, bool) {
	idx := s.indexOf(current)
	if idx < 0 {
		return nil, false
	}

	nextIdx := idx + 1
	if s.Loop && nextIdx >= len(s.Scenes) {
		nextIdx = 0
		s.startIdx = 0
	}
	if nextIdx < len(s.Scenes) {
		return s.Scenes[nextIdx], true
	}

	return nil, false
}

func (s *SequencialFlow) indexOf(scene Scene) int {
	for i := s.startIdx; i < len(s.Scenes); i++ {
		if s.Scenes[i] == scene {
			s.startIdx = i
			return i
		}
	}

	return -1
}

// CompositFlow is a collection of Flow. It executes the element's NextScene methods in sequence, returning the first Scene found as the next Scene.
type CompositFlow []Flow

func (c CompositFlow) Init() {
	for _, f := range c {
		f.Init()
	}
}

func (c CompositFlow) NextScene(current Scene) (Scene, bool) {
	for _, f := range c {
		s, ok := f.NextScene(current)
		if ok {
			return s, true
		}
	}

	return nil, false
}
