package vocabulary

import "testing"

// Test against all recognized characters in the vocab syntax
func TestBasicTokenScan(t *testing.T) {
	testCases := map[string]Token{
		"Hello": TokenTextLiteral,
		"234":   TokenNumericLiteral,
		">":     TokenGreaterThan,
		">>":    TokenDoubleGreaterThan,
		",":     TokenComma,
		"`":     TokenBacktick,
		"(":     TokenOpenBracket,
		")":     TokenCloseBracket,
		"/":     TokenSlash,
		"```":   TokenMarkdownCodefence,
		"<!--":  TokenMarkdownCommentStart,
		"-->":   TokenMarkdownCommentEnd,
	}

	scanner := NewScanner("xxx")
	actual, actualText := scanner.Scan()

	if actualText != "" {
		t.Fatalf("Before scanning, actual text should be empty")
	}

	if actual != TokenUnknown {
		t.Fatalf("Before scanning, scanner token should be unknown")
	}

	for tokenText := range testCases {
		scanner := NewScanner(tokenText)

		actual, actualText = scanner.Scan()
		expected := testCases[tokenText]

		if actualText != tokenText {
			t.Fatalf("Text does not match: expected %s, got %s", tokenText, actualText)
		}
		if actual != testCases[tokenText] {
			t.Fatalf("Token does not match, expected %d, got %d", actual, expected)
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
