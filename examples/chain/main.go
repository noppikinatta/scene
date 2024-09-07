package main

import (
	"fmt"
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

	s1.leftClickNext = &s2
	s1.rightClickNext = &s3

	s2.leftClickNext = &s4
	s2.rightClickNext = &s1

	s3.leftClickNext = &s5
	s3.rightClickNext = &s4

	s4.leftClickNext = &s2
	s4.rightClickNext = &s1

	s5.leftClickNext = &s3
	s5.rightClickNext = nil

	return scene.NewChain(&s1)
}

type exampleScene struct {
	waitFrames     int
	color          color.Color
	name           string
	ended          bool
	leftClickNext  *exampleScene
	rightClickNext *exampleScene
	next           scene.Scene
}

func (s *exampleScene) Init() {
	s.waitFrames = 15
	s.ended = false
}

func (s *exampleScene) Update() error {
	if s.waitFrames > 0 {
		s.waitFrames--
		return nil
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		s.ended = true
		if s.leftClickNext != nil {
			s.next = s.leftClickNext
		} else {
			s.next = nil
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		s.ended = true
		if s.rightClickNext != nil {
			s.next = s.rightClickNext
		} else {
			s.next = nil
		}
	}
	return nil
}

func (s *exampleScene) Draw(screen *ebiten.Image) {
	screen.Fill(s.color)

	name := func(s *exampleScene) string {
		if s == nil {
			return "EXIT GAME"
		}
		return s.name
	}

	txt := fmt.Sprintf(
		"current: %s\nleft click: %s\nright click: %s",
		s.name,
		name(s.leftClickNext),
		name(s.rightClickNext),
	)

	ebitenutil.DebugPrint(screen, txt)
}

func (s *exampleScene) Ended() bool {
	return s.ended
}

func (s *exampleScene) NextScene() (scene.Scene, bool) {
	if !s.Ended() {
		return nil, false
	}

	return s.next, s.next != nil
}

func (s *exampleScene) Dispose() {
}
