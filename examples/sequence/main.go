package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/noppikinatta/scene"
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
	s1 := exampleScene{color: color.RGBA{R: 200, A: 255}, name: "red"}
	s2 := exampleScene{color: color.RGBA{G: 180, A: 255}, name: "green"}
	s3 := exampleScene{color: color.RGBA{B: 200, A: 255}, name: "blue"}
	s4 := exampleScene{color: color.RGBA{R: 200, G: 180, A: 255}, name: "yellow"}
	s5 := exampleScene{color: color.RGBA{R: 200, B: 200, A: 255}, name: "purple"}

	scenes := make([]scene.Scene, 0, 5)

	scenes = append(scenes, scene.WithSimpleFade(&s1, 15, color.Black))
	scenes = append(scenes, scene.WithSimpleFade(&s2, 15, color.Black))
	scenes = append(scenes, scene.WithSimpleFade(&s3, 15, color.Black))
	scenes = append(scenes, scene.WithSimpleFade(&s4, 15, color.Black))
	scenes = append(scenes, scene.WithSimpleFade(&s5, 15, color.Black))

	seq := scene.NewSequence(scenes...)
	seq.AddTransition(scenes[4], scenes[0])

	return seq
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
