package bamenn_test

import (
	"fmt"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/bamenn"
)

func TestTransition(t *testing.T) {
	r := recorder{}

	s1 := eventsForTest{gameForTest: gameForTest{Name: "s1"}}
	s2 := eventsForTest{gameForTest: gameForTest{Name: "s2"}}

	seq := bamenn.NewSequence(&s1)

	tran := bamenn.NewLinearTransition(2, 5, &linearTransitionDrawerForTest{Recorder: &r})

	s1.UpdateFn = func() error {
		seq.SwitchWithTransition(&s2, tran)
		return nil
	}

	canEndS2 := false
	s2.OnArrivalFn = func() {
		canEndS2 = true
	}
	s2.UpdateFn = func() error {
		if canEndS2 {
			return ebiten.Termination
		}
		return nil
	}

	runForTest(t, seq)

	compareLogs(t, []string{
		"t:0 5 2 0.0",
		"t:1 5 2 0.2",
		"t:2 5 2 0.4",
		"t:3 5 2 0.6",
		"t:4 5 2 0.8",
	}, r.Log)
}

type linearTransitionDrawerForTest struct {
	Recorder *recorder
}

// Draw draws as the LinearTransition progresses.
func (d *linearTransitionDrawerForTest) Draw(screen *ebiten.Image, progress bamenn.LinearTransitionProgress) {
	d.Recorder.Append(
		"t",
		fmt.Sprintf(
			"%d %d %d %.1f",
			progress.CurrentFrame,
			progress.MaxFrames,
			progress.FrameToSwitch,
			progress.Rate(),
		),
	)
}
