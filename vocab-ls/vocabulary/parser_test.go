package vocabulary

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func printJSON(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func testParseExpectation(t *testing.T, text string, sections []*VocabularySection, expectedErrors []ParsingError) {
	t.Helper()

	split := strings.Split(text, "\n")
	for i := range split {
		split[i] = strings.TrimSpace(split[i])
	}
	joined := strings.Join(split, "\n")

	var parsedErrors []any
	parser := NewParser(t.Context(), "xxx", NewScanner(joined), func(a any) {
		parsedErrors = append(parsedErrors, a)
	})

	// act
	parser.Parse()

	var result *VocabAst = parser.ast

	got := printJSON(result)

	newAst := &VocabAst{uri: "xxx", Sections: sections}
	expected := printJSON(newAst)

	if expected != got {
		t.Fatalf("Parsing test failed. Expected:\n%s \nGot:\n%s", expected, got)
	}

	if len(expectedErrors) != len(parsedErrors) {
		t.Fatalf("Parsing test failed. Expected error length != parsedErrors")
	}

	for i := range expectedErrors {
		if expectedErrors[i] != parsedErrors[i] {
			t.Fatalf("Parsing test failed. Expected error: %+v, got %+v", expectedErrors[i], parsedErrors[i])
		}
	}
}

func TestFullSectionParsing(t *testing.T) {
	testParseExpectation(t,
		`
			20/08/2025
			> (it) la magia, bene
			> (de) anlegen
			Ho una magia molto speciale. Non ti conviene metterti contro di me!
			Ne, will gar nicht mit ihm anlegen.
			21/08/2025
			> (it) brillare
			>> (it) la maga
			C'era una volta un piccolo villaggio in Italia. In questo villaggio, viveva una giovane maga. La maga si chiamava Luna, e il suo potere era molto semplice: poteva far brillare le stelle nel cielo. 
		`, []*VocabularySection{
			{
				Date: &DateSection{Text: "20/08/2025", Time: time.Date(2025, time.August, 20, 0, 0, 0, 0, time.Local), Start: 0, End: 10},
				NewWords: []*WordsSection{
					{
						Language: Italiano,
						Line:     1,
						Words: []*Word{
							{Text: "magia", FullText: "la magia", Start: 7, End: 14},
							{Text: "bene", FullText: "bene", Start: 17, End: 20},
						},
					},
					{
						Language: Deutsch,
						Line:     2,
						Words: []*Word{
							{Text: "anlegen", FullText: "anlegen", Start: 7, End: 13},
						},
					},
				},
				ReviewedWords: []*WordsSection{},
				Sentences: []*SentenceSection{
					{
						StartLine: 3,
						EndLine:   3,
						StartPos:  0,
						EndPos:    len("Ho una magia molto speciale. Non ti conviene metterti contro di me!"),
						Text:      "Ho una magia molto speciale. Non ti conviene metterti contro di me!",
					},
					{
						StartLine: 4,
						EndLine:   4,
						StartPos:  0,
						EndPos:    len("Ne, will gar nicht mit ihm anlegen."),
						Text:      "Ne, will gar nicht mit ihm anlegen.",
					},
				},
			},
			{
				Date: &DateSection{Text: "21/08/2025", Time: time.Date(2025, time.August, 20, 0, 0, 0, 0, time.Local), Start: 0, End: 10},
				NewWords: []*WordsSection{
					{
						Language: Italiano,
						Line:     5,
						Words: []*Word{
							{Text: "brillare", FullText: "brillare", Start: 7, End: 14},
						},
					},
				},
				ReviewedWords: []*WordsSection{
					{
						Language: Italiano,
						Line:     5,
						Words: []*Word{
							{Text: "maga", FullText: "la maga", Start: 7, End: 14},
						},
					},
				},
				Sentences: []*SentenceSection{
					{
						StartLine: 6,
						EndLine:   6,
						StartPos:  0,
						EndPos:    len("C'era una volta un piccolo villaggio in Italia. In questo villaggio, viveva una giovane maga. La maga si chiamava Luna, e il suo potere era molto semplice: poteva far brillare le stelle nel cielo."),
						Text:      "C'era una volta un piccolo villaggio in Italia. In questo villaggio, viveva una giovane maga. La maga si chiamava Luna, e il suo potere era molto semplice: poteva far brillare le stelle nel cielo. ",
					},
				},
			},
		}, []ParsingError{})
}

// Incomplete sections don't necessarily emit diagnostics error as missing vocabulary is already covered by the compiler.
func TestOnlyDateSection(t *testing.T) {
	type Expectation struct {
		Input      string
		ParsedDate time.Time
		Start      int
		End        int
		Error      ParsingError
	}
	expectations := []Expectation{
		{
			Input:      "20/08/2025",
			ParsedDate: time.Date(2025, time.August, 20, 0, 0, 0, 0, time.Local),
			Start:      0,
			End:        10,
		},
		{
			Input:      " 20/08/2025 ",
			ParsedDate: time.Date(2025, time.August, 20, 0, 0, 0, 0, time.Local),
			Start:      1,
			End:        11,
		},
		{
			Input:      "00/00/0000",
			ParsedDate: time.Date(1, 1, 1, 0, 0, 0, 0, time.Local),
			Start:      0,
			End:        10,
			Error:      MalformedDate,
		},
	}

	for _, expectation := range expectations {
		var parsedError ParsingError = ""
		parser := NewParser(t.Context(), "xxx", NewScanner(expectation.Input), func(a any) {
			parsedError = a.(ParsingError)
		})
		parser.Parse()
		gotTime := parser.ast.Sections[0].Date.Time
		if expectation.ParsedDate != gotTime {
			t.Fatalf("Date mismatched, expected: %+v, got %+v", expectation.ParsedDate, gotTime)
		}

		gotStart := parser.ast.Sections[0].Date.Start
		gotEnd := parser.ast.Sections[0].Date.End
		if gotStart != expectation.Start || gotEnd != expectation.End {
			t.Fatalf("Start and end don't match, expected %d, %d -- got %d, %d", expectation.Start, expectation.End, gotStart, gotEnd)
		}

		if expectation.Error != "" && expectation.Error != parsedError {
			t.Fatalf("Error mismatch, expected: %s, got %s", expectation.Error, parsedError)
		}
	}

	// var result *VocabAst = parser.ast
	// testParseExpectation(t,
	// 	` 20/08/2025`, []*VocabularySection{
	// 		{
	// 			Date: &DateSection{Text: "20/08/2025", Time: time.Date(2025, time.August, 20, 0, 0, 0, 0, time.Local), Start: 0, End: 10},
	// 		},
	// 	}, []ParsingError{})
	// testParseExpectation(t,
	// 	``, []*VocabularySection{
	// 		{
	// 			Date: nil,
	// 		},
	// 	}, []ParsingError{})
	// testParseExpectation(t,
	// 	`20/`, []*VocabularySection{
	// 		{Date: nil},
	// 	}, []ParsingError{
	// 		MalformedDate,
	// 	})
	// testParseExpectation(t,
	// 	` 00/00/0000`, []*VocabularySection{
	// 		{Date: nil},
	// 	}, []ParsingError{
	// 		MalformedDate,
	// 	})
}

// TODO: incomplete date, incomplete language, incomlpete word, etc.
func TestSyntacticError(t *testing.T) {
}
