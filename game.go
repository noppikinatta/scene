package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// ToGame wraps the Scene with ebiten.Game.
func ToGame(scene Scene, layoutFn func(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)) ebiten.Game {
	return &game{
		scene:    scene,
		layoutFn: layoutFn,
	}
}

type game struct {
	inited   bool
	scene    Scene
	layoutFn func(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
}

func (g *game) Update() error {
	if !g.inited {
		if o, ok := g.scene.(OnSceneStarter); ok {
			o.OnSceneStart()
		}
		g.inited = true
	}
	if g.scene.CanEnd() {
		if o, ok := g.scene.(OnSceneEnder); ok {
			o.OnSceneEnd()
		}
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
	return g.layoutFn(outsideWidth, outsideHeight)
}
