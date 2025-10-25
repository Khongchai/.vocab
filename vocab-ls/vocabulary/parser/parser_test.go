package parser

import (
	"fmt"
	"testing"
	"time"
	lsproto "vocab/lsp"
	test "vocab/vocab_testing"
)

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

		section := parser.Ast.Sections[0]
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
	text := test.TrimLines(`
		20/08/2025
		> (it) la magia, bene,scorprire
		>> (de) was
	`)

	parser := NewParser(t.Context(), "xxx", NewScanner(text), func(a any) {})
	parser.Parse()

	test.Expect(t, 1, len(parser.Ast.Sections))
	test.Expect(t, 2, parser.line)
	test.Expect(t, 0, len(parser.currentVocabSection().Diagnostics))

	section := parser.Ast.Sections[0]
	test.Expect(t, time.Date(2025, time.August, 20, 0, 0, 0, 0, time.Local), section.Date.Time)
	test.Expect(t, 0, section.Date.Line)
	test.Expect(t, 0, section.Date.Start)
	test.Expect(t, 10, section.Date.End)
	test.Expect(t, 1, len(section.NewWords))

	newWords := section.NewWords
	test.Expect(t, 1, len(newWords))
	test.Expect(t, Italiano, newWords[0].Language)
	test.Expect(t, false, newWords[0].Reviewed)
	test.Expect(t, 1, newWords[0].Line)

	n := newWords[0].Words
	test.Expect(t, "la magia", n[0].Text)
	test.Expect(t, 1, n[0].Line)
	test.Expect(t, 7, n[0].Start)
	test.Expect(t, 7+len("la magia"), n[0].End)
	test.Expect(t, false, n[0].Literally)
	test.Expect(t, "bene", n[1].Text)
	test.Expect(t, false, n[1].Literally)
	test.Expect(t, "scorprire", n[2].Text)
	test.Expect(t, false, n[2].Literally)

	reviewedWords := section.ReviewedWords
	r := reviewedWords[0].Words
	test.Expect(t, 1, len(reviewedWords))
	test.Expect(t, Deutsch, reviewedWords[0].Language)
	test.Expect(t, 2, reviewedWords[0].Line)
	test.Expect(t, true, reviewedWords[0].Reviewed)
	test.Expect(t, "was", r[0].Text)
	test.Expect(t, false, r[0].Literally)
}

func TestWordSectionWithoutDate(t *testing.T) {
	text := test.TrimLines(`
		> (it) la magia, bene,scorprire
	`)

	parser := NewParser(t.Context(), "xxx", NewScanner(text), func(any) {})
	parser.Parse()

	test.Expect(t, 1, len(parser.Ast.Sections)) // upon error, create an empty vocab section even if there is none
	diag := parser.currentVocabSection().Diagnostics[0]
	test.Expect(t, ExpectDateSection, diag.Message)
	test.Expect(t, 0, diag.Range.Start.Character)
	test.Expect(t, 1, diag.Range.End.Character)
	test.Expect(t, 0, diag.Range.Start.Line)
	test.Expect(t, 0, diag.Range.End.Line)
	test.Expect(t, lsproto.DiagnosticsSeverityError, diag.Severity)

}

func TestWordExpression(t *testing.T) {
	text := test.TrimLines(fmt.Sprintf(`
		20/08/2025
		> (it) %sla magia%s, bene
	`, "`", "`"))

	parser := NewParser(t.Context(), "xxx", NewScanner(text), func(any) {})
	parser.Parse()

	section := parser.currentVocabSection()
	words := section.NewWords[0]
	test.Expect(t, Italiano, words.Language)
	test.Expect(t, 1, words.Line)
	test.Expect(t, "la magia", words.Words[0].Text)
	test.Expect(t, true, words.Words[0].Literally)
}

