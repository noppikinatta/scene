package scene

import "github.com/hajimehoshi/ebiten/v2"

type FinalScreenDrawerConvertible interface {
	drawFinalScreenFunc() func(screen ebiten.FinalScreen, offScreen *ebiten.Image, geoM ebiten.GeoM)
}

func FinalScreenDrawer(f FinalScreenDrawerConvertible) ebiten.FinalScreenDrawer {
	return &finalScreenDrawer{fn: f.drawFinalScreenFunc()}
}

type finalScreenDrawer struct {
	fn func(screen ebiten.FinalScreen, offScreen *ebiten.Image, geoM ebiten.GeoM)
}

func (f *finalScreenDrawer) DrawFinalScreen(screen ebiten.FinalScreen, offScreen *ebiten.Image, geoM ebiten.GeoM) {
	f.fn(screen, offScreen, geoM)
}

type LayoutFConvertible interface {
	layoutFFunc() func(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64)
}

func LayoutFer(l LayoutFConvertible) ebiten.LayoutFer {
	return &layoutFer{fn: l.layoutFFunc()}
}

type layoutFer struct {
	fn func(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64)
}

func (l *layoutFer) LayoutF(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	return l.fn(outsideWidth, outsideHeight)
}
