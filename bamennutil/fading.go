package bamennutil

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/bamenn"
)

// LinearFillFadingDrawer can be used to draw LinearTransitions. It performs a fade-in/fade-out that fills in the specified color.
type LinearFillFadingDrawer struct {
	Color color.Color
}

func (d LinearFillFadingDrawer) Draw(screen *ebiten.Image, progress bamenn.LinearTransitionProgress) {
	screenSize := screen.Bounds().Size()

	o := ebiten.DrawImageOptions{}
	o.ColorScale.ScaleWithColor(d.Color)
	o.ColorScale.ScaleAlpha(float32(d.alpha(progress)))
	o.GeoM.Scale(
		float64(screenSize.X),
		float64(screenSize.Y),
	)

	screen.DrawImage(dummyWhitePixel, &o)
}

func (d LinearFillFadingDrawer) alpha(progress bamenn.LinearTransitionProgress) float64 {
	if progress.FrameToSwitch {
		return 1
	}

	rate := progress.Rate()
	if rate < 0.5 {
		return rate * 2
	} else {
		prevFrame := progress
		prevFrame.CurrentFrame--
		return (1 - prevFrame.Rate()) * 2
	}
}