func TestWordExpressionMissingClosingBacktickShouldAutoClose(t *testing.T) {
	text := test.TrimLines(fmt.Sprintf(`
		20/08/2025
		> (it) %sla magia, bene
		21/08/2025
		> (de) %sder Inhalt
	`, "`", "`"))

	parser := NewParser(t.Context(), "xxx", NewScanner(text), func(any) {})
	parser.Parse()

	// (de)
	section2 := parser.currentVocabSection()
	words2 := section2.NewWords[0]
	test.Expect(t, Deutsch, words2.Language)
	test.Expect(t, 3, words2.Line)
	test.Expect(t, 1, len(words2.Words))
	test.Expect(t, "der Inhalt", words2.Words[0].Text)
	test.Expect(t, true, words2.Words[0].Literally)

	// (it)
	section1 := parser.Ast.Sections[0]
	words1 := section1.NewWords[0]
	test.Expect(t, Italiano, words1.Language)
	test.Expect(t, 1, words1.Line)
	test.Expect(t, 1, len(words1.Words))
	test.Expect(t, "la magia, bene", words1.Words[0].Text)
	test.Expect(t, true, words1.Words[0].Literally)

}

func TestMultipleWordSection(t *testing.T) {
	text := test.TrimLines(`
		> (it) la magia, bene,scorprire
		>> (de) was
	`)

	parser := NewParser(t.Context(), "xxx", NewScanner(text), func(any) {})
	parser.Parse()

	test.Expect(t, 1, len(parser.Ast.Sections))
	diag := parser.currentVocabSection().Diagnostics[0]
	test.Expect(t, ExpectDateSection, diag.Message)
	test.Expect(t, 0, diag.Range.Start.Character)
	test.Expect(t, 1, diag.Range.End.Character)
	test.Expect(t, 0, diag.Range.Start.Line)
	test.Expect(t, 0, diag.Range.End.Line)
	test.Expect(t, lsproto.DiagnosticsSeverityError, diag.Severity)

}

func TestUtteranceSection(t *testing.T) {
	text := test.TrimLines(`
		01/08/1997
		> (de) ablenken, ansprechen
		Das lenkt mich wirklich ab!
		Sag einfach Bescheid, was dir gerade am meisten anspricht!
	`)

	parser := NewParser(t.Context(), "xxx", NewScanner(text), func(any) {})
	parser.Parse()

	utterance := parser.Ast.Sections[0].Utterance
	test.Expect(t, 2, len(utterance))
	test.Expect(t, "Das lenkt mich wirklich ab!", utterance[0].Text)
	test.Expect(t, 2, utterance[0].Line)
	test.Expect(t, 0, utterance[0].Start)
	test.Expect(t, len("Das lenkt mich wirklich ab!"), utterance[0].End)

	test.Expect(t, "Sag einfach Bescheid, was dir gerade am meisten anspricht!", utterance[1].Text)
	test.Expect(t, 3, utterance[1].Line)
	test.Expect(t, 0, utterance[1].Start)
	test.Expect(t, len("Sag einfach Bescheid, was dir gerade am meisten anspricht!"), utterance[1].End)
}

