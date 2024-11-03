package scene_test

import (
	"errors"
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/scene"
)

func TestParallelProcessesAllGames(t *testing.T) {
	cases := []struct {
		Name        string
		Fn          func(p *scene.Parallel)
		ExpectedLog []string
	}{
		{
			Name:        "update",
			Fn:          func(p *scene.Parallel) { p.Update() },
			ExpectedLog: []string{"s1:update", "s2:update", "s3:update"},
		},
		{
			Name:        "draw",
			Fn:          func(p *scene.Parallel) { p.Draw(nil) },
			ExpectedLog: []string{"s1:draw", "s2:draw", "s3:draw"},
		},
		{
			Name:        "layout",
			Fn:          func(p *scene.Parallel) { p.Layout(1, 1) },
			ExpectedLog: []string{"s1:layout", "s2:layout", "s3:layout"},
		},
		{
			Name:        "layoutf",
			Fn:          func(p *scene.Parallel) { p.LayoutF(1, 1) },
			ExpectedLog: []string{"s1:layout", "s2:layout", "s3:layoutf"},
		},
		{
			Name:        "onstart",
			Fn:          func(p *scene.Parallel) { p.OnStart() },
			ExpectedLog: []string{"s1:onstart", "s2:onstart"},
		},
		{
			Name:        "onarrival",
			Fn:          func(p *scene.Parallel) { p.OnArrival() },
			ExpectedLog: []string{"s1:onarrival", "s2:onarrival"},
		},
		{
			Name:        "ondeparture",
			Fn:          func(p *scene.Parallel) { p.OnDeparture() },
			ExpectedLog: []string{"s1:ondeparture", "s2:ondeparture"},
		},
		{
			Name:        "onend",
			Fn:          func(p *scene.Parallel) { p.OnEnd() },
			ExpectedLog: []string{"s1:onend", "s2:onend"},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			recorder := recorder{}

			s1 := eventsForTest{gameForTest: gameForTest{Name: "s1", Recorder: &recorder}}
			s2 := eventsForTest{gameForTest: gameForTest{Name: "s2", Recorder: &recorder}}
			s3 := layoutFerForTest{gameForTest: gameForTest{Name: "s3", Recorder: &recorder}}

			p := scene.NewParallel(&s1, &s2, &s3)
			c.Fn(p)

			compareLogs(t, c.ExpectedLog, recorder.Log)
		})
	}
}

func TestParallelMergesErrors(t *testing.T) {
	err1 := errors.New("err1")
	err3 := errors.New("err3")

	s1 := gameForTest{UpdateFn: func() error { return err1 }}
	s2 := gameForTest{UpdateFn: func() error { return nil }}
	s3 := gameForTest{UpdateFn: func() error { return err3 }}

	p := scene.NewParallel(&s1, &s2, &s3)
	err := p.Update()

	if !errors.Is(err, err1) {
		t.Errorf("err should contain %v", err1)
	}
	if !errors.Is(err, err3) {
		t.Errorf("err should contain %v", err3)
	}
}

func TestParallelLayoutReturnValue(t *testing.T) {
	s1 := gameForTest{LayoutW: 7, LayoutH: 2}
	s2 := gameForTest{LayoutW: 6, LayoutH: 9}
	s3 := gameForTest{LayoutW: 5, LayoutH: 8}

	p := scene.NewParallel(&s1, &s2, &s3)

	w, h := p.Layout(1, 1)
	if w != 7 || h != 9 {
		t.Errorf("expected (7,9), but got (%d,%d)", w, h)
	}
}

