package vocabulary

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
	test "vocab/vocab_testing"
)

func printJSON(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func trimLines(text string) string {
	split := strings.Split(text, "\n")
	filtered := []string{}
	for i := range split {
		trimmed := strings.TrimSpace(split[i])
		if trimmed == "" {
			continue
		}
		filtered = append(filtered, trimmed)
	}
	joined := strings.Join(filtered, "\n")
	return joined
}

func testParseExpectation(t *testing.T, text string, sections []*VocabularySection, expectedErrors []string) {
	t.Helper()

	joined := trimLines(text)

	var parsedErrors []any
	parser := NewParser(t.Context(), "xxx", NewScanner(joined),
		func(a any) {
			parsedErrors = append(parsedErrors, a)
		}, func(a any) {})

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
							{Text: "la magia", Start: 7, End: 14},
							{Text: "bene", Start: 17, End: 20},
						},
					},
					{
						Language: Deutsch,
						Line:     2,
						Words: []*Word{
							{Text: "anlegen", Start: 7, End: 13},
						},
					},
				},
				ReviewedWords: []*WordsSection{},
				Utterance: []*UtteranceSection{
					{
						Line:  3,
						Start: 0,
						End:   len("Ho una magia molto speciale. Non ti conviene metterti contro di me!"),
						Text:  "Ho una magia molto speciale. Non ti conviene metterti contro di me!",
					},
					{
						Line:  4,
						Start: 0,
						End:   len("Ne, will gar nicht mit ihm anlegen."),
						Text:  "Ne, will gar nicht mit ihm anlegen.",
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
							{Text: "brillare", Start: 7, End: 14},
						},
					},
				},
				ReviewedWords: []*WordsSection{
					{
						Language: Italiano,
						Line:     5,
						Words: []*Word{
							{Text: "maga", Start: 7, End: 14},
						},
					},
				},
				Utterance: []*UtteranceSection{
					{
						Line:  6,
						Start: 0,
						End:   len("C'era una volta un piccolo villaggio in Italia. In questo villaggio, viveva una giovane maga. La maga si chiamava Luna, e il suo potere era molto semplice: poteva far brillare le stelle nel cielo."),
						Text:  "C'era una volta un piccolo villaggio in Italia. In questo villaggio, viveva una giovane maga. La maga si chiamava Luna, e il suo potere era molto semplice: poteva far brillare le stelle nel cielo. ",
					},
				},
			},
		}, []string{})
}

// Incomplete sections don't necessarily emit diagnostics error as missing vocabulary is already covered by the compiler.
func TestOnlyDateSection(t *testing.T) {
	type Expectation struct {
		Input      string
		ParsedDate time.Time
		Start      int
		End        int
		Error      string
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
		var parsedError string = ""
		parser := NewParser(t.Context(), "xxx", NewScanner(expectation.Input), func(a any) {
			parsedError = a.(string)
		}, func(a any) {})
		parser.Parse()

		sectionDate := parser.ast.Sections[0].Date

		if expectation.Error != "" && expectation.Error != parsedError {
			t.Fatalf("Error mismatch, expected: %s, got %s", expectation.Error, parsedError)
		}

		gotTime := sectionDate.Time
		if expectation.ParsedDate != gotTime {
			t.Fatalf("Date mismatched, expected: %+v, got %+v", expectation.ParsedDate, gotTime)
		}

		gotStart := sectionDate.Start
		gotEnd := sectionDate.End
		if gotStart != expectation.Start || gotEnd != expectation.End {
			t.Fatalf("Start and end don't match, expected %d, %d -- got %d, %d", expectation.Start, expectation.End, gotStart, gotEnd)
		}
	}
}

func TestIncompleteDateSection(t *testing.T) {
	var parsedError string = ""
	parser := NewParser(t.Context(), "xxx", NewScanner("23/00"), func(a any) {
		parsedError = a.(string)
	}, func(a any) {})
	parser.Parse()

	if parsedError != ExpectVocabSection {
		t.Errorf("Expected parsed error to be MalformedDate, instead got %s", parsedError)
	}

	if parser.tokenStart != 0 && parser.tokenEnd != 5 {
		t.Errorf("Token start and end not 0 and 5, :%d, %d", parser.tokenStart, parser.tokenEnd)
	}
}

