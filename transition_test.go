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

	s2Counter := 0
	s2.UpdateFn = func() error {
		s2Counter++
		if s2Counter <= 10 {
			return nil
		}

		return ebiten.Termination
	}

	runForTest(t, seq)

	compareLogs(t, []string{
		"t:0 5 false",
		"t:1 5 false",
		"t:2 5 true",
		"t:3 5 false",
		"t:4 5 false",
		"t:5 5 false",
	}, r.Log)
}

type linearTransitionDrawerForTest struct {
	Recorder *recorder
}

// Draw draws as the LinearTransition progresses.
func (d *linearTransitionDrawerForTest) Draw(screen *ebiten.Image, progress bamenn.LinearTransitionProgress) {
	d.Recorder.Append("t", fmt.Sprint(progress.CurrentFrame, progress.MaxFrames, progress.FrameToSwitch))
}
