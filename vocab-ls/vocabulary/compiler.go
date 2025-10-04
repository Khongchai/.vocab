package vocabulary

import (
	"context"
	lsproto "vocab/lsp"
)

type WordBranch struct {
	Branch              *VocabularySection
	CompiledDiagnostics *lsproto.Diagnostic
}

type WordTree struct {
	Branches map[string]*[]WordBranch
}

// The program design is to allow fast lookup of word: "Given a word, is it time to review this?"
//
// `Entries` is therefore a hash map of words to a linked list of sorted dates.
type Compiler struct {
	ctx  context.Context
	tree *WordTree
}

func NewCompiler(ctx context.Context) *Compiler {
	return &Compiler{
		ctx:  ctx,
		tree: &WordTree{},
	}
}

// Accept new document uri, turn it into a branch, put into word tree for quick lookup LATER and compile it NOW.
func (p *Compiler) Accept(documentUri string, text string, changeRange *lsproto.Range) {
	ast := NewAst(p.ctx, documentUri, text, changeRange)
	p.Asts = append(p.Asts, ast)
}

func (p *Compiler) compile() {
	p.compile()
}
