package compiler

import (
	"maps"
	"slices"
	"strings"
	"testing"
	"time"
	lsproto "vocab/lsp"
	test "vocab/vocab_testing"
	parser "vocab/vocabulary/parser"
)

func fakeWord(text string, grade int, parent *parser.WordsSection) *parser.Word {
	return &parser.Word{
		Line:      0,
		Text:      text,
		Literally: true,
		Start:     0,
		End:       len(text),
		Grade:     grade,
		Parent:    parent,
	}
}
func fakeTree(zeit time.Time, dateLine int, newWordDeutschSection1 string, newWordItalianoSection2 string, reviewedWordDeutschSection3 string) *WordTree {
	tree := NewWordTree()
	section := parser.NewVocabularySection("xxx")
	section.Date = &parser.DateSection{Time: zeit, Line: dateLine}
	newWordsSection1 := &parser.WordsSection{
		Parent:   section,
		Reviewed: false,
		Language: parser.Deutsch,
	}
	tree.AddTwig(parser.Deutsch, fakeWord(newWordDeutschSection1, 5, newWordsSection1), "xxx", section, []*lsproto.Diagnostic{})
	section.NewWords = append(section.NewWords, newWordsSection1)

	newWordsSection2 := &parser.WordsSection{
		Parent:   section,
		Reviewed: false,
		Language: parser.Italiano,
	}
	tree.AddTwig(parser.Italiano, fakeWord(newWordItalianoSection2, 5, newWordsSection2), "xxx", section, []*lsproto.Diagnostic{})
	section.NewWords = append(section.NewWords, newWordsSection2)

	reviewedWordsSection1 := &parser.WordsSection{
		Parent:   section,
		Reviewed: true,
		Language: parser.Deutsch,
	}
	tree.AddTwig(parser.Deutsch, fakeWord(reviewedWordDeutschSection3, 5, reviewedWordsSection1), "xxx", section, []*lsproto.Diagnostic{})
	section.ReviewedWords = append(section.ReviewedWords, reviewedWordsSection1)
	return tree
}

func TestCreatingANewTreeShouldNotThrowError(t *testing.T) {
	NewWordTree()
}

func TestAddTwigToEmptyTree(t *testing.T) {
	tree := NewWordTree()
	vocabSection := parser.NewVocabularySection("xxx")
	vocabSection.Date = &parser.DateSection{Time: time.Now()}
	wordSection := &parser.WordsSection{
		Parent: vocabSection,
	}
	vocabSection.NewWords = append(vocabSection.NewWords, wordSection)
	word := fakeWord("Testen", 5, wordSection)
	startingDiags := []*lsproto.Diagnostic{
		lsproto.MakeDiagnostics("test diagnostics", 1, 2, 3, lsproto.DiagnosticsSeverityError),
	}

	//act
	tree.AddTwig(parser.Deutsch, word, "xxx", vocabSection, startingDiags)

	wordNormalized := strings.ToLower(word.Text)
	test.Expect(t, true, tree.branches[string(parser.Deutsch)] != nil)
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
	section1 := parser.NewVocabularySection("xxx")
	section1.Date = &parser.DateSection{Time: time.Now()}
	newWordSection := &parser.WordsSection{
		Parent:   section1,
		Reviewed: false,
		Language: parser.Deutsch,
	}
	section1.NewWords = append(section1.NewWords, newWordSection)
	section2 := parser.NewVocabularySection("xxx")
	section2.Date = &parser.DateSection{Time: time.Now()}
	reviewedWordSection := &parser.WordsSection{
		Parent:   section1,
		Reviewed: true,
		Language: parser.Italiano,
	}
	section2.ReviewedWords = append(section2.ReviewedWords, reviewedWordSection)

	//act
	tree.AddTwig(parser.Deutsch, fakeWord("Testen", 5, newWordSection), "xxx", section1, []*lsproto.Diagnostic{})
	tree.AddTwig(parser.Italiano, fakeWord("Test", 3, newWordSection), "xxx", section1, []*lsproto.Diagnostic{})
	tree.AddTwig(parser.Deutsch, fakeWord("testen", 4, reviewedWordSection), "xxx", section2, []*lsproto.Diagnostic{})

	test.Expect(t, true, tree.branches[string(parser.Deutsch)] != nil)
	test.Expect(t, true, tree.branches[string(parser.Italiano)] != nil)
	branches := slices.Collect(maps.Values(tree.branches))
	test.Expect(t, 2, len(branches))

	// german
	derAst := tree.branches[string(parser.Deutsch)]
	test.Expect(t, 2, len(derAst.twigs["testen"]))
	test.Expect(t, 5, derAst.twigs["testen"][0].grade)
	test.Expect(t, section1, derAst.twigs["testen"][0].section)
	test.Expect(t, 4, derAst.twigs["testen"][1].grade)
	test.Expect(t, section2, derAst.twigs["testen"][1].section)

	// italian
	unRamo := tree.branches[string(parser.Italiano)]
	test.Expect(t, 1, len(unRamo.twigs["test"]))
	test.Expect(t, 3, unRamo.twigs["test"][0].grade)
	test.Expect(t, section1, unRamo.twigs["test"][0].section)
}

