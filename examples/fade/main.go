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
	showFPS := showFPSScene{}

	s := scene.NewParallel(
		&showFPS,
		scene.NewSequence(
			scene.NewFade(60, scene.ProgressDrawerFadeInFill(color.Black)),
			scene.NewWait(&showFPS),
			scene.NewFade(60, scene.ProgressDrawerFadeOutFill(color.Black)),
		),
	)
	g := scene.ToGame(s, func(outsideWidth, outsideHeight int) (screenWidth int, screenHeight int) {
		return outsideWidth, outsideHeight
	})

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}

type showFPSScene struct {
	ended bool
}

func (s *showFPSScene) Init() {
	s.ended = false
}

func (s *showFPSScene) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		s.ended = true
	}
	return nil
}

func (s *showFPSScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{B: 200, A: 255})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.1f / TPS: %.1f", ebiten.ActualFPS(), ebiten.ActualTPS()))
	fmt.Println("DRAWING")
}

func (s *showFPSScene) Ended() bool {
	return s.ended
}

func (s *showFPSScene) Dispose() {
}
