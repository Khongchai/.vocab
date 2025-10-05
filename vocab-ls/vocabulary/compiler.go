package vocabulary

import (
	"context"
	"fmt"
	lsproto "vocab/lsp"
	data "vocab/vocabulary/data"
)

// The program design is to allow fast lookup of word: "Given a word, is it time to review this?"
//
// `Entries` is therefore a hash map of words to a an array of date sections they appear in.
type Compiler struct {
	ctx  context.Context
	tree *data.WordTree
	log  func(any)
}

func NewCompiler(ctx context.Context, log func(any)) *Compiler {
	return &Compiler{
		ctx:  ctx,
		tree: &data.WordTree{},
		log:  log,
	}
}

func (c *Compiler) Accept(documentUri string, text string, changeRange *lsproto.Range) {
	scanner := NewScanner(text)
	parser := NewParser(c.ctx, documentUri, scanner, c.log)
	parser.Parse()

	newWordTree := c.astToWordTree(parser.ast)
	if c.tree == nil {
		c.tree = newWordTree
	} else {
		c.tree.Graft(newWordTree)
	}
}

func (c *Compiler) Compile() []lsproto.Diagnostic {
	return c.tree.Harvest()
}

func (c *Compiler) astToWordTree(ast *data.VocabAst) *data.WordTree {
	tree := data.NewWordTree()
	for _, section := range ast.Sections {

		// add new words
		for _, newWordSection := range section.NewWords {
			lang := newWordSection.Language
			for _, word := range newWordSection.Words {
				diag := []*lsproto.Diagnostic{}
				text := word.Text
				existingTwigs := tree.GetTwigs(lang, text)

				if len(existingTwigs) != 0 {
					message := fmt.Sprintf("This has been seen before. First occurence in %s", existingTwigs[0].GetLocation())
					diag = append(diag, &lsproto.Diagnostic{
						Message:  message,
						Severity: lsproto.DiagnosticsSeverityWarning,
						Range: lsproto.Range{
							Start: lsproto.Position{
								Line:      newWordSection.Line,
								Character: word.Start,
							},
							End: lsproto.Position{
								Line:      newWordSection.Line,
								Character: word.End,
							},
						},
					})
				}

				tree.AddTwig(lang, text, ast.Uri, section, diag)
			}
		}

		// add reviewed words
		for _, reviewedWordSection := range section.ReviewedWords {
			lang := reviewedWordSection.Language
			for _, word := range reviewedWordSection.Words {
				diag := []*lsproto.Diagnostic{}
				text := word.Text
				existingTwigs := tree.GetTwigs(lang, text)
				if len(existingTwigs) == 0 {
					message := "This has never been seen before."
					diag = append(diag, &lsproto.Diagnostic{
						Message:  message,
						Severity: lsproto.DiagnosticsSeverityWarning,
						Range: lsproto.Range{
							Start: lsproto.Position{
								Line:      reviewedWordSection.Line,
								Character: word.Start,
							},
							End: lsproto.Position{
								Line:      reviewedWordSection.Line,
								Character: word.End,
							},
						},
					})
				}

				tree.AddTwig(lang, text, ast.Uri, section, diag)
			}
		}
	}

	return tree
}
