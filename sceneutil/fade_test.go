package sceneutil_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/scene/sceneutil"
)

func TestFade(t *testing.T) {
	pd := dummyProgressDrawer{}
	fade := sceneutil.NewFade(3, &pd)

	fade.Init()
	fade.Update()
	fade.Draw(nil)

	if pd.progress < 0.332 || pd.progress > 0.334 {
		t.Errorf("expected progress: 0.332~0.334, actual progress: %f", pd.progress)
	}
	if fade.Ended() {
		t.Error("falde should not end before complete fading 1/3")
	}

	fade.Update()
	fade.Draw(nil)

	if pd.progress < 0.665 || pd.progress > 0.667 {
		t.Errorf("expected progress: 0.665~0.667, actual progress: %f", pd.progress)
	}
	if fade.Ended() {
		t.Error("falde should not end before complete fading 2/3")
	}

	fade.Update()
	fade.Draw(nil)

	if pd.progress < 1 {
		t.Errorf("expected progress: 1, actual progress: %f", pd.progress)
	}
	if !fade.Ended() {
		t.Error("falde should end after complete fading 3/3")
	}
}

type dummyProgressDrawer struct {
	progress float64
}

func (d *dummyProgressDrawer) Draw(screen *ebiten.Image, progress float64) {
	d.progress = progress
}
