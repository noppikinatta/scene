package scene

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: replace it when Ebitrngine v2.9.0 is released
// https://github.com/hajimehoshi/ebiten/blob/366f4899a2c54082afa47f7a095a1795c35d8117/gameforui.go#L145
func defaultDrawFinalScreenTemporaryImplRemoveItWhenEbitengineV290Released(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
	scale := geoM.Element(0, 0)
	switch {
	case math.Floor(scale) == scale:
		op := &ebiten.DrawImageOptions{}
		op.GeoM = geoM
		screen.DrawImage(offscreen, op)
	default:
		op := &ebiten.DrawImageOptions{}
		op.GeoM = geoM
		op.Filter = ebiten.FilterLinear
		screen.DrawImage(offscreen, op)
	}
}
