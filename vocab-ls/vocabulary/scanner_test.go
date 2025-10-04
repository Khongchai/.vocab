package vocabulary

import (
	"testing"
)

type ScanExpect struct {
	TextValue  string
	TokenValue Token
	Line       int
	LineOffset int
	Pos        int
}

type ScanCase struct {
	Text         string
	Expectations []ScanExpect
}

func testScanExpectations(t *testing.T, testCases []ScanCase) {
	t.Helper()

	for i := range testCases {

		currentCase := testCases[i]
		scanner := NewScanner(currentCase.Text)

		for j := range currentCase.Expectations {
			// fmt.Printf("Testing case %+v\n", currentCase)

			expectation := currentCase.Expectations[j]
			actualToken, actualText := scanner.Scan()
			actualLine := scanner.line
			actualOffset := scanner.lineOffset
			actualPos := scanner.pos

			if actualLine != expectation.Line {
				t.Fatalf("Line does not match expectation, expected %d, got %d", expectation.Line, actualLine)
			}

			if actualToken != expectation.TokenValue {
				t.Fatalf("Token does not match, expected %+v, got %+v", expectation.TokenValue, actualToken)
			}

			if actualText != expectation.TextValue {
				t.Fatalf("Text does not match: expected %s, got %s", expectation.TextValue, actualText)
			}

			if actualOffset != expectation.LineOffset {
				t.Fatalf("Line offset does not match the offset of expected token text. Expected %d, got %d", expectation.LineOffset, actualOffset)
			}

			if actualPos != expectation.Pos {
				t.Fatalf("Offset does not match the length of expected token text. Expected %d, got %d", expectation.Pos, actualPos)
			}
		}

	}
}

// Test against all recognized characters in the vocab syntax
func TestBasicTokenScan(t *testing.T) {
	testScanExpectations(t, []ScanCase{
		{Text: "ÄäöÖé", Expectations: []ScanExpect{{TextValue: "ÄäöÖé", TokenValue: TokenText, LineOffset: 10, Pos: 10}}},
		{Text: "Hello", Expectations: []ScanExpect{{TextValue: "Hello", TokenValue: TokenText, LineOffset: 5, Pos: 5}}},
		{Text: "20/08/2025", Expectations: []ScanExpect{{TextValue: "20/08/2025", TokenValue: TokenDateExpression, LineOffset: 10, Pos: 10}}},
		{Text: "23/00", Expectations: []ScanExpect{{TextValue: "23/00", TokenValue: TokenText, LineOffset: 5, Pos: 5}}},
		{Text: ">", Expectations: []ScanExpect{{TextValue: ">", TokenValue: TokenGreaterThan, LineOffset: 1, Pos: 1}}},
		{Text: ">>", Expectations: []ScanExpect{{TextValue: ">>", TokenValue: TokenDoubleGreaterThan, LineOffset: 2, Pos: 2}}},
		{Text: ",", Expectations: []ScanExpect{{TextValue: ",", TokenValue: TokenComma, LineOffset: 1, Pos: 1}}},
		{Text: "`foo`", Expectations: []ScanExpect{{TextValue: "foo", TokenValue: TokenWordLiteral, LineOffset: 5, Pos: 5}}},
		{Text: "`foo", Expectations: []ScanExpect{{TextValue: "foo", TokenValue: TokenWordLiteral, LineOffset: 4, Pos: 4}}},
		{Text: "(it)", Expectations: []ScanExpect{{TextValue: "it", TokenValue: TokenLanguageLiteral, LineOffset: 4, Pos: 4}}},
		{Text: "(i)", Expectations: []ScanExpect{{TextValue: "(i)", TokenValue: TokenText, LineOffset: 3, Pos: 3}}},
		{Text: "/", Expectations: []ScanExpect{{TextValue: "/", TokenValue: TokenSlash, LineOffset: 1, Pos: 1}}},
		{Text: "<!--", Expectations: []ScanExpect{{TextValue: "<!--", TokenValue: TokenMarkdownCommentStart, LineOffset: 4, Pos: 4}}},
		{Text: "-->", Expectations: []ScanExpect{{TextValue: "-->", TokenValue: TokenMarkdownCommentEnd, LineOffset: 3, Pos: 3}}},
		{Text: "-", Expectations: []ScanExpect{{TextValue: "-", TokenValue: TokenMinus, LineOffset: 1, Pos: 1}}},
	})

}

func TestNewline(t *testing.T) {
	testScanExpectations(t, []ScanCase{
		{Text: "\n", Expectations: []ScanExpect{{TextValue: "\n", TokenValue: TokenLineBreak, LineOffset: 0, Line: 1, Pos: 1}}},
		{Text: "\r", Expectations: []ScanExpect{{TextValue: "\r", TokenValue: TokenLineBreak, LineOffset: 0, Line: 1, Pos: 1}}},
		{Text: "Hello \nWorld!",
			Expectations: []ScanExpect{
				{TextValue: "Hello", TokenValue: TokenText, LineOffset: 5, Line: 0, Pos: 5},
				{TextValue: " ", TokenValue: TokenWhitespace, LineOffset: 6, Line: 0, Pos: 6},
				{TextValue: "\n", TokenValue: TokenLineBreak, LineOffset: 0, Line: 1, Pos: 7},
				{TextValue: "World", TokenValue: TokenText, LineOffset: 5, Line: 1, Pos: 12},
				{TextValue: "!", TokenValue: TokenText, LineOffset: 6, Line: 1, Pos: 13},
			},
		},
	})
}

