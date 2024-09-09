package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func NewSequence(scenes ...Scene) Scene {
	if len(scenes) == 0 {
		return &nopScene{}
	}

	ss := withNext(scenes...)
	return NewChain(ss[0])
}

type nopScene struct {
}

func (s *nopScene) Init() {}

func (s *nopScene) Update() error {
	return nil
}

func (s *nopScene) Draw(screen *ebiten.Image) {}

func (s *nopScene) Ended() bool {
	return true
}

func (s *nopScene) Dispose() {}

type SeqenceScene struct {
	Scene
	next Scene
}

func (s *SeqenceScene) NextScene() (Scene, bool) {
	return s.next, true
}

func withNext(scenes ...Scene) []Scene {
	ss := make([]Scene, len(scenes))
	for i := range scenes {
		if _, ok := scenes[i].(NextScener); ok || i == len(scenes)-1 {
			ss[i] = scenes[i]
		} else {
			ss[i] = &SeqenceScene{Scene: scenes[i]}
		}
	}

	for i := range ss {
		if s, ok := ss[i].(*SeqenceScene); ok {
			s.next = ss[i+1]
		}
	}

	return ss
}
