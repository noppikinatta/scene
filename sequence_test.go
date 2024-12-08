package bamenn_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/bamenn"
)

func TestSequence(t *testing.T) {
	cases := []struct {
		Name        string
		GameFn      func() (*bamenn.Sequence, *recorder)
		ExpectedLog []string
	}{
		{
			Name: "simple",
			GameFn: func() (*bamenn.Sequence, *recorder) {
				r := recorder{}

				s1 := gameForTest{Name: "s1", Recorder: &r}
				s2 := gameForTest{Name: "s2", Recorder: &r}

				seq := &bamenn.Sequence{}
				seq.SetFirst(&s1)

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
			GameFn: func() (*bamenn.Sequence, *recorder) {
				r := recorder{}

				s1 := eventsForTest{gameForTest: gameForTest{Name: "s1", Recorder: &r}}
				s2 := eventsForTest{gameForTest: gameForTest{Name: "s2", Recorder: &r}}

				seq := bamenn.NewSequence(&s1)

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
			},
		},
		{
			Name: "finalscreendrawer",
			GameFn: func() (*bamenn.Sequence, *recorder) {
				r := recorder{}

				s1 := gameForTest{Name: "s1", Recorder: &r}
				s2 := finalScreenDrawerForTest{gameForTest: gameForTest{Name: "s2", Recorder: &r}}

				seq := bamenn.NewSequence(&s1)

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
			GameFn: func() (*bamenn.Sequence, *recorder) {
				r := recorder{}

				s1 := gameForTest{Name: "s1", Recorder: &r}
				s2 := layoutFerForTest{gameForTest: gameForTest{Name: "s2", Recorder: &r}}

				seq := bamenn.NewSequence(&s1)

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
			GameFn: func() (*bamenn.Sequence, *recorder) {
				r := recorder{}

				s1 := eventsForTest{gameForTest: gameForTest{Name: "s1", Recorder: &r}}
				s2 := eventsForTest{gameForTest: gameForTest{Name: "s2", Recorder: &r}}

				seq := bamenn.NewSequence(&s1)

				tran := transitionForTest{Name: "t1", SwitchFrames: 3, MaxFrames: 5, Recorder: &r}

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
			},
		},
		{
			Name: "nested-sequence",
			GameFn: func() (*bamenn.Sequence, *recorder) {
				r := recorder{}

				s1 := eventsForTest{gameForTest: gameForTest{Name: "s1", Recorder: &r}}
				s21 := eventsForTest{gameForTest: gameForTest{Name: "s21", Recorder: &r}}
				s22 := eventsForTest{gameForTest: gameForTest{Name: "s22", Recorder: &r}}
				seq2 := bamenn.NewSequence(&s21)
				s3 := eventsForTest{gameForTest: gameForTest{Name: "s3", Recorder: &r}}

				seq := bamenn.NewSequence(&s1)

				s1count := 0
				s1.UpdateFn = func() error {
					if s1count < 1 {
						s1count++
						return nil
					}
					seq.Switch(seq2)
					return nil
				}

				s21count := 0
				s21.UpdateFn = func() error {
					if s21count < 1 {
						s21count++
						return nil
					}
					seq2.Switch(&s22)
					return nil
				}

				s22count := 0
				s22.UpdateFn = func() error {
					if s22count < 1 {
						s22count++
						return nil
					}
					seq.Switch(&s3)
					return nil
				}

				s3.UpdateFn = func() error {
					return ebiten.Termination
				}

				return seq, &r
			},
			ExpectedLog: []string{
				"s1:layout",
				"s1:onstart",
				"s1:onarrival",
				"s1:update",
				"s1:draw",
				"s1:layout",
				"s1:update",
				"s1:ondeparture",
				"s1:draw",
				"s1:layout",
				"s1:onend",
				"s21:onstart",
				"s21:onarrival",
				"s21:update",
				"s21:draw",
				"s21:layout",
				"s21:update",
				"s21:ondeparture",
				"s21:draw",
				"s21:layout",
				"s21:onend",
				"s22:onstart",
				"s22:onarrival",
				"s22:update",
				"s22:draw",
				"s22:layout",
				"s22:update",
				"s22:ondeparture",
				"s22:draw",
				"s22:layout",
				"s22:onend",
				"s3:onstart",
				"s3:onarrival",
				"s3:update",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			seq, recorder := c.GameFn()

			runForTest(t, seq)

			compareLogs(t, c.ExpectedLog, recorder.Log)
		})
	}
}

func TestSequenceLayoutReturnValue(t *testing.T) {
	s := gameForTest{LayoutW: 3, LayoutH: 3}
	seq := bamenn.NewSequence(&s)

	w, h := seq.Layout(1, 1)
	if w != 3 || h != 3 {
		t.Errorf("layout expected (3,3), but got (%d,%d)", w, h)
	}
}

func TestSequenceLayoutFReturnValue(t *testing.T) {
	cases := []struct {
		Name      string
		Game      ebiten.Game
		ExpectedW float64
		ExpectedH float64
	}{
		{
			Name:      "with-layoutf",
			Game:      &layoutFerForTest{gameForTest: gameForTest{LayoutW: 3, LayoutH: 3}, layoutFW: 4, layoutFH: 4},
			ExpectedW: 4,
			ExpectedH: 4,
		},
		{
			Name:      "without-layoutf",
			Game:      &gameForTest{LayoutW: 3, LayoutH: 3},
			ExpectedW: 3,
			ExpectedH: 3,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			seq := bamenn.NewSequence(c.Game)

			w, h := seq.LayoutF(1, 1)
			if w != c.ExpectedW || h != c.ExpectedH {
				t.Errorf("layout expected (%f,%f), but got (%f,%f)", c.ExpectedW, c.ExpectedH, w, h)
			}
		})
	}
}
