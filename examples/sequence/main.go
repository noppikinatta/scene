package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/noppikinatta/scene"
	"github.com/noppikinatta/scene/sceneutil"
)

func main() {
	s := createScenes()
	g := scene.ToGame(s, func(outsideWidth, outsideHeight int) (screenWidth int, screenHeight int) {
		return outsideWidth, outsideHeight
	})

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}

func createScenes() scene.Scene {
	s1 := sceneutil.WithSimpleFade(&exampleScene{color: color.RGBA{R: 200, A: 255}, name: "red"}, 15, color.Black)
	s2 := sceneutil.WithSimpleFade(&exampleScene{color: color.RGBA{G: 180, A: 255}, name: "green"}, 15, color.Black)
	s3 := sceneutil.WithSimpleFade(&exampleScene{color: color.RGBA{B: 200, A: 255}, name: "blue"}, 15, color.Black)
	s4 := sceneutil.WithSimpleFade(&exampleScene{color: color.RGBA{R: 200, G: 180, A: 255}, name: "yellow"}, 15, color.Black)
	s5 := sceneutil.WithSimpleFade(&exampleScene{color: color.RGBA{R: 200, B: 200, A: 255}, name: "purple"}, 15, color.Black)

	nextScener := scene.NewSequencialLoopNextScener(s1, s2, s3, s4, s5)

	return scene.NewChain(s1, nextScener)
}

type exampleScene struct {
	color color.Color
	name  string
	ended bool
}

func (s *exampleScene) Init() {
	s.ended = false
}

func (s *exampleScene) Update() error {

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		s.ended = true
	}
	return nil
}

func (s *exampleScene) Draw(screen *ebiten.Image) {
	screen.Fill(s.color)
	ebitenutil.DebugPrint(screen, s.name)
}

func (s *exampleScene) Ended() bool {
	return s.ended
}

func (s *exampleScene) Dispose() {
}