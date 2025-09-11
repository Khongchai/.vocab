package vocabulary

import (
	"context"
	"time"
	lsproto "vocab/lsp"
)

type Language string

const (
	Deutsch Language = "Tedesco"
	Italian Language = "Italienisch"
)

type Word struct {
	Text string
}

type Sentence struct {
	Text string
	// future positions for reviewed and such will be here.
}

type Date struct {
	Text      string
	Time      time.Time
	TextRange lsproto.Range
}

type VocabularySection struct {
	Date          *Date
	NewWords      []*Word
	ReviewedWords []*Word
	Sentences     []*Sentence
	Language      Language
}

type Document struct {
	Sections []*VocabularySection
}

type VocabAst struct {
	ctx       context.Context
	documents map[string][]*Document
	scanner   *Scanner
}

func NewAst(ctx context.Context, uri string, text string, changeRange *lsproto.Range) *VocabAst {
	if changeRange != nil {
		panic("Partial update not yet handled")
	}

	ast := &VocabAst{
		ctx:       ctx,
		documents: map[string][]*Document{},
	}

	return ast
}

// TODO: Yield diagnostics result from scanner and all of its children.
func (ast *VocabAst) GetCurrentDiagnostics(uri string) []lsproto.Diagnostic {
	scannerDiagnostics := ast.scanner.GetCurrentDiagnostics(uri)
	return []lsproto.Diagnostic{}
}
