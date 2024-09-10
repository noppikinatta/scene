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
	const (
		n1 = "red"
		n2 = "green"
		n3 = "blue"
		n4 = "yellow"
		n5 = "purple"
	)

	recorder := mouseButtonDownRecorder{}

	s1 := &exampleScene{
		name:       n1,
		color:      color.RGBA{R: 200, A: 255},
		recorder:   &recorder,
		transition: &mouseButtonTransition{Left: n2, Right: n3},
	}
	s2 := &exampleScene{
		name:       n2,
		color:      color.RGBA{G: 180, A: 255},
		recorder:   &recorder,
		transition: &mouseButtonTransition{Left: n4, Right: n1},
	}
	s3 := &exampleScene{
		name:       n3,
		color:      color.RGBA{B: 200, A: 255},
		recorder:   &recorder,
		transition: &mouseButtonTransition{Left: n5, Right: n4},
	}
	s4 := &exampleScene{
		name:       n4,
		color:      color.RGBA{R: 200, G: 180, A: 255},
		recorder:   &recorder,
		transition: &mouseButtonTransition{Left: n2, Right: n1},
	}
	s5 := &exampleScene{
		name:       n5,
		color:      color.RGBA{R: 200, B: 200, A: 255},
		recorder:   &recorder,
		transition: &mouseButtonTransition{Left: n3, Right: "EXIT GAME"},
	}

	fs1 := scene.WithSimpleFade(s1, 15, color.Black)
	fs2 := scene.WithSimpleFade(s2, 15, color.Black)
	fs3 := scene.WithSimpleFade(s3, 15, color.Black)
	fs4 := scene.WithSimpleFade(s4, 15, color.Black)
	fs5 := scene.WithSimpleFade(s5, 15, color.Black)

	scenes := map[string]scene.Scene{
		n1: fs1,
		n2: fs2,
		n3: fs3,
		n4: fs4,
		n5: fs5,
	}

	transitions := map[scene.Scene]*mouseButtonTransition{
		fs1: s1.transition,
		fs2: s2.transition,
		fs3: s3.transition,
		fs4: s4.transition,
		fs5: s5.transition,
	}

	ns := mouseButtonNextScener{
		scenes:      scenes,
		transitions: transitions,
		recorder:    &recorder,
	}

	return scene.NewChain(fs1, &ns)
}

type exampleScene struct {
	name       string
	waitFrames int
	color      color.Color
	ended      bool
	transition *mouseButtonTransition
	recorder   *mouseButtonDownRecorder
}

func (s *exampleScene) Init() {
	s.waitFrames = 15
	s.recorder.Clear()
	s.ended = false
}

func (s *exampleScene) Update() error {
	if s.waitFrames > 0 {
		s.waitFrames--
		return nil
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		s.ended = true
		s.recorder.Set(ebiten.MouseButtonLeft)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		s.ended = true
		s.recorder.Set(ebiten.MouseButtonRight)
	}
	return nil
}

func (s *exampleScene) Draw(screen *ebiten.Image) {
	screen.Fill(s.color)

	txt := fmt.Sprintf(
		"current: %s\nleft click: %s\nright click: %s",
		s.name,
		s.transition.Left,
		s.transition.Right,
	)

	ebitenutil.DebugPrint(screen, txt)
}

func (s *exampleScene) Ended() bool {
	return s.ended
}

func (s *exampleScene) Dispose() {
}

type mouseButtonDownRecorder struct {
	Pressed bool
	Button  ebiten.MouseButton
}

func (r *mouseButtonDownRecorder) Set(button ebiten.MouseButton) {
	r.Pressed = true
	r.Button = button
}

func (r *mouseButtonDownRecorder) Clear() {
	r.Pressed = false
}

func (r *mouseButtonDownRecorder) IsButtonDown(button ebiten.MouseButton) bool {
	if !r.Pressed {
		return false
	}
	return button == r.Button
}

type mouseButtonTransition struct {
	Left  string
	Right string
}

type mouseButtonNextScener struct {
	scenes      map[string]scene.Scene
	transitions map[scene.Scene]*mouseButtonTransition
	recorder    *mouseButtonDownRecorder
}

func (m *mouseButtonNextScener) NextScene(current scene.Scene) (scene.Scene, bool) {
	t, ok := m.transitions[current]
	if !ok {
		return nil, false
	}

	if m.recorder.IsButtonDown(ebiten.MouseButtonLeft) {
		return m.scene(t.Left)
	}
	if m.recorder.IsButtonDown(ebiten.MouseButtonRight) {
		return m.scene(t.Right)
	}

	return nil, false
}

func (m *mouseButtonNextScener) scene(name string) (scene.Scene, bool) {
	s, ok := m.scenes[name]
	return s, ok
}