func TestParallelLayoutFReturnValue(t *testing.T) {
	cases := []struct {
		Name      string
		Games     []ebiten.Game
		ExpectedW float64
		ExpectedH float64
	}{
		{
			Name: "all-without-layoutf",
			Games: []ebiten.Game{
				&gameForTest{LayoutW: 3, LayoutH: 3},
				&gameForTest{LayoutW: 1, LayoutH: 2},
				&gameForTest{LayoutW: 2, LayoutH: 1},
			},
			ExpectedW: 3,
			ExpectedH: 3,
		},
		{
			Name: "all-with-layoutf",
			Games: []ebiten.Game{
				&layoutFerForTest{layoutFW: 4, layoutFH: 9, gameForTest: gameForTest{LayoutW: 3, LayoutH: 3}},
				&layoutFerForTest{layoutFW: 7, layoutFH: 8, gameForTest: gameForTest{LayoutW: 3, LayoutH: 3}},
				&layoutFerForTest{layoutFW: 6, layoutFH: 5, gameForTest: gameForTest{LayoutW: 3, LayoutH: 3}},
			},
			ExpectedW: 7,
			ExpectedH: 9,
		},
		{
			Name: "parcial-with-layoutf-1",
			Games: []ebiten.Game{
				&layoutFerForTest{layoutFW: 4, layoutFH: 9, gameForTest: gameForTest{LayoutW: 3, LayoutH: 3}},
				&gameForTest{LayoutW: 5, LayoutH: 7},
				&layoutFerForTest{layoutFW: 6, layoutFH: 5, gameForTest: gameForTest{LayoutW: 3, LayoutH: 3}},
			},
			ExpectedW: 6,
			ExpectedH: 9,
		},
		{
			Name: "parcial-with-layoutf-2",
			Games: []ebiten.Game{
				&layoutFerForTest{layoutFW: 4, layoutFH: 9, gameForTest: gameForTest{LayoutW: 3, LayoutH: 3}},
				&gameForTest{LayoutW: 10, LayoutH: 20},
				&layoutFerForTest{layoutFW: 6, layoutFH: 5, gameForTest: gameForTest{LayoutW: 3, LayoutH: 3}},
			},
			ExpectedW: 10,
			ExpectedH: 20,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			p := scene.NewParallel(c.Games...)

			w, h := p.LayoutF(1, 1)
			if w != c.ExpectedW || h != c.ExpectedH {
				t.Errorf("expected (%f,%f), but got (%f,%f)", c.ExpectedW, c.ExpectedH, w, h)
			}
		})
	}
}

func TestParallelDrawFinalScreen(t *testing.T) {
	cases := []struct {
		Name          string
		Games         []ebiten.Game
		ExpectedColor color.Color
	}{
		{
			Name: "all-without-finalscreendrawer",
			Games: []ebiten.Game{
				&gameForTest{},
				&gameForTest{},
				&gameForTest{},
			},
			ExpectedColor: color.White,
		},
		{
			Name: "all-with-finalscreendrawer",
			Games: []ebiten.Game{
				&finalScreenDrawerForTest{
					gameForTest: gameForTest{},
					drawFn: func(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
						offscreen.Fill(color.RGBA{R: 255, A: 255})
						screen.DrawImage(offscreen, &ebiten.DrawImageOptions{GeoM: geoM})
					},
				},
				&finalScreenDrawerForTest{
					gameForTest: gameForTest{},
					drawFn: func(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
						offscreen.Fill(color.RGBA{G: 255, A: 255})
						screen.DrawImage(offscreen, &ebiten.DrawImageOptions{GeoM: geoM})
					},
				},
				&finalScreenDrawerForTest{
					gameForTest: gameForTest{},
					drawFn: func(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
						offscreen.Fill(color.RGBA{B: 255, A: 255})
						screen.DrawImage(offscreen, &ebiten.DrawImageOptions{GeoM: geoM})
					},
				},
			},
			ExpectedColor: color.RGBA{R: 255, A: 255},
		},
		{
			Name: "partial-with-finalscreendrawer",
			Games: []ebiten.Game{
				&gameForTest{},
				&gameForTest{},
				&finalScreenDrawerForTest{
					gameForTest: gameForTest{},
					drawFn: func(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
						offscreen.Fill(color.RGBA{B: 255, A: 255})
						screen.DrawImage(offscreen, &ebiten.DrawImageOptions{GeoM: geoM})
					},
				},
			},
			ExpectedColor: color.RGBA{B: 255, A: 255},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			screen := ebiten.NewImage(3, 3)
			finalScreen := ebiten.NewImage(3, 3)

			screen.Fill(color.White)

			p := scene.NewParallel(c.Games...)
			p.DrawFinalScreen(finalScreen, screen, ebiten.GeoM{})

			clr := finalScreen.At(1, 1)
			ar, ag, ab, aa := clr.RGBA()
			er, eg, eb, ea := c.ExpectedColor.RGBA()
			if ar != er || ag != eg || ab != eb || aa != ea {
				t.Errorf("expected RGBA=%d,%d,%d,%d, but got %d,%d,%d,%d", er, eg, eb, ea, ar, ag, ab, aa)
			}
		})
	}
}
