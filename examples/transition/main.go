package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/noppikinatta/scene"
)

func main() {
	s := &showFPSScene{}

	tran := scene.NewLinearTransition(30, 61, verticalLineTransitionDrawer{})

	seq := scene.NewSequence(s)

	handler := func() {
		seq.SwitchWithTransition(s, tran)
	}
	s.handler = handler

	ebiten.SetWindowSize(640, 480)

	err := ebiten.RunGame(seq)
	if err != nil {
		log.Fatal(err)
	}
}

type showFPSScene struct {
	handler func()
}

func (s *showFPSScene) Update() error {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		s.handler()
	}
	return nil
}

func (s *showFPSScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{B: 200, A: 255})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.1f / TPS: %.1f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (s *showFPSScene) Layout(ow, oh int) (int, int) {
	return ow, oh
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
	if progress.FrameToSwitch {
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
