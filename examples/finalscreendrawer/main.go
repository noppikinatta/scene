package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/noppikinatta/scene"
	"github.com/noppikinatta/scene/sceneutil"
)

func main() {
	// Create ebiten.Games as scenes.
	s1 := &gameWithDinalScreenDrawer{}
	s2 := &gameWithoutDinalScreenDrawer{}

	// Create Sequence.
	seq := scene.NewSequence(s1)

	// Add handlers to switch scenes.
	tran := scene.NewLinearTransition(5, 10, sceneutil.LinearFillFadingDrawer{Color: color.Black})
	s1.handler = func() { seq.SwitchWithTransition(s2, tran) }
	s2.handler = func() { seq.SwitchWithTransition(s1, tran) }

	ebiten.SetWindowSize(600, 600)

	err := ebiten.RunGame(seq)
	if err != nil {
		log.Fatal(err)
	}
}

// gameWithDinalScreenDrawer is example game scene with FinalScreenDrawer implementation.
type gameWithDinalScreenDrawer struct {
	handler func()
}

func (g *gameWithDinalScreenDrawer) Update() error {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.handler()
	}
	return nil
}

func (g *gameWithDinalScreenDrawer) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 200, G: 100, A: 255})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.1f / TPS: %.1f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *gameWithDinalScreenDrawer) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *gameWithDinalScreenDrawer) DrawFinalScreen(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
	ebitenutil.DebugPrintAt(offscreen, "FINAL SCREEN DRAWER", 100, 100)

	opt := ebiten.DrawImageOptions{}
	opt.GeoM = geoM
	screen.DrawImage(offscreen, &opt)
}

// gameWithoutDinalScreenDrawer is example game scene without FinalScreenDrawer.
type gameWithoutDinalScreenDrawer struct {
	handler func()
}

func (g *gameWithoutDinalScreenDrawer) Update() error {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.handler()
	}
	return nil
}

func (g *gameWithoutDinalScreenDrawer) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{G: 100, B: 200, A: 255})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.1f / TPS: %.1f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *gameWithoutDinalScreenDrawer) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth, outsideHeight
}
