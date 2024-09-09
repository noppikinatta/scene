package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Chain struct {
	current Scene
}

func NewChain(first Scene) *Chain {
	return &Chain{current: first}
}

func (c *Chain) Init() {
	c.current.Init()
}

func (c *Chain) Update() error {
	if c.current.Ended() {
		c.goToNext()
	}
	return c.current.Update()
}

func (c *Chain) Draw(screen *ebiten.Image) {
	c.current.Draw(screen)
}

func (c *Chain) goToNext() {
	s, ok := c.nextScene()
	if ok {
		c.current.Dispose()
		c.current = s
		c.current.Init()
	}
}

func (c *Chain) nextScene() (Scene, bool) {
	ns, ok := c.current.(NextScener)
	if !ok {
		return nil, false
	}

	return ns.NextScene()
}

func (c *Chain) Ended() bool {
	if !c.current.Ended() {
		return false
	}

	// if next scene is exists,
	// Chain is not ended
	if _, ok := c.nextScene(); ok {
		return false
	}

	return true
}

func (c *Chain) Dispose() {
	c.current.Dispose()
}
