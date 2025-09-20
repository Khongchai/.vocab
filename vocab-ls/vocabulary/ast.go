package vocabulary

import (
	"time"
)

type Language string

const (
	Deutsch  Language = "Tedesco"
	Italiano Language = "Italienisch"
)

type Word struct {
	// Text represent the actual string value of a word without its article. For example, l'uccello should be normalized to uccello
	Text string
	// Full text is the full content including its article.
	FullText string
	// the start of FullText
	// "hello" start = 0
	Start int
	// the end of FullText
	// "hello" end = 4
	End int
}

type Sentence struct {
	StartLine int
	EndLine   int
	StartPos  int
	EndPos    int
	Text      string
	// future positions for reviewed and such will be here.
}

type Date struct {
	Text  string
	Time  time.Time
	Start int
	End   int
}

type WordsSection struct {
	Words    []*Word
	Language Language
	Line     int
}

type VocabularySection struct {
	Date          *Date
	NewWords      []*WordsSection
	ReviewedWords []*WordsSection
	Sentences     []*Sentence
}

type Document struct {
	Sections []*VocabularySection
	uri      string
}

// func (d *Document) String() string {
// 	var sb strings.Builder
// 	for i, section := range d.Sections {
// 		sectionString := section.String()
// 		sb.WriteString(sectionString)
// 		if i != len(d.Sections)-1 {
// 			sb.WriteString(", ")
// 		}
// 	}
// 	return sb.String()
// }

type VocabAst struct {
	Documents []*Document
}
