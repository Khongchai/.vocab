package vocabulary

import (
	"context"
	lsproto "vocab/lsp"
)

type DiagnosticsArray []lsproto.Diagnostic
type DiagnosticsMap map[string]DiagnosticsArray

type Program struct {
	ctx         context.Context
	diagnostics DiagnosticsMap
	Asts        []*VocabAst
}

func NewProgram(ctx context.Context) *Program {
	return &Program{
		ctx:  ctx,
		Asts: []*VocabAst{},
	}
}

func (p *Program) BuildAst(documentUri string, text string, changeRange *lsproto.Range) {
	ast := NewAst(p.ctx, documentUri, text, changeRange)
	p.Asts = append(p.Asts, ast)
}

func (p *Program) Compile() {
	p.compile()
}

func (p *Program) GetDiagnostics(documentUri string) DiagnosticsArray {
	return p.diagnostics[documentUri]
}

func (p *Program) compile() {
	// TODO
	// This would compile the ast into diagnostics. This includes both syntactic and spaced-repetition-related errors.
}
