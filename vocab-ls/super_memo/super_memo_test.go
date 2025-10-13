package super_memo

import (
	"fmt"
	"testing"
	test "vocab/vocab_testing"
)

const initialEase = 2.5

func TestPlayground(t *testing.T) {
	r, i, e := int(0.0), float64(0.0), float64(initialEase)
	grade := 5
	for range 10 {
		r, i, e = Sm2(grade, r, i, e)
		fmt.Printf("repetition: %d, interval: %f, easiness: %f\n", r, i, e)
	}
}

func TestKnownInputsAndOutputs_IntervalAlwaysOneAndRepetitionZero_IfGradeLessThan3(t *testing.T) {
	r1, i1, e1 := Sm2(0, 0, 0, initialEase)
	r2, i2, e2 := Sm2(1, r1, i1, e1)
	_, i3, _ := Sm2(2, r2, i2, e2)

	test.Expect(t, true, i3 == i2)
	test.Expect(t, true, i2 == i1)
}

func TestKnownInputsAndOutputs_DistanceKeepsGoingUpOnceGrade3OrMore(t *testing.T) {
	r1, i1, e1 := Sm2(3, 0, 0, initialEase)
	r2, i2, e2 := Sm2(4, r1, i1, e1)
	r3, i3, e3 := Sm2(5, r2, i2, e2)
	_, i4, _ := Sm2(5, r3, i3, e3)

	test.Expect(t, true, i4 > i3)
	test.Expect(t, true, i3 > i2)
	test.Expect(t, true, i2 > i1)
}
