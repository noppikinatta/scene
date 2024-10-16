package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/noppikinatta/scene"
	"github.com/noppikinatta/scene/sceneutil"
)

func main() {
	s := &showFPSScene{}
	flow := scene.NewSequencialLoopFlow(s)

	tran := scene.NewLinearTransition(61, verticalLineTransitionDrawer{})
	traner := scene.NewFixedTransitioner(tran)

	chain := scene.NewChain(s, flow, traner)

	g := scene.ToGame(chain, sceneutil.SimpleLayoutFunc())

	ebiten.SetWindowSize(640, 480)

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}

type showFPSScene struct {
	ended bool
}

func (s *showFPSScene) OnSceneStart() {
	s.ended = false
}

func (s *showFPSScene) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		s.ended = true
	}
	return nil
}

func (s *showFPSScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{B: 200, A: 255})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.1f / TPS: %.1f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (s *showFPSScene) CanEnd() bool {
	return s.ended
}

var (
	dummyImageBase = ebiten.NewImage(3, 3)

	// dummyWhitePixel is a 1x1 white pixel image.
	dummyWhitePixel = dummyImageBase.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	dummyImageBase.Fill(color.White)
}

type verticalLineTransitionDrawer struct{}

func (d verticalLineTransitionDrawer) Draw(screen *ebiten.Image, progress scene.LinearTransitionProgress) {
	if progress.Halfway() {
		screen.Fill(color.Black)
		return
	}
	d.drawVerticalLines(screen, progress)
}

func (d verticalLineTransitionDrawer) drawVerticalLines(screen *ebiten.Image, progress scene.LinearTransitionProgress) {
	const wCount = 8

	var l, w int
	wBase := screen.Bounds().Dx() / wCount

	if progress.Rate() < 0.5 {
		w = int(float64(wBase) * progress.Rate() * 2)
		l = 0
	} else {
		w = int(float64(wBase) * (1 - progress.Rate()) * 2)
		l = wBase - w
	}

	if w == 0 {
		return
	}

	screenSize := screen.Bounds().Size()

	for i := 0; i < wCount; i++ {
		o := ebiten.DrawImageOptions{}
		o.ColorScale.ScaleWithColor(color.Black)
		o.GeoM.Scale(
			float64(w),
			float64(screenSize.Y),
		)
		o.GeoM.Translate(float64(l+i*wBase), 0)

		screen.DrawImage(dummyWhitePixel, &o)
	}
}