func TestAddedTwigsShouldBeSorted(t *testing.T) {
	tree := NewWordTree()
	sectionFromTime := func(time time.Time) *parser.VocabularySection {
		dateText := time.Format("2006-01-02")
		section := &parser.VocabularySection{
			Date: &parser.DateSection{Time: time, Text: dateText},
			Uri:  "xxx",
		}
		wordsSection := &parser.WordsSection{
			Parent:   section,
			Reviewed: false,
			Language: parser.Deutsch,
		}
		section.NewWords = append(section.NewWords, wordsSection)
		return section
	}
	now := time.Now()
	nowSection := sectionFromTime(now)
	yesterdaySection := sectionFromTime(now.AddDate(0, 0, -1))
	tomorrowSection := sectionFromTime(now.AddDate(0, 0, 1))

	//act
	tree.AddTwig(parser.Deutsch,
		fakeWord("poopy", 5, nowSection.NewWords[0]),
		"xxx",
		nowSection,
		[]*lsproto.Diagnostic{},
	)
	tree.AddTwig(parser.Deutsch,
		fakeWord("poopy", 5, tomorrowSection.NewWords[0]),
		"xxx",
		tomorrowSection,
		[]*lsproto.Diagnostic{},
	)
	tree.AddTwig(parser.Deutsch,
		fakeWord("poopy", 5, yesterdaySection.NewWords[0]),
		"xxx",
		yesterdaySection,
		[]*lsproto.Diagnostic{},
	)

	äste := tree.branches[string(parser.Deutsch)]
	test.Expect(t, 1, len(äste.twigs))
	zweige := äste.twigs["poopy"]
	test.Expect(t, yesterdaySection, zweige[0].section)
	test.Expect(t, nowSection, zweige[1].section)
	test.Expect(t, tomorrowSection, zweige[2].section)
}

func TestAddTwigsWithInvalidGrade_ShouldProduceExtraDiagnosticsError(t *testing.T) {
	tree := NewWordTree()
	section := parser.NewVocabularySection("xxx")
	section.Date = &parser.DateSection{Time: time.Now()}
	newWordSection := &parser.WordsSection{
		Parent:   section,
		Reviewed: false,
		Language: parser.Italiano,
	}
	section.NewWords = append(section.NewWords, newWordSection)

	//act
	tree.AddTwig(parser.Deutsch,
		fakeWord("ding", -1, section.NewWords[0]),
		"xxx",
		section,
		[]*lsproto.Diagnostic{},
	)
	tree.AddTwig(parser.Deutsch,
		fakeWord("cosa", 6, section.NewWords[0]),
		"xxx",
		section,
		[]*lsproto.Diagnostic{
			lsproto.MakeDiagnostics("test diagnostics", 1, 2, 3, lsproto.DiagnosticsSeverityError),
		},
	)

	äste := tree.branches[string(parser.Deutsch)]
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

	germanBranch := mergedTree.branches[string(parser.Deutsch)]

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

	italianBranch := mergedTree.branches[string(parser.Italiano)]

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

	test.Expect(t, 1, len(mergedTree.branches[string(parser.Deutsch)].twigs["ein"]))
	test.Expect(t, 1, len(mergedTree.branches[string(parser.Italiano)].twigs["due"]))
	test.Expect(t, 1, len(mergedTree.branches[string(parser.Deutsch)].twigs["drei"]))
}

func TestHarvest_ShouldOutputWordFruitsWithCorrectInterval(t *testing.T) {
	now := time.Now()
	lastWeek := now.AddDate(0, 0, -7)
	tomorrow := now.AddDate(0, 0, 1)

	tree1 := fakeTree(now, 0, "halo", "ciao", "halo")
	tree2 := fakeTree(lastWeek, 5, "schon", "già", "schon")
	tree3 := fakeTree(tomorrow, 10, "oft", "spesso", "oft")
	merged := tree1.Graft(tree2).Graft(tree3)

	harvested := merged.Harvest()
	test.Expect(t, 6, len(harvested))
	print(harvested)

}
