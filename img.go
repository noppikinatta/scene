package scene

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	dummyImageBase = ebiten.NewImage(3, 3)

	// DummyWhitePixel is a 1x1 white pixel image.
	DummyWhitePixel = dummyImageBase.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	dummyImageBase.Fill(color.White)
}
