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
	s := sceneutil.WithFade(&showFPSScene{}, 60, 60, &exampleProgressDrawer{true}, &exampleProgressDrawer{false})
	flow := scene.NewSequencialLoopFlow(s)
	chain := scene.NewChain(s, flow)

	g := scene.ToGame(chain, func(outsideWidth, outsideHeight int) (screenWidth int, screenHeight int) {
		return outsideWidth, outsideHeight
	})

	ebiten.SetWindowSize(640, 480)

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}

type showFPSScene struct {
	ended bool
}

func (s *showFPSScene) Init() {
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

func (s *showFPSScene) Ended() bool {
	return s.ended
}

func (s *showFPSScene) Dispose() {
}

var (
	dummyImageBase = ebiten.NewImage(3, 3)

	// dummyWhitePixel is a 1x1 white pixel image.
	dummyWhitePixel = dummyImageBase.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	dummyImageBase.Fill(color.White)
}

type exampleProgressDrawer struct {
	in bool
}

func (d *exampleProgressDrawer) Draw(screen *ebiten.Image, progress float64) {
	const wCount = 8

	var l, w int
	wBase := screen.Bounds().Dx() / wCount

	if d.in {
		w = int(float64(wBase) * (1 - progress))
		l = wBase - w
	} else {
		w = int(float64(wBase) * progress)
		l = 0
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
