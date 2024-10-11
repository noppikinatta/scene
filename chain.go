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
	tryCall(c.current, func(o OnSceneStarter) { o.OnSceneStart() })
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
		tryCall(c.current, func(o OnTransitionStarter) { o.OnTransitionStart() })
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

func (c *Chain) goToNext() {
	s, ok := c.nextScene()
	if !ok {
		return
	}

	tryCall(c.current, func(o OnSceneEnder) { o.OnSceneEnd() })
	c.current = s
	tryCall(c.current, func(o OnSceneStarter) { o.OnSceneStart() })
}

func (c *Chain) nextScene() (Scene, bool) {
	if c.current == nil {
		return nil, false
	}
	return c.flow.NextScene(c.current)
}

func (c *Chain) Draw(screen *ebiten.Image) {
	if c.current == nil {
		return
	}
	c.current.Draw(screen)
}

func (c *Chain) CanEnd() bool {
	if c.current == nil {
		return true
	}

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
	tryCall(c.current, func(o OnSceneEnder) { o.OnSceneEnd() })
	c.current = nil
}
