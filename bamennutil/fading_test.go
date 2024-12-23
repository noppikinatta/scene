package bamennutil_test

import (
	"errors"
	"fmt"
	"image/color"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/bamenn"
	"github.com/noppikinatta/bamenn/bamennutil"
)

func TestLinearFillFadingDrawer(t *testing.T) {
	s1 := dummyScene{
		drawFn: func(screen *ebiten.Image) {
			screen.Fill(color.White)
		},
	}

	var canEndS2 bool
	s2 := dummyScene{
		updateFn: func() error {
			if canEndS2 {
				return ebiten.Termination
			}
			return nil
		},
		drawFn: func(screen *ebiten.Image) {
			screen.Fill(color.White)
		},
		onArrivalFn: func() { canEndS2 = true },
	}

	game := bamenn.NewSequence(&s1)
	fading := bamennutil.LinearFillFadingDrawer{color.Black}
	tran := bamenn.NewLinearTransition(3, 7, fading)

	switched := false
	s1.updateFn = func() error {
		if switched {
			return nil
		}
		if game.SwitchWithTransition(&s2, tran) {
			switched = true
		}
		return nil
	}

	records := make([]string, 0)
	recordFn := func(screen *ebiten.Image) {
		c := screen.At(1, 1)
		rcd := fmt.Sprintf("%v", c)
		records = append(records, rcd)
	}

	fakeScreen := ebiten.NewImage(3, 3)

	for range 100 { // loop 100 times to avoid inf loop
		err := game.Update()

		if errors.Is(err, ebiten.Termination) {
			break
		}
		if err != nil {
			t.Fatalf("unexpected err on Game.Update(): %v", err)
		}

		game.Draw(fakeScreen)
		recordFn(fakeScreen)
	}

	expecteds := []string{
		fmt.Sprint(color.RGBA{191, 191, 191, 255}),
		fmt.Sprint(color.RGBA{128, 128, 128, 255}),
		fmt.Sprint(color.RGBA{64, 64, 64, 255}),
		fmt.Sprint(color.RGBA{0, 0, 0, 255}),
		fmt.Sprint(color.RGBA{64, 64, 64, 255}),
		fmt.Sprint(color.RGBA{128, 128, 128, 255}),
		fmt.Sprint(color.RGBA{191, 191, 191, 255}),
	}

	if len(expecteds) != len(records) {
		t.Fatalf("record len expected %d but got %d", len(expecteds), len(records))
	}

	for i := range expecteds {
		e := expecteds[i]
		r := records[i]

		if e != r {
			t.Errorf("%d: records are different:\n%s\n%s", i, e, r)
		}
	}
}

type dummyScene struct {
	updateFn      func() error
	drawFn        func(screen *ebiten.Image)
	onDepartureFn func()
	onArrivalFn   func()
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

func (s *dummyScene) Layout(ow, oh int) (int, int) {
	return ow, oh
}

func (s *dummyScene) OnDeparture() {
	if s.onDepartureFn == nil {
		return
	}
	s.onDepartureFn()
}

func (s *dummyScene) OnArrival() {
	if s.onArrivalFn == nil {
		return
	}
	s.onArrivalFn()
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
	return 3, 3
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
