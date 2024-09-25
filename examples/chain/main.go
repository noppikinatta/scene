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
	const (
		name1 = "red"
		name2 = "green"
		name3 = "blue"
		name4 = "yellow"
		name5 = "purple"
	)

	nextScener := exampleNextScener{}

	scene1 := newExampleScene(name1, color.RGBA{R: 200, A: 255}, &nextScener, name2, name3)
	scene2 := newExampleScene(name2, color.RGBA{G: 160, A: 255}, &nextScener, name4, name1)
	scene3 := newExampleScene(name3, color.RGBA{B: 200, A: 255}, &nextScener, name5, name4)
	scene4 := newExampleScene(name4, color.RGBA{R: 200, G: 160, A: 255}, &nextScener, name2, name1)
	scene5 := newExampleScene(name5, color.RGBA{R: 200, B: 200, A: 255}, &nextScener, name3, "EXIT")

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

	nextScener.Scenes = scenes

	return scene.NewChain(scenes[scene1.Name], &nextScener)
}

type exampleScene struct {
	scene.Scene
	Name    string
	Buttons []*exampleButton
	color   color.Color
	ended   bool
}

func newExampleScene(name string, color color.Color, nextScener *exampleNextScener, nextSceneNames ...string) *exampleScene {
	s := &exampleScene{
		Name:  name,
		color: color,
	}

	setNextScene := func(args any) {
		sceneName := args.(string)
		nextScener.SetNextSceneName(sceneName)
	}

	setEnded := func(any) {
		s.SetEnded(true)
	}

	buttons := make([]*exampleButton, 0, len(nextSceneNames))
	for _, n := range nextSceneNames {
		buttons = append(buttons, &exampleButton{
			NextSceneName: n,
			EventHandlers: []func(any){
				setNextScene,
				setEnded,
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

	for i, b := range s.Buttons {
		if b.Bounds.Empty() {
			b.Bounds = image.Rect(i*100, 30, i*100+80, 80)
		}
		b.Draw(screen)
	}
}

func (s *exampleScene) Ended() bool {
	return s.ended
}

func (s *exampleScene) SetEnded(value bool) {
	s.ended = value
}

func (s *exampleScene) Dispose() {

}

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

type exampleNextScener struct {
	Scenes        map[string]scene.Scene
	NextSceneName string
}

func (n *exampleNextScener) SetNextSceneName(name string) {
	n.NextSceneName = name
}

func (n *exampleNextScener) NextScene(current scene.Scene) (scene.Scene, bool) {
	s, ok := n.Scenes[n.NextSceneName]
	return s, ok
}
