package scene_test

import (
	"testing"

	"github.com/noppikinatta/scene"
)

// Sequence

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
				"s1:update",
				"s1:draw",
				"s2:update",
			},
		},
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
