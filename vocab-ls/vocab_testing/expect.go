package vocab_testing

import (
	"testing"
)

func Expect[T comparable](t *testing.T, expectation T, got T) {
	t.Helper()

	if expectation == got {
		return
	}
	t.Fatalf("Error: expected %+v, got %+v", expectation, got)
}
