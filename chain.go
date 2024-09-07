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
	if c.Ended() {
		return nil
	}

	return c.current.Update()
}

func (c *Chain) Draw(screen *ebiten.Image) {
	c.current.Draw(screen)
	if c.current.Ended() {
		c.goToNext()
	}
}

func (c *Chain) goToNext() {
	ns, ok := c.current.(NextScener)
	if !ok {
		return
	}

	s, ok := ns.NextScene()
	if ok {
		c.current = s
		c.current.Init()
	}
}

func (c *Chain) Ended() bool {
	if !c.current.Ended() {
		return false
	}

	// if current Scene does not implement NextScener,
	// Chain is ended
	ns, ok := c.current.(NextScener)
	if !ok {
		return true
	}

	// if next scene is exists,
	// Chain is not ended
	if _, ok = ns.NextScene(); ok {
		return false
	}

	return true
}

func (c *Chain) Dispose() {
	c.current.Dispose()
}
