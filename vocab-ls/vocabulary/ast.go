package vocabulary

import (
	"context"
	"time"
	lsproto "vocab/lsp"
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
}

type Document struct {
	Sections []*VocabularySection
}

type VocabAst struct {
	ctx       context.Context
	documents map[string][]*Document
}

func NewAst(ctx context.Context) *VocabAst {
	return &VocabAst{
		ctx:       ctx,
		documents: map[string][]*Document{},
	}
}

func (ast *VocabAst) Update(uri string, text string, changeRange *lsproto.Range) {
	if changeRange != nil {
		panic("Partial update not yet handled")
	}
}

func (ast *VocabAst) GetCurrentDiagnostics(uri string) []lsproto.Diagnostic {
	return []lsproto.Diagnostic{}
}
