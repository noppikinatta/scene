package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Chain struct {
	first      Scene
	current    Scene
	nextScener NextScener
}

func NewChain(first Scene, nextScener NextScener) *Chain {
	return &Chain{first: first, nextScener: nextScener}
}

func (c *Chain) Init() {
	if c.current != nil {
		c.current.Dispose()
	}
	c.current = c.first
	c.current.Init()
}

func (c *Chain) Update() error {
	if c.current == nil {
		return nil
	}
	if c.current.Ended() {
		c.goToNext()
	}
	return c.current.Update()
}

func (c *Chain) Draw(screen *ebiten.Image) {
	if c.current == nil {
		return
	}
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
	return c.nextScener.NextScene(c.current)
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

type NextScener interface {
	NextScene(current Scene) (Scene, bool)
}

type SequencialNextScener struct {
	Scenes []Scene
	Loop   bool
}

func NewSequencialNextScener(first Scene, rest ...Scene) *SequencialNextScener {
	ss := make([]Scene, 0, len(rest)+1)
	ss = append(ss, first)
	ss = append(ss, rest...)
	return &SequencialNextScener{Scenes: ss, Loop: false}
}

func NewSequencialLoopNextScener(first Scene, rest ...Scene) *SequencialNextScener {
	ns := NewSequencialNextScener(first, rest...)
	ns.Loop = true
	return ns
}

func (s *SequencialNextScener) NextScene(current Scene) (Scene, bool) {
	idx := s.indexOf(current)
	if idx < 0 {
		return nil, false
	}

	nextIdx := idx + 1
	if s.Loop {
		nextIdx = nextIdx % len(s.Scenes)
	}
	if nextIdx < len(s.Scenes) {
		return s.Scenes[nextIdx], true
	}

	return nil, false
}

func (s *SequencialNextScener) indexOf(scene Scene) int {
	for i := range s.Scenes {
		if s.Scenes[i] == scene {
			return i
		}
	}

	return -1
}

type CompositNextScener []NextScener

func (c CompositNextScener) NextScene(current Scene) (Scene, bool) {
	for _, ns := range c {
		s, ok := ns.NextScene(current)
		if ok {
			return s, true
		}
	}

	return nil, false
}
