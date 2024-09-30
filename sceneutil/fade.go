package sceneutil

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/scene"
)

// Fade is a Scene for drawing in progress.
type Fade struct {
	currentFrame int
	maxFrames    int
	drawer       ProgressDrawer
}

// NewFade creates a new Fade instance.
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

// Progress returns a value from 0.0 to 1.0 depending on progress.
func (f *Fade) Progress() float64 {
	return float64(f.currentFrame) / float64(f.maxFrames)
}

func (f *Fade) Ended() bool {
	return f.currentFrame >= f.maxFrames
}

func (f *Fade) Dispose() {}

// ProgressDrawer draws according to progress.
type ProgressDrawer interface {
	// Draw draws according to progress.
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

	screen.DrawImage(dummyWhitePixel, &o)
}

// ProgressDrawerFadeInFill fills a single color that can be used for fade-ins.
func ProgressDrawerFadeInFill(color color.Color) ProgressDrawer {
	return &progressDrawerFadeFill{
		color:  color,
		fadeIn: true,
	}
}

// ProgressDrawerFadeOutFill fills a single color that can be used for fade-outs.
func ProgressDrawerFadeOutFill(color color.Color) ProgressDrawer {
	return &progressDrawerFadeFill{
		color:  color,
		fadeIn: false,
	}
}

// WithSimpleFade wraps the passed Scene with a simple fade-in and fade-out.
func WithSimpleFade(s scene.Scene, frames int, color color.Color) scene.Scene {
	fadeIn := NewFade(frames, ProgressDrawerFadeInFill(color))
	fadeOut := NewFade(frames, ProgressDrawerFadeOutFill(color))

	seq := scene.NewSequencialFlow(fadeIn, scene.NewBarrier(s.Ended), fadeOut)

	return scene.NewParallel(
		s,
		scene.NewChain(fadeIn, seq),
	)
}
