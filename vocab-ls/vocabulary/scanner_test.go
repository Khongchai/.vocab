package vocabulary

import "testing"

type Expectation struct {
	TextValue  string
	TokenValue Token
	Line       int
	LineOffset int
	Pos        int
}

type TestCase struct {
	Text         string
	Expectations []Expectation
}

// Test against all recognized characters in the vocab syntax
func TestBasicTokenScan(t *testing.T) {
	testCases := []TestCase{
		{
			Text: "Hello",
			Expectations: []Expectation{
				{
					TextValue:  "Hello",
					TokenValue: TokenTextLiteral,
					LineOffset: 5,
					Pos:        5,
				},
			},
		},
		{
			Text: "234",
			Expectations: []Expectation{
				{
					TextValue:  "234",
					TokenValue: TokenNumericLiteral,
					LineOffset: 3,
					Pos:        3,
				},
			},
		},
		{
			Text: ">",
			Expectations: []Expectation{
				{
					TextValue:  ">",
					TokenValue: TokenGreaterThan,
					LineOffset: 1,
					Pos:        1,
				},
			},
		},
		{
			Text: ">>",
			Expectations: []Expectation{
				{
					TextValue:  ">>",
					TokenValue: TokenDoubleGreaterThan,
					LineOffset: 2,
					Pos:        2,
				},
			},
		},
		{
			Text: ",",
			Expectations: []Expectation{
				{
					TextValue:  ",",
					TokenValue: TokenComma,
					LineOffset: 1,
					Pos:        1,
				},
			},
		},
		{
			Text: "`",
			Expectations: []Expectation{
				{
					TextValue:  "`",
					TokenValue: TokenBacktick,
					LineOffset: 1,
					Pos:        1,
				},
			},
		},
		{
			Text: "(",
			Expectations: []Expectation{
				{
					TextValue:  "(",
					TokenValue: TokenLeftParen,
					LineOffset: 1,
					Pos:        1,
				},
			},
		},
		{
			Text: ")",
			Expectations: []Expectation{
				{
					TextValue:  ")",
					TokenValue: TokenRightParen,
					LineOffset: 1,
					Pos:        1,
				},
			},
		},
		{
			Text: "/",
			Expectations: []Expectation{
				{
					TextValue:  "/",
					TokenValue: TokenSlash,
					LineOffset: 1,
					Pos:        1,
				},
			},
		},
		{
			Text: "```",
			Expectations: []Expectation{
				{
					TextValue:  "```",
					TokenValue: TokenMarkdownCodefence,
					LineOffset: 3,
					Pos:        3,
				},
			},
		},
		{
			Text: "<!--",
			Expectations: []Expectation{
				{
					TextValue:  "<!--",
					TokenValue: TokenMarkdownCommentStart,
					LineOffset: 4,
					Pos:        4,
				},
			},
		},
		{
			Text: "-->",
			Expectations: []Expectation{
				{
					TextValue:  "-->",
					TokenValue: TokenMarkdownCommentEnd,
					LineOffset: 3,
					Pos:        3,
				},
			},
		},
		{
			Text: "-",
			Expectations: []Expectation{
				{
					TextValue:  "-",
					TokenValue: TokenMinus,
					LineOffset: 1,
					Pos:        1,
				},
			},
		},
	}

	for i := range testCases {
		currentCase := testCases[i]
		scanner := NewScanner(currentCase.Text)

		for j := range currentCase.Expectations {
			expectation := currentCase.Expectations[j]
			actualToken, actualText := scanner.Scan()
			actualLine := scanner.line
			actualOffset := scanner.lineOffset
			actualPos := scanner.pos

			if actualLine != expectation.Line {
				t.Fatalf("Line does not match expectation, expected %d, got %d", expectation.Line, actualLine)
			}

			if actualToken != expectation.TokenValue {
				t.Fatalf("Token does not match, expected %d, got %d", expectation.TokenValue, actualToken)
			}

			if actualText != expectation.TextValue {
				t.Fatalf("Text does not match: expected %s, got %s", expectation.TextValue, actualText)
			}

			if actualOffset != expectation.LineOffset {
				t.Fatalf("Line offset does not match the length of expected token text. Expected %d, got %d", expectation.LineOffset, actualOffset)
			}

			if actualPos != expectation.Pos {
				t.Fatalf("Offset does not match the length of expected token text. Expected %d, got %d", expectation.Pos, actualPos)
			}
		}

	}
}

func TestNewline(t *testing.T) {
	// testCases := testCase{
	// 	"\n":       {TokenLineBreak},
	// 	"Hello \n": {TokenTextLiteral, TokenLineBreak},
	// 	"12345 \n": {TokenNumericLiteral, TokenLineBreak},
	// 	"\r\n":     {TokenLineBreak},
	// }

	// for tokenText := range testCases {
	// 	scanner := NewScanner(tokenText)

	// 	for i := range testCases[tokenText] {
	// 		actualToken, actualText := scanner.Scan()
	// 		expectedToken = tokenText

	// 		if scanner.line != 0 {
	// 			t.Fatalf("Line should still be zero")
	// 		}

	// 		if actualToken != testCases[tokenText][i] {
	// 			t.Fatalf("Token does not match, expected %d, got %d", actualToken, expected)
	// 		}

	// 		if actualText != tokenText {
	// 			t.Fatalf("Text does not match: expected %s, got %s", tokenText, actualText)
	// 		}

	// 		if scanner.lineOffset != len(tokenText) {
	// 			t.Fatalf("Line offset does not match the length of expected token text. Expected %d, got %d", scanner.lineOffset, len(tokenText))
	// 		}

	// 		if scanner.pos != len(tokenText) {
	// 			t.Fatalf("Offset does not match the length of expected token text. Expected %d, got %d", scanner.lineOffset, len(tokenText))
	// 		}
	// 	}

	// }
}

func TestDateScan(t *testing.T) {
	// dateCapture := []Token{
	// 	TokenNumericLiteral,
	// 	TokenSlash,
	// 	TokenNumericLiteral,
	// 	TokenSlash,
	// 	TokenNumericLiteral,
	// }
	// testCases := map[string][]Token{
	// 	"00/00/0000":  dateCapture,
	// 	"#20/05/2025": append([]Token{TokenTextLiteral}, dateCapture...),
	// 	"# 20/05/2025": append([]Token{TokenTextLiteral}, dateCapture...),
	// }

	// for tokenText := range testCases {
	// 	scanner := NewScanner(tokenText)

	// 	actual, actualText := scanner.Scan()
	// 	expected := testCases[tokenText]

	// 	if actual != testCases[tokenText] {
	// 		t.Fatalf("Token does not match, expected %d, got %d", actual, expected)
	// 	}

	// 	if actualText != tokenText {
	// 		t.Fatalf("Text does not match: expected %s, got %s", tokenText, actualText)
	// 	}
	// }
}

func TestNewVocabScan(t *testing.T) {

}

func TestReviewedVocabScan(t *testing.T) {

}

func TestFullSectionScan(t *testing.T) {

}
