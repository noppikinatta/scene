package sceneutil

import "github.com/noppikinatta/scene"

func SimpleLayoutFunc() scene.Layouter {
	return scene.NewLayouterFromFunc(layout)
}

func layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
