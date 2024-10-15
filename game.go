package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// ToGame wraps the Scene with ebiten.Game.
func ToGame(scene Scene, layouter Layouter) ebiten.Game {
	return &game{
		scene:    scene,
		layouter: layouter,
	}
}

type Layouter interface {
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
}

type layouterFn func(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)

func (l layouterFn) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return l(outsideWidth, outsideHeight)
}

func NewLayouterFromFunc(fn func(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)) Layouter {
	return layouterFn(fn)
}

type game struct {
	inited   bool
	scene    Scene
	layouter Layouter
}

func (g *game) Update() error {
	if !g.inited {
		tryCall(g.scene, func(o OnSceneStarter) { o.OnSceneStart() })
		g.inited = true
	}
	if g.scene.CanEnd() {
		tryCall(g.scene, func(o OnSceneEnder) { o.OnSceneEnd() })
		return ebiten.Termination
	}
	return g.scene.Update()
}

func (g *game) Draw(screen *ebiten.Image) {
	if !g.inited {
		return
	}
	if g.scene.CanEnd() {
		return
	}
	g.scene.Draw(screen)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.layouter.Layout(outsideWidth, outsideHeight)
}
