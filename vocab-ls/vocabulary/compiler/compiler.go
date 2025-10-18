package compiler

import (
	"context"
	"fmt"
	lsproto "vocab/lsp"
	"vocab/vocabulary/parser"
)

type Compiler struct {
	ctx  context.Context
	tree *WordTree
	log  func(any)
}

func NewCompiler(ctx context.Context, log func(any)) *Compiler {
	return &Compiler{
		ctx:  ctx,
		tree: &WordTree{},
		log:  log,
	}
}

func (c *Compiler) Accept(documentUri string, text string, changeRange *lsproto.Range) {
	scanner := parser.NewScanner(text)
	parser := parser.NewParser(c.ctx, documentUri, scanner, c.log)
	parser.Parse()

	newWordTree := c.astToWordTree(parser.Ast)
	if c.tree == nil {
		c.tree = newWordTree
	} else {
		c.tree.Graft(newWordTree)
	}
}

func (c *Compiler) Compile() []lsproto.Diagnostic {
	details := c.tree.Harvest()
	diags := []lsproto.Diagnostic{}

	addDiagToWords := func(timeRemaining float64, severitiy lsproto.DiagnosticsSeverity, words []*parser.Word) {
		for _, word := range words {
			err := lsproto.MakeDiagnostics(
				fmt.Sprintf("Needs review within: %f day(s)", timeRemaining),
				word.Line,
				word.Start,
				word.End,
				severitiy,
			)

			diags = append(diags, *err)
		}
	}

	for _, detail := range details {
		for _, starting := range detail.StartingDiagnostics {
			diags = append(diags, *starting)
		}

		severity := func() lsproto.DiagnosticsSeverity {
			if detail.TimeRemaining <= 1 {
				return lsproto.DiagnosticsSeverityError
			} else if detail.TimeRemaining < 3 {
				return lsproto.DiagnosticsSeverityHint
			}

			return lsproto.DiagnosticsSeverityInformation
		}()

		addDiagToWords(detail.TimeRemaining, severity, detail.Words)
	}

	return diags
}

func (c *Compiler) astToWordTree(ast *parser.VocabAst) *WordTree {
	tree := NewWordTree()

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
					warn := lsproto.MakeDiagnostics(
						warnText,
						wordSection.Line,
						word.Start,
						word.End,
						lsproto.DiagnosticsSeverityWarning,
					)
					diag = append(diag, warn)
				}

				tree.AddTwig(lang, word, ast.Uri, section, diag)
			}
		}
	}
	return tree
}