func TestNewVocabScan(t *testing.T) {
	testScanExpectations(t, []ScanCase{
		{
			Text: ">Hello, World",
			Expectations: []ScanExpect{
				{TextValue: ">", TokenValue: TokenGreaterThan, Line: 0, LineOffset: 1, Pos: 1},
				{TextValue: "Hello", TokenValue: TokenText, Line: 0, LineOffset: 6, Pos: 6},
				{TextValue: ",", TokenValue: TokenComma, Line: 0, LineOffset: 7, Pos: 7},
				{TextValue: " ", TokenValue: TokenWhitespace, LineOffset: 8, Line: 0, Pos: 8},
				{TextValue: "World", TokenValue: TokenText, Line: 0, LineOffset: 13, Pos: 13},
			},
		},
		{
			Text: "> Hello, World",
			Expectations: []ScanExpect{
				{TextValue: ">", TokenValue: TokenGreaterThan, Line: 0, LineOffset: 1, Pos: 1},
				{TextValue: " ", TokenValue: TokenWhitespace, Line: 0, LineOffset: 2, Pos: 2},
				{TextValue: "Hello", TokenValue: TokenText, Line: 0, LineOffset: 7, Pos: 7},
				{TextValue: ",", TokenValue: TokenComma, Line: 0, LineOffset: 8, Pos: 8},
				{TextValue: " ", TokenValue: TokenWhitespace, Line: 0, LineOffset: 9, Pos: 9},
				{TextValue: "World", TokenValue: TokenText, Line: 0, LineOffset: 14, Pos: 14},
			},
		},
		{
			Text: "> (it) `ciao`, bello!",
			Expectations: []ScanExpect{
				{TextValue: ">", TokenValue: TokenGreaterThan, Line: 0, LineOffset: 1, Pos: 1},
				{TextValue: " ", TokenValue: TokenWhitespace, Line: 0, LineOffset: 2, Pos: 2},
				{TextValue: "it", TokenValue: TokenLanguageLiteral, Line: 0, LineOffset: 6, Pos: 6},
				{TextValue: " ", TokenValue: TokenWhitespace, Line: 0, LineOffset: 7, Pos: 7},
				{TextValue: "ciao", TokenValue: TokenWordLiteral, Line: 0, LineOffset: 13, Pos: 13},
				{TextValue: ",", TokenValue: TokenComma, Line: 0, LineOffset: 14, Pos: 14},
				{TextValue: " ", TokenValue: TokenWhitespace, Line: 0, LineOffset: 15, Pos: 15},
				{TextValue: "bello", TokenValue: TokenText, Line: 0, LineOffset: 20, Pos: 20},
				{TextValue: "!", TokenValue: TokenText, Line: 0, LineOffset: 21, Pos: 21},
			},
		},
	})
}

func TestReviewedVocabScan(t *testing.T) {
	testScanExpectations(t, []ScanCase{
		{
			Text: ">>Hello, World",
			Expectations: []ScanExpect{
				{TextValue: ">>", TokenValue: TokenDoubleGreaterThan, Line: 0, LineOffset: 2, Pos: 2},
				{TextValue: "Hello", TokenValue: TokenText, Line: 0, LineOffset: 7, Pos: 7},
				{TextValue: ",", TokenValue: TokenComma, Line: 0, LineOffset: 8, Pos: 8},
				{TextValue: " ", TokenValue: TokenWhitespace, Line: 0, LineOffset: 9, Pos: 9},
				{TextValue: "World", TokenValue: TokenText, Line: 0, LineOffset: 14, Pos: 14},
			},
		},
		{
			Text: ">> Hello, World",
			Expectations: []ScanExpect{
				{TextValue: ">>", TokenValue: TokenDoubleGreaterThan, Line: 0, LineOffset: 2, Pos: 2},
				{TextValue: " ", TokenValue: TokenWhitespace, Line: 0, LineOffset: 3, Pos: 3},
				{TextValue: "Hello", TokenValue: TokenText, Line: 0, LineOffset: 8, Pos: 8},
				{TextValue: ",", TokenValue: TokenComma, Line: 0, LineOffset: 9, Pos: 9},
				{TextValue: " ", TokenValue: TokenWhitespace, Line: 0, LineOffset: 10, Pos: 10},
				{TextValue: "World", TokenValue: TokenText, Line: 0, LineOffset: 15, Pos: 15},
			},
		},
		{
			Text: ">> (de) halo, schön!",
			Expectations: []ScanExpect{
				{TextValue: ">>", TokenValue: TokenDoubleGreaterThan, Line: 0, LineOffset: 2, Pos: 2},
				{TextValue: " ", TokenValue: TokenWhitespace, Line: 0, LineOffset: 3, Pos: 3},
				{TextValue: "de", TokenValue: TokenLanguageLiteral, Line: 0, LineOffset: 7, Pos: 7},
				{TextValue: " ", TokenValue: TokenWhitespace, Line: 0, LineOffset: 8, Pos: 8},
				{TextValue: "halo", TokenValue: TokenText, Line: 0, LineOffset: 12, Pos: 12},
				{TextValue: ",", TokenValue: TokenComma, Line: 0, LineOffset: 13, Pos: 13},
				{TextValue: " ", TokenValue: TokenWhitespace, Line: 0, LineOffset: 14, Pos: 14},
				{TextValue: "schön", TokenValue: TokenText, Line: 0, LineOffset: 20, Pos: 20},
				{TextValue: "!", TokenValue: TokenText, Line: 0, LineOffset: 21, Pos: 21},
			},
		},
	})
}
