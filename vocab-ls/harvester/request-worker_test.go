package harvester

import "testing"

func TestFilePathTransform(t *testing.T) {
	input := "c:\\Users\\world\\Desktop\\vocab\\test.vocab"
	expect := "file:///c%3A/Users/world/Desktop/vocab/test.vocab"

	result := TransformWindowsPathToLspUri(input)
	if result != expect {
		t.Fatalf("Invalid file transform. Expected '%s', got '%s'", expect, result)
	}
}
