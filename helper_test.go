package bamenn_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func runForTest(t *testing.T, game ebiten.Game) {
	t.Helper()

	dummyScreen := ebiten.NewImage(3, 3)
	dummyFinalScreen := ebiten.NewImage(3, 3)

	for range 100 { // loop 100 times to avoid inf loop
		if l, ok := game.(ebiten.LayoutFer); ok {
			l.LayoutF(0, 0)
		} else {
			game.Layout(0, 0)
		}

		err := game.Update()

		if errors.Is(err, ebiten.Termination) {
			break
		}
		if err != nil {
			t.Fatalf("unexpected err on Game.Update(): %v", err)
		}

		game.Draw(dummyScreen)
		if f, ok := game.(ebiten.FinalScreenDrawer); ok {
			f.DrawFinalScreen(dummyFinalScreen, dummyScreen, ebiten.GeoM{})
		}
	}
}

func compareLogs(t *testing.T, expectedLog, actualLog []string) {
	t.Helper()

	exLen := len(expectedLog)
	acLen := len(actualLog)

	if exLen != acLen {
		t.Errorf("expected log length %d, but got %d", exLen, acLen)
	}

	l := exLen
	if acLen < exLen {
		l = acLen
	}

	for i := 0; i < l; i++ {
		exLog := expectedLog[i]
		acLog := actualLog[i]

		if exLog != acLog {
			t.Errorf("%d: log different\nex: %s\nac: %s", i, exLog, acLog)
		}
	}

	if acLen < exLen {
		for i := acLen; i < exLen; i++ {
			exLog := expectedLog[i]
			t.Errorf("%d: log different\nex: %s\nac: NO ITEM", i, exLog)
		}
	}
	if acLen > exLen {
		for i := exLen; i < acLen; i++ {
			acLog := actualLog[i]
			t.Errorf("%d: log different\nex: NO ITEM\nac: %s", i, acLog)
		}
	}
}

type recorder struct {
	Log []string
}

func (r *recorder) Append(name, logType string) {
	r.Log = append(r.Log, fmt.Sprintf("%s:%s", name, logType))
}

type gameForTest struct {
	Name             string
	UpdateFn         func() error
	OnArrivalFn      func()
	Recorder         *recorder
	LayoutW, LayoutH int
}

func (g *gameForTest) append(logType string) {
	if g.Recorder == nil {
		return
	}
	g.Recorder.Append(g.Name, logType)
}

func (g *gameForTest) Update() error {
	g.append("update")
	if g.UpdateFn == nil {
		return ebiten.Termination
	}
	return g.UpdateFn()
}

func (g *gameForTest) Draw(screen *ebiten.Image) {
	g.append("draw")
}

func (g *gameForTest) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	g.append("layout")
	return g.LayoutW, g.LayoutH
}

type eventsForTest struct {
	gameForTest
}

func (e *eventsForTest) OnStart() {
	e.gameForTest.append("onstart")
}

func (e *eventsForTest) OnArrival() {
	e.gameForTest.append("onarrival")
	if e.OnArrivalFn != nil {
		e.OnArrivalFn()
	}
}

func (e *eventsForTest) OnDeparture() {
	e.gameForTest.append("ondeparture")
}

func (e *eventsForTest) OnEnd() {
	e.gameForTest.append("onend")
}

type finalScreenDrawerForTest struct {
	gameForTest
	drawFn func(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM)
}

func (f *finalScreenDrawerForTest) DrawFinalScreen(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
	f.gameForTest.append("drawfinalscreen")
	if f.drawFn != nil {
		f.drawFn(screen, offscreen, geoM)
	}
}

type layoutFerForTest struct {
	gameForTest
	layoutFW, layoutFH float64
}

func (l *layoutFerForTest) LayoutF(outsideWidth float64, outsideHeight float64) (screenWidth float64, screenHeight float64) {
	l.gameForTest.append("layoutf")
	return l.layoutFW, l.layoutFH
}

type transitionForTest struct {
	Name         string
	Recorder     *recorder
	SwitchFrames int
	MaxFrames    int
	currentFrame int
}

func (t *transitionForTest) append(logType string) {
	if t.Recorder == nil {
		return
	}
	t.Recorder.Append(t.Name, logType)
}

func (t *transitionForTest) Reset() {
	t.currentFrame = 0
	t.append("reset")
}

func (t *transitionForTest) Update() error {
	t.append("update")
	if t.currentFrame < t.MaxFrames {
		t.currentFrame++
	}
	return nil
}

func (t *transitionForTest) Draw(screen *ebiten.Image) {
	t.append("draw")
}

func (t *transitionForTest) Completed() bool {
	return t.currentFrame >= t.MaxFrames
}

func (t *transitionForTest) CanSwitchScenes() bool {
	return t.currentFrame >= t.SwitchFrames
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
