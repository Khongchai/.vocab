package entity

import (
	"fmt"
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
	// grade parsed after word -> word(5)
	Grade int
}

type Section interface {
	Identity() string
}

type UtteranceSection struct {
	Line  int
	Text  string
	Start int
	End   int
}

func (d *UtteranceSection) Identity() string {
	return fmt.Sprintf("%s::%d", d.Text, d.Line)
}

type DateSection struct {
	Line  int
	Text  string
	Time  time.Time
	Start int
	End   int
}

func (d *DateSection) Identity() string {
	return fmt.Sprintf("%s::%d", d.Text, d.Line)
}

type WordsSection struct {
	Words    []*Word
	Reviewed bool
	Language Language
	Line     int
}

func (w *WordsSection) Identity() string {
	return fmt.Sprintf("%s::%d", string(w.Language), w.Line)
}

type VocabularySection struct {
	Date          *DateSection
	NewWords      []*WordsSection
	ReviewedWords []*WordsSection
	Utterance     []*UtteranceSection
	Diagnostics   []*lsproto.Diagnostic
	Uri           string
}

func NewVocabularySection(uri string) *VocabularySection {
	return &VocabularySection{Uri: uri}
}

func (v *VocabularySection) Identity() string {
	return fmt.Sprintf("%s::%s", v.Uri, v.Date.Identity())
}

type VocabAst struct {
	// Might make this an array later, we'll see
	Sections []*VocabularySection
	Uri      string
}
