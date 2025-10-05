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

	diagnosticsForWord := func(message string, line, start, end int) *lsproto.Diagnostic {
		diag := &lsproto.Diagnostic{
			Message:  message,
			Severity: lsproto.DiagnosticsSeverityWarning,
			Range: lsproto.Range{
				Start: lsproto.Position{
					Line:      line,
					Character: start,
				},
				End: lsproto.Position{
					Line:      line,
					Character: end,
				},
			},
		}
		return diag
	}

	for _, section := range ast.Sections {
		// add new words
		for _, newWordSection := range section.NewWords {
			lang := newWordSection.Language
			for _, word := range newWordSection.Words {
				var diag []*lsproto.Diagnostic
				existingTwigs := tree.GetTwigs(lang, word.Text)

				if len(existingTwigs) != 0 {
					warn := diagnosticsForWord(
						fmt.Sprintf("This has been seen before. First occurence in %s", existingTwigs[0].GetLocation()),
						newWordSection.Line,
						word.Start,
						word.End,
					)
					diag = append(diag, warn)
				}

				tree.AddTwig(lang, word.Text, ast.Uri, section, diag)
			}
		}

		for _, reviewedWordSection := range section.ReviewedWords {
			lang := reviewedWordSection.Language
			for _, word := range reviewedWordSection.Words {
				var diag []*lsproto.Diagnostic
				existingTwigs := tree.GetTwigs(lang, word.Text)

				if len(existingTwigs) == 0 {
					warn := diagnosticsForWord(
						"This has never been seen before",
						reviewedWordSection.Line,
						word.Start,
						word.End,
					)
					diag = append(diag, warn)
				}

				tree.AddTwig(lang, word.Text, ast.Uri, section, diag)
			}
		}
	}

	return tree
}