func TestInvalidDateSectionUnexpectedToken(t *testing.T) {
	var parsedError string = ""
	parser := NewParser(t.Context(), "xxx", NewScanner("08/09/2025 foo"), func(a any) {
		parsedError = a.(string)
	}, func(a any) {})
	parser.Parse()

	if parsedError != ExpectVocabSection {
		t.Errorf("Expected parsed error to be UnexpectedToken, instead got %s", parsedError)
	}

	if parser.line != 0 {
		t.Errorf("Expect line to be 0, instead got: %d", parser.line)
	}

	if parser.tokenStart != 11 && parser.tokenEnd != 14 {
		t.Errorf("Token start and end not 0 and 5, :%d, %d", parser.tokenStart, parser.tokenEnd)
	}
}

func TestSingleWordSection(t *testing.T) {
	text := trimLines(`
		20/08/2025
		> (it) la magia, bene,scorprire
	`)

	var parsedError string = ""
	parser := NewParser(t.Context(), "xxx", NewScanner(text), func(a any) {
		parsedError = a.(string)
	}, func(a any) {})
	parser.Parse()

	test.Expect(t, "", parsedError)
	test.Expect(t, 1, parser.line)
	test.Expect(t, 1, len(parser.ast.Sections))

	section := parser.ast.Sections[0]
	test.Expect(t, time.Date(2025, time.August, 20, 0, 0, 0, 0, time.Local), section.Date.Time)
	test.Expect(t, 0, section.Date.Line)
	test.Expect(t, 0, section.Date.Start)
	test.Expect(t, 10, section.Date.End)
	test.Expect(t, 1, len(section.NewWords))

	newWords := section.NewWords
	test.Expect(t, 1, len(newWords))
	test.Expect(t, Italiano, newWords[0].Language)
	test.Expect(t, 1, newWords[0].Line)

	words := newWords[0].Words
	test.Expect(t, "la magia", words[0].Text)
	test.Expect(t, "bene", words[1].Text)
	test.Expect(t, "scorprire", words[2].Text)
}

// TODO
// word section missing language identifier
// word section missing date
// word section with malformed date
func TestWordSections(t *testing.T) {
	// text := `
	// 	20/08/2025
	// 	> (it) la magia, bene
	// `
	// var parsedError string = ""
	// parser := NewParser(t.Context(), "xxx", NewScanner(text), func(a any) {
	// 	parsedError = a.(string)
	// }, func(a any) {})
	// parser.Parse()

	// if parsedError != "" {
	// 	t.Errorf("Expected no error, got: %s", parsedError)
	// }

	// if parser.line != 1 {
	// 	t.Errorf("Expect line to be 1, instead got: %d", parser.line)
	// }

	// // correct word section
	// testParseExpectation(t,
	// 	`
	// 		20/08/2025
	// 		> (it) la magia, bene
	// 		> (de) die Gelegenheit
	// 		>> (de) anlegen
	// 	`, []*VocabularySection{
	// 		{
	// 			Date: &DateSection{Text: "20/08/2025", Time: time.Date(2025, time.August, 20, 0, 0, 0, 0, time.Local), Start: 0, End: 10},
	// 			NewWords: []*WordsSection{
	// 				{
	// 					Language: Italiano,
	// 					Line:     1,
	// 					Words: []*Word{
	// 						{Text: "magia", Text: "la magia", Start: 7, End: 14},
	// 						{Text: "bene", Text: "bene", Start: 17, End: 20},
	// 					},
	// 				},
	// 				{
	// 					Language: Deutsch,
	// 					Line:     2,
	// 					Words: []*Word{
	// 						{Text: "anlegen", Text: "anlegen", Start: 7, End: 13},
	// 					},
	// 				},
	// 			},
	// 			ReviewedWords: []*WordsSection{},
	// 			Utterance: []*UtteranceSection{
	// 				{
	// 					Start: 0,
	// 					End:   len("Ho una magia molto speciale. Non ti conviene metterti contro di me!"),
	// 					Text:  "Ho una magia molto speciale. Non ti conviene metterti contro di me!",
	// 				},
	// 				{
	// 					Start: 0,
	// 					End:   len("Ne, will gar nicht mit ihm anlegen."),
	// 					Text:  "Ne, will gar nicht mit ihm anlegen.",
	// 				},
	// 			},
	// 		},
	// 		{},
	// 	}, []ParsingError{})
}

// Test one paragraph
// Test multiple paragraph
// Test sentence in Italian
// Test sentence in German
// Test sentence without Date
// func TestParagraphsSection() {

// }
