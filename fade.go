package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Fade struct {
	currentFrame int
	maxFrames    int
	drawer       ProgressDrawer
}

func NewFade(frames int, drawer ProgressDrawer) *Fade {
	return &Fade{
		maxFrames: frames,
		drawer:    drawer,
	}
}

func (f *Fade) Init() {
	f.currentFrame = 0
}

func (f *Fade) Update() error {
	if f.currentFrame < f.maxFrames {
		f.currentFrame++
	}
	return nil
}

func (f *Fade) Draw(screen *ebiten.Image) {
	f.drawer.Draw(screen, f.Progress())
}

func (f *Fade) Progress() float64 {
	return float64(f.currentFrame) / float64(f.maxFrames)
}

func (f *Fade) Ended() bool {
	return f.currentFrame >= f.maxFrames
}

func (f *Fade) Dispose() {}

type ProgressDrawer interface {
	Draw(screen *ebiten.Image, progress float64)
}

type progressDrawerFadeFill struct {
	color  color.Color
	fadeIn bool
}

func (d *progressDrawerFadeFill) Draw(screen *ebiten.Image, progress float64) {
	alpha := progress
	if d.fadeIn {
		alpha = 1 - alpha
	}

	screenSize := screen.Bounds().Size()

	o := ebiten.DrawImageOptions{}
	o.ColorScale.ScaleWithColor(d.color)
	o.ColorScale.ScaleAlpha(float32(alpha))
	o.GeoM.Scale(
		float64(screenSize.X),
		float64(screenSize.Y),
	)

	screen.DrawImage(DummyWhitePixel, &o)
}

func ProgressDrawerFadeInFill(color color.Color) ProgressDrawer {
	return &progressDrawerFadeFill{
		color:  color,
		fadeIn: true,
	}
}

func ProgressDrawerFadeOutFill(color color.Color) ProgressDrawer {
	return &progressDrawerFadeFill{
		color:  color,
		fadeIn: false,
	}
}

// WithSimpleFade wraps the passed Scene with a simple fade-in and fade-out.
func WithSimpleFade(s Scene, frames int, color color.Color) Scene {
	fadeIn := NewFade(frames, ProgressDrawerFadeInFill(color))

	seq := NewSequencialNextScener(
		fadeIn,
		NewWait(s.Ended),
		NewFade(frames, ProgressDrawerFadeOutFill(color)),
	)

	return NewParallel(
		s,
		NewChain(
			fadeIn,
			seq,
		),
	)
}
