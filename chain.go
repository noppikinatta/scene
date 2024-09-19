package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Chain runs Scenes in sequence. The next Scene when one Scene is finished is controlled using NextScener.
type Chain struct {
	first      Scene
	current    Scene
	nextScener NextScener
}

// NewChain creates a new Chain instance.
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

// NextScener determines the next Scene.
type NextScener interface {
	// NextScene returns the next Scene. If the next Scene does not exist, nil and false are returned.
	NextScene(current Scene) (Scene, bool)
}

// SequencialNextScener is an implementation of NextScener that connects Scenes in the order of the given Scene slices.
type SequencialNextScener struct {
	Scenes []Scene
	Loop   bool
}

// NewSequencialNextScener creates a new SequencialNextScener instance.
func NewSequencialNextScener(first Scene, rest ...Scene) *SequencialNextScener {
	ss := make([]Scene, 0, len(rest)+1)
	ss = append(ss, first)
	ss = append(ss, rest...)
	return &SequencialNextScener{Scenes: ss, Loop: false}
}

// NewSequencialLoopNextScener creates a new SequencialNextScener with loop.
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

// CompositNextScener is a collection of NextScener. It executes the element's NextScene methods in sequence, returning the first Scene found as the next Scene.
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
