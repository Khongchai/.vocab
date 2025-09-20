package vocabulary

import (
	"context"
	lsproto "vocab/lsp"
)

type DiagnosticsArray []lsproto.Diagnostic
type DiagnosticsMap map[string]DiagnosticsArray

// The program design is to allow fast lookup of word: "Given a word, is it time to review this?"
//
// `Entries` is therefore a hash map of words to a linked list of sorted dates.
type Compiler struct {
	ctx         context.Context
	diagnostics DiagnosticsMap
	Asts        []*VocabAst
	Entries     map[string][]*VocabularySection
}

func NewCompiler(ctx context.Context) *Compiler {
	return &Compiler{
		ctx:  ctx,
		Asts: []*VocabAst{},
	}
}

func (p *Compiler) ParseDocument(documentUri string, text string, changeRange *lsproto.Range) {
	// ast := NewAst(p.ctx, documentUri, text, changeRange)
	// p.Asts = append(p.Asts, ast)
}

func (p *Compiler) Compile() {
	p.compile()
}

func (p *Compiler) GetDiagnostics(documentUri string) DiagnosticsArray {
	return p.diagnostics[documentUri]
}

func (p *Compiler) compile() {
	// TODO
	// This would compile the ast into diagnostics. This includes both syntactic and spaced-repetition-related errors.
}
