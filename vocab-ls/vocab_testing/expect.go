package vocab_testing

import (
	"testing"
)

func Expect[T comparable](t *testing.T, expectation T, gots ...T) {
	t.Helper()

	for i, got := range gots {
		if expectation == got {
			continue
		}

		t.Fatalf("Error: expected %+v, got %+v (index %d)", expectation, got, i)
	}
}
