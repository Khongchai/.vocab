package vocabulary

import (
	"context"
	"time"
	lsproto "vocab/lsp"
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
}

type VocabAst struct {
	ctx       context.Context
	uri       string
	documents *Document
	scanner   *Scanner
}

func NewAst(ctx context.Context, uri string, text string, changeRange *lsproto.Range) *VocabAst {
	if changeRange != nil {
		panic("Partial update not yet handled")
	}

	ast := &VocabAst{
		ctx:       ctx,
		documents: &Document{},
	}

	ast.scanner = NewScanner(uri)

	return ast
}

// TODO: Yield diagnostics result from scanner and all of its children.
func (ast *VocabAst) GetCurrentDiagnostics(uri string) []lsproto.Diagnostic {
	return []lsproto.Diagnostic{}
}
