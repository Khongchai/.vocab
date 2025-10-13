package entity

import (
	"maps"
	"slices"
	"time"
	lsproto "vocab/lsp"
	"vocab/super_memo"
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
func (wt *WordTree) AddTwig(language Language, word *Word, uri string, section *VocabularySection, startingDiagnostics []*lsproto.Diagnostic) {
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

func (wt *WordTree) Harvest() []*WordFruit {
	details := []*WordFruit{}
	// für jede LanguageBranch
	// für jede WordTwig auf LanguageBranch (Wir gehen davon aus, dass die Twigs schon sortiert sind.)
	for lang, langBranch := range wt.branches {
		for word, twigs := range langBranch.twigs {
			detail := &WordFruit{
				Words:               []*Word{},
				TimeRemaining:       0,
				StartingDiagnostics: []*lsproto.Diagnostic{},
				Lang:                Language(lang),
				Text:                word,
			}

			repetitionNumber := 0
			easinessFactor := super_memo.InitialEasinessFactor

			// interval is the final output we want
			var interval float64
			var lastSeenDate *time.Time
			for _, twig := range twigs {
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

			// In the future, this can be moved to a different layer.
			// If max(0, today - last reviewed date) is more than interval, produce diagnostics
			diffHours := time.Since(*lastSeenDate).Hours()
			diffDays := diffHours / 24
			remaining := diffDays - interval
			detail.TimeRemaining = remaining

			details = append(details, detail)
		}
	}

	return details
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
	grade               int
	section             *VocabularySection
	startingDiagnostics []*lsproto.Diagnostic
	// Document location file name)
	location string
}

// The final review detail of a word
type WordFruit struct {
	Lang  Language
	Words []*Word
	Text  string
	// Remaining time as the number of days
	TimeRemaining       float64
	StartingDiagnostics []*lsproto.Diagnostic
}

func (wb *WordTwig) GetLocation() string {
	return wb.location
}
