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
	alpha := 0.0

	switch f := progress.CurrentFrame - progress.FrameToSwitch; {
	case f < 0:
		alpha = float64(progress.CurrentFrame+1) / float64(progress.FrameToSwitch+1)
	case f == 0:
		alpha = 1
	case f > 0:
		alpha = float64(progress.MaxFrames-progress.CurrentFrame) / float64(progress.MaxFrames-progress.FrameToSwitch)
	}

	return alpha
}
