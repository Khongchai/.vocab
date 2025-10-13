package entity

import (
	"maps"
	"slices"
	"strings"
	"testing"
	lsproto "vocab/lsp"
	test "vocab/vocab_testing"
)

func TestCreatingANewTreeShouldNotThrowError(t *testing.T) {
	NewWordTree()
}

func TestAddTwigToEmptyTree(t *testing.T) {
	tree := NewWordTree()
	vocabSection := &VocabularySection{}
	wordSection := &WordsSection{
		Parent: vocabSection,
	}
	vocabSection.NewWords = append(vocabSection.NewWords, wordSection)
	word := &Word{
		Line:      0,
		Text:      "Testen",
		Literally: true,
		Start:     0,
		End:       len("Testen"),
		Grade:     5,
		Parent:    wordSection,
	}
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
	section1 := &VocabularySection{}
	newWordSection := &WordsSection{
		Parent:   section1,
		Reviewed: false,
		Language: Deutsch,
	}
	section1.NewWords = append(section1.NewWords, newWordSection)
	section2 := &VocabularySection{}
	reviewedWordSection := &WordsSection{
		Parent:   section1,
		Reviewed: true,
		Language: Italiano,
	}
	section2.ReviewedWords = append(section2.ReviewedWords, reviewedWordSection)

	//act
	tree.AddTwig(Deutsch, &Word{
		Line:      0,
		Text:      "Testen",
		Literally: true,
		Start:     0,
		End:       len("Testen"),
		Grade:     5,
		Parent:    newWordSection,
	}, "xxx", section1, []*lsproto.Diagnostic{})
	tree.AddTwig(Italiano, &Word{
		Line:      0,
		Text:      "Test",
		Literally: true,
		Start:     0,
		End:       len("Test"),
		Grade:     3,
		Parent:    newWordSection,
	}, "xxx", section1, []*lsproto.Diagnostic{})
	tree.AddTwig(Deutsch, &Word{
		Line:      0,
		Text:      "testen",
		Literally: true,
		Start:     0,
		End:       len("testen"),
		Grade:     4,
		Parent:    reviewedWordSection,
	}, "xxx", section2, []*lsproto.Diagnostic{})

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

}

func TestAddTwigsWithInvalidGrade_ShouldProduceExtraDiagnosticsError(t *testing.T) {

}

func TestGraftingTrees(t *testing.T) {

}
