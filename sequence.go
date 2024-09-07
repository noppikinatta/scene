package scene

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sequence struct {
	scenes      []Scene
	currentIdx  int
	transitions map[int]int
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
	s.currentIdx = 0
	for i := range s.scenes {
		s.scenes[i].Init()
	}
}

func (s *Sequence) Update() error {
	if s.Ended() {
		return nil
	}

	current, ok := s.current()
	if !ok {
		return nil
	}

	return current.Update()
}

func (s *Sequence) Draw(screen *ebiten.Image) {
	current, ok := s.current()
	if !ok {
		return
	}
	current.Draw(screen)
	if current.Ended() {
		s.goToNext()
	}
}

func (s *Sequence) Ended() bool {
	return s.currentIdx >= len(s.scenes)
}

func (s *Sequence) current() (Scene, bool) {
	if s.currentIdx < 0 || s.currentIdx >= len(s.scenes) {
		return nil, false
	}

	return s.scenes[s.currentIdx], true
}

func (s *Sequence) goToNext() {
	prev, ok := s.current()
	if ok {
		prev.Dispose()
	}
	s.currentIdx = s.getNext()
	next, ok := s.current()
	if ok {
		next.Init()
	}
}

func (s *Sequence) getNext() int {
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
