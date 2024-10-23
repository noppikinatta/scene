package main

import (
	"errors"
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

	ebiten.SetWindowSize(600, 600)

	err := ebiten.RunGame(s)
	if err != nil {
		log.Fatal(err)
	}
}

func createScenes() ebiten.Game {
	// Create scene instances.
	scene1 := newExampleScene("red", color.RGBA{R: 200, A: 255})
	scene2 := newExampleScene("green", color.RGBA{G: 160, A: 255})
	scene3 := newExampleScene("blue", color.RGBA{B: 200, A: 255})
	scene4 := newExampleScene("yellow", color.RGBA{R: 200, G: 160, A: 255})
	scene5 := newExampleScene("purple", color.RGBA{R: 200, B: 200, A: 255})

	// Create Director.
	director := scene.NewSequence(scene1)

	// Add buttons to switch scenes.
	tran := scene.NewLinearTransition(10, sceneutil.LinearFillFadingDrawer{Color: color.Black})
	addButton := func(scene, nextScene *exampleScene) {
		scene.AddButton(nextScene.Name, func() error {
			director.SwitchWithTransition(nextScene, tran)
			return nil
		})
	}

	addButton(scene1, scene2)
	addButton(scene1, scene3)
	addButton(scene2, scene4)
	addButton(scene2, scene1)
	addButton(scene3, scene5)
	addButton(scene3, scene4)
	addButton(scene4, scene1)
	addButton(scene4, scene3)
	addButton(scene5, scene2)

	// Add exit button.
	scene5.AddButton("EXIT", func() error {
		return ebiten.Termination
	})

	return director
}

// example Scene implementation
type exampleScene struct {
	Name    string
	Buttons []*exampleButton
	color   color.Color
	ended   bool
}

func newExampleScene(name string, color color.Color) *exampleScene {
	return &exampleScene{
		Name:  name,
		color: color,
	}
}

func (s *exampleScene) AddButton(name string, handlers ...func() error) {
	buttonBounds := func(num int) image.Rectangle {
		left := 10 + num*100
		top := 20
		width := 80
		height := 60
		return image.Rect(
			left,
			top,
			left+width,
			top+height,
		)
	}

	s.Buttons = append(s.Buttons, &exampleButton{
		Name:          name,
		Bounds:        buttonBounds(len(s.Buttons)),
		EventHandlers: handlers,
	})
}

func (s *exampleScene) OnSceneStart() {
	s.ended = false
}

func (s *exampleScene) Update() error {
	for _, b := range s.Buttons {
		err := b.Update()
		if err != nil {
			return err
		}
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

func (s *exampleScene) Layout(ow, oh int) (int, int) {
	return ow, oh
}

// simple button for example
type exampleButton struct {
	Name          string
	Bounds        image.Rectangle
	EventHandlers []func() error
}

func (b *exampleButton) Update() error {
	if !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return nil
	}

	if !b.contains(ebiten.CursorPosition()) {
		return nil
	}

	errs := make([]error, 0)
	for _, h := range b.EventHandlers {
		errs = append(errs, h())
	}

	return errors.Join(errs...)
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
	ebitenutil.DebugPrintAt(screen, b.Name, b.Bounds.Min.X+4, b.Bounds.Min.Y+4)
}
