package compiler

import (
	"maps"
	"slices"
	"time"
	lsproto "vocab/lsp"
	"vocab/super_memo"
	"vocab/vocabulary/parser"
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

func (wt *WordTree) GetTwigs(language parser.Language, word string) []*WordTwig {
	existingBranch := wt.branches[string(language)]
	if existingBranch == nil {
		return []*WordTwig{}
	}
	twigs := existingBranch.twigs[word]
	return twigs
}

func (wt *WordTree) GetBranches() []*LanguageBranch {
	return slices.Collect(maps.Values(wt.branches))
}

// Add a new word to the tree. If language branch does not exists, one is created.
func (wt *WordTree) AddTwig(language parser.Language, word *parser.Word, uri string, section *parser.VocabularySection, startingDiagnostics []*lsproto.Diagnostic) {
	lang := string(language)
	branch := wt.branches[lang]

	if branch == nil {
		branch = &LanguageBranch{twigs: map[string][]*WordTwig{}}
		wt.branches[lang] = branch
	}

	clamped_grade := max(super_memo.MemoBlackout, min(word.Grade, super_memo.MemoPerfect))
	if clamped_grade != word.Grade {
		newDiag := lsproto.MakeDiagnostics(
			"Expect grade to be from 0 to 5. Can also leave empty for the default 0",
			word.Line,
			word.Start,
			word.End,
			lsproto.DiagnosticsSeverityError,
		)
		startingDiagnostics = append(startingDiagnostics, newDiag)
	}
	twig := &WordTwig{
		grade:               clamped_grade,
		section:             section,
		startingDiagnostics: startingDiagnostics,
	}

	norm := word.GetNormalizedText()
	branch.twigs[norm] = append(branch.twigs[norm], twig)
	branch.sortTwigs(norm)
}

func (wt *WordTree) Graft(other *WordTree) *WordTree {
	for key, value := range other.branches {
		if wt.branches[key] == nil {
			wt.branches[key] = value
			continue
		}

		wt.branches[key].Graft(value)
	}
	return wt
}

func (wt *WordTree) Harvest() []*WordFruit {
	details := []*WordFruit{}
	// für jede LanguageBranch
	// für jede WordTwig auf LanguageBranch (Wir gehen davon aus, dass die Twigs schon sortiert sind.)
	for lang, langBranch := range wt.branches {
		for word, twigs := range langBranch.twigs {
			detail := &WordFruit{
				Words:               []*parser.Word{},
				Interval:            0,
				LastSeenDate:        time.Time{},
				StartingDiagnostics: []*lsproto.Diagnostic{},
				Lang:                parser.Language(lang),
				Text:                word,
			}

			repetitionNumber := 0
			easinessFactor := super_memo.InitialEasinessFactor

			// interval is the final output we want
			var interval float64
			var lastSeenDate *time.Time
			for _, twig := range twigs {
				detail.Words = append(detail.Words, twig.word)
				detail.StartingDiagnostics = append(detail.StartingDiagnostics, twig.startingDiagnostics...)
				currentInterval := func() float64 {
					if lastSeenDate == nil {
						return 0
					}
					diff := twig.section.Date.Time.Sub(*lastSeenDate)
					diffDays := diff.Hours() / 24
					return diffDays
				}()
				repetitionNumber, interval, easinessFactor = super_memo.Sm2(twig.grade, repetitionNumber, currentInterval, easinessFactor)

				lastSeenDate = &twig.section.Date.Time
			}

			detail.Interval = interval
			detail.LastSeenDate = *lastSeenDate

			details = append(details, detail)
		}
	}

	return details
}

type LanguageBranch struct {
	// Map of word to places they appear in.
	twigs map[string][]*WordTwig
}

func (lb *LanguageBranch) Graft(other *LanguageBranch) *LanguageBranch {
	for word, twigs := range other.twigs {
		lb.twigs[word] = append(lb.twigs[word], twigs...)
	}

	// If grafting is called more than once
	// this makes sure no section is repeated twice...
	for word := range lb.twigs {
		uniques := make(map[string]*WordTwig)
		for _, section := range lb.twigs[word] {
			ident := section.section.Identity()
			uniques[ident] = section
		}

		uniqued := slices.Collect(maps.Values(uniques))
		lb.twigs[word] = uniqued
		lb.sortTwigs(word)
	}

	return lb
}

func (lb *LanguageBranch) sortTwigs(word string) {
	unsorted := slices.Values(lb.twigs[word])
	sorted := slices.SortedFunc(unsorted, func(a, b *WordTwig) int {
		if a.section.Date.Time.Before(b.section.Date.Time) {
			return -1
		}

		if a.section.Date.Time.After(b.section.Date.Time) {
			return 1
		}

		return 0
	})
	lb.twigs[word] = sorted
}

type WordTwig struct {
	grade               int
	word                *parser.Word
	section             *parser.VocabularySection // word.Parent.Parent
	startingDiagnostics []*lsproto.Diagnostic
	// Document location file name)
	location string
}

// The final review detail of a word
type WordFruit struct {
	Lang parser.Language
	// All word object tied to normalized `text`
	Words []*parser.Word
	// The normalized text.
	Text                string
	Interval            float64
	LastSeenDate        time.Time
	StartingDiagnostics []*lsproto.Diagnostic
}

func (wb *WordTwig) GetLocation() string {
	return wb.location
}
