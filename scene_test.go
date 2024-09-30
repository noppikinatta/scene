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

	b.Init() // can call without any panics
	err := b.Update()
	if err != nil {
		t.Error("Barrier should not return err")
	}
	b.Draw(nil) // can call without any panics
	b.Dispose() // can call without any panics

	if b.Ended() {
		t.Error("Ended should return false")
	}

	ended = true

	if !b.Ended() {
		t.Error("Ended should return true")
	}
}

func TestChain(t *testing.T) {
	logger := dummyCountSceneLogger{}

	s1 := newDummyCountScene("s1", 1, &logger)
	s2 := newDummyCountScene("s2", 2, &logger)
	s3 := newDummyCountScene("s3", 3, &logger)

	orderIdx := 0
	order := []scene.Scene{s1, s2, s1, s3, s2, s3}
	f := dummyFlow{func(current scene.Scene) (scene.Scene, bool) {
		for i := orderIdx; i < len(order)-1; i++ {
			if order[i] == current {
				orderIdx = i
				return order[i+1], true
			}
		}

		return nil, false
	}}

	chain := scene.NewChain(order[0], &f)

	chain.Init()

	// i for avoid inf loop
	for i := 0; i < 1000; i++ {
		if err := chain.Update(); err != nil {
			t.Fatalf("err in Update(): %s", err.Error())
		}
		chain.Draw(nil)
		if chain.Ended() {
			break
		}
	}

	chain.Dispose()

	expectedLog := []string{
		"s1:init",
		"s1:update",
		"s1:draw",
		"s1:dispose",
		"s2:init",
		"s2:update",
		"s2:draw",
		"s2:update",
		"s2:draw",
		"s2:dispose",
		"s1:init",
		"s1:update",
		"s1:draw",
		"s1:dispose",
		"s3:init",
		"s3:update",
		"s3:draw",
		"s3:update",
		"s3:draw",
		"s3:update",
		"s3:draw",
		"s3:dispose",
		"s2:init",
		"s2:update",
		"s2:draw",
		"s2:update",
		"s2:draw",
		"s2:dispose",
		"s3:init",
		"s3:update",
		"s3:draw",
		"s3:update",
		"s3:draw",
		"s3:update",
		"s3:draw",
		"s3:dispose",
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

func TestNoPanicEvenIfCurrentIsNil(t *testing.T) {
	logger := dummyCountSceneLogger{}

	s1 := newDummyCountScene("s1", 1, &logger)
	s2 := newDummyCountScene("s2", 1, &logger)
	s3 := newDummyCountScene("s3", 1, &logger)

	f := scene.NewSequencialFlow(s1, s2, s3)

	chain := scene.NewChain(s1, f)

	// no panic
	chain.Update()
	chain.Draw(nil)
}

func TestChainNest(t *testing.T) {
	logger := dummyCountSceneLogger{}

	s1 := newDummyCountScene("s1", 1, &logger)
	s2 := newDummyCountScene("s2", 1, &logger)
	ns1 := scene.NewSequencialFlow(s1, s2)
	c1 := scene.NewChain(s1, ns1)

	s3 := newDummyCountScene("s3", 1, &logger)
	s4 := newDummyCountScene("s4", 1, &logger)
	ns2 := scene.NewSequencialFlow(s3, s4)
	c2 := scene.NewChain(s3, ns2)

	f := scene.NewSequencialFlow(c1, c2, c1, c2)
	chain := scene.NewChain(c1, f)

	chain.Init()

	// i for avoid inf loop
	for i := 0; i < 1000; i++ {
		if err := chain.Update(); err != nil {
			t.Fatalf("err in Update(): %s", err.Error())
		}
		chain.Draw(nil)
		if chain.Ended() {
			break
		}
	}

	chain.Dispose()

	expectedLog := []string{
		"s1:init",
		"s1:update",
		"s1:draw",
		"s1:dispose",
		"s2:init",
		"s2:update",
		"s2:draw",
		"s2:dispose",
		"s3:init",
		"s3:update",
		"s3:draw",
		"s3:dispose",
		"s4:init",
		"s4:update",
		"s4:draw",
		"s4:dispose",
		"s1:init",
		"s1:update",
		"s1:draw",
		"s1:dispose",
		"s2:init",
		"s2:update",
		"s2:draw",
		"s2:dispose",
		"s3:init",
		"s3:update",
		"s3:draw",
		"s3:dispose",
		"s4:init",
		"s4:update",
		"s4:draw",
		"s4:dispose",
	}

	if len(expectedLog) != len(logger.log) {
		t.Errorf("expected len:%d, actual len: %d", len(expectedLog), len(logger.log))
		t.Fatalf("expected logs: %v, actual logs: %v", expectedLog, logger.log)
	}

	for i := range expectedLog {
		if expectedLog[i] != logger.log[i] {
			t.Errorf("expected log: %s, actual log: %s", expectedLog[i], logger.log[i])
		}
	}
}

func TestChainSequencial(t *testing.T) {
	logger := dummyCountSceneLogger{}

	s1 := newDummyCountScene("s1", 1, &logger)
	s2 := newDummyCountScene("s2", 2, &logger)
	s3 := newDummyCountScene("s3", 3, &logger)

	f := scene.NewSequencialFlow(s1, s2, s3)

	chain := scene.NewChain(s1, f)
	chain.Init()

	// i for avoid inf loop
	for i := 0; i < 1000; i++ {
		if err := chain.Update(); err != nil {
			t.Fatalf("err in Update(): %s", err.Error())
		}
		chain.Draw(nil)
		if chain.Ended() {
			break
		}
	}

	chain.Dispose()

	expectedLog := []string{
		"s1:init",
		"s1:update",
		"s1:draw",
		"s1:dispose",
		"s2:init",
		"s2:update",
		"s2:draw",
		"s2:update",
		"s2:draw",
		"s2:dispose",
		"s3:init",
		"s3:update",
		"s3:draw",
		"s3:update",
		"s3:draw",
		"s3:update",
		"s3:draw",
		"s3:dispose",
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

func TestChainequencialLoop(t *testing.T) {
	logger := dummyCountSceneLogger{}

	s1 := newDummyCountScene("s1", 1, &logger)
	s2 := newDummyCountScene("s2", 2, &logger)
	s3 := newDummyCountScene("s3", 3, &logger)

	f := scene.NewSequencialLoopFlow(s1, s2, s3)

	chain := scene.NewChain(s1, f)
	chain.Init()

	// i for avoid inf loop
	for i := 0; i < 1000; i++ {
		if err := chain.Update(); err != nil {
			t.Fatalf("err in Update(): %s", err.Error())
		}
		chain.Draw(nil)
		if chain.Ended() {
			break
		}

		s1InitCount := 0
		for _, l := range logger.log {
			if l == "s1:init" {
				s1InitCount++
			}
		}
		if s1InitCount >= 2 {
			break
		}
	}

	chain.Dispose()

	expectedLog := []string{
		"s1:init",
		"s1:update",
		"s1:draw",
		"s1:dispose",
		"s2:init",
		"s2:update",
		"s2:draw",
		"s2:update",
		"s2:draw",
		"s2:dispose",
		"s3:init",
		"s3:update",
		"s3:draw",
		"s3:update",
		"s3:draw",
		"s3:update",
		"s3:draw",
		"s3:dispose",
		"s1:init",
		"s1:update",
		"s1:draw",
		"s1:dispose",
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

	chain := scene.NewChain(s1, f)

	chain.Init()

	// i for avoid inf loop
	for i := 0; i < 1000; i++ {
		if err := chain.Update(); err != nil {
			t.Fatalf("err in Update(): %s", err.Error())
		}
		chain.Draw(nil)
		if chain.Ended() {
			break
		}
	}

	chain.Dispose()

	expectedLog := []string{
		"s1:init",
		"s1:update",
		"s1:draw",
		"s1:dispose",
		"s2:init",
		"s2:update",
		"s2:draw",
		"s2:update",
		"s2:draw",
		"s2:dispose",
		"s3:init",
		"s3:update",
		"s3:draw",
		"s3:update",
		"s3:draw",
		"s3:update",
		"s3:draw",
		"s3:dispose",
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

	g := scene.ToGame(s, func(outsideWidth, outsideHeight int) (screenWidth int, screenHeight int) {
		return outsideWidth, outsideHeight
	})

	var err error
	for i := 0; i < 10; i++ { // loop 10 times to avoid inf loop
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
		"s:init",
		"s:update",
		"s:draw",
		"s:update",
		"s:draw",
		"s:update",
		"s:dispose",
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

type dummyCountSceneLogger struct {
	log []string
}

func (l *dummyCountSceneLogger) Append(log string) {
	l.log = append(l.log, log)
}

func newDummyCountScene(name string, maxCount int, logger *dummyCountSceneLogger) *dummyScene {
	counter := 0
	return &dummyScene{
		initFn: func() {
			counter = 0
			logger.Append(name + ":init")
		},
		updateFn: func() error {
			counter++
			logger.Append(name + ":update")
			return nil
		},
		drawFn: func(screen *ebiten.Image) {
			logger.Append(name + ":draw")
		},
		endedFn: func() bool {
			return counter >= maxCount
		},
		disposeFn: func() {
			logger.Append(name + ":dispose")
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
