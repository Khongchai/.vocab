package entity

import (
	"maps"
	"slices"
	"strings"
	"testing"
	"time"
	lsproto "vocab/lsp"
	test "vocab/vocab_testing"
)

func fakeWord(text string, grade int, parent *WordsSection) *Word {
	return &Word{
		Line:      0,
		Text:      text,
		Literally: true,
		Start:     0,
		End:       len(text),
		Grade:     grade,
		Parent:    parent,
	}
}

func TestCreatingANewTreeShouldNotThrowError(t *testing.T) {
	NewWordTree()
}

func TestAddTwigToEmptyTree(t *testing.T) {
	tree := NewWordTree()
	vocabSection := NewVocabularySection("xxx")
	vocabSection.Date = &DateSection{Time: time.Now()}
	wordSection := &WordsSection{
		Parent: vocabSection,
	}
	vocabSection.NewWords = append(vocabSection.NewWords, wordSection)
	word := fakeWord("Testen", 5, wordSection)
	startingDiags := []*lsproto.Diagnostic{
		lsproto.MakeDiagnostics("test diagnostics", 1, 2, 3, lsproto.DiagnosticsSeverityError),
	}

	//act
	tree.AddTwig(Deutsch, word, "xxx", vocabSection, startingDiags)

	wordNormalized := strings.ToLower(word.Text)
	test.Expect(t, true, tree.branches[string(Deutsch)] != nil)
	branches := slices.Collect(maps.Values(tree.branches))
	test.Expect(t, 1, len(branches))
	test.Expect(t, 1, len(branches[0].twigs))
	test.Expect(t, true, branches[0].twigs[wordNormalized] != nil)
	test.Expect(t, 1, len(branches[0].twigs[wordNormalized]))
	twig := branches[0].twigs[wordNormalized][0]
	test.Expect(t, 1, len(twig.startingDiagnostics))
	test.Expect(t, 5, twig.grade)
	test.Expect(t, vocabSection, twig.section)
}

func TestAddTwigToNonEmptyTree_WithMultipleSections(t *testing.T) {
	tree := NewWordTree()
	section1 := NewVocabularySection("xxx")
	section1.Date = &DateSection{Time: time.Now()}
	newWordSection := &WordsSection{
		Parent:   section1,
		Reviewed: false,
		Language: Deutsch,
	}
	section1.NewWords = append(section1.NewWords, newWordSection)
	section2 := NewVocabularySection("xxx")
	section2.Date = &DateSection{Time: time.Now()}
	reviewedWordSection := &WordsSection{
		Parent:   section1,
		Reviewed: true,
		Language: Italiano,
	}
	section2.ReviewedWords = append(section2.ReviewedWords, reviewedWordSection)

	//act
	tree.AddTwig(Deutsch, fakeWord("Testen", 5, newWordSection), "xxx", section1, []*lsproto.Diagnostic{})
	tree.AddTwig(Italiano, fakeWord("Test", 3, newWordSection), "xxx", section1, []*lsproto.Diagnostic{})
	tree.AddTwig(Deutsch, fakeWord("testen", 4, reviewedWordSection), "xxx", section2, []*lsproto.Diagnostic{})

	test.Expect(t, true, tree.branches[string(Deutsch)] != nil)
	test.Expect(t, true, tree.branches[string(Italiano)] != nil)
	branches := slices.Collect(maps.Values(tree.branches))
	test.Expect(t, 2, len(branches))

	// german
	derAst := tree.branches[string(Deutsch)]
	test.Expect(t, 2, len(derAst.twigs["testen"]))
	test.Expect(t, 5, derAst.twigs["testen"][0].grade)
	test.Expect(t, section1, derAst.twigs["testen"][0].section)
	test.Expect(t, 4, derAst.twigs["testen"][1].grade)
	test.Expect(t, section2, derAst.twigs["testen"][1].section)

	// italian
	unRamo := tree.branches[string(Italiano)]
	test.Expect(t, 1, len(unRamo.twigs["test"]))
	test.Expect(t, 3, unRamo.twigs["test"][0].grade)
	test.Expect(t, section1, unRamo.twigs["test"][0].section)
}

func TestAddedTwigsShouldBeSorted(t *testing.T) {
	tree := NewWordTree()
	sectionFromTime := func(time time.Time) *VocabularySection {
		dateText := time.Format("2006-01-02")
		section := &VocabularySection{
			Date: &DateSection{Time: time, Text: dateText},
			Uri:  "xxx",
		}
		wordsSection := &WordsSection{
			Parent:   section,
			Reviewed: false,
			Language: Deutsch,
		}
		section.NewWords = append(section.NewWords, wordsSection)
		return section
	}
	now := time.Now()
	nowSection := sectionFromTime(now)
	yesterdaySection := sectionFromTime(now.AddDate(0, 0, -1))
	tomorrowSection := sectionFromTime(now.AddDate(0, 0, 1))

	//act
	tree.AddTwig(Deutsch,
		fakeWord("poopy", 5, nowSection.NewWords[0]),
		"xxx",
		nowSection,
		[]*lsproto.Diagnostic{},
	)
	tree.AddTwig(Deutsch,
		fakeWord("poopy", 5, tomorrowSection.NewWords[0]),
		"xxx",
		tomorrowSection,
		[]*lsproto.Diagnostic{},
	)
	tree.AddTwig(Deutsch,
		fakeWord("poopy", 5, yesterdaySection.NewWords[0]),
		"xxx",
		yesterdaySection,
		[]*lsproto.Diagnostic{},
	)

	äste := tree.branches[string(Deutsch)]
	test.Expect(t, 1, len(äste.twigs))
	zweige := äste.twigs["poopy"]
	test.Expect(t, yesterdaySection, zweige[0].section)
	test.Expect(t, nowSection, zweige[1].section)
	test.Expect(t, tomorrowSection, zweige[2].section)
}

func TestAddTwigsWithInvalidGrade_ShouldProduceExtraDiagnosticsError(t *testing.T) {
	tree := NewWordTree()
	section := NewVocabularySection("xxx")
	section.Date = &DateSection{Time: time.Now()}
	newWordSection := &WordsSection{
		Parent:   section,
		Reviewed: false,
		Language: Italiano,
	}
	section.NewWords = append(section.NewWords, newWordSection)

	//act
	tree.AddTwig(Deutsch,
		fakeWord("ding", -1, section.NewWords[0]),
		"xxx",
		section,
		[]*lsproto.Diagnostic{},
	)
	tree.AddTwig(Deutsch,
		fakeWord("cosa", 6, section.NewWords[0]),
		"xxx",
		section,
		[]*lsproto.Diagnostic{
			lsproto.MakeDiagnostics("test diagnostics", 1, 2, 3, lsproto.DiagnosticsSeverityError),
		},
	)

	äste := tree.branches[string(Deutsch)]
	test.Expect(t, 2, len(äste.twigs))
	// expect clamping
	test.Expect(t, 0, äste.twigs["ding"][0].grade)
	test.Expect(t, 1, len(äste.twigs["ding"][0].startingDiagnostics))
	test.Expect(t, 5, äste.twigs["cosa"][0].grade)
	test.Expect(t, 2, len(äste.twigs["cosa"][0].startingDiagnostics))
}

func TestGraftingTrees(t *testing.T) {

}
