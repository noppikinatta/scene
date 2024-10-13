package sceneutil

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/scene"
)

type LinearFillFadingDrawer struct {
	Color color.Color
}

func (d LinearFillFadingDrawer) Draw(screen *ebiten.Image, progress scene.LinearTransitionProgress) bool {
	screenSize := screen.Bounds().Size()

	o := ebiten.DrawImageOptions{}
	o.ColorScale.ScaleWithColor(d.Color)
	o.ColorScale.ScaleAlpha(float32(d.alpha(progress)))
	o.GeoM.Scale(
		float64(screenSize.X),
		float64(screenSize.Y),
	)

	screen.DrawImage(dummyWhitePixel, &o)
	return progress.JustHalf()
}

func (d LinearFillFadingDrawer) alpha(progress scene.LinearTransitionProgress) float64 {
	if progress.JustHalf() {
		return 1
	}

	rate := progress.Rate()
	if rate < 0.5 {
		return rate * 2
	} else {
		return (1 - rate) * 2
	}
}