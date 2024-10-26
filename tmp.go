package scene

import (
	"fmt"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	theScreenShader     *ebiten.Shader
	theScreenShaderOnce sync.Once
)

var screenShaderSource = []byte(`//kage:unit pixels

package main

func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
	// Blend source colors in a square region, which size is 1/scale.
	scale := imageDstSize()/imageSrc0Size()
	p0 := srcPos - 1/2.0/scale
	p1 := srcPos + 1/2.0/scale

	// Texels must be in the source rect, so it is not necessary to check.
	c0 := imageSrc0UnsafeAt(p0)
	c1 := imageSrc0UnsafeAt(vec2(p1.x, p0.y))
	c2 := imageSrc0UnsafeAt(vec2(p0.x, p1.y))
	c3 := imageSrc0UnsafeAt(p1)

	rate := clamp(fract(p1)*scale, 0, 1)
	return mix(mix(c0, c1, rate.x), mix(c2, c3, rate.x), rate.y)
}
`)

// TODO: replace to it when Ebitrngine v2.9.0 is released
// https://github.com/hajimehoshi/ebiten/blob/366f4899a2c54082afa47f7a095a1795c35d8117/gameforui.go#L145
func defaultDrawFinalScreenTemporaryImplRemoveItWhenEbitengineV290Released(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
	theScreenShaderOnce.Do(func() {
		s, err := ebiten.NewShader(screenShaderSource)
		if err != nil {
			panic(fmt.Sprintf("ebiten: compiling the screen shader failed: %v", err))
		}
		theScreenShader = s
	})

	scale := geoM.Element(0, 0)
	switch {
	case !ebiten.IsScreenFilterEnabled(), math.Floor(scale) == scale:
		op := &ebiten.DrawImageOptions{}
		op.GeoM = geoM
		screen.DrawImage(offscreen, op)
	case scale < 1:
		op := &ebiten.DrawImageOptions{}
		op.GeoM = geoM
		op.Filter = ebiten.FilterLinear
		screen.DrawImage(offscreen, op)
	default:
		op := &ebiten.DrawRectShaderOptions{}
		op.Images[0] = offscreen
		op.GeoM = geoM
		w, h := offscreen.Bounds().Dx(), offscreen.Bounds().Dy()
		screen.DrawRectShader(w, h, theScreenShader, op)
	}
}
