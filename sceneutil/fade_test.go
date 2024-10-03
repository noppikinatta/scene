package sceneutil_test

import (
	"image/color"
	"os"
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
	ds := dummyScene{
		drawFn: func(screen *ebiten.Image) {
			screen.Fill(color.Black)
		},
		endedFn: func() bool {
			return true
		},
	}

	s := sceneutil.WithSimpleFade(&ds, 4, color.RGBA{R: 20, G: 40, B: 80, A: 255})

	img := ebiten.NewImage(3, 3)

	s.Init()

	colors := make([]color.Color, 0)

	// i for avoid inf loop
	for i := 0; i < 1000; i++ {
		err := s.Update()
		if err != nil {
			t.Fatalf("expected nil err, but got %s", err.Error())
		}

		s.Draw(img)
		colors = append(colors, img.At(1, 1)) // get center pixel
		if s.Ended() {
			break
		}
	}

	expectedColors := []color.Color{
		color.RGBA{R: 15, G: 30, B: 60, A: 255}, // 0.25
		color.RGBA{R: 10, G: 20, B: 40, A: 255}, // 0.50
		color.RGBA{R: 5, G: 10, B: 20, A: 255},  // 0.75
		color.Black,                             // 1.00
		color.Black,
		color.RGBA{R: 5, G: 10, B: 20, A: 255},  // 0.25
		color.RGBA{R: 10, G: 20, B: 40, A: 255}, // 0.50
		color.RGBA{R: 15, G: 30, B: 60, A: 255}, // 0.75
		color.RGBA{R: 20, G: 40, B: 80, A: 255}, // 1.00
	}

	if len(expectedColors) != len(colors) {
		t.Fatalf("expected color len %d but got %d", len(expectedColors), len(colors))
	}

	rgba := func(c color.Color) color.Color {
		r, g, b, a := c.RGBA()
		return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	}

	for i := range colors {
		ec := rgba(expectedColors[i])
		ac := rgba(colors[i])

		if ec != ac {
			t.Errorf("%d:expected color %v but got %v", i, ec, ac)
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

// from: https://github.com/hajimehoshi/ebiten/blob/main/internal/testing/testing.go
type game struct {
	m    *testing.M
	code int
}

func (g *game) Update() error {
	g.code = g.m.Run()
	return ebiten.Termination
}

func (*game) Draw(*ebiten.Image) {
}

func (*game) Layout(int, int) (int, int) {
	return 320, 240
}

func MainWithRunLoop(m *testing.M) {
	// Run an Ebiten process so that (*Image).At is available.
	g := &game{
		m:    m,
		code: 1,
	}
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
	if g.code != 0 {
		os.Exit(g.code)
	}
}

func TestMain(m *testing.M) {
	MainWithRunLoop(m)
}
