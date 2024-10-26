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
	g := createGame()

	ebiten.SetWindowSize(600, 600)

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}

func createGame() ebiten.Game {
	// Create scene instances.
	scene1 := &gameWithDinalScreenDrawer{}
	scene2 := &gameWithoutDinalScreenDrawer{}

	// Create Sequence.
	sequence := scene.NewSequence(scene1)

	// Add handlers to switch scenes.
	tran := scene.NewLinearTransition(10, sceneutil.LinearFillFadingDrawer{Color: color.Black})
	scene1.handler = func() { sequence.SwitchWithTransition(scene2, tran) }
	scene2.handler = func() { sequence.SwitchWithTransition(scene1, tran) }

	finalScreenDrawer := scene.ToFinalScreenDrawer(sequence)

	game := gameAndFinalScreenDrawer{
		Game:              sequence,
		FinalScreenDrawer: finalScreenDrawer,
	}

	return game
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

// Game to run on this example.
type gameAndFinalScreenDrawer struct {
	ebiten.Game
	ebiten.FinalScreenDrawer
}
