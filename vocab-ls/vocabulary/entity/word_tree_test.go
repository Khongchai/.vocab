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
func fakeTree(zeit time.Time, dateLine int, word1 string, word2 string, word3 string) *WordTree {
	tree := NewWordTree()
	section := NewVocabularySection("xxx")
	section.Date = &DateSection{Time: zeit, Line: dateLine}
	newWordsSection1 := &WordsSection{
		Parent:   section,
		Reviewed: false,
		Language: Deutsch,
	}
	tree.AddTwig(Deutsch, fakeWord(word1, 5, newWordsSection1), "xxx", section, []*lsproto.Diagnostic{})
	section.NewWords = append(section.NewWords, newWordsSection1)

	newWordsSection2 := &WordsSection{
		Parent:   section,
		Reviewed: false,
		Language: Italiano,
	}
	tree.AddTwig(Italiano, fakeWord(word2, 5, newWordsSection2), "xxx", section, []*lsproto.Diagnostic{})
	section.NewWords = append(section.NewWords, newWordsSection2)

	reviewedWordsSection1 := &WordsSection{
		Parent:   section,
		Reviewed: true,
		Language: Deutsch,
	}
	tree.AddTwig(Deutsch, fakeWord(word3, 5, reviewedWordsSection1), "xxx", section, []*lsproto.Diagnostic{})
	section.ReviewedWords = append(section.ReviewedWords, reviewedWordsSection1)
	return tree
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

func TestGraftingTrees_ShouldCombineCorrectLanguageBranch_AndSortTwigs(t *testing.T) {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)

	// Resulting tree:
	// date:now
	// > (de) ein, gestern
	// > (it) due, oggi
	// >> (de) drei, morgen
	// date:tomorrow
	// > (de) links
	// > (it) dritto
	// >> (de) ein
	tree1 := fakeTree(now, 0, "ein", "due", "drei")
	tree2 := fakeTree(tomorrow, 10, "links", "dritto", "ein")
	tree3 := fakeTree(now, 5, "gestern", "oggi", "morgen")

	//act
	mergedTree := NewWordTree().Graft(tree1).Graft(tree2).Graft(tree3)

	test.Expect(t, true, mergedTree != nil)

	germanBranch := mergedTree.branches[string(Deutsch)]

	einTwigs := germanBranch.twigs["ein"]
	test.Expect(t, 2, len(einTwigs))

	test.Expect(t, now.YearDay(), einTwigs[0].section.Date.Time.YearDay())
	test.Expect(t, 0, einTwigs[0].section.Date.Line)
	test.Expect(t, tomorrow.YearDay(), einTwigs[1].section.Date.Time.YearDay())
	test.Expect(t, 10, einTwigs[1].section.Date.Line)

	dreiTwigs := germanBranch.twigs["drei"]
	test.Expect(t, 1, len(dreiTwigs))
	test.Expect(t, now.YearDay(), dreiTwigs[0].section.Date.Time.YearDay())
	test.Expect(t, 0, dreiTwigs[0].section.Date.Line)

	gesternTwigs := germanBranch.twigs["gestern"]
	test.Expect(t, 1, len(gesternTwigs))
	test.Expect(t, now.YearDay(), gesternTwigs[0].section.Date.Time.YearDay())
	test.Expect(t, 5, gesternTwigs[0].section.Date.Line)

	morgenTwigs := germanBranch.twigs["morgen"]
	test.Expect(t, 1, len(morgenTwigs))
	test.Expect(t, now.YearDay(), morgenTwigs[0].section.Date.Time.YearDay())
	test.Expect(t, 5, morgenTwigs[0].section.Date.Line)

	linksTwigs := germanBranch.twigs["links"]
	test.Expect(t, 1, len(linksTwigs))
	test.Expect(t, tomorrow.YearDay(), linksTwigs[0].section.Date.Time.YearDay())
	test.Expect(t, 10, linksTwigs[0].section.Date.Line)

	italianBranch := mergedTree.branches[string(Italiano)]

	dueTwigs := italianBranch.twigs["due"]
	test.Expect(t, 1, len(dueTwigs))
	test.Expect(t, now.YearDay(), dueTwigs[0].section.Date.Time.YearDay())
	test.Expect(t, 0, dueTwigs[0].section.Date.Line)

	oggiTwigs := italianBranch.twigs["oggi"]
	test.Expect(t, 1, len(oggiTwigs))
	test.Expect(t, now.YearDay(), oggiTwigs[0].section.Date.Time.YearDay())
	test.Expect(t, 5, oggiTwigs[0].section.Date.Line)

	drittoTwigs := italianBranch.twigs["dritto"]
	test.Expect(t, 1, len(drittoTwigs))
	test.Expect(t, tomorrow.YearDay(), drittoTwigs[0].section.Date.Time.YearDay())
	test.Expect(t, 10, drittoTwigs[0].section.Date.Line)
}

func TestGraftingTrees_ShouldNotRecombineTreesWithSameIdentity(t *testing.T) {

	now := time.Now()

	tree1 := fakeTree(now, 0, "ein", "due", "drei")
	tree2 := fakeTree(now, 0, "ein", "due", "drei")

	//act
	mergedTree := NewWordTree().Graft(tree1).Graft(tree2)

	test.Expect(t, true, mergedTree != nil)

	test.Expect(t, 1, len(mergedTree.branches[string(Deutsch)].twigs["ein"]))
	test.Expect(t, 1, len(mergedTree.branches[string(Italiano)].twigs["due"]))
	test.Expect(t, 1, len(mergedTree.branches[string(Deutsch)].twigs["drei"]))
}

func TestHarvest_(t *testing.T) {
	// t = today
	// given
	// t-10 (3 words)
	// t-5 (2 words)
	// expect number of days output by WordFruit
}
