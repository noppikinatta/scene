package sceneutil

import (
	"image/color"

	"github.com/noppikinatta/scene"
)

// SimpleFadeAdderFunction returns a function adds fade functionality to the given Scene.
func SimpleFadeAdderFunction(frames int, color color.Color) func(s scene.Scene) scene.Scene {
	return func(s scene.Scene) scene.Scene {
		p := scene.NewParallel(
			s,
			scene.NewSequence(
				scene.NewFade(frames, scene.ProgressDrawerFadeInFill(color)),
				scene.NewWait(s),
				scene.NewFade(frames, scene.ProgressDrawerFadeOutFill(color)),
			),
		)

		return p
	}
}
