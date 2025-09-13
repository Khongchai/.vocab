package vocabulary

import "testing"

// Test against all recognized characters in the vocab syntax
func TestBasicTokenScan(t *testing.T) {
	testCases := map[string]Token{
		">":    TokenGreaterThan,
		">>":   TokenDoubleGreaterThan,
		",":    TokenComma,
		"`":    TokenBacktick,
		"(":    TokenOpenBracket,
		")":    TokenCloseBracket,
		"(it)": TokenItalianKeyword,
		"(de)": TokenGermanKeyword,
		"/":    TokenSlash,
	}

	scanner := NewScanner("xxx")
	if scanner.CurrentToken() != nil {
		t.Errorf("Before scanning, scanner token should be nil")
	}

	for tokenText := range testCases {
		scanner := NewScanner(tokenText)
		scanner.NextToken()

		actual := *scanner.CurrentToken()
		expected := testCases[tokenText]
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
