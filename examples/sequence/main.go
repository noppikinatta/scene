package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/noppikinatta/scene"
	"github.com/noppikinatta/scene/sceneutil"
)

func main() {
	s := createScenes()
	g := scene.ToGame(s, sceneutil.SimpleLayoutFunc())

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}

func createScenes() scene.Scene {
	s1 := exampleScene{color: color.RGBA{R: 200, A: 255}, name: "red"}
	s2 := exampleScene{color: color.RGBA{G: 180, A: 255}, name: "green"}
	s3 := exampleScene{color: color.RGBA{B: 200, A: 255}, name: "blue"}
	s4 := exampleScene{color: color.RGBA{R: 200, G: 180, A: 255}, name: "yellow"}
	s5 := exampleScene{color: color.RGBA{R: 200, B: 200, A: 255}, name: "purple"}

	flow := scene.NewSequencialLoopFlow(&s1, &s2, &s3, &s4, &s5)

	tran := scene.NewLinearTransition(10, sceneutil.LinearFillFadingDrawer{Color: color.Black})
	traner := scene.NewFixedTransitioner(tran)

	return scene.NewChain(&s1, flow, traner)
}

type exampleScene struct {
	color color.Color
	name  string
	ended bool
}

func (s *exampleScene) OnSceneStart() {
	s.ended = false
}

func (s *exampleScene) Update() error {

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		s.ended = true
	}
	return nil
}

func (s *exampleScene) Draw(screen *ebiten.Image) {
	screen.Fill(s.color)
	ebitenutil.DebugPrint(screen, s.name)
}

func (s *exampleScene) CanEnd() bool {
	return s.ended
}
