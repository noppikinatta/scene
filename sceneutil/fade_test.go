package sceneutil_test

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/scene/sceneutil"
)

func TestFade(t *testing.T) {
	pd := dummyProgressDrawer{}
	fade := sceneutil.NewFade(3, &pd)

	fade.Init()
	fade.Update()
	fade.Draw(nil)

	if pd.progress < 0.332 || pd.progress > 0.334 {
		t.Errorf("expected progress: 0.332~0.334, actual progress: %f", pd.progress)
	}
	if fade.Ended() {
		t.Error("falde should not end before complete fading 1/3")
	}

	fade.Update()
	fade.Draw(nil)

	if pd.progress < 0.665 || pd.progress > 0.667 {
		t.Errorf("expected progress: 0.665~0.667, actual progress: %f", pd.progress)
	}
	if fade.Ended() {
		t.Error("falde should not end before complete fading 2/3")
	}

	fade.Update()
	fade.Draw(nil)

	if pd.progress < 1 {
		t.Errorf("expected progress: 1, actual progress: %f", pd.progress)
	}
	if !fade.Ended() {
		t.Error("falde should end after complete fading 3/3")
	}

	fade.Dispose()
}

func TestWithSimpleFade(t *testing.T) {
	ds := dummyScene{drawFn: func(screen *ebiten.Image) {
		screen.Fill(color.Black)
	}}

	s := sceneutil.WithSimpleFade(&ds, 4, color.RGBA{R: 20, G: 40, B: 80, A: 255})

	img := ebiten.NewImage(10, 10)

	s.Init()

	// i for avoid inf loop
	for i := 0; i < 1000; i++ {
		err := s.Update()
		if err != nil {
			t.Fatalf("expected nil err, but got %s", err.Error())
		}

		s.Draw(img)
		if s.Ended() {
			break
		}
	}

	s.Dispose()
}

type dummyProgressDrawer struct {
	progress float64
}

func (d *dummyProgressDrawer) Draw(screen *ebiten.Image, progress float64) {
	d.progress = progress
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
