package super_memo

import (
	"fmt"
	"testing"
)

const initialEase = 2.5

func TestPlayground(t *testing.T) {
	r1, i1, e1 := Sm2(3, 0, 1, initialEase)
	fmt.Printf("%d, %f, %f\n", r1, i1, e1)

	r2, i2, e2 := Sm2(3, r1, i1, e1)
	fmt.Printf("%d, %f, %f\n", r2, i2, e2)

	r3, i3, e3 := Sm2(3, r2, i2, e2)
	fmt.Printf("%d, %f, %f\n", r3, i3, e3)

	r4, i4, e4 := Sm2(3, r3, i3, e3)
	fmt.Printf("%d, %f, %f\n", r4, i4, e4)
}
