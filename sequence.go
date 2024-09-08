package scene

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sequence struct {
	scenes      []Scene
	currentIdx  int
	transitions map[int]int
	ended       bool
}

func NewSequence(scenes ...Scene) *Sequence {
	s := Sequence{
		scenes:      scenes,
		currentIdx:  0,
		transitions: make(map[int]int),
	}
	return &s
}

func (s *Sequence) AddTransition(from, to Scene) error {
	fi, ok := s.indexOf(from)
	if !ok {
		return errors.New("from scene for transition is not in container")
	}

	ti, ok := s.indexOf(to)
	if !ok {
		return errors.New("to scene for transition is not in container")
	}

	s.transitions[fi] = ti
	return nil
}

func (s *Sequence) indexOf(scene Scene) (int, bool) {
	for i := range s.scenes {
		if scene == s.scenes[i] {
			return i, true
		}
	}

	return 0, false
}

func (s *Sequence) Init() {
	s.ended = false
	s.currentIdx = 0
	for i := range s.scenes {
		s.scenes[i].Init()
	}
}

func (s *Sequence) Update() error {
	current, ok := s.current()
	if !ok {
		return nil
	}

	if current.Ended() {
		s.goToNext()
	}

	return current.Update()
}

func (s *Sequence) Draw(screen *ebiten.Image) {
	current, ok := s.current()
	if !ok {
		return
	}
	current.Draw(screen)
}

func (s *Sequence) Ended() bool {
	return s.ended
}

func (s *Sequence) current() (Scene, bool) {
	return s.sceneFromIdx(s.currentIdx)
}

func (s *Sequence) sceneFromIdx(idx int) (Scene, bool) {
	if idx < 0 || idx >= len(s.scenes) {
		return nil, false
	}

	return s.scenes[idx], true
}

func (s *Sequence) goToNext() {
	prev, ok := s.current()
	if ok {
		prev.Dispose()
	}

	nextIdx := s.getNextIdx()
	next, ok := s.sceneFromIdx(nextIdx)
	if ok {
		s.currentIdx = nextIdx
		next.Init()
	} else {
		s.ended = true
	}
}

func (s *Sequence) getNextIdx() int {
	idx, ok := s.transitions[s.currentIdx]
	if ok {
		return idx
	}

	idx = s.currentIdx + 1
	return idx
}

func (s *Sequence) Dispose() {
	for i := range s.scenes {
		s.scenes[i].Dispose()
	}
}
