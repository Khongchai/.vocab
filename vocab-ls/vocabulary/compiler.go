package vocabulary

import (
	"context"
	"fmt"
	lsproto "vocab/lsp"
	entity "vocab/vocabulary/entity"
)

// The program design is to allow fast lookup of word: "Given a word, is it time to review this?"
//
// `Entries` is therefore a hash map of words to a an array of date sections they appear in.
type Compiler struct {
	ctx  context.Context
	tree *entity.WordTree
	log  func(any)
}

func NewCompiler(ctx context.Context, log func(any)) *Compiler {
	return &Compiler{
		ctx:  ctx,
		tree: &entity.WordTree{},
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

func (c *Compiler) astToWordTree(ast *entity.VocabAst) *entity.WordTree {
	tree := entity.NewWordTree()

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
		allWordSections := append(section.NewWords, section.ReviewedWords...)
		for _, wordSection := range allWordSections {
			lang := wordSection.Language
			reviewed := wordSection.Reviewed
			for _, word := range wordSection.Words {
				var diag []*lsproto.Diagnostic
				existingTwigs := tree.GetTwigs(lang, word.Text)

				warnText := func() string {
					if !reviewed && len(existingTwigs) != 0 {
						return fmt.Sprintf("This has been seen before. First occurence in %s", existingTwigs[0].GetLocation())
					} else if reviewed && len(existingTwigs) == 0 {
						return "This has never been seen before"
					}
					return ""
				}()

				if warnText != "" {
					warn := diagnosticsForWord(
						warnText,
						wordSection.Line,
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