func TestFullSectionParsing(t *testing.T) {
	text := test.TrimLines(`
		02/10/2025
		>> (it) la notizia, chiacchierare
		> (de) aufschlüsseln
		Guardando le notizie italiane. Che tipo di accento è questo?
		Oggi, passegiamo e chiacchieriamo in italiano.
		Stamattina ho un po'sonno.
		I didn't sleep much either
		Anch'io ho dormito poco.
		Anch'io sono un' po stanco, sono andato a letto tardi.
		A mezzanotte, ma che per me è tardissimo.
		Camminiamo un altro pocchino?
		Kannst du mir diesen Satz aufschlüsseln?
		03/10/2025
		> (de) ansprechen, schnappen, ausfragen
		>> (de) anlegen,
		Sag einfach was dir so im Kopf rumgehen, und wir plaudern ein bisschen.
		Ich werde nicht mit ihm anlegen, das ist mein Vorschlag für dich.
		Na klar, ich helfe dir gerne dabei. Wir können zum Beispiel über deinen Tag reden, über irgendwelche Hobbys oder einfach über etwas, das dich interessiert, wie Reisen, Bücher oder Technik. Sag einfach Bescheid, was dich gerade am meisten anspricht.
		Na gut, dann schnapp ich mir einfach mal ein paar technische Themen für dich. Wie wäre es, wenn ich dich ein bisschen dazu ausfrage, welche Technik dich im Moment so fasziniert?
	`)

	parser := NewParser(t.Context(), "xxx", NewScanner(text), func(any) {})
	parser.Parse()

	test.Expect(t, 2, len(parser.Ast.Sections))

	// ======== SECTION 1: 02/10/2025 ========
	section1 := parser.Ast.Sections[0]
	test.Expect(t, time.Date(2025, time.October, 2, 0, 0, 0, 0, time.Local), section1.Date.Time)

	// Reviewed words (>>)
	test.Expect(t, 1, len(section1.ReviewedWords))
	reviewed := section1.ReviewedWords[0]
	test.Expect(t, Italiano, reviewed.Language)
	test.Expect(t, 1, reviewed.Line)
	test.Expect(t, 2, len(reviewed.Words))
	test.Expect(t, "la notizia", reviewed.Words[0].Text)
	test.Expect(t, false, reviewed.Words[0].Literally)
	test.Expect(t, "chiacchierare", reviewed.Words[1].Text)

	// New words (>)
	test.Expect(t, 1, len(section1.NewWords))
	newWords := section1.NewWords[0]
	test.Expect(t, Deutsch, newWords.Language)
	test.Expect(t, 2, newWords.Line)
	test.Expect(t, 1, len(newWords.Words))
	test.Expect(t, "aufschlüsseln", newWords.Words[0].Text)

	// Utterances
	test.Expect(t, 9, len(section1.Utterance))
	test.Expect(t, "Guardando le notizie italiane. Che tipo di accento è questo?", section1.Utterance[0].Text)
	test.Expect(t, "Camminiamo un altro pocchino?", section1.Utterance[7].Text)
	test.Expect(t, "Kannst du mir diesen Satz aufschlüsseln?", section1.Utterance[8].Text)

	// ======== SECTION 2: 03/10/2025 ========
	section2 := parser.Ast.Sections[1]
	test.Expect(t, time.Date(2025, time.October, 3, 0, 0, 0, 0, time.Local), section2.Date.Time)
	test.Expect(t, 12, section2.Date.Line)

	// New words (>)
	test.Expect(t, 1, len(section2.NewWords))
	words := section2.NewWords[0]
	test.Expect(t, Deutsch, words.Language)
	test.Expect(t, 13, words.Line)
	test.Expect(t, 3, len(words.Words))
	test.Expect(t, "ansprechen", words.Words[0].Text)
	test.Expect(t, "schnappen", words.Words[1].Text)
	test.Expect(t, "ausfragen", words.Words[2].Text)

	// Reviewed words (>>)
	test.Expect(t, 1, len(section2.ReviewedWords))
	rw := section2.ReviewedWords[0]
	test.Expect(t, Deutsch, rw.Language)
	test.Expect(t, 14, rw.Line)
	test.Expect(t, 1, len(rw.Words))
	test.Expect(t, "anlegen", rw.Words[0].Text)

	// Utterances (German + mixed)
	test.Expect(t, 4, len(section2.Utterance))
	test.Expect(t, "Ich werde nicht mit ihm anlegen, das ist mein Vorschlag für dich.", section2.Utterance[1].Text)
	test.Expect(t, "Na gut, dann schnapp ich mir einfach mal ein paar technische Themen für dich. Wie wäre es, wenn ich dich ein bisschen dazu ausfrage, welche Technik dich im Moment so fasziniert?", section2.Utterance[3].Text)
}

func TestGrading(t *testing.T) {
	text := test.TrimLines(fmt.Sprintf(`
		20/08/2025
		> (it) %sla magia%s(1), chiacchierare, caminare(0), cosa(10)
	`, "`", "`"))

	parser := NewParser(t.Context(), "xxx", NewScanner(text), func(any) {})
	parser.Parse()

	words := parser.Ast.Sections[0].NewWords[0].Words
	test.Expect(t, 4, len(words))
	test.Expect(t, 1, words[0].Grade)
	test.Expect(t, 0, words[1].Grade) // default, no score = 0
	test.Expect(t, 0, words[2].Grade)
	test.Expect(t, 10, words[3].Grade)
}

