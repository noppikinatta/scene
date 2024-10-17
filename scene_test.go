package scene_test

import (
	"errors"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/scene"
	"github.com/noppikinatta/scene/sceneutil"
)

func TestParallelAllScenesAreProcessed(t *testing.T) {
	log := make([]string, 0)

	newScene := func(name string) *dummyScene {
		counter := 0
		return &dummyScene{
			onSceneStartFn: func() {
				log = append(log, name+":onSceneStart")
			},
			updateFn: func() error {
				log = append(log, name+":update")
				counter++
				return nil
			},
			drawFn: func(screen *ebiten.Image) {
				log = append(log, name+":draw")
			},
			onSceneEndFn: func() {
				log = append(log, name+":onSceneEnd")
			},
			canEndFn: func() bool {
				return counter >= 1
			},
		}
	}

	p := scene.NewParallel(
		newScene("1"),
		newScene("2"),
		newScene("3"),
	)

	g := scene.ToGame(p, sceneutil.SimpleLayoutFunc())
	loopGame(t, g)

	expected := []string{
		"1:onSceneStart",
		"2:onSceneStart",
		"3:onSceneStart",
		"1:update",
		"2:update",
		"3:update",
		"1:draw",
		"2:draw",
		"3:draw",
		"1:onSceneEnd",
		"2:onSceneEnd",
		"3:onSceneEnd",
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
			canEndFn: func() bool {
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
		if p.CanEnd() {
			t.Errorf("CanEnd() should be false when counter = %d", counter)
		}
	}

	counter = 3

	if !p.CanEnd() {
		t.Errorf("CanEnd() should be true when counter = %d", counter)
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
				t.Errorf("err should be nil but: %s", err.Error())
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

func TestBarrier(t *testing.T) {
	ended := false
	endedFn := func() bool {
		return ended
	}

	b := scene.NewBarrier(endedFn)

	err := b.Update()
	if err != nil {
		t.Error("Barrier should not return err")
	}
	b.Draw(nil) // can call without any panics

	if b.CanEnd() {
		t.Error("CanEnd() should return false")
	}

	ended = true

	if !b.CanEnd() {
		t.Error("CanEnd() should return true")
	}
}

func TestChain(t *testing.T) {
	cases := []struct {
		Name        string
		Flower      func(s1, s2, s3 scene.Scene) scene.Flow
		ExpectedLog []string
	}{
		{
			Name: "own-flow",
			Flower: func(s1, s2, s3 scene.Scene) scene.Flow {
				orderIdx := 0
				order := []scene.Scene{s1, s2, s1, s3, s2, s3}
				return &dummyFlow{func(current scene.Scene) (scene.Scene, bool) {
					for i := orderIdx; i < len(order)-1; i++ {
						if order[i] == current {
							orderIdx = i
							return order[i+1], true
						}
					}

					return nil, false
				}}
			},
			ExpectedLog: []string{
				"s1:onSceneStart",
				"s1:update",
				"s1:draw",
				"s1:onSceneEnd",
				"s2:onSceneStart",
				"s2:update",
				"s2:draw",
				"s2:update",
				"s2:draw",
				"s2:onSceneEnd",
				"s1:onSceneStart",
				"s1:update",
				"s1:draw",
				"s1:onSceneEnd",
				"s3:onSceneStart",
				"s3:update",
				"s3:draw",
				"s3:update",
				"s3:draw",
				"s3:update",
				"s3:draw",
				"s3:onSceneEnd",
				"s2:onSceneStart",
				"s2:update",
				"s2:draw",
				"s2:update",
				"s2:draw",
				"s2:onSceneEnd",
				"s3:onSceneStart",
				"s3:update",
				"s3:draw",
				"s3:update",
				"s3:draw",
				"s3:update",
				"s3:draw",
				"s3:onSceneEnd",
			},
		},
		{
			Name: "sequencial-normal",
			Flower: func(s1, s2, s3 scene.Scene) scene.Flow {
				return scene.NewSequencialFlow(s1, s2, s3)
			},
			ExpectedLog: []string{
				"s1:onSceneStart",
				"s1:update",
				"s1:draw",
				"s1:onSceneEnd",
				"s2:onSceneStart",
				"s2:update",
				"s2:draw",
				"s2:update",
				"s2:draw",
				"s2:onSceneEnd",
				"s3:onSceneStart",
				"s3:update",
				"s3:draw",
				"s3:update",
				"s3:draw",
				"s3:update",
				"s3:draw",
				"s3:onSceneEnd",
			},
		},
		{
			Name: "sequencial-no-inf-loop-if-same-scene-is-set-twice",
			Flower: func(s1, s2, s3 scene.Scene) scene.Flow {
				return scene.NewSequencialFlow(s1, s2, s1, s3)
			},
			ExpectedLog: []string{
				"s1:onSceneStart",
				"s1:update",
				"s1:draw",
				"s1:onSceneEnd",
				"s2:onSceneStart",
				"s2:update",
				"s2:draw",
				"s2:update",
				"s2:draw",
				"s2:onSceneEnd",
				"s1:onSceneStart",
				"s1:update",
				"s1:draw",
				"s1:onSceneEnd",
				"s3:onSceneStart",
				"s3:update",
				"s3:draw",
				"s3:update",
				"s3:draw",
				"s3:update",
				"s3:draw",
				"s3:onSceneEnd",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			logger := dummyCountSceneLogger{}

			s1 := newDummyCountScene("s1", 1, &logger)
			s2 := newDummyCountScene("s2", 2, &logger)
			s3 := newDummyCountScene("s3", 3, &logger)

			f := c.Flower(s1, s2, s3)

			tran := scene.NopTransition
			traner := scene.NewFixedTransitioner(tran)

			chain := scene.NewChain(s1, f, traner)

			g := scene.ToGame(chain, sceneutil.SimpleLayoutFunc())
			loopGame(t, g)

			expectedLogLen := len(c.ExpectedLog)
			actualLogLen := len(logger.log)

			if expectedLogLen != actualLogLen {
				t.Errorf("expected log len: %d, actual log len: %d", expectedLogLen, actualLogLen)
			}

			logLen := expectedLogLen
			if actualLogLen < expectedLogLen {
				logLen = actualLogLen
			}

			for i := 0; i < logLen; i++ {
				if c.ExpectedLog[i] != logger.log[i] {
					t.Errorf("%d: expected log: %s, actual log: %s", i, c.ExpectedLog[i], logger.log[i])
				}
			}
		})
	}
}

func TestChainNoPanic(t *testing.T) {
	cases := []struct {
		Name string
		Fn   func(c *scene.Chain)
	}{
		{
			Name: "update-without-init",
			Fn:   func(c *scene.Chain) { c.Update() },
		},
		{
			Name: "draw-without-init",
			Fn:   func(c *scene.Chain) { c.Draw(nil) },
		},
		{
			Name: "on-scene-start-twice",
			Fn: func(c *scene.Chain) {
				c.OnSceneStart()
				c.OnSceneStart()
			},
		},
		{
			Name: "on-scene-end-twice",
			Fn: func(c *scene.Chain) {
				c.OnSceneStart()
				c.OnSceneEnd()
				c.OnSceneEnd()
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			logger := dummyCountSceneLogger{}

			s1 := newDummyCountScene("s1", 1, &logger)
			s2 := newDummyCountScene("s2", 1, &logger)
			s3 := newDummyCountScene("s3", 1, &logger)

			f := scene.NewSequencialFlow(s1, s2, s3)
			traner := scene.NewFixedTransitioner(scene.NopTransition)

			chain := scene.NewChain(s1, f, traner)

			// expected no panic
			c.Fn(chain)
		})
	}
}

func TestChainSequencialLoop(t *testing.T) {
	logger := dummyCountSceneLogger{}

	s1 := newDummyCountScene("s1", 1, &logger)
	s2 := newDummyCountScene("s2", 2, &logger)
	s3 := newDummyCountScene("s3", 3, &logger)

	// make sure the game terminates when s1 started 2 times
	s1StartCount := 0
	s1StartFn := s1.onSceneStartFn
	s1.onSceneStartFn = func() {
		s1StartFn()
		s1StartCount++
	}
	s1UpdateFn := s1.updateFn
	s1.updateFn = func() error {
		if err := s1UpdateFn(); err != nil {
			return err
		}
		if s1StartCount >= 2 {
			return ebiten.Termination
		}
		return nil
	}

	f := scene.NewSequencialLoopFlow(s1, s2, s3)

	traner := scene.NewFixedTransitioner(scene.NopTransition)

	chain := scene.NewChain(s1, f, traner)
	g := scene.ToGame(chain, sceneutil.SimpleLayoutFunc())
	loopGame(t, g)

	expectedLog := []string{
		"s1:onSceneStart",
		"s1:update",
		"s1:draw",
		"s1:onSceneEnd",
		"s2:onSceneStart",
		"s2:update",
		"s2:draw",
		"s2:update",
		"s2:draw",
		"s2:onSceneEnd",
		"s3:onSceneStart",
		"s3:update",
		"s3:draw",
		"s3:update",
		"s3:draw",
		"s3:update",
		"s3:draw",
		"s3:onSceneEnd",
		"s1:onSceneStart",
		"s1:update",
		"s1:onSceneEnd",
	}

	if len(expectedLog) != len(logger.log) {
		t.Fatalf("expected logs: %v, actual logs: %v", expectedLog, logger.log)
	}

	for i := range expectedLog {
		if expectedLog[i] != logger.log[i] {
			t.Errorf("expected log: %s, actual log: %s", expectedLog[i], logger.log[i])
		}
	}
}

func TestCompositFlow(t *testing.T) {
	logger := dummyCountSceneLogger{}

	s1 := newDummyCountScene("s1", 1, &logger)
	s2 := newDummyCountScene("s2", 2, &logger)
	s3 := newDummyCountScene("s3", 3, &logger)

	f1 := dummyFlow{func(current scene.Scene) (scene.Scene, bool) {
		if current == s1 {
			return s2, true
		}

		return nil, false
	}}

	f2 := dummyFlow{func(current scene.Scene) (scene.Scene, bool) {
		if current == s3 {
			return nil, false
		}
		return s3, true
	}}

	f := scene.CompositFlow{&f1, &f2}
	traner := scene.NewFixedTransitioner(scene.NopTransition)

	chain := scene.NewChain(s1, f, traner)

	g := scene.ToGame(chain, sceneutil.SimpleLayoutFunc())
	loopGame(t, g)

	expectedLog := []string{
		"s1:onSceneStart",
		"s1:update",
		"s1:draw",
		"s1:onSceneEnd",
		"s2:onSceneStart",
		"s2:update",
		"s2:draw",
		"s2:update",
		"s2:draw",
		"s2:onSceneEnd",
		"s3:onSceneStart",
		"s3:update",
		"s3:draw",
		"s3:update",
		"s3:draw",
		"s3:update",
		"s3:draw",
		"s3:onSceneEnd",
	}

	if len(expectedLog) != len(logger.log) {
		t.Fatalf("expected logs: %v, actual logs: %v", expectedLog, logger.log)
	}

	for i := range expectedLog {
		if expectedLog[i] != logger.log[i] {
			t.Errorf("expected log: %s, actual log: %s", expectedLog[i], logger.log[i])
		}
	}
}

func TestToGame(t *testing.T) {
	logger := dummyCountSceneLogger{}
	s := newDummyCountScene("s", 3, &logger)

	l := scene.NewLayouterFromFunc(func(outsideWidth, outsideHeight int) (screenWidth int, screenHeight int) {
		return outsideWidth * 2, outsideHeight * 2
	})
	g := scene.ToGame(s, l)

	w, h := g.Layout(100, 100)
	if w != 200 || h != 200 {
		t.Errorf("expected w,h=200,200, but got %d,%d", w, h)
	}

	// calling Draw() before first Update() is OK
	g.Draw(nil)

	var err error
	for range 100 { // loop 100 times to avoid inf loop
		err = g.Update()
		if err != nil {
			break
		}
		g.Draw(nil)
	}

	if !errors.Is(err, ebiten.Termination) {
		t.Errorf("expected ebiten.Termination, but got %v", err)
	}

	expectedLog := []string{
		"s:onSceneStart",
		"s:update",
		"s:draw",
		"s:update",
		"s:draw",
		"s:update",
		"s:draw",
		"s:onSceneEnd",
	}

	if len(expectedLog) != len(logger.log) {
		t.Fatalf("expected logs: %v, actual logs: %v", expectedLog, logger.log)
	}

	for i := range expectedLog {
		if expectedLog[i] != logger.log[i] {
			t.Errorf("expected log: %s, actual log: %s", expectedLog[i], logger.log[i])
		}
	}
}

type dummyScene struct {
	onSceneStartFn      func()
	onTransitionEndFn   func()
	updateFn            func() error
	drawFn              func(screen *ebiten.Image)
	canEndFn            func() bool
	onTransitionStartFn func()
	onSceneEndFn        func()
}

func (s *dummyScene) OnSceneStart() {
	if s.onSceneStartFn == nil {
		return
	}
	s.onSceneStartFn()
}

func (s *dummyScene) OnTransitionEnd() {
	if s.onTransitionEndFn == nil {
		return
	}
	s.onTransitionEndFn()
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

func (s *dummyScene) CanEnd() bool {
	if s.canEndFn == nil {
		return false
	}
	return s.canEndFn()
}

func (s *dummyScene) OnTransitionStart() {
	if s.onTransitionStartFn == nil {
		return
	}
	s.onTransitionStartFn()
}

func (s *dummyScene) OnSceneEnd() {
	if s.onSceneEndFn == nil {
		return
	}
	s.onSceneEndFn()
}

type dummyCountSceneLogger struct {
	log []string
}

func (l *dummyCountSceneLogger) Append(log string) {
	l.log = append(l.log, log)
}

func newDummyCountScene(name string, maxCount int, logger *dummyCountSceneLogger) *dummyScene {
	counter := 0
	return &dummyScene{
		onSceneStartFn: func() {
			counter = 0
			logger.Append(name + ":onSceneStart")
		},
		updateFn: func() error {
			counter++
			logger.Append(name + ":update")
			return nil
		},
		drawFn: func(screen *ebiten.Image) {
			logger.Append(name + ":draw")
		},
		canEndFn: func() bool {
			return counter >= maxCount
		},
		onSceneEndFn: func() {
			logger.Append(name + ":onSceneEnd")
		},
	}
}

type dummyFlow struct {
	fn func(current scene.Scene) (scene.Scene, bool)
}

func (f *dummyFlow) Init() {}

func (f *dummyFlow) NextScene(current scene.Scene) (scene.Scene, bool) {
	return f.fn(current)
}

func loopGame(t *testing.T, g ebiten.Game) {
	t.Helper()

	fakeScreen := ebiten.NewImage(3, 3)
	for range 100 { // loop 100 times to avoid inf loop
		err := g.Update()

		if errors.Is(err, ebiten.Termination) {
			break
		}
		if err != nil {
			t.Fatalf("unexpected err on Game.Update(): %v", err)
		}

		g.Draw(fakeScreen)
	}
}
