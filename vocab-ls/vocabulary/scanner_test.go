package vocabulary

import "testing"

// Test against all recognized characters in the vocab syntax
func TestBasicTokenScan(t *testing.T) {
	testCases := map[string]Token{
		">":   TokenGreaterThan,
		">>":  TokenDoubleGreaterThan,
		",":   TokenComma,
		"`":   TokenBacktick,
		"(":   TokenOpenBracket,
		")":   TokenCloseBracket,
		"/":   TokenSlash,
		"```": TokenMarkdownCodefence,
	}

	scanner := NewScanner("xxx")
	actual, actualText := scanner.Scan()

	if actualText != "" {
		t.Errorf("Before scanning, actual text should be empty")
	}
	if actual != TokenUnknown {
		t.Errorf("Before scanning, scanner token should be unknown")
	}

	for tokenText := range testCases {
		scanner := NewScanner(tokenText)

		actual, actualText = scanner.Scan()
		expected := testCases[tokenText]

		if actualText != tokenText {
			t.Errorf("Text does not match: expected %s, got %s", tokenText, actualText)
		}
		if actual != testCases[tokenText] {
			t.Errorf("Token does not match, expected %d, got %d", actual, expected)
		}
	}
}

func TestDateScan(t *testing.T) {
}

func TestNewVocabScan(t *testing.T) {

}

func TestReviewedVocabScan(t *testing.T) {

}

func TestFullSectionScan(t *testing.T) {

}