func TestInvalidGradeShouldIgnoreWordsAfterCompletely(t *testing.T) {
	text := test.TrimLines(fmt.Sprintf(`
		20/08/2025
		> (it) %sla magia%s(xxx), these, should, not, count
		21/08/2025
		> (it) chiacchierare(4j2)
		22/08/2025
		> (it) chiacchierare()
	`, "`", "`"))

	parser := NewParser(t.Context(), "xxx", NewScanner(text), func(any) {})
	parser.Parse()

	words1 := parser.Ast.Sections[0].NewWords[0].Words
	test.Expect(t, 1, len(words1))
	diag1 := parser.Ast.Sections[0].Diagnostics[0]
	test.Expect(t, InvalidScore, diag1.Message)
	test.Expect(t, 17, diag1.Range.Start.Character)
	test.Expect(t, 22, diag1.Range.End.Character) // remember, it's [start, ...end)
	test.Expect(t, 1, diag1.Range.Start.Line, diag1.Range.End.Line)

	words2 := parser.Ast.Sections[1].NewWords[0].Words
	test.Expect(t, 0, len(words2))
	diag2 := parser.Ast.Sections[1].Diagnostics[0]
	test.Expect(t, InvalidScore, diag2.Message)
	test.Expect(t, 20, diag2.Range.Start.Character)
	test.Expect(t, 25, diag2.Range.End.Character)
	test.Expect(t, 3, diag2.Range.Start.Line, diag2.Range.End.Line)

	words3 := parser.Ast.Sections[2].NewWords[0].Words
	test.Expect(t, 0, len(words3))
	diag3 := parser.Ast.Sections[2].Diagnostics[0]
	test.Expect(t, InvalidScore, diag3.Message)
	test.Expect(t, 20, diag3.Range.Start.Character)
	test.Expect(t, 22, diag3.Range.End.Character)
	test.Expect(t, 5, diag3.Range.Start.Line, diag3.Range.End.Line)
}

func TestWordWithNumber_ShouldWorkBecauseItSuperUsefulWhenTesting(t *testing.T) {
	text := test.TrimLines(`
		20/05/2025
		> (it) it_word1, it_word2
		lorem ipsum...
	`)
	ast := NewParser(t.Context(), "xxx", NewScanner(text), func(any) {}).Parse().Ast
	test.Expect(t, 1, len(ast.Sections))
	test.Expect(t, 1, len(ast.Sections[0].NewWords))
	test.Expect(t, 2, len(ast.Sections[0].NewWords[0].Words))
	test.Expect(t, "it_word1", ast.Sections[0].NewWords[0].Words[0].Text)
	test.Expect(t, "it_word2", ast.Sections[0].NewWords[0].Words[1].Text)
}

func TestIgnoreComment(t *testing.T) {
	text := test.TrimLines(`
		20/05/2025
		| hola
		> (it) it_word1, it_word2
		| amigo!
		lorem ipsum...
	`)
	ast := NewParser(t.Context(), "xxx", NewScanner(text), func(any) {}).Parse().Ast
	test.Expect(t, 1, len(ast.Sections))
	test.Expect(t, 1, len(ast.Sections[0].NewWords))
	test.Expect(t, 2, len(ast.Sections[0].NewWords[0].Words))
	test.Expect(t, "it_word1", ast.Sections[0].NewWords[0].Words[0].Text)
	test.Expect(t, "it_word2", ast.Sections[0].NewWords[0].Words[1].Text)
}

func TestRepeatedToken(t *testing.T) {
	text := test.TrimLines(`
		20/05/2025
		> (it) la magia, maga, la magia
	`)
	ast := NewParser(t.Context(), "xxx", NewScanner(text), func(any) {}).Parse().Ast
	test.Expect(t, 1, len(ast.Sections))
	test.Expect(t, 1, len(ast.Sections[0].NewWords))
	test.Expect(t, 2, len(ast.Sections[0].NewWords[0].Words))
	test.Expect(t, 1, len(ast.Sections[0].Diagnostics))
	test.Expect(t, DuplicateToken, ast.Sections[0].Diagnostics[0].Message)
	test.Expect(t, 23, ast.Sections[0].Diagnostics[0].Range.Start.Character)
	test.Expect(t, 31, ast.Sections[0].Diagnostics[0].Range.End.Character)
}
