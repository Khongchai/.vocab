package vocabulary

import (
	"time"
	lsproto "vocab/lsp"
)

type Language string

const (
	Unrecognized Language = ""
	Deutsch      Language = "Tedesco"
	Italiano     Language = "Italienisch"
)

type Word struct {
	// Text represent the actual string value of a word with or without its article.
	Text string
	// Literally means the text was wrapped with backticks. The compiler may choose to do something
	// differently with this information.
	Literally bool
	// the start of Text
	// "hello" start = 0
	Start int
	// the end of Text
	// "hello" end = 4
	End int
}

type Section interface {
	SectionName() string
}

type UtteranceSection struct {
	Line  int
	Text  string
	Start int
	End   int
}

func (d *UtteranceSection) SectionName() string { return "Utterance" }

type DateSection struct {
	Line  int
	Text  string
	Time  time.Time
	Start int
	End   int
}

func (d *DateSection) SectionName() string { return "Date" }

type WordsSection struct {
	Words    []*Word
	Language Language
	Line     int
}

func (w *WordsSection) SectionName() string { return "Words" }

type VocabularySection struct {
	Date          *DateSection
	NewWords      []*WordsSection
	ReviewedWords []*WordsSection
	Utterance     []*UtteranceSection
	Diagnostics   []*lsproto.Diagnostic
}

func (v *VocabularySection) SectionName() string { return "Vocabulary" }

type VocabAst struct {
	// Might make this an array later, we'll see
	Sections []*VocabularySection
	uri      string
}
