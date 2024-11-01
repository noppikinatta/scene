package scene_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/scene"
)

func TestSequence(t *testing.T) {
	cases := []struct {
		Name        string
		GameFn      func() (*scene.Sequence, *recorder)
		ExpectedLog []string
	}{
		{
			Name: "simple",
			GameFn: func() (*scene.Sequence, *recorder) {
				r := recorder{}

				s1 := gameForTest{Name: "s1", recorder: &r}
				s2 := gameForTest{Name: "s2", recorder: &r}

				seq := scene.NewSequence(&s1)

				s1.UpdateFn = func() error {
					seq.Switch(&s2)
					return nil
				}

				return seq, &r
			},
			ExpectedLog: []string{
				"s1:layout",
				"s1:update",
				"s1:draw",
				"s1:layout",
				"s2:update",
			},
		},
		{
			Name: "event-handlers",
			GameFn: func() (*scene.Sequence, *recorder) {
				r := recorder{}

				s1 := eventsForTest{gameForTest: gameForTest{Name: "s1", recorder: &r}}
				s2 := eventsForTest{gameForTest: gameForTest{Name: "s2", recorder: &r}}

				seq := scene.NewSequence(&s1)

				s1.UpdateFn = func() error {
					seq.Switch(&s2)
					return nil
				}

				return seq, &r
			},
			ExpectedLog: []string{
				"s1:layout",
				"s1:onstart",
				"s1:onarrival",
				"s1:update",
				"s1:ondeparture",
				"s1:draw",
				"s1:layout",
				"s1:onend",
				"s2:onstart",
				"s2:onarrival",
				"s2:update",
				"s2:ondeparture",
				"s2:onend",
			},
		},
		{
			Name: "finalscreendrawer",
			GameFn: func() (*scene.Sequence, *recorder) {
				r := recorder{}

				s1 := gameForTest{Name: "s1", recorder: &r}
				s2 := finalScreenDrawerForTest{gameForTest: gameForTest{Name: "s2", recorder: &r}}

				seq := scene.NewSequence(&s1)

				s1.UpdateFn = func() error {
					seq.Switch(&s2)
					return nil
				}

				s2Counter := 0
				s2.UpdateFn = func() error {
					s2Counter++
					if s2Counter <= 1 {
						return nil
					}

					return ebiten.Termination
				}

				return seq, &r
			},
			ExpectedLog: []string{
				"s1:layout",
				"s1:update",
				"s1:draw",
				"s1:layout",
				"s2:update",
				"s2:draw",
				"s2:drawfinalscreen",
				"s2:layout",
				"s2:update",
			},
		},
		{
			Name: "layoutfer",
			GameFn: func() (*scene.Sequence, *recorder) {
				r := recorder{}

				s1 := gameForTest{Name: "s1", recorder: &r}
				s2 := layoutFerForTest{gameForTest: gameForTest{Name: "s2", recorder: &r}}

				seq := scene.NewSequence(&s1)

				s1.UpdateFn = func() error {
					seq.Switch(&s2)
					return nil
				}

				s2Counter := 0
				s2.UpdateFn = func() error {
					s2Counter++
					if s2Counter <= 1 {
						return nil
					}

					return ebiten.Termination
				}

				return seq, &r
			},
			ExpectedLog: []string{
				"s1:layout",
				"s1:update",
				"s1:draw",
				"s1:layout",
				"s2:update",
				"s2:draw",
				"s2:layoutf",
				"s2:update",
			},
		},
		{
			Name: "transition",
			GameFn: func() (*scene.Sequence, *recorder) {
				r := recorder{}

				s1 := eventsForTest{gameForTest: gameForTest{Name: "s1", recorder: &r}}
				s2 := eventsForTest{gameForTest: gameForTest{Name: "s2", recorder: &r}}

				seq := scene.NewSequence(&s1)

				tran := transitionForTest{Name: "t1", switchFrames: 3, maxFrames: 5, recorder: &r}

				s1.UpdateFn = func() error {
					seq.SwitchWithTransition(&s2, &tran)
					return nil
				}

				s2Counter := 0
				s2.UpdateFn = func() error {
					s2Counter++
					if s2Counter <= 3 {
						return nil
					}

					return ebiten.Termination
				}

				return seq, &r
			},
			ExpectedLog: []string{
				"s1:layout",
				"s1:onstart",
				"s1:onarrival",
				"s1:update",
				"t1:reset",
				"s1:ondeparture",
				"s1:draw",
				"t1:draw",
				"s1:layout",
				"t1:update",
				"s1:update",
				"s1:draw",
				"t1:draw",
				"s1:layout",
				"t1:update",
				"s1:update",
				"s1:draw",
				"t1:draw",
				"s1:layout",
				"t1:update",
				"s1:onend",
				"s2:onstart",
				"s2:update",
				"s2:draw",
				"t1:draw",
				"s2:layout",
				"t1:update",
				"s2:update",
				"s2:draw",
				"t1:draw",
				"s2:layout",
				"t1:update",
				"s2:onarrival",
				"s2:update",
				"s2:draw",
				"s2:layout",
				"s2:update",
				"s2:ondeparture",
				"s2:onend",
			},
		},
		{
			Name: "nested-sequence",
			GameFn: func() (*scene.Sequence, *recorder) {
				r := recorder{}

				s1 := gameForTest{Name: "s1", recorder: &r}
				s2 := gameForTest{Name: "s2", recorder: &r}

				seq := scene.NewSequence(&s1)

				s1.UpdateFn = func() error {
					seq.Switch(&s2)
					return nil
				}

				return seq, &r
			},
			ExpectedLog: []string{
				"s1:layout",
				"s1:update",
				"s1:draw",
				"s1:layout",
				"s2:update",
			},
		},
		// nested sequence
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			seq, recorder := c.GameFn()

			runForTest(t, seq)

			exLen := len(c.ExpectedLog)
			acLen := len(recorder.logs)

			if exLen != acLen {
				t.Errorf("expected log length %d, but got %d", exLen, acLen)
			}

			l := exLen
			if acLen < exLen {
				l = acLen
			}

			for i := 0; i < l; i++ {
				exLog := c.ExpectedLog[i]
				acLog := recorder.logs[i]

				if exLog != acLog {
					t.Errorf("%d: log different\nex: %s\nac: %s", i, exLog, acLog)
				}
			}
		})
	}
}

func TestSequenceLayoutReturnValue(t *testing.T) {

}

func TestSequenceLayoutFReturnValue(t *testing.T) {

}
