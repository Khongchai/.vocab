package vocabulary

import (
	"strings"
	"testing"
	"time"
	lsproto "vocab/lsp"
	test "vocab/vocab_testing"
)

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
		parser := NewParser(t.Context(), "xxx", NewScanner(expectation.Input), func(a any) {})
		parser.Parse()

		section := parser.ast.Sections[0]
		sectionDate := section.Date

		if expectation.Error != "" && section.Diagnostics[0].Message != expectation.Error {
			t.Fatalf("Error mismatch, expected: %s, got %+v", expectation.Error, section.Diagnostics[0].Message)
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
	parser := NewParser(t.Context(), "xxx", NewScanner("23/00"), func(a any) {})
	parser.Parse()

	test.Expect(t, 1, len(parser.currentVocabSection().Diagnostics))
	diag := parser.currentVocabSection().Diagnostics[0]
	test.Expect(t, ExpectDateSection, diag.Message)
	test.Expect(t, 0, diag.Range.Start.Character)
	test.Expect(t, 5, diag.Range.End.Character)
}

func TestInvalidDateSectionUnexpectedToken(t *testing.T) {
	parser := NewParser(t.Context(), "xxx", NewScanner("08/09/2025 foo"), func(a any) {})
	parser.Parse()

	test.Expect(t, 1, len(parser.currentVocabSection().Diagnostics))
	diag := parser.currentVocabSection().Diagnostics[0]
	test.Expect(t, ExpectVocabSection, diag.Message)
	test.Expect(t, 11, diag.Range.Start.Character)
	test.Expect(t, 14, diag.Range.End.Character)
}

func TestSingleWordSection(t *testing.T) {
	text := trimLines(`
		20/08/2025
		> (it) la magia, bene,scorprire
	`)

	parser := NewParser(t.Context(), "xxx", NewScanner(text), func(a any) {})
	parser.Parse()

	test.Expect(t, 1, parser.line)
	test.Expect(t, 1, len(parser.ast.Sections))
	test.Expect(t, 0, len(parser.currentVocabSection().Diagnostics))

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

func TestWordSectionWithoutDate(t *testing.T) {
	text := trimLines(`
		> (it) la magia, bene,scorprire
	`)

	parser := NewParser(t.Context(), "xxx", NewScanner(text), func(any) {})
	parser.Parse()

	test.Expect(t, 1, len(parser.ast.Sections)) // upon error, create an empty vocab section even if there is none
	diag := parser.currentVocabSection().Diagnostics[0]
	test.Expect(t, ExpectDateSection, diag.Message)
	test.Expect(t, 0, diag.Range.Start.Character)
	test.Expect(t, 1, diag.Range.End.Character)
	test.Expect(t, 0, diag.Range.Start.Line)
	test.Expect(t, 0, diag.Range.End.Line)
	test.Expect(t, lsproto.DiagnosticsSeverityError, diag.Severity)

}

func TestMultipleWordSection(t *testing.T) {
}

func TestUtteranceSection(t *testing.T) {}

func TestUtteranceSectionWithoutDate(t *testing.T) {}

func TestUtteranceSectionWithoutVocab(t *testing.T) {}
func TestUtteranceSectionAsStart(t *testing.T)      {}

// func TestFullSectionParsing(t *testing.T) {
// testParseExpectation(t,
// 	`
// 		20/08/2025
// 		> (it) la magia, bene
// 		> (de) anlegen
// 		Ho una magia molto speciale. Non ti conviene metterti contro di me!
// 		Ne, will gar nicht mit ihm anlegen.
// 		21/08/2025
// 		> (it) brillare
// 		>> (it) la maga
// 		C'era una volta un piccolo villaggio in Italia. In questo villaggio, viveva una giovane maga. La maga si chiamava Luna, e il suo potere era molto semplice: poteva far brillare le stelle nel cielo.
// 	`, []*VocabularySection{
// 		{
// 			Date: &DateSection{Text: "20/08/2025", Time: time.Date(2025, time.August, 20, 0, 0, 0, 0, time.Local), Start: 0, End: 10},
// 			NewWords: []*WordsSection{
// 				{
// 					Language: Italiano,
// 					Line:     1,
// 					Words: []*Word{
// 						{Text: "la magia", Start: 7, End: 14},
// 						{Text: "bene", Start: 17, End: 20},
// 					},
// 				},
// 				{
// 					Language: Deutsch,
// 					Line:     2,
// 					Words: []*Word{
// 						{Text: "anlegen", Start: 7, End: 13},
// 					},
// 				},
// 			},
// 			ReviewedWords: []*WordsSection{},
// 			Utterance: []*UtteranceSection{
// 				{
// 					Line:  3,
// 					Start: 0,
// 					End:   len("Ho una magia molto speciale. Non ti conviene metterti contro di me!"),
// 					Text:  "Ho una magia molto speciale. Non ti conviene metterti contro di me!",
// 				},
// 				{
// 					Line:  4,
// 					Start: 0,
// 					End:   len("Ne, will gar nicht mit ihm anlegen."),
// 					Text:  "Ne, will gar nicht mit ihm anlegen.",
// 				},
// 			},
// 		},
// 		{
// 			Date: &DateSection{Text: "21/08/2025", Time: time.Date(2025, time.August, 20, 0, 0, 0, 0, time.Local), Start: 0, End: 10},
// 			NewWords: []*WordsSection{
// 				{
// 					Language: Italiano,
// 					Line:     5,
// 					Words: []*Word{
// 						{Text: "brillare", Start: 7, End: 14},
// 					},
// 				},
// 			},
// 			ReviewedWords: []*WordsSection{
// 				{
// 					Language: Italiano,
// 					Line:     5,
// 					Words: []*Word{
// 						{Text: "maga", Start: 7, End: 14},
// 					},
// 				},
// 			},
// 			Utterance: []*UtteranceSection{
// 				{
// 					Line:  6,
// 					Start: 0,
// 					End:   len("C'era una volta un piccolo villaggio in Italia. In questo villaggio, viveva una giovane maga. La maga si chiamava Luna, e il suo potere era molto semplice: poteva far brillare le stelle nel cielo."),
// 					Text:  "C'era una volta un piccolo villaggio in Italia. In questo villaggio, viveva una giovane maga. La maga si chiamava Luna, e il suo potere era molto semplice: poteva far brillare le stelle nel cielo. ",
// 				},
// 			},
// 		},
// 	}, []string{})
// }
