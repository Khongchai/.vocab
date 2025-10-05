package entity

import (
	"maps"
	"slices"
	lsproto "vocab/lsp"
)

// A WordTree is a map of the word, or exact literal string to
// a branch of languages, which is a map of the language literal name to an
// array of `twigs` of sections they are in
type WordTree struct {
	branches map[string]*LanguageBranch
}

func NewWordTree() *WordTree {
	tree := &WordTree{branches: map[string]*LanguageBranch{}}
	return tree
}

func (wt *WordTree) GetTwigs(language Language, word string) []*WordTwig {
	found := wt.branches[string(language)].twigs[word]
	return found
}

// Add a new word to the tree. If language branch does not exists, one is created.
func (wt *WordTree) AddTwig(language Language, word string, uri string, section *VocabularySection, startingDiagnostics []*lsproto.Diagnostic) {
	lang := string(language)
	branch := wt.branches[lang]

	if branch == nil {
		branch = &LanguageBranch{twigs: map[string][]*WordTwig{}}
		wt.branches[lang] = branch
	}

	twig := &WordTwig{
		section:             section,
		startingDiagnostics: startingDiagnostics,
	}

	branch.twigs[word] = append(branch.twigs[word], twig)
}

func (wt *WordTree) Graft(other *WordTree) {
	for key, value := range other.branches {
		if wt.branches[key] == nil {
			wt.branches[key] = value
			continue
		}

		wt.branches[key].Graft(value)
	}
}

type LanguageBranch struct {
	twigs map[string][]*WordTwig
}

func (wb *LanguageBranch) Graft(other *LanguageBranch) {
	for lang, twigs := range other.twigs {
		wb.twigs[lang] = append(wb.twigs[lang], twigs...)
	}

	// If grafting is called more than once
	// this makes sure no section is repeated twice...
	for lang := range wb.twigs {
		uniques := make(map[string]*WordTwig)
		for _, section := range wb.twigs[lang] {
			ident := section.section.Identity()
			uniques[ident] = section
		}

		unsorted := maps.Values(uniques)
		sorted := slices.SortedFunc(unsorted, func(a, b *WordTwig) int {
			if a.section.Date.Time.Before(b.section.Date.Time) {
				return -1
			}

			if a.section.Date.Time.After(b.section.Date.Time) {
				return 1
			}

			return 0
		})
		wb.twigs[lang] = sorted
	}
}

type WordTwig struct {
	section             *VocabularySection
	startingDiagnostics []*lsproto.Diagnostic
	// Document location (file name)
	location string
}

func (wb *WordTwig) GetLocation() string {
	return wb.location
}

// Use a fixed-grade array as input to super memo2
func (*WordTree) Harvest() []lsproto.Diagnostic {

}
