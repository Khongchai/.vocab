package entity

import (
	"maps"
	"slices"
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

	test.Expect(t, true, tree.branches[string(Deutsch)] != nil)
	branches := slices.Collect(maps.Values(tree.branches))
	test.Expect(t, 1, len(branches))
	test.Expect(t, 1, len(branches[0].twigs))
	test.Expect(t, true, branches[0].twigs[word.Text] != nil)
	test.Expect(t, 1, len(branches[0].twigs[word.Text]))
	twig := branches[0].twigs[word.Text][0]
	test.Expect(t, 5, twig.grade)
	test.Expect(t, vocabSection, twig.section)
}

func TestAddTwigToNonEmptyTree(t *testing.T) {
}

func TestAddTwigToNonEmptyLanguageBranch(t *testing.T) {
}

func TestGraftingTrees(t *testing.T) {

}
