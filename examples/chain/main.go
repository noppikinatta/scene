package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/noppikinatta/scene"
	"github.com/noppikinatta/scene/sceneutil"
)

func main() {
	s := createScenes()
	g := scene.ToGame(s, func(outsideWidth, outsideHeight int) (screenWidth int, screenHeight int) {
		return outsideWidth, outsideHeight
	})

	ebiten.SetWindowSize(600, 600)

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}

func createScenes() scene.Scene {
	// Define scene names.
	const (
		name1 = "red"
		name2 = "green"
		name3 = "blue"
		name4 = "yellow"
		name5 = "purple"
	)

	flow := exampleFlow{}

	// Create scene instances.
	scene1 := newExampleScene(name1, color.RGBA{R: 200, A: 255}, &flow, name2, name3)
	scene2 := newExampleScene(name2, color.RGBA{G: 160, A: 255}, &flow, name4, name1)
	scene3 := newExampleScene(name3, color.RGBA{B: 200, A: 255}, &flow, name5, name4)
	scene4 := newExampleScene(name4, color.RGBA{R: 200, G: 160, A: 255}, &flow, name2, name1)
	scene5 := newExampleScene(name5, color.RGBA{R: 200, B: 200, A: 255}, &flow, name3, "EXIT")

	// Create scene map and registar to Flow.
	withFade := func(s scene.Scene) scene.Scene {
		return sceneutil.WithSimpleFade(s, 15, color.Black)
	}

	scenes := map[string]scene.Scene{
		scene1.Name: withFade(scene1),
		scene2.Name: withFade(scene2),
		scene3.Name: withFade(scene3),
		scene4.Name: withFade(scene4),
		scene5.Name: withFade(scene5),
	}

	flow.Scenes = scenes

	// Create Chain instance.
	return scene.NewChain(scenes[scene1.Name], &flow)
}

// example Scene implementation
type exampleScene struct {
	Name    string
	Buttons []*exampleButton
	color   color.Color
	ended   bool
}

func newExampleScene(name string, color color.Color, flow *exampleFlow, nextSceneNames ...string) *exampleScene {
	s := &exampleScene{
		Name:  name,
		color: color,
	}

	// button event handler
	buttonClickHandler := func(args any) {
		sceneName := args.(string)
		flow.SetNextSceneName(sceneName)
		s.SetEnded(true)
	}

	buttons := make([]*exampleButton, 0, len(nextSceneNames))
	for i, n := range nextSceneNames {
		buttons = append(buttons, &exampleButton{
			NextSceneName: n,
			Bounds:        image.Rect(i*100, 30, i*100+80, 80),
			EventHandlers: []func(any){
				buttonClickHandler,
			},
		})
	}

	s.Buttons = buttons
	return s
}

func (s *exampleScene) Init() {
	s.ended = false
}

func (s *exampleScene) Update() error {
	for _, b := range s.Buttons {
		b.Update()
	}
	return nil
}

func (s *exampleScene) Draw(screen *ebiten.Image) {
	screen.Fill(s.color)
	ebitenutil.DebugPrint(screen, s.Name)

	for _, b := range s.Buttons {
		b.Draw(screen)
	}
}

func (s *exampleScene) Ended() bool {
	return s.ended
}

func (s *exampleScene) SetEnded(value bool) {
	// In this example, this function is called via the button's event handler.
	s.ended = value
}

func (s *exampleScene) Dispose() {}

// simple button for example
type exampleButton struct {
	NextSceneName string
	Bounds        image.Rectangle
	EventHandlers []func(args any)
}

func (b *exampleButton) Update() {
	if !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return
	}

	if !b.contains(ebiten.CursorPosition()) {
		return
	}

	for _, h := range b.EventHandlers {
		h(b.NextSceneName)
	}
}

func (b *exampleButton) contains(x, y int) bool {
	if x < b.Bounds.Min.X {
		return false
	}
	if x > b.Bounds.Max.X {
		return false
	}
	if y < b.Bounds.Min.Y {
		return false
	}
	if y > b.Bounds.Max.Y {
		return false
	}

	return true
}

func (b *exampleButton) Draw(screen *ebiten.Image) {
	var x, y, w, h float32
	x = float32(b.Bounds.Min.X)
	y = float32(b.Bounds.Min.Y)
	w = float32(b.Bounds.Dx())
	h = float32(b.Bounds.Dy())

	vector.DrawFilledRect(screen, x, y, w, h, color.RGBA{B: 128, A: 128}, false)
	ebitenutil.DebugPrintAt(screen, b.NextSceneName, b.Bounds.Min.X+4, b.Bounds.Min.Y+4)
}

// example Flow implementation
type exampleFlow struct {
	Scenes        map[string]scene.Scene
	NextSceneName string
}

func (f *exampleFlow) Init() {
	f.NextSceneName = ""
}

func (f *exampleFlow) SetNextSceneName(name string) {
	// In this example, this function is called via the button's event handler.
	f.NextSceneName = name
}

func (f *exampleFlow) NextScene(current scene.Scene) (scene.Scene, bool) {
	// You don't have to use current parameter.
	s, ok := f.Scenes[f.NextSceneName]
	return s, ok
}
