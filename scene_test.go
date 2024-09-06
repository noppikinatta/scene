package scene_test

import (
	"errors"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/scene"
)

func TestParallelAllScenesAreProcessed(t *testing.T) {
	log := make([]string, 0)

	newScene := func(name string) *dummyScene {
		return &dummyScene{
			initFn: func() {
				log = append(log, name+":init")
			},
			updateFn: func() error {
				log = append(log, name+":update")
				return nil
			},
			drawFn: func(screen *ebiten.Image) {
				log = append(log, name+":draw")
			},
			disposeFn: func() {
				log = append(log, name+":dispose")
			},
		}
	}

	p := scene.NewParallel(
		newScene("1"),
		newScene("2"),
		newScene("3"),
	)

	p.Init()
	p.Update()
	p.Draw(nil)
	p.Dispose()

	expected := []string{
		"1:init",
		"2:init",
		"3:init",
		"1:update",
		"2:update",
		"3:update",
		"1:draw",
		"2:draw",
		"3:draw",
		"1:dispose",
		"2:dispose",
		"3:dispose",
	}

	max := len(expected)
	if max > len(log) {
		max = len(log)
	}

	for i := 0; i < max; i++ {
		e := expected[i]
		a := log[i]

		if e != a {
			t.Errorf("%d: elements are different: expected:%s / actual:%s", i, e, a)
		}
	}

	if len(expected) > len(log) {
		t.Fatal("some methods in Parallel are not called")
	}
	if len(expected) < len(log) {
		t.Fatal("too many methods in Parallel are called")
	}
}

func TestParallelEnded(t *testing.T) {
	counter := 0

	newScene := func(max int) *dummyScene {
		return &dummyScene{
			endedFn: func() bool {
				return counter >= max
			},
		}
	}

	scenes := []scene.Scene{
		newScene(1),
		newScene(2),
		newScene(3),
	}

	p := scene.NewParallel(scenes...)

	for i := 0; i < 3; i++ {
		counter = i
		if p.Ended() {
			t.Errorf("Ended() should be false when counter = %d", counter)
		}
	}

	counter = 3

	if !p.Ended() {
		t.Errorf("Ended() should be true when counter = %d", counter)
	}
}

func TestParallelUpdateErrorMerged(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")

	newScene := func(err error) *dummyScene {
		return &dummyScene{
			updateFn: func() error {
				return err
			},
		}
	}

	cases := []struct {
		Name           string
		ReturnsErr1    bool
		ReturnsErr2    bool
		ErrShouldBeNil bool
	}{
		{
			Name:           "no-err",
			ErrShouldBeNil: true,
		},
		{
			Name:        "err1",
			ReturnsErr1: true,
			ReturnsErr2: false,
		},
		{
			Name:        "err2",
			ReturnsErr1: false,
			ReturnsErr2: true,
		},
		{
			Name:        "err1-and-err2",
			ReturnsErr1: true,
			ReturnsErr2: true,
		},
	}

	for _, c := range cases {
		var s1, s2 scene.Scene

		if c.ReturnsErr1 {
			s1 = newScene(err1)
		} else {
			s1 = newScene(nil)
		}
		if c.ReturnsErr2 {
			s2 = newScene(err2)
		} else {
			s2 = newScene(nil)
		}

		p := scene.NewParallel(s1, s2)

		err := p.Update()

		if c.ErrShouldBeNil {
			if err != nil {
				t.Error("err should be nil")
			}
		} else {
			if errors.Is(err, err1) != c.ReturnsErr1 {
				t.Errorf(
					"errors.Is %v should be %t: actual err is %v",
					err1,
					c.ReturnsErr1,
					err,
				)
			}
			if errors.Is(err, err2) != c.ReturnsErr2 {
				t.Errorf(
					"errors.Is %v should be %t: actual err is %v",
					err2,
					c.ReturnsErr2,
					err,
				)
			}
		}
	}
}

type dummyScene struct {
	initFn    func()
	updateFn  func() error
	drawFn    func(screen *ebiten.Image)
	endedFn   func() bool
	disposeFn func()
}

func (s *dummyScene) Init() {
	if s.initFn == nil {
		return
	}
	s.initFn()
}

func (s *dummyScene) Update() error {
	if s.updateFn == nil {
		return nil
	}
	return s.updateFn()
}

func (s *dummyScene) Draw(screen *ebiten.Image) {
	if s.drawFn == nil {
		return
	}
	s.drawFn(screen)
}

func (s *dummyScene) Ended() bool {
	if s.endedFn == nil {
		return false
	}
	return s.endedFn()
}

func (s *dummyScene) Dispose() {
	if s.disposeFn == nil {
		return
	}
	s.disposeFn()
}
